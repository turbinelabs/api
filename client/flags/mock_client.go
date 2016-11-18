// Automatically generated by MockGen. DO NOT EDIT!
// Source: client.go

package flags

import (
	gomock "github.com/golang/mock/gomock"
	service "github.com/turbinelabs/api/service"
)

// Mock of ClientFromFlags interface
type MockClientFromFlags struct {
	ctrl     *gomock.Controller
	recorder *_MockClientFromFlagsRecorder
}

// Recorder for MockClientFromFlags (not exported)
type _MockClientFromFlagsRecorder struct {
	mock *MockClientFromFlags
}

func NewMockClientFromFlags(ctrl *gomock.Controller) *MockClientFromFlags {
	mock := &MockClientFromFlags{ctrl: ctrl}
	mock.recorder = &_MockClientFromFlagsRecorder{mock}
	return mock
}

func (_m *MockClientFromFlags) EXPECT() *_MockClientFromFlagsRecorder {
	return _m.recorder
}

func (_m *MockClientFromFlags) Make() (service.All, error) {
	ret := _m.ctrl.Call(_m, "Make")
	ret0, _ := ret[0].(service.All)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockClientFromFlagsRecorder) Make() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Make")
}