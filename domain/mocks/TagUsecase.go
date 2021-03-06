// Code generated by mockery v2.10.4. DO NOT EDIT.

package mocks

import (
	domain "github.com/nrmadi02/mini-project/domain"
	mock "github.com/stretchr/testify/mock"

	request "github.com/nrmadi02/mini-project/web/request"
)

// TagUsecase is an autogenerated mock type for the TagUsecase type
type TagUsecase struct {
	mock.Mock
}

// CreateNewTag provides a mock function with given fields: _a0
func (_m *TagUsecase) CreateNewTag(_a0 request.CreateTagRequest) (domain.Tag, error) {
	ret := _m.Called(_a0)

	var r0 domain.Tag
	if rf, ok := ret.Get(0).(func(request.CreateTagRequest) domain.Tag); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(domain.Tag)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(request.CreateTagRequest) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteTag provides a mock function with given fields: id
func (_m *TagUsecase) DeleteTag(id string) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAllTags provides a mock function with given fields:
func (_m *TagUsecase) GetAllTags() (domain.Tags, error) {
	ret := _m.Called()

	var r0 domain.Tags
	if rf, ok := ret.Get(0).(func() domain.Tags); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(domain.Tags)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
