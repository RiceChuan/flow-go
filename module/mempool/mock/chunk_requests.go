// Code generated by mockery v1.0.0. DO NOT EDIT.

package mempool

import (
	flow "github.com/onflow/flow-go/model/flow"

	mock "github.com/stretchr/testify/mock"

	verification "github.com/onflow/flow-go/model/verification"
)

// ChunkRequests is an autogenerated mock type for the ChunkRequests type
type ChunkRequests struct {
	mock.Mock
}

// Add provides a mock function with given fields: request
func (_m *ChunkRequests) Add(request *verification.ChunkDataPackRequest) bool {
	ret := _m.Called(request)

	var r0 bool
	if rf, ok := ret.Get(0).(func(*verification.ChunkDataPackRequest) bool); ok {
		r0 = rf(request)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// All provides a mock function with given fields:
func (_m *ChunkRequests) All() []*verification.ChunkDataPackRequest {
	ret := _m.Called()

	var r0 []*verification.ChunkDataPackRequest
	if rf, ok := ret.Get(0).(func() []*verification.ChunkDataPackRequest); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*verification.ChunkDataPackRequest)
		}
	}

	return r0
}

// ByID provides a mock function with given fields: chunkID
func (_m *ChunkRequests) ByID(chunkID flow.Identifier) (*verification.ChunkDataPackRequest, int, bool) {
	ret := _m.Called(chunkID)

	var r0 *verification.ChunkDataPackRequest
	if rf, ok := ret.Get(0).(func(flow.Identifier) *verification.ChunkDataPackRequest); ok {
		r0 = rf(chunkID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*verification.ChunkDataPackRequest)
		}
	}

	var r1 int
	if rf, ok := ret.Get(1).(func(flow.Identifier) int); ok {
		r1 = rf(chunkID)
	} else {
		r1 = ret.Get(1).(int)
	}

	var r2 bool
	if rf, ok := ret.Get(2).(func(flow.Identifier) bool); ok {
		r2 = rf(chunkID)
	} else {
		r2 = ret.Get(2).(bool)
	}

	return r0, r1, r2
}

// IncrementAttempt provides a mock function with given fields: chunkID
func (_m *ChunkRequests) IncrementAttempt(chunkID flow.Identifier) bool {
	ret := _m.Called(chunkID)

	var r0 bool
	if rf, ok := ret.Get(0).(func(flow.Identifier) bool); ok {
		r0 = rf(chunkID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Rem provides a mock function with given fields: chunkID
func (_m *ChunkRequests) Rem(chunkID flow.Identifier) bool {
	ret := _m.Called(chunkID)

	var r0 bool
	if rf, ok := ret.Get(0).(func(flow.Identifier) bool); ok {
		r0 = rf(chunkID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}
