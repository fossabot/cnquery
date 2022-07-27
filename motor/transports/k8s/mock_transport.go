// Code generated by MockGen. DO NOT EDIT.
// Source: ./transport.go

// Package k8s is a generated GoMock package.
package k8s

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	afero "github.com/spf13/afero"
	platform "go.mondoo.io/mondoo/motor/platform"
	transports "go.mondoo.io/mondoo/motor/transports"
	resources "go.mondoo.io/mondoo/motor/transports/k8s/resources"
	v1 "k8s.io/api/core/v1"
	version "k8s.io/apimachinery/pkg/version"
)

// MockTransport is a mock of Transport interface.
type MockTransport struct {
	ctrl     *gomock.Controller
	recorder *MockTransportMockRecorder
}

// MockTransportMockRecorder is the mock recorder for MockTransport.
type MockTransportMockRecorder struct {
	mock *MockTransport
}

// NewMockTransport creates a new mock instance.
func NewMockTransport(ctrl *gomock.Controller) *MockTransport {
	mock := &MockTransport{ctrl: ctrl}
	mock.recorder = &MockTransportMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTransport) EXPECT() *MockTransportMockRecorder {
	return m.recorder
}

// Capabilities mocks base method.
func (m *MockTransport) Capabilities() transports.Capabilities {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Capabilities")
	ret0, _ := ret[0].(transports.Capabilities)
	return ret0
}

// Capabilities indicates an expected call of Capabilities.
func (mr *MockTransportMockRecorder) Capabilities() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Capabilities", reflect.TypeOf((*MockTransport)(nil).Capabilities))
}

// Close mocks base method.
func (m *MockTransport) Close() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close.
func (mr *MockTransportMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockTransport)(nil).Close))
}

// FS mocks base method.
func (m *MockTransport) FS() afero.Fs {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FS")
	ret0, _ := ret[0].(afero.Fs)
	return ret0
}

// FS indicates an expected call of FS.
func (mr *MockTransportMockRecorder) FS() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FS", reflect.TypeOf((*MockTransport)(nil).FS))
}

// FileInfo mocks base method.
func (m *MockTransport) FileInfo(path string) (transports.FileInfoDetails, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FileInfo", path)
	ret0, _ := ret[0].(transports.FileInfoDetails)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FileInfo indicates an expected call of FileInfo.
func (mr *MockTransportMockRecorder) FileInfo(path interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FileInfo", reflect.TypeOf((*MockTransport)(nil).FileInfo), path)
}

// ID mocks base method.
func (m *MockTransport) ID() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ID")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ID indicates an expected call of ID.
func (mr *MockTransportMockRecorder) ID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ID", reflect.TypeOf((*MockTransport)(nil).ID))
}

// Identifier mocks base method.
func (m *MockTransport) Identifier() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Identifier")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Identifier indicates an expected call of Identifier.
func (mr *MockTransportMockRecorder) Identifier() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Identifier", reflect.TypeOf((*MockTransport)(nil).Identifier))
}

// Kind mocks base method.
func (m *MockTransport) Kind() transports.Kind {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Kind")
	ret0, _ := ret[0].(transports.Kind)
	return ret0
}

// Kind indicates an expected call of Kind.
func (mr *MockTransportMockRecorder) Kind() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Kind", reflect.TypeOf((*MockTransport)(nil).Kind))
}

// Name mocks base method.
func (m *MockTransport) Name() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Name indicates an expected call of Name.
func (mr *MockTransportMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockTransport)(nil).Name))
}

// Namespaces mocks base method.
func (m *MockTransport) Namespaces() ([]v1.Namespace, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Namespaces")
	ret0, _ := ret[0].([]v1.Namespace)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Namespaces indicates an expected call of Namespaces.
func (mr *MockTransportMockRecorder) Namespaces() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Namespaces", reflect.TypeOf((*MockTransport)(nil).Namespaces))
}

