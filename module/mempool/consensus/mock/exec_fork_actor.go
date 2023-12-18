// Code generated by mockery v2.21.4. DO NOT EDIT.

package mock

import (
	flow "github.com/onflow/flow-go/model/flow"
	mock "github.com/stretchr/testify/mock"
)

// ExecForkActorMock is an autogenerated mock type for the ExecForkActor type
type ExecForkActorMock struct {
	mock.Mock
}

// OnExecFork provides a mock function with given fields: _a0
func (_m *ExecForkActorMock) OnExecFork(_a0 []*flow.IncorporatedResultSeal) {
	_m.Called(_a0)
}

type mockConstructorTestingTNewExecForkActorMock interface {
	mock.TestingT
	Cleanup(func())
}

// NewExecForkActorMock creates a new instance of ExecForkActorMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewExecForkActorMock(t mockConstructorTestingTNewExecForkActorMock) *ExecForkActorMock {
	mock := &ExecForkActorMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
