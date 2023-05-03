// Code generated by mockery v2.21.4. DO NOT EDIT.

package mockp2p

import (
	flow "github.com/onflow/flow-go/model/flow"
	mock "github.com/stretchr/testify/mock"
)

// ClusterIDUpdateConsumer is an autogenerated mock type for the ClusterIDUpdateConsumer type
type ClusterIDUpdateConsumer struct {
	mock.Mock
}

// OnClusterIDSUpdate provides a mock function with given fields: _a0
func (_m *ClusterIDUpdateConsumer) OnClusterIDSUpdate(_a0 flow.ChainIDList) {
	_m.Called(_a0)
}

type mockConstructorTestingTNewClusterIDUpdateConsumer interface {
	mock.TestingT
	Cleanup(func())
}

// NewClusterIDUpdateConsumer creates a new instance of ClusterIDUpdateConsumer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewClusterIDUpdateConsumer(t mockConstructorTestingTNewClusterIDUpdateConsumer) *ClusterIDUpdateConsumer {
	mock := &ClusterIDUpdateConsumer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
