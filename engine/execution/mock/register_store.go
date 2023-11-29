// Code generated by mockery v2.21.4. DO NOT EDIT.

package mock

import (
	flow "github.com/onflow/flow-go/model/flow"
	mock "github.com/stretchr/testify/mock"
)

// RegisterStore is an autogenerated mock type for the RegisterStore type
type RegisterStore struct {
	mock.Mock
}

// GetRegister provides a mock function with given fields: height, blockID, register
func (_m *RegisterStore) GetRegister(height uint64, blockID flow.Identifier, register flow.RegisterID) ([]byte, error) {
	ret := _m.Called(height, blockID, register)

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(uint64, flow.Identifier, flow.RegisterID) ([]byte, error)); ok {
		return rf(height, blockID, register)
	}
	if rf, ok := ret.Get(0).(func(uint64, flow.Identifier, flow.RegisterID) []byte); ok {
		r0 = rf(height, blockID, register)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(uint64, flow.Identifier, flow.RegisterID) error); ok {
		r1 = rf(height, blockID, register)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsBlockExecuted provides a mock function with given fields: height, blockID
func (_m *RegisterStore) IsBlockExecuted(height uint64, blockID flow.Identifier) (bool, error) {
	ret := _m.Called(height, blockID)

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(uint64, flow.Identifier) (bool, error)); ok {
		return rf(height, blockID)
	}
	if rf, ok := ret.Get(0).(func(uint64, flow.Identifier) bool); ok {
		r0 = rf(height, blockID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(uint64, flow.Identifier) error); ok {
		r1 = rf(height, blockID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// LastFinalizedAndExecutedHeight provides a mock function with given fields:
func (_m *RegisterStore) LastFinalizedAndExecutedHeight() uint64 {
	ret := _m.Called()

	var r0 uint64
	if rf, ok := ret.Get(0).(func() uint64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint64)
	}

	return r0
}

// OnBlockFinalized provides a mock function with given fields:
func (_m *RegisterStore) OnBlockFinalized() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SaveRegisters provides a mock function with given fields: header, registers
func (_m *RegisterStore) SaveRegisters(header *flow.Header, registers flow.RegisterEntries) error {
	ret := _m.Called(header, registers)

	var r0 error
	if rf, ok := ret.Get(0).(func(*flow.Header, flow.RegisterEntries) error); ok {
		r0 = rf(header, registers)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewRegisterStore interface {
	mock.TestingT
	Cleanup(func())
}

// NewRegisterStore creates a new instance of RegisterStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewRegisterStore(t mockConstructorTestingTNewRegisterStore) *RegisterStore {
	mock := &RegisterStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}