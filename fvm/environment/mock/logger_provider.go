// Code generated by mockery v2.43.2. DO NOT EDIT.

package mock

import (
	zerolog "github.com/rs/zerolog"
	mock "github.com/stretchr/testify/mock"
)

// LoggerProvider is an autogenerated mock type for the LoggerProvider type
type LoggerProvider struct {
	mock.Mock
}

// Logger provides a mock function with given fields:
func (_m *LoggerProvider) Logger() zerolog.Logger {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Logger")
	}

	var r0 zerolog.Logger
	if rf, ok := ret.Get(0).(func() zerolog.Logger); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(zerolog.Logger)
	}

	return r0
}

// NewLoggerProvider creates a new instance of LoggerProvider. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewLoggerProvider(t interface {
	mock.TestingT
	Cleanup(func())
}) *LoggerProvider {
	mock := &LoggerProvider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}