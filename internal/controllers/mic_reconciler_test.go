package controllers

import (
	"context"
	"errors"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	kmmv1beta1 "github.com/rh-ecosystem-edge/kernel-module-management/api/v1beta1"
	"github.com/rh-ecosystem-edge/kernel-module-management/internal/client"
	"github.com/rh-ecosystem-edge/kernel-module-management/internal/mbsc"
	"github.com/rh-ecosystem-edge/kernel-module-management/internal/mic"
	"github.com/rh-ecosystem-edge/kernel-module-management/internal/pod"
	"go.uber.org/mock/gomock"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var _ = Describe("MicReconciler_Reconcile", func() {
	var (
		ctrl               *gomock.Controller
		mockImagePuller    *pod.MockImagePuller
		mockMicReconHelper *MockmicReconcilerHelper
		mr                 *micReconciler
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockImagePuller = pod.NewMockImagePuller(ctrl)
		mockMicReconHelper = NewMockmicReconcilerHelper(ctrl)

		mr = &micReconciler{
			micReconHelper: mockMicReconHelper,
			imagePullerAPI: mockImagePuller,
		}
	})

	ctx := context.Background()
	testMic := kmmv1beta1.ModuleImagesConfig{}

	DescribeTable("check good and error flows", func(listPullPodsError,
		updateStatusByPodsError,
		updateStatusByMBSCError,
		processImagesSpecsError bool) {

		returnedError := errors.New("some error")
		expectedErr := returnedError
		pullPods := []v1.Pod{}
		if listPullPodsError {
			mockImagePuller.EXPECT().ListPullPods(ctx, &testMic).Return(nil, returnedError)
			goto executeTestFunction
		}
		mockImagePuller.EXPECT().ListPullPods(ctx, &testMic).Return(pullPods, nil)
		if updateStatusByPodsError {
			mockMicReconHelper.EXPECT().updateStatusByPullPods(ctx, &testMic, pullPods).Return(returnedError)
			goto executeTestFunction
		}
		mockMicReconHelper.EXPECT().updateStatusByPullPods(ctx, &testMic, pullPods).Return(nil)
		if updateStatusByMBSCError {
			mockMicReconHelper.EXPECT().updateStatusByMBSC(ctx, &testMic).Return(returnedError)
			goto executeTestFunction
		}
		mockMicReconHelper.EXPECT().updateStatusByMBSC(ctx, &testMic).Return(nil)
		if processImagesSpecsError {
			mockMicReconHelper.EXPECT().processImagesSpecs(ctx, &testMic, pullPods).Return(returnedError)
			goto executeTestFunction
		}
		mockMicReconHelper.EXPECT().processImagesSpecs(ctx, &testMic, pullPods).Return(nil)
		expectedErr = nil

	executeTestFunction:
		res, err := mr.Reconcile(ctx, &testMic)

		Expect(res).To(Equal(reconcile.Result{}))
		if expectedErr != nil {
			Expect(err).To(HaveOccurred())
		} else {
			Expect(err).To(BeNil())
		}
	},
		Entry("listPullPods failed", true, false, false, false),
		Entry("updateStatusByPullPods failed", false, true, false, false),
		Entry("updateStatusByMBSC failed", false, false, true, false),
		Entry("processImagesSpecs failed", false, false, false, true),
		Entry("everything worked", false, false, false, false),
	)
})

