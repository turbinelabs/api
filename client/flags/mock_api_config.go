// Automatically generated by MockGen. DO NOT EDIT!
// Source: api_config.go

package flags

import (
	gomock "github.com/golang/mock/gomock"
	http0 "github.com/turbinelabs/api/http"
	http "net/http"
)

// Mock of APIConfigFromFlags interface
type MockAPIConfigFromFlags struct {
	ctrl     *gomock.Controller
	recorder *_MockAPIConfigFromFlagsRecorder
}

// Recorder for MockAPIConfigFromFlags (not exported)
type _MockAPIConfigFromFlagsRecorder struct {
	mock *MockAPIConfigFromFlags
}

func NewMockAPIConfigFromFlags(ctrl *gomock.Controller) *MockAPIConfigFromFlags {
	mock := &MockAPIConfigFromFlags{ctrl: ctrl}
	mock.recorder = &_MockAPIConfigFromFlagsRecorder{mock}
	return mock
}

func (_m *MockAPIConfigFromFlags) EXPECT() *_MockAPIConfigFromFlagsRecorder {
	return _m.recorder
}

func (_m *MockAPIConfigFromFlags) MakeClient() *http.Client {
	ret := _m.ctrl.Call(_m, "MakeClient")
	ret0, _ := ret[0].(*http.Client)
	return ret0
}

func (_mr *_MockAPIConfigFromFlagsRecorder) MakeClient() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "MakeClient")
}

func (_m *MockAPIConfigFromFlags) MakeEndpoint() (http0.Endpoint, error) {
	ret := _m.ctrl.Call(_m, "MakeEndpoint")
	ret0, _ := ret[0].(http0.Endpoint)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockAPIConfigFromFlagsRecorder) MakeEndpoint() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "MakeEndpoint")
}

func (_m *MockAPIConfigFromFlags) APIKey() string {
	ret := _m.ctrl.Call(_m, "APIKey")
	ret0, _ := ret[0].(string)
	return ret0
}

func (_mr *_MockAPIConfigFromFlagsRecorder) APIKey() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "APIKey")
}

func (_m *MockAPIConfigFromFlags) APIAuthKeyFromFlags() APIAuthKeyFromFlags {
	ret := _m.ctrl.Call(_m, "APIAuthKeyFromFlags")
	ret0, _ := ret[0].(APIAuthKeyFromFlags)
	return ret0
}

func (_mr *_MockAPIConfigFromFlagsRecorder) APIAuthKeyFromFlags() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "APIAuthKeyFromFlags")
}
