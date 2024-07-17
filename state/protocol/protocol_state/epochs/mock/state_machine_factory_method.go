// Code generated by mockery v2.43.2. DO NOT EDIT.

package mock

import (
	flow "github.com/onflow/flow-go/model/flow"
	epochs "github.com/onflow/flow-go/state/protocol/protocol_state/epochs"

	mock "github.com/stretchr/testify/mock"
)

// StateMachineFactoryMethod is an autogenerated mock type for the StateMachineFactoryMethod type
type StateMachineFactoryMethod struct {
	mock.Mock
}

// Execute provides a mock function with given fields: candidateView, parentState
func (_m *StateMachineFactoryMethod) Execute(candidateView uint64, parentState *flow.RichEpochStateEntry) (epochs.StateMachine, error) {
	ret := _m.Called(candidateView, parentState)

	if len(ret) == 0 {
		panic("no return value specified for Execute")
	}

	var r0 epochs.StateMachine
	var r1 error
	if rf, ok := ret.Get(0).(func(uint64, *flow.RichEpochStateEntry) (epochs.StateMachine, error)); ok {
		return rf(candidateView, parentState)
	}
	if rf, ok := ret.Get(0).(func(uint64, *flow.RichEpochStateEntry) epochs.StateMachine); ok {
		r0 = rf(candidateView, parentState)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(epochs.StateMachine)
		}
	}

	if rf, ok := ret.Get(1).(func(uint64, *flow.RichEpochStateEntry) error); ok {
		r1 = rf(candidateView, parentState)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewStateMachineFactoryMethod creates a new instance of StateMachineFactoryMethod. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewStateMachineFactoryMethod(t interface {
	mock.TestingT
	Cleanup(func())
}) *StateMachineFactoryMethod {
	mock := &StateMachineFactoryMethod{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
