// Code generated by MockGen. DO NOT EDIT.
// Source: proxy.go

// Package flags is a generated GoMock package.
package flags

import (
	gomock "github.com/golang/mock/gomock"
	service "github.com/turbinelabs/api/service"
	reflect "reflect"
)

// MockProxyFromFlags is a mock of ProxyFromFlags interface
type MockProxyFromFlags struct {
	ctrl     *gomock.Controller
	recorder *MockProxyFromFlagsMockRecorder
}

// MockProxyFromFlagsMockRecorder is the mock recorder for MockProxyFromFlags
type MockProxyFromFlagsMockRecorder struct {
	mock *MockProxyFromFlags
}

// NewMockProxyFromFlags creates a new mock instance
func NewMockProxyFromFlags(ctrl *gomock.Controller) *MockProxyFromFlags {
	mock := &MockProxyFromFlags{ctrl: ctrl}
	mock.recorder = &MockProxyFromFlagsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockProxyFromFlags) EXPECT() *MockProxyFromFlagsMockRecorder {
	return m.recorder
}

// Name mocks base method
func (m *MockProxyFromFlags) Name() string {
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name
func (mr *MockProxyFromFlagsMockRecorder) Name() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockProxyFromFlags)(nil).Name))
}

// Ref mocks base method
func (m *MockProxyFromFlags) Ref(arg0 service.ZoneRef) service.ProxyRef {
	ret := m.ctrl.Call(m, "Ref", arg0)
	ret0, _ := ret[0].(service.ProxyRef)
	return ret0
}

// Ref indicates an expected call of Ref
func (mr *MockProxyFromFlagsMockRecorder) Ref(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ref", reflect.TypeOf((*MockProxyFromFlags)(nil).Ref), arg0)
}
