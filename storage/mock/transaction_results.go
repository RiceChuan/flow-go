// Code generated by mockery v1.0.0. DO NOT EDIT.

package mock

import (
	flow "github.com/dapperlabs/flow-go/model/flow"
	mock "github.com/stretchr/testify/mock"
)

// TransactionResults is an autogenerated mock type for the TransactionResults type
type TransactionResults struct {
	mock.Mock
}

// BatchStore provides a mock function with given fields: blockID, transactionResults
func (_m *TransactionResults) BatchStore(blockID flow.Identifier, transactionResults []flow.TransactionResult) error {
	ret := _m.Called(blockID, transactionResults)

	var r0 error
	if rf, ok := ret.Get(0).(func(flow.Identifier, []flow.TransactionResult) error); ok {
		r0 = rf(blockID, transactionResults)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ByBlockIDTransactionID provides a mock function with given fields: blockID, transactionID
func (_m *TransactionResults) ByBlockIDTransactionID(blockID flow.Identifier, transactionID flow.Identifier) (*flow.TransactionResult, error) {
	ret := _m.Called(blockID, transactionID)

	var r0 *flow.TransactionResult
	if rf, ok := ret.Get(0).(func(flow.Identifier, flow.Identifier) *flow.TransactionResult); ok {
		r0 = rf(blockID, transactionID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*flow.TransactionResult)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(flow.Identifier, flow.Identifier) error); ok {
		r1 = rf(blockID, transactionID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Store provides a mock function with given fields: blockID, transactionResult
func (_m *TransactionResults) Store(blockID flow.Identifier, transactionResult *flow.TransactionResult) error {
	ret := _m.Called(blockID, transactionResult)

	var r0 error
	if rf, ok := ret.Get(0).(func(flow.Identifier, *flow.TransactionResult) error); ok {
		r0 = rf(blockID, transactionResult)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
