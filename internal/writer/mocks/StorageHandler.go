// Code generated by mockery v2.29.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// StorageHandler is an autogenerated mock type for the StorageHandler type
type StorageHandler struct {
	mock.Mock
}

// Add provides a mock function with given fields: _a0, _a1
func (_m *StorageHandler) Add(_a0 string, _a1 []byte) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, []byte) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetCurrent provides a mock function with given fields: _a0
func (_m *StorageHandler) GetCurrent(_a0 string) ([]byte, error) {
	ret := _m.Called(_a0)

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(string) ([]byte, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(string) []byte); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Initialize provides a mock function with given fields:
func (_m *StorageHandler) Initialize() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// IsClean provides a mock function with given fields:
func (_m *StorageHandler) IsClean() (bool, error) {
	ret := _m.Called()

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func() (bool, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Store provides a mock function with given fields:
func (_m *StorageHandler) Store() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewStorageHandler interface {
	mock.TestingT
	Cleanup(func())
}

// NewStorageHandler creates a new instance of StorageHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewStorageHandler(t mockConstructorTestingTNewStorageHandler) *StorageHandler {
	mock := &StorageHandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
