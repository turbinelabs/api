// Automatically generated by MockGen. DO NOT EDIT!
// Source: api_key.go

package flags

import (
	gomock "github.com/golang/mock/gomock"
)

// Mock of APIAuthKeyFromFlags interface
type MockAPIAuthKeyFromFlags struct {
	ctrl     *gomock.Controller
	recorder *_MockAPIAuthKeyFromFlagsRecorder
}

// Recorder for MockAPIAuthKeyFromFlags (not exported)
type _MockAPIAuthKeyFromFlagsRecorder struct {
	mock *MockAPIAuthKeyFromFlags
}

func NewMockAPIAuthKeyFromFlags(ctrl *gomock.Controller) *MockAPIAuthKeyFromFlags {
	mock := &MockAPIAuthKeyFromFlags{ctrl: ctrl}
	mock.recorder = &_MockAPIAuthKeyFromFlagsRecorder{mock}
	return mock
}

func (_m *MockAPIAuthKeyFromFlags) EXPECT() *_MockAPIAuthKeyFromFlagsRecorder {
	return _m.recorder
}

func (_m *MockAPIAuthKeyFromFlags) Make() string {
	ret := _m.ctrl.Call(_m, "Make")
	ret0, _ := ret[0].(string)
	return ret0
}

func (_mr *_MockAPIAuthKeyFromFlagsRecorder) Make() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Make")
}
