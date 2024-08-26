// Code generated by mockery v2.43.2. DO NOT EDIT.

package mock

import mock "github.com/stretchr/testify/mock"

// EVMMetrics is an autogenerated mock type for the EVMMetrics type
type EVMMetrics struct {
	mock.Mock
}

// EVMBlockExecuted provides a mock function with given fields: txCount, totalGasUsed, totalSupplyInFlow
func (_m *EVMMetrics) EVMBlockExecuted(txCount int, totalGasUsed uint64, totalSupplyInFlow float64) {
	_m.Called(txCount, totalGasUsed, totalSupplyInFlow)
}

// EVMTransactionExecuted provides a mock function with given fields: gasUsed, isDirectCall, failed
func (_m *EVMMetrics) EVMTransactionExecuted(gasUsed uint64, isDirectCall bool, failed bool) {
	_m.Called(gasUsed, isDirectCall, failed)
}

// SetNumberOfDeployedCOAs provides a mock function with given fields: count
func (_m *EVMMetrics) SetNumberOfDeployedCOAs(count uint64) {
	_m.Called(count)
}

// NewEVMMetrics creates a new instance of EVMMetrics. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewEVMMetrics(t interface {
	mock.TestingT
	Cleanup(func())
}) *EVMMetrics {
	mock := &EVMMetrics{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}