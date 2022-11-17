// Code generated by mockery v2.15.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	bus "github.com/vitalii-tkachuk/verification-service/internal/application/shared/bus"
)

// Query is an autogenerated mock type for the Query type
type Query struct {
	mock.Mock
}

// Type provides a mock function with given fields:
func (_m *Query) Type() bus.QueryType {
	ret := _m.Called()

	var r0 bus.QueryType
	if rf, ok := ret.Get(0).(func() bus.QueryType); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bus.QueryType)
	}

	return r0
}

type mockConstructorTestingTNewQuery interface {
	mock.TestingT
	Cleanup(func())
}

// NewQuery creates a new instance of Query. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewQuery(t mockConstructorTestingTNewQuery) *Query {
	mock := &Query{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}