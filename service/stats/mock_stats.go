// Code generated by MockGen. DO NOT EDIT.
// Source: stats.go

// Package stats is a generated GoMock package.
package stats

import (
	gomock "github.com/golang/mock/gomock"
	v2 "github.com/turbinelabs/api/service/stats/v2"
	reflect "reflect"
)

// MockStatsService is a mock of StatsService interface
type MockStatsService struct {
	ctrl     *gomock.Controller
	recorder *MockStatsServiceMockRecorder
}

// MockStatsServiceMockRecorder is the mock recorder for MockStatsService
type MockStatsServiceMockRecorder struct {
	mock *MockStatsService
}

// NewMockStatsService creates a new mock instance
func NewMockStatsService(ctrl *gomock.Controller) *MockStatsService {
	mock := &MockStatsService{ctrl: ctrl}
	mock.recorder = &MockStatsServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockStatsService) EXPECT() *MockStatsServiceMockRecorder {
	return m.recorder
}

// ForwardV2 mocks base method
func (m *MockStatsService) ForwardV2(arg0 *v2.Payload) (*v2.ForwardResult, error) {
	ret := m.ctrl.Call(m, "ForwardV2", arg0)
	ret0, _ := ret[0].(*v2.ForwardResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ForwardV2 indicates an expected call of ForwardV2
func (mr *MockStatsServiceMockRecorder) ForwardV2(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ForwardV2", reflect.TypeOf((*MockStatsService)(nil).ForwardV2), arg0)
}

// QueryV2 mocks base method
func (m *MockStatsService) QueryV2(arg0 *v2.Query) (*v2.QueryResult, error) {
	ret := m.ctrl.Call(m, "QueryV2", arg0)
	ret0, _ := ret[0].(*v2.QueryResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryV2 indicates an expected call of QueryV2
func (mr *MockStatsServiceMockRecorder) QueryV2(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryV2", reflect.TypeOf((*MockStatsService)(nil).QueryV2), arg0)
}

// Close mocks base method
func (m *MockStatsService) Close() error {
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockStatsServiceMockRecorder) Close() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockStatsService)(nil).Close))
}
