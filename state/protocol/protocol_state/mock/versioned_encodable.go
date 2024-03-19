// Code generated by mockery v2.21.4. DO NOT EDIT.

package mock

import mock "github.com/stretchr/testify/mock"

// VersionedEncodable is an autogenerated mock type for the VersionedEncodable type
type VersionedEncodable struct {
	mock.Mock
}

// VersionedEncode provides a mock function with given fields:
func (_m *VersionedEncodable) VersionedEncode() (uint64, []byte, error) {
	ret := _m.Called()

	var r0 uint64
	var r1 []byte
	var r2 error
	if rf, ok := ret.Get(0).(func() (uint64, []byte, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() uint64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint64)
	}

	if rf, ok := ret.Get(1).(func() []byte); ok {
		r1 = rf()
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).([]byte)
		}
	}

	if rf, ok := ret.Get(2).(func() error); ok {
		r2 = rf()
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

type mockConstructorTestingTNewVersionedEncodable interface {
	mock.TestingT
	Cleanup(func())
}

// NewVersionedEncodable creates a new instance of VersionedEncodable. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewVersionedEncodable(t mockConstructorTestingTNewVersionedEncodable) *VersionedEncodable {
	mock := &VersionedEncodable{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
