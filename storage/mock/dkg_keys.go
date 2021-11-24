// Code generated by mockery v1.0.0. DO NOT EDIT.

package mock

import (
	encodable "github.com/onflow/flow-go/model/encodable"
	mock "github.com/stretchr/testify/mock"
)

// DKGKeys is an autogenerated mock type for the DKGKeys type
type DKGKeys struct {
	mock.Mock
}

// InsertMyDKGPrivateInfo provides a mock function with given fields: epochCounter, key
func (_m *DKGKeys) InsertMyDKGPrivateInfo(epochCounter uint64, key *encodable.RandomBeaconPrivKey) error {
	ret := _m.Called(epochCounter, key)

	var r0 error
	if rf, ok := ret.Get(0).(func(uint64, *encodable.RandomBeaconPrivKey) error); ok {
		r0 = rf(epochCounter, key)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RetrieveMyDKGPrivateInfo provides a mock function with given fields: epochCounter
func (_m *DKGKeys) RetrieveMyDKGPrivateInfo(epochCounter uint64) (*encodable.RandomBeaconPrivKey, error) {
	ret := _m.Called(epochCounter)

	var r0 *encodable.RandomBeaconPrivKey
	if rf, ok := ret.Get(0).(func(uint64) *encodable.RandomBeaconPrivKey); ok {
		r0 = rf(epochCounter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*encodable.RandomBeaconPrivKey)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint64) error); ok {
		r1 = rf(epochCounter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
