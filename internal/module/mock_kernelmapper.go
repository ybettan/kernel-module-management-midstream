// Code generated by MockGen. DO NOT EDIT.
// Source: kernelmapper.go
//
// Generated by this command:
//
//	mockgen -source=kernelmapper.go -package=module -destination=mock_kernelmapper.go KernelMapper,kernelMapperHelperAPI
//
// Package module is a generated GoMock package.
package module

import (
	reflect "reflect"

	v1beta1 "github.com/rh-ecosystem-edge/kernel-module-management/api/v1beta1"
	api "github.com/rh-ecosystem-edge/kernel-module-management/internal/api"
	gomock "go.uber.org/mock/gomock"
)

// MockKernelMapper is a mock of KernelMapper interface.
type MockKernelMapper struct {
	ctrl     *gomock.Controller
	recorder *MockKernelMapperMockRecorder
}

// MockKernelMapperMockRecorder is the mock recorder for MockKernelMapper.
type MockKernelMapperMockRecorder struct {
	mock *MockKernelMapper
}

// NewMockKernelMapper creates a new mock instance.
func NewMockKernelMapper(ctrl *gomock.Controller) *MockKernelMapper {
	mock := &MockKernelMapper{ctrl: ctrl}
	mock.recorder = &MockKernelMapperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockKernelMapper) EXPECT() *MockKernelMapperMockRecorder {
	return m.recorder
}

// GetModuleLoaderDataForKernel mocks base method.
func (m *MockKernelMapper) GetModuleLoaderDataForKernel(mod *v1beta1.Module, kernelVersion string) (*api.ModuleLoaderData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetModuleLoaderDataForKernel", mod, kernelVersion)
	ret0, _ := ret[0].(*api.ModuleLoaderData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetModuleLoaderDataForKernel indicates an expected call of GetModuleLoaderDataForKernel.
func (mr *MockKernelMapperMockRecorder) GetModuleLoaderDataForKernel(mod, kernelVersion any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetModuleLoaderDataForKernel", reflect.TypeOf((*MockKernelMapper)(nil).GetModuleLoaderDataForKernel), mod, kernelVersion)
}

// MockkernelMapperHelperAPI is a mock of kernelMapperHelperAPI interface.
type MockkernelMapperHelperAPI struct {
	ctrl     *gomock.Controller
	recorder *MockkernelMapperHelperAPIMockRecorder
}

// MockkernelMapperHelperAPIMockRecorder is the mock recorder for MockkernelMapperHelperAPI.
type MockkernelMapperHelperAPIMockRecorder struct {
	mock *MockkernelMapperHelperAPI
}

// NewMockkernelMapperHelperAPI creates a new mock instance.
func NewMockkernelMapperHelperAPI(ctrl *gomock.Controller) *MockkernelMapperHelperAPI {
	mock := &MockkernelMapperHelperAPI{ctrl: ctrl}
	mock.recorder = &MockkernelMapperHelperAPIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockkernelMapperHelperAPI) EXPECT() *MockkernelMapperHelperAPIMockRecorder {
	return m.recorder
}

// findKernelMapping mocks base method.
func (m *MockkernelMapperHelperAPI) findKernelMapping(mappings []v1beta1.KernelMapping, kernelVersion string) (*v1beta1.KernelMapping, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "findKernelMapping", mappings, kernelVersion)
	ret0, _ := ret[0].(*v1beta1.KernelMapping)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// findKernelMapping indicates an expected call of findKernelMapping.
func (mr *MockkernelMapperHelperAPIMockRecorder) findKernelMapping(mappings, kernelVersion any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "findKernelMapping", reflect.TypeOf((*MockkernelMapperHelperAPI)(nil).findKernelMapping), mappings, kernelVersion)
}

// getRelevantBuild mocks base method.
func (m *MockkernelMapperHelperAPI) getRelevantBuild(moduleBuild, mappingBuild *v1beta1.Build) *v1beta1.Build {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "getRelevantBuild", moduleBuild, mappingBuild)
	ret0, _ := ret[0].(*v1beta1.Build)
	return ret0
}

// getRelevantBuild indicates an expected call of getRelevantBuild.
func (mr *MockkernelMapperHelperAPIMockRecorder) getRelevantBuild(moduleBuild, mappingBuild any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "getRelevantBuild", reflect.TypeOf((*MockkernelMapperHelperAPI)(nil).getRelevantBuild), moduleBuild, mappingBuild)
}

// getRelevantSign mocks base method.
func (m *MockkernelMapperHelperAPI) getRelevantSign(moduleSign, mappingSign *v1beta1.Sign, kernel string) (*v1beta1.Sign, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "getRelevantSign", moduleSign, mappingSign, kernel)
	ret0, _ := ret[0].(*v1beta1.Sign)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// getRelevantSign indicates an expected call of getRelevantSign.
func (mr *MockkernelMapperHelperAPIMockRecorder) getRelevantSign(moduleSign, mappingSign, kernel any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "getRelevantSign", reflect.TypeOf((*MockkernelMapperHelperAPI)(nil).getRelevantSign), moduleSign, mappingSign, kernel)
}

// prepareModuleLoaderData mocks base method.
func (m *MockkernelMapperHelperAPI) prepareModuleLoaderData(mapping *v1beta1.KernelMapping, mod *v1beta1.Module, kernelVersion string) (*api.ModuleLoaderData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "prepareModuleLoaderData", mapping, mod, kernelVersion)
	ret0, _ := ret[0].(*api.ModuleLoaderData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// prepareModuleLoaderData indicates an expected call of prepareModuleLoaderData.
func (mr *MockkernelMapperHelperAPIMockRecorder) prepareModuleLoaderData(mapping, mod, kernelVersion any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "prepareModuleLoaderData", reflect.TypeOf((*MockkernelMapperHelperAPI)(nil).prepareModuleLoaderData), mapping, mod, kernelVersion)
}

// replaceTemplates mocks base method.
func (m *MockkernelMapperHelperAPI) replaceTemplates(mld *api.ModuleLoaderData) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "replaceTemplates", mld)
	ret0, _ := ret[0].(error)
	return ret0
}

// replaceTemplates indicates an expected call of replaceTemplates.
func (mr *MockkernelMapperHelperAPIMockRecorder) replaceTemplates(mld any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "replaceTemplates", reflect.TypeOf((*MockkernelMapperHelperAPI)(nil).replaceTemplates), mld)
}
