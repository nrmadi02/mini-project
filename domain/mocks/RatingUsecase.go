// Code generated by mockery v2.10.4. DO NOT EDIT.

package mocks

import (
	domain "github.com/nrmadi02/mini-project/domain"
	mock "github.com/stretchr/testify/mock"
)

// RatingUsecase is an autogenerated mock type for the RatingUsecase type
type RatingUsecase struct {
	mock.Mock
}

// AddNewRanting provides a mock function with given fields: id, userid, value
func (_m *RatingUsecase) AddNewRanting(id string, userid string, value int) (domain.RatingEnterprise, error) {
	ret := _m.Called(id, userid, value)

	var r0 domain.RatingEnterprise
	if rf, ok := ret.Get(0).(func(string, string, int) domain.RatingEnterprise); ok {
		r0 = rf(id, userid, value)
	} else {
		r0 = ret.Get(0).(domain.RatingEnterprise)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, int) error); ok {
		r1 = rf(id, userid, value)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteRating provides a mock function with given fields: id, userid
func (_m *RatingUsecase) DeleteRating(id string, userid string) error {
	ret := _m.Called(id, userid)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(id, userid)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindRating provides a mock function with given fields: id, userid
func (_m *RatingUsecase) FindRating(id string, userid string) (domain.RatingEnterprise, error) {
	ret := _m.Called(id, userid)

	var r0 domain.RatingEnterprise
	if rf, ok := ret.Get(0).(func(string, string) domain.RatingEnterprise); ok {
		r0 = rf(id, userid)
	} else {
		r0 = ret.Get(0).(domain.RatingEnterprise)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(id, userid)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAllRatingByEnterpriseID provides a mock function with given fields: id
func (_m *RatingUsecase) GetAllRatingByEnterpriseID(id string) (domain.RatingEnterprises, error) {
	ret := _m.Called(id)

	var r0 domain.RatingEnterprises
	if rf, ok := ret.Get(0).(func(string) domain.RatingEnterprises); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(domain.RatingEnterprises)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateRating provides a mock function with given fields: id, userid, value
func (_m *RatingUsecase) UpdateRating(id string, userid string, value int) (domain.RatingEnterprise, error) {
	ret := _m.Called(id, userid, value)

	var r0 domain.RatingEnterprise
	if rf, ok := ret.Get(0).(func(string, string, int) domain.RatingEnterprise); ok {
		r0 = rf(id, userid, value)
	} else {
		r0 = ret.Get(0).(domain.RatingEnterprise)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, int) error); ok {
		r1 = rf(id, userid, value)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
