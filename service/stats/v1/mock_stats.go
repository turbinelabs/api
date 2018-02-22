// Code generated by MockGen. DO NOT EDIT.
// Source: stats.go

// Package v1 is a generated GoMock package.
package v1

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockStatsQueryService is a mock of StatsQueryService interface
type MockStatsQueryService struct {
	ctrl     *gomock.Controller
	recorder *MockStatsQueryServiceMockRecorder
}

// MockStatsQueryServiceMockRecorder is the mock recorder for MockStatsQueryService
type MockStatsQueryServiceMockRecorder struct {
	mock *MockStatsQueryService
}

// NewMockStatsQueryService creates a new mock instance
func NewMockStatsQueryService(ctrl *gomock.Controller) *MockStatsQueryService {
	mock := &MockStatsQueryService{ctrl: ctrl}
	mock.recorder = &MockStatsQueryServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockStatsQueryService) EXPECT() *MockStatsQueryServiceMockRecorder {
	return m.recorder
}

// Query mocks base method
func (m *MockStatsQueryService) Query(arg0 *Query) (*QueryResult, error) {
	ret := m.ctrl.Call(m, "Query", arg0)
	ret0, _ := ret[0].(*QueryResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Query indicates an expected call of Query
func (mr *MockStatsQueryServiceMockRecorder) Query(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Query", reflect.TypeOf((*MockStatsQueryService)(nil).Query), arg0)
}