// PlatformIdDetectors mocks base method.
func (m *MockTransport) PlatformIdDetectors() []transports.PlatformIdDetector {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PlatformIdDetectors")
	ret0, _ := ret[0].([]transports.PlatformIdDetector)
	return ret0
}

// PlatformIdDetectors indicates an expected call of PlatformIdDetectors.
func (mr *MockTransportMockRecorder) PlatformIdDetectors() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PlatformIdDetectors", reflect.TypeOf((*MockTransport)(nil).PlatformIdDetectors))
}

// PlatformIdentifier mocks base method.
func (m *MockTransport) PlatformIdentifier() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PlatformIdentifier")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PlatformIdentifier indicates an expected call of PlatformIdentifier.
func (mr *MockTransportMockRecorder) PlatformIdentifier() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PlatformIdentifier", reflect.TypeOf((*MockTransport)(nil).PlatformIdentifier))
}

// PlatformInfo mocks base method.
func (m *MockTransport) PlatformInfo() *platform.Platform {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PlatformInfo")
	ret0, _ := ret[0].(*platform.Platform)
	return ret0
}

// PlatformInfo indicates an expected call of PlatformInfo.
func (mr *MockTransportMockRecorder) PlatformInfo() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PlatformInfo", reflect.TypeOf((*MockTransport)(nil).PlatformInfo))
}

// Pod mocks base method.
func (m *MockTransport) Pod(namespace, name string) (*v1.Pod, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Pod", namespace, name)
	ret0, _ := ret[0].(*v1.Pod)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Pod indicates an expected call of Pod.
func (mr *MockTransportMockRecorder) Pod(namespace, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Pod", reflect.TypeOf((*MockTransport)(nil).Pod), namespace, name)
}

// Pods mocks base method.
func (m *MockTransport) Pods(namespace v1.Namespace) ([]v1.Pod, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Pods", namespace)
	ret0, _ := ret[0].([]v1.Pod)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Pods indicates an expected call of Pods.
func (mr *MockTransportMockRecorder) Pods(namespace interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Pods", reflect.TypeOf((*MockTransport)(nil).Pods), namespace)
}

// Resources mocks base method.
func (m *MockTransport) Resources(kind, name string) (*ResourceResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Resources", kind, name)
	ret0, _ := ret[0].(*ResourceResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Resources indicates an expected call of Resources.
func (mr *MockTransportMockRecorder) Resources(kind, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Resources", reflect.TypeOf((*MockTransport)(nil).Resources), kind, name)
}

// RunCommand mocks base method.
func (m *MockTransport) RunCommand(command string) (*transports.Command, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RunCommand", command)
	ret0, _ := ret[0].(*transports.Command)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RunCommand indicates an expected call of RunCommand.
func (mr *MockTransportMockRecorder) RunCommand(command interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RunCommand", reflect.TypeOf((*MockTransport)(nil).RunCommand), command)
}

// Runtime mocks base method.
func (m *MockTransport) Runtime() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Runtime")
	ret0, _ := ret[0].(string)
	return ret0
}

// Runtime indicates an expected call of Runtime.
func (mr *MockTransportMockRecorder) Runtime() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Runtime", reflect.TypeOf((*MockTransport)(nil).Runtime))
}

// ServerVersion mocks base method.
func (m *MockTransport) ServerVersion() *version.Info {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ServerVersion")
	ret0, _ := ret[0].(*version.Info)
	return ret0
}

// ServerVersion indicates an expected call of ServerVersion.
func (mr *MockTransportMockRecorder) ServerVersion() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ServerVersion", reflect.TypeOf((*MockTransport)(nil).ServerVersion))
}

// SupportedResourceTypes mocks base method.
func (m *MockTransport) SupportedResourceTypes() (*resources.ApiResourceIndex, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SupportedResourceTypes")
	ret0, _ := ret[0].(*resources.ApiResourceIndex)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SupportedResourceTypes indicates an expected call of SupportedResourceTypes.
func (mr *MockTransportMockRecorder) SupportedResourceTypes() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SupportedResourceTypes", reflect.TypeOf((*MockTransport)(nil).SupportedResourceTypes))
}