var _ = Describe("updateStatusByPullPods", func() {
	var (
		ctrl            *gomock.Controller
		clnt            *client.MockClient
		statusWriter    *client.MockStatusWriter
		mockImagePuller *pod.MockImagePuller
		micHelper       *mic.MockMIC
		mrh             micReconcilerHelper
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		clnt = client.NewMockClient(ctrl)
		statusWriter = client.NewMockStatusWriter(ctrl)
		mockImagePuller = pod.NewMockImagePuller(ctrl)
		micHelper = mic.NewMockMIC(ctrl)
		mrh = newMICReconcilerHelper(clnt, mockImagePuller, micHelper, nil, nil)
	})

	ctx := context.Background()
	testMic := kmmv1beta1.ModuleImagesConfig{
		Spec: kmmv1beta1.ModuleImagesConfigSpec{
			Images: []kmmv1beta1.ModuleImageSpec{
				{
					Image: "image 1",
					Build: &kmmv1beta1.Build{},
				},
				{

					Image: "image 2",
					Sign:  &kmmv1beta1.Sign{},
				},
				{

					Image: "image 3",
				},
			},
		},
	}

	It("zero pull pods", func() {
		pullPods := []v1.Pod{}
		err := mrh.updateStatusByPullPods(ctx, &testMic, pullPods)
		Expect(err).To(BeNil())
	})

	It("pod's image is not in spec", func() {
		pullPod := v1.Pod{}
		gomock.InOrder(
			mockImagePuller.EXPECT().GetPullPodImage(pullPod).Return("image 3"),
			micHelper.EXPECT().GetModuleImageSpec(&testMic, "image 3").Return(nil),
			clnt.EXPECT().Status().Return(statusWriter),
			statusWriter.EXPECT().Patch(ctx, &testMic, gomock.Any()),
			mockImagePuller.EXPECT().DeletePod(ctx, &pullPod).Return(nil),
		)
		err := mrh.updateStatusByPullPods(ctx, &testMic, []v1.Pod{pullPod})
		Expect(err).To(BeNil())
	})

	It("pod failed, build config present", func() {
		pullPod := v1.Pod{
			Status: v1.PodStatus{
				Phase: v1.PodFailed,
			},
		}
		gomock.InOrder(
			mockImagePuller.EXPECT().GetPullPodImage(pullPod).Return("image 1"),
			micHelper.EXPECT().GetModuleImageSpec(&testMic, "image 1").Return(&testMic.Spec.Images[0]),
			micHelper.EXPECT().SetImageStatus(&testMic, "image 1", kmmv1beta1.ImageNeedsBuilding),
			clnt.EXPECT().Status().Return(statusWriter),
			statusWriter.EXPECT().Patch(ctx, &testMic, gomock.Any()),
			mockImagePuller.EXPECT().DeletePod(ctx, &pullPod).Return(nil),
		)
		err := mrh.updateStatusByPullPods(ctx, &testMic, []v1.Pod{pullPod})
		Expect(err).To(BeNil())
	})

	It("pod failed, build config not present, sign config present", func() {
		pullPod := v1.Pod{
			Status: v1.PodStatus{
				Phase: v1.PodFailed,
			},
		}
		gomock.InOrder(
			mockImagePuller.EXPECT().GetPullPodImage(pullPod).Return("image 2"),
			micHelper.EXPECT().GetModuleImageSpec(&testMic, "image 2").Return(&testMic.Spec.Images[1]),
			micHelper.EXPECT().SetImageStatus(&testMic, "image 2", kmmv1beta1.ImageNeedsSigning),
			clnt.EXPECT().Status().Return(statusWriter),
			statusWriter.EXPECT().Patch(ctx, &testMic, gomock.Any()),
			mockImagePuller.EXPECT().DeletePod(ctx, &pullPod).Return(nil),
		)
		err := mrh.updateStatusByPullPods(ctx, &testMic, []v1.Pod{pullPod})
		Expect(err).To(BeNil())
	})

	It("pod failed, build or sign config is not present", func() {
		pullPod := v1.Pod{
			Status: v1.PodStatus{
				Phase: v1.PodFailed,
			},
		}
		gomock.InOrder(
			mockImagePuller.EXPECT().GetPullPodImage(pullPod).Return("image 3"),
			micHelper.EXPECT().GetModuleImageSpec(&testMic, "image 3").Return(&testMic.Spec.Images[2]),
			clnt.EXPECT().Status().Return(statusWriter),
			statusWriter.EXPECT().Patch(ctx, &testMic, gomock.Any()),
			mockImagePuller.EXPECT().DeletePod(ctx, &pullPod).Return(nil),
		)
		err := mrh.updateStatusByPullPods(ctx, &testMic, []v1.Pod{pullPod})
		Expect(err).To(BeNil())
	})

	It("pod succeeded", func() {
		pullPod := v1.Pod{
			Status: v1.PodStatus{
				Phase: v1.PodSucceeded,
			},
		}
		gomock.InOrder(
			mockImagePuller.EXPECT().GetPullPodImage(pullPod).Return("image 2"),
			micHelper.EXPECT().GetModuleImageSpec(&testMic, "image 2").Return(&testMic.Spec.Images[1]),
			micHelper.EXPECT().SetImageStatus(&testMic, "image 2", kmmv1beta1.ImageExists),
			clnt.EXPECT().Status().Return(statusWriter),
			statusWriter.EXPECT().Patch(ctx, &testMic, gomock.Any()),
			mockImagePuller.EXPECT().DeletePod(ctx, &pullPod).Return(nil),
		)
		err := mrh.updateStatusByPullPods(ctx, &testMic, []v1.Pod{pullPod})
		Expect(err).To(BeNil())
	})
})

