// Code generated by mockery v2.43.2. DO NOT EDIT.

package mock

import (
	io "io"

	storage "github.com/onflow/flow-go/storage"
	mock "github.com/stretchr/testify/mock"
)

// Reader is an autogenerated mock type for the Reader type
type Reader struct {
	mock.Mock
}

// Get provides a mock function with given fields: key
func (_m *Reader) Get(key []byte) ([]byte, io.Closer, error) {
	ret := _m.Called(key)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 []byte
	var r1 io.Closer
	var r2 error
	if rf, ok := ret.Get(0).(func([]byte) ([]byte, io.Closer, error)); ok {
		return rf(key)
	}
	if rf, ok := ret.Get(0).(func([]byte) []byte); ok {
		r0 = rf(key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func([]byte) io.Closer); ok {
		r1 = rf(key)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(io.Closer)
		}
	}

	if rf, ok := ret.Get(2).(func([]byte) error); ok {
		r2 = rf(key)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// NewIter provides a mock function with given fields: startPrefix, endPrefix, ops
func (_m *Reader) NewIter(startPrefix []byte, endPrefix []byte, ops storage.IteratorOption) (storage.Iterator, error) {
	ret := _m.Called(startPrefix, endPrefix, ops)

	if len(ret) == 0 {
		panic("no return value specified for NewIter")
	}

	var r0 storage.Iterator
	var r1 error
	if rf, ok := ret.Get(0).(func([]byte, []byte, storage.IteratorOption) (storage.Iterator, error)); ok {
		return rf(startPrefix, endPrefix, ops)
	}
	if rf, ok := ret.Get(0).(func([]byte, []byte, storage.IteratorOption) storage.Iterator); ok {
		r0 = rf(startPrefix, endPrefix, ops)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(storage.Iterator)
		}
	}

	if rf, ok := ret.Get(1).(func([]byte, []byte, storage.IteratorOption) error); ok {
		r1 = rf(startPrefix, endPrefix, ops)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewReader creates a new instance of Reader. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewReader(t interface {
	mock.TestingT
	Cleanup(func())
}) *Reader {
	mock := &Reader{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}