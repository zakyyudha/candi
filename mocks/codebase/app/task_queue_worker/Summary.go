// Code generated by mockery v2.13.1. DO NOT EDIT.

package mocks

import (
	context "context"

	taskqueueworker "github.com/golangid/candi/codebase/app/task_queue_worker"
	mock "github.com/stretchr/testify/mock"
)

// Summary is an autogenerated mock type for the Summary type
type Summary struct {
	mock.Mock
}

// FindAllSummary provides a mock function with given fields: ctx, filter
func (_m *Summary) FindAllSummary(ctx context.Context, filter *taskqueueworker.Filter) []taskqueueworker.TaskSummary {
	ret := _m.Called(ctx, filter)

	var r0 []taskqueueworker.TaskSummary
	if rf, ok := ret.Get(0).(func(context.Context, *taskqueueworker.Filter) []taskqueueworker.TaskSummary); ok {
		r0 = rf(ctx, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]taskqueueworker.TaskSummary)
		}
	}

	return r0
}

// FindDetailSummary provides a mock function with given fields: ctx, taskName
func (_m *Summary) FindDetailSummary(ctx context.Context, taskName string) taskqueueworker.TaskSummary {
	ret := _m.Called(ctx, taskName)

	var r0 taskqueueworker.TaskSummary
	if rf, ok := ret.Get(0).(func(context.Context, string) taskqueueworker.TaskSummary); ok {
		r0 = rf(ctx, taskName)
	} else {
		r0 = ret.Get(0).(taskqueueworker.TaskSummary)
	}

	return r0
}

// IncrementSummary provides a mock function with given fields: ctx, taskName, incr
func (_m *Summary) IncrementSummary(ctx context.Context, taskName string, incr map[string]int64) {
	_m.Called(ctx, taskName, incr)
}

// UpdateSummary provides a mock function with given fields: ctx, taskName, updated
func (_m *Summary) UpdateSummary(ctx context.Context, taskName string, updated map[string]interface{}) {
	_m.Called(ctx, taskName, updated)
}

type mockConstructorTestingTNewSummary interface {
	mock.TestingT
	Cleanup(func())
}

// NewSummary creates a new instance of Summary. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewSummary(t mockConstructorTestingTNewSummary) *Summary {
	mock := &Summary{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