var _ = Describe("updateStatusByMBSC", func() {
	var (
		ctrl         *gomock.Controller
		clnt         *client.MockClient
		statusWriter *client.MockStatusWriter
		mbscHelper   *mbsc.MockMBSC
		micHelper    *mic.MockMIC
		mrh          micReconcilerHelper
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		clnt = client.NewMockClient(ctrl)
		statusWriter = client.NewMockStatusWriter(ctrl)
		micHelper = mic.NewMockMIC(ctrl)
		mbscHelper = mbsc.NewMockMBSC(ctrl)
		mrh = newMICReconcilerHelper(clnt, nil, micHelper, mbscHelper, nil)
	})

	ctx := context.Background()
	testMic := kmmv1beta1.ModuleImagesConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "some name",
			Namespace: "some namespace",
		},
		Spec: kmmv1beta1.ModuleImagesConfigSpec{
			Images: []kmmv1beta1.ModuleImageSpec{
				{
					Image: "image 1",
					Build: &kmmv1beta1.Build{},
				},
				{

					Image: "image 2",
				},
			},
		},
	}

	It("failed to get MBSC", func() {
		mbscHelper.EXPECT().Get(ctx, testMic.Name, testMic.Namespace).Return(nil, fmt.Errorf("some error"))
		err := mrh.updateStatusByMBSC(ctx, &testMic)
		Expect(err).To(HaveOccurred())
	})

	It("MBSC does not exists", func() {
		mbscHelper.EXPECT().Get(ctx, testMic.Name, testMic.Namespace).Return(nil, nil)
		err := mrh.updateStatusByMBSC(ctx, &testMic)
		Expect(err).To(BeNil())
	})

	It("Image in MBSC status does not exists in MIC spec", func() {
		testMBSC := kmmv1beta1.ModuleBuildSignConfig{
			Status: kmmv1beta1.ModuleImagesConfigStatus{
				ImagesStates: []kmmv1beta1.ModuleImageState{
					{
						Image: "some image",
					},
				},
			},
		}
		gomock.InOrder(
			mbscHelper.EXPECT().Get(ctx, testMic.Name, testMic.Namespace).Return(&testMBSC, nil),
			micHelper.EXPECT().GetModuleImageSpec(&testMic, testMBSC.Status.ImagesStates[0].Image).Return(nil),
			clnt.EXPECT().Status().Return(statusWriter),
			statusWriter.EXPECT().Patch(ctx, &testMic, gomock.Any()),
		)
		err := mrh.updateStatusByMBSC(ctx, &testMic)
		Expect(err).To(BeNil())
	})

	DescribeTable("images in MBSC status exist in MIC spec",
		func(signExists bool, mbscImageStatus kmmv1beta1.ImageState, expectedImageState kmmv1beta1.ImageState) {
			testMBSC := kmmv1beta1.ModuleBuildSignConfig{
				Status: kmmv1beta1.ModuleImagesConfigStatus{
					ImagesStates: []kmmv1beta1.ModuleImageState{
						{
							Image:  "some image",
							Status: mbscImageStatus,
						},
					},
				},
			}
			imageSpec := kmmv1beta1.ModuleImageSpec{}
			if signExists {
				imageSpec.Sign = &kmmv1beta1.Sign{}
			}
			gomock.InOrder(
				mbscHelper.EXPECT().Get(ctx, testMic.Name, testMic.Namespace).Return(&testMBSC, nil),
				micHelper.EXPECT().GetModuleImageSpec(&testMic, "some image").Return(&imageSpec),
				micHelper.EXPECT().SetImageStatus(&testMic, "some image", expectedImageState),
				clnt.EXPECT().Status().Return(statusWriter),
				statusWriter.EXPECT().Patch(ctx, &testMic, gomock.Any()),
			)

			err := mrh.updateStatusByMBSC(ctx, &testMic)
			Expect(err).To(BeNil())
		},
		Entry("sign config does not exists, status is ImageBuildFailed", false, kmmv1beta1.ImageBuildFailed, kmmv1beta1.ImageDoesNotExist),
		Entry("sign config does not exists, status is ImageBuildSucceeded", false, kmmv1beta1.ImageBuildSucceeded, kmmv1beta1.ImageExists),
		Entry("sign config exists, status is ImageBuildFailed", true, kmmv1beta1.ImageBuildFailed, kmmv1beta1.ImageDoesNotExist),
		Entry("sign config exists, status is ImageBuildSucceeded", true, kmmv1beta1.ImageBuildSucceeded, kmmv1beta1.ImageNeedsSigning),
		Entry("status is ImageSignFailed", false, kmmv1beta1.ImageSignFailed, kmmv1beta1.ImageDoesNotExist),
		Entry("status is ImageSignSucceeded", false, kmmv1beta1.ImageSignSucceeded, kmmv1beta1.ImageExists),
	)
})

