/*
Copyright 2017 Turbine Labs, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Automatically generated by MockGen. DO NOT EDIT!
// Source: zone_key.go

package flags

import (
	gomock "github.com/golang/mock/gomock"
	api "github.com/turbinelabs/api"
	service "github.com/turbinelabs/api/service"
)

// Mock of ZoneKeyFromFlags interface
type MockZoneKeyFromFlags struct {
	ctrl     *gomock.Controller
	recorder *_MockZoneKeyFromFlagsRecorder
}

// Recorder for MockZoneKeyFromFlags (not exported)
type _MockZoneKeyFromFlagsRecorder struct {
	mock *MockZoneKeyFromFlags
}

func NewMockZoneKeyFromFlags(ctrl *gomock.Controller) *MockZoneKeyFromFlags {
	mock := &MockZoneKeyFromFlags{ctrl: ctrl}
	mock.recorder = &_MockZoneKeyFromFlagsRecorder{mock}
	return mock
}

func (_m *MockZoneKeyFromFlags) EXPECT() *_MockZoneKeyFromFlagsRecorder {
	return _m.recorder
}

func (_m *MockZoneKeyFromFlags) Get(_param0 service.All) (api.ZoneKey, error) {
	ret := _m.ctrl.Call(_m, "Get", _param0)
	ret0, _ := ret[0].(api.ZoneKey)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockZoneKeyFromFlagsRecorder) Get(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Get", arg0)
}

func (_m *MockZoneKeyFromFlags) ZoneName() string {
	ret := _m.ctrl.Call(_m, "ZoneName")
	ret0, _ := ret[0].(string)
	return ret0
}

func (_mr *_MockZoneKeyFromFlagsRecorder) ZoneName() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ZoneName")
}
