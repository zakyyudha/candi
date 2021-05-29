// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	taskqueueworker "pkg.agungdp.dev/candi/codebase/app/task_queue_worker"
)

// QueueStorage is an autogenerated mock type for the QueueStorage type
type QueueStorage struct {
	mock.Mock
}

// Clear provides a mock function with given fields: taskName
func (_m *QueueStorage) Clear(taskName string) {
	_m.Called(taskName)
}

// GetAllJobs provides a mock function with given fields: taskName
func (_m *QueueStorage) GetAllJobs(taskName string) []*taskqueueworker.Job {
	ret := _m.Called(taskName)

	var r0 []*taskqueueworker.Job
	if rf, ok := ret.Get(0).(func(string) []*taskqueueworker.Job); ok {
		r0 = rf(taskName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*taskqueueworker.Job)
		}
	}

	return r0
}

// NextJob provides a mock function with given fields: taskName
func (_m *QueueStorage) NextJob(taskName string) *taskqueueworker.Job {
	ret := _m.Called(taskName)

	var r0 *taskqueueworker.Job
	if rf, ok := ret.Get(0).(func(string) *taskqueueworker.Job); ok {
		r0 = rf(taskName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*taskqueueworker.Job)
		}
	}

	return r0
}

// PopJob provides a mock function with given fields: taskName
func (_m *QueueStorage) PopJob(taskName string) taskqueueworker.Job {
	ret := _m.Called(taskName)

	var r0 taskqueueworker.Job
	if rf, ok := ret.Get(0).(func(string) taskqueueworker.Job); ok {
		r0 = rf(taskName)
	} else {
		r0 = ret.Get(0).(taskqueueworker.Job)
	}

	return r0
}

// PushJob provides a mock function with given fields: job
func (_m *QueueStorage) PushJob(job *taskqueueworker.Job) {
	_m.Called(job)
}