var _ = Describe("processImagesSpecs", func() {
	var (
		ctrl            *gomock.Controller
		clnt            *client.MockClient
		mockImagePuller *pod.MockImagePuller
		mbscHelper      *mbsc.MockMBSC
		micHelper       *mic.MockMIC
		mrh             micReconcilerHelper
		testMic         kmmv1beta1.ModuleImagesConfig
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		clnt = client.NewMockClient(ctrl)
		mockImagePuller = pod.NewMockImagePuller(ctrl)
		micHelper = mic.NewMockMIC(ctrl)
		mbscHelper = mbsc.NewMockMBSC(ctrl)
		mrh = newMICReconcilerHelper(clnt, mockImagePuller, micHelper, mbscHelper, scheme)
		testMic = kmmv1beta1.ModuleImagesConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "some name",
				Namespace: "some namespace",
			},
			Spec: kmmv1beta1.ModuleImagesConfigSpec{
				Images: []kmmv1beta1.ModuleImageSpec{
					{
						Image: "image 1",
						Build: &kmmv1beta1.Build{},
					},
				},
			},
		}

	})

	ctx := context.Background()
	pullPods := []v1.Pod{}
	testMic = kmmv1beta1.ModuleImagesConfig{
		Spec: kmmv1beta1.ModuleImagesConfigSpec{
			Images: []kmmv1beta1.ModuleImageSpec{
				{
					Image: "image 1",
				},
			},
		},
	}

	It("image status empty, pull pod does not exists, need to create a pull pod", func() {
		gomock.InOrder(
			micHelper.EXPECT().GetImageState(&testMic, "image 1").Return(kmmv1beta1.ImageState("")),
			mockImagePuller.EXPECT().GetPullPodForImage(pullPods, "image 1").Return(nil),
			mockImagePuller.EXPECT().CreatePullPod(ctx, &testMic.Spec.Images[0], &testMic).Return(nil),
		)
		err := mrh.processImagesSpecs(ctx, &testMic, pullPods)
		Expect(err).To(BeNil())
	})

	It("image status empty, pull pod exists, nothing to do", func() {
		gomock.InOrder(
			micHelper.EXPECT().GetImageState(&testMic, "image 1").Return(kmmv1beta1.ImageState("")),
			mockImagePuller.EXPECT().GetPullPodForImage(pullPods, "image 1").Return(&v1.Pod{}),
		)
		err := mrh.processImagesSpecs(ctx, &testMic, pullPods)
		Expect(err).To(BeNil())
	})

	DescribeTable("images in MBSC status exist in MIC spec",
		func(imageState kmmv1beta1.ImageState, buildExists, signExists, updateMSBC bool, msbcAction kmmv1beta1.BuildOrSignAction) {
			testMic := kmmv1beta1.ModuleImagesConfig{
				Spec: kmmv1beta1.ModuleImagesConfigSpec{
					Images: []kmmv1beta1.ModuleImageSpec{
						{
							Image: "image 1",
						},
					},
				},
			}
			if buildExists {
				testMic.Spec.Images[0].Build = &kmmv1beta1.Build{}
			}
			if signExists {
				testMic.Spec.Images[0].Sign = &kmmv1beta1.Sign{}
			}
			micHelper.EXPECT().GetImageState(&testMic, "image 1").Return(imageState)

			if updateMSBC {
				mbscHelper.EXPECT().CreateOrPatch(ctx, &testMic, &testMic.Spec.Images[0], msbcAction).Return(nil)
			}

			err := mrh.processImagesSpecs(ctx, &testMic, pullPods)
			Expect(err).To(BeNil())
		},
		Entry("image state ImageDoesNotExist, no build or sign configs, do nothing",
			kmmv1beta1.ImageDoesNotExist, false, false, false, kmmv1beta1.SignImage),
		Entry("image state ImageDoesNotExist, build config exists, sign does not exists, update MBSC to build action",
			kmmv1beta1.ImageDoesNotExist, true, false, true, kmmv1beta1.BuildImage),
		Entry("image state ImageDoesNotExist, build config does not exists, sign config exists, update MBSC to build action",
			kmmv1beta1.ImageDoesNotExist, false, true, true, kmmv1beta1.BuildImage),
		Entry("image state ImageNeedsBuilding, build/sign config not important, update MBSC to build action",
			kmmv1beta1.ImageNeedsBuilding, false, false, true, kmmv1beta1.BuildImage),
		Entry("image state ImageNeedsSigning, build/sign config not important, update MBSC to sign action",
			kmmv1beta1.ImageNeedsSigning, false, false, true, kmmv1beta1.SignImage),
	)
})
