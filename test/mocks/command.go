// Code generated by mockery v2.15.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	bus "github.com/vitalii-tkachuk/verification-service/internal/application/shared/bus"
)

// Command is an autogenerated mock type for the Command type
type Command struct {
	mock.Mock
}

// Type provides a mock function with given fields:
func (_m *Command) Type() bus.CommandType {
	ret := _m.Called()

	var r0 bus.CommandType
	if rf, ok := ret.Get(0).(func() bus.CommandType); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bus.CommandType)
	}

	return r0
}

type mockConstructorTestingTNewCommand interface {
	mock.TestingT
	Cleanup(func())
}

// NewCommand creates a new instance of Command. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewCommand(t mockConstructorTestingTNewCommand) *Command {
	mock := &Command{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
