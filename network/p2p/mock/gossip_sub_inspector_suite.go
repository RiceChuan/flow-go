// Code generated by mockery v2.21.4. DO NOT EDIT.

package mockp2p

import (
	irrecoverable "github.com/onflow/flow-go/module/irrecoverable"
	mock "github.com/stretchr/testify/mock"

	p2p "github.com/onflow/flow-go/network/p2p"

	peer "github.com/libp2p/go-libp2p/core/peer"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

// GossipSubInspectorSuite is an autogenerated mock type for the GossipSubInspectorSuite type
type GossipSubInspectorSuite struct {
	mock.Mock
}

// AddInvCtrlMsgNotifConsumer provides a mock function with given fields: _a0
func (_m *GossipSubInspectorSuite) AddInvCtrlMsgNotifConsumer(_a0 p2p.GossipSubInvCtrlMsgNotifConsumer) {
	_m.Called(_a0)
}

// Done provides a mock function with given fields:
func (_m *GossipSubInspectorSuite) Done() <-chan struct{} {
	ret := _m.Called()

	var r0 <-chan struct{}
	if rf, ok := ret.Get(0).(func() <-chan struct{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan struct{})
		}
	}

	return r0
}

// InspectFunc provides a mock function with given fields:
func (_m *GossipSubInspectorSuite) InspectFunc() func(peer.ID, *pubsub.RPC) error {
	ret := _m.Called()

	var r0 func(peer.ID, *pubsub.RPC) error
	if rf, ok := ret.Get(0).(func() func(peer.ID, *pubsub.RPC) error); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(func(peer.ID, *pubsub.RPC) error)
		}
	}

	return r0
}

// Ready provides a mock function with given fields:
func (_m *GossipSubInspectorSuite) Ready() <-chan struct{} {
	ret := _m.Called()

	var r0 <-chan struct{}
	if rf, ok := ret.Get(0).(func() <-chan struct{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan struct{})
		}
	}

	return r0
}

// Start provides a mock function with given fields: _a0
func (_m *GossipSubInspectorSuite) Start(_a0 irrecoverable.SignalerContext) {
	_m.Called(_a0)
}

type mockConstructorTestingTNewGossipSubInspectorSuite interface {
	mock.TestingT
	Cleanup(func())
}

// NewGossipSubInspectorSuite creates a new instance of GossipSubInspectorSuite. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewGossipSubInspectorSuite(t mockConstructorTestingTNewGossipSubInspectorSuite) *GossipSubInspectorSuite {
	mock := &GossipSubInspectorSuite{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
