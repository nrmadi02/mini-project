// Code generated by mockery v2.10.4. DO NOT EDIT.

package mocks

import (
	domain "github.com/nrmadi02/mini-project/domain"
	mock "github.com/stretchr/testify/mock"
)

// RoleRepository is an autogenerated mock type for the RoleRepository type
type RoleRepository struct {
	mock.Mock
}

// FindByName provides a mock function with given fields: name
func (_m *RoleRepository) FindByName(name string) (domain.Role, error) {
	ret := _m.Called(name)

	var r0 domain.Role
	if rf, ok := ret.Get(0).(func(string) domain.Role); ok {
		r0 = rf(name)
	} else {
		r0 = ret.Get(0).(domain.Role)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
