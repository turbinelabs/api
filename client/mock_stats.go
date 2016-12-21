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
// Source: stats.go

package client

import (
	gomock "github.com/golang/mock/gomock"
	stats "github.com/turbinelabs/api/service/stats"
	executor "github.com/turbinelabs/nonstdlib/executor"
)

// Mock of internalStatsClient interface
type MockinternalStatsClient struct {
	ctrl     *gomock.Controller
	recorder *_MockinternalStatsClientRecorder
}

// Recorder for MockinternalStatsClient (not exported)
type _MockinternalStatsClientRecorder struct {
	mock *MockinternalStatsClient
}

func NewMockinternalStatsClient(ctrl *gomock.Controller) *MockinternalStatsClient {
	mock := &MockinternalStatsClient{ctrl: ctrl}
	mock.recorder = &_MockinternalStatsClientRecorder{mock}
	return mock
}

func (_m *MockinternalStatsClient) EXPECT() *_MockinternalStatsClientRecorder {
	return _m.recorder
}

func (_m *MockinternalStatsClient) Forward(_param0 *stats.Payload) (*stats.ForwardResult, error) {
	ret := _m.ctrl.Call(_m, "Forward", _param0)
	ret0, _ := ret[0].(*stats.ForwardResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockinternalStatsClientRecorder) Forward(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Forward", arg0)
}

func (_m *MockinternalStatsClient) Query(_param0 *stats.Query) (*stats.QueryResult, error) {
	ret := _m.ctrl.Call(_m, "Query", _param0)
	ret0, _ := ret[0].(*stats.QueryResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockinternalStatsClientRecorder) Query(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Query", arg0)
}

func (_m *MockinternalStatsClient) Close() error {
	ret := _m.ctrl.Call(_m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockinternalStatsClientRecorder) Close() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Close")
}

func (_m *MockinternalStatsClient) ForwardWithCallback(_param0 *stats.Payload, _param1 executor.CallbackFunc) error {
	ret := _m.ctrl.Call(_m, "ForwardWithCallback", _param0, _param1)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockinternalStatsClientRecorder) ForwardWithCallback(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ForwardWithCallback", arg0, arg1)
}
