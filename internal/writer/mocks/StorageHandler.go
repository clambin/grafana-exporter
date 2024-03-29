// Code generated by mockery v2.32.4. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// StorageHandler is an autogenerated mock type for the StorageHandler type
type StorageHandler struct {
	mock.Mock
}

type StorageHandler_Expecter struct {
	mock *mock.Mock
}

func (_m *StorageHandler) EXPECT() *StorageHandler_Expecter {
	return &StorageHandler_Expecter{mock: &_m.Mock}
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

// StorageHandler_Add_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Add'
type StorageHandler_Add_Call struct {
	*mock.Call
}

// Add is a helper method to define mock.On call
//   - _a0 string
//   - _a1 []byte
func (_e *StorageHandler_Expecter) Add(_a0 interface{}, _a1 interface{}) *StorageHandler_Add_Call {
	return &StorageHandler_Add_Call{Call: _e.mock.On("Add", _a0, _a1)}
}

func (_c *StorageHandler_Add_Call) Run(run func(_a0 string, _a1 []byte)) *StorageHandler_Add_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].([]byte))
	})
	return _c
}

func (_c *StorageHandler_Add_Call) Return(_a0 error) *StorageHandler_Add_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *StorageHandler_Add_Call) RunAndReturn(run func(string, []byte) error) *StorageHandler_Add_Call {
	_c.Call.Return(run)
	return _c
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

// StorageHandler_GetCurrent_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetCurrent'
type StorageHandler_GetCurrent_Call struct {
	*mock.Call
}

// GetCurrent is a helper method to define mock.On call
//   - _a0 string
func (_e *StorageHandler_Expecter) GetCurrent(_a0 interface{}) *StorageHandler_GetCurrent_Call {
	return &StorageHandler_GetCurrent_Call{Call: _e.mock.On("GetCurrent", _a0)}
}

func (_c *StorageHandler_GetCurrent_Call) Run(run func(_a0 string)) *StorageHandler_GetCurrent_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *StorageHandler_GetCurrent_Call) Return(_a0 []byte, _a1 error) *StorageHandler_GetCurrent_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *StorageHandler_GetCurrent_Call) RunAndReturn(run func(string) ([]byte, error)) *StorageHandler_GetCurrent_Call {
	_c.Call.Return(run)
	return _c
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

// StorageHandler_Initialize_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Initialize'
type StorageHandler_Initialize_Call struct {
	*mock.Call
}

// Initialize is a helper method to define mock.On call
func (_e *StorageHandler_Expecter) Initialize() *StorageHandler_Initialize_Call {
	return &StorageHandler_Initialize_Call{Call: _e.mock.On("Initialize")}
}

func (_c *StorageHandler_Initialize_Call) Run(run func()) *StorageHandler_Initialize_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *StorageHandler_Initialize_Call) Return(_a0 error) *StorageHandler_Initialize_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *StorageHandler_Initialize_Call) RunAndReturn(run func() error) *StorageHandler_Initialize_Call {
	_c.Call.Return(run)
	return _c
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

// StorageHandler_IsClean_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IsClean'
type StorageHandler_IsClean_Call struct {
	*mock.Call
}

// IsClean is a helper method to define mock.On call
func (_e *StorageHandler_Expecter) IsClean() *StorageHandler_IsClean_Call {
	return &StorageHandler_IsClean_Call{Call: _e.mock.On("IsClean")}
}

func (_c *StorageHandler_IsClean_Call) Run(run func()) *StorageHandler_IsClean_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *StorageHandler_IsClean_Call) Return(_a0 bool, _a1 error) *StorageHandler_IsClean_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *StorageHandler_IsClean_Call) RunAndReturn(run func() (bool, error)) *StorageHandler_IsClean_Call {
	_c.Call.Return(run)
	return _c
}

// Store provides a mock function with given fields: msg
func (_m *StorageHandler) Store(msg string) error {
	ret := _m.Called(msg)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(msg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// StorageHandler_Store_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Store'
type StorageHandler_Store_Call struct {
	*mock.Call
}

// Store is a helper method to define mock.On call
//   - msg string
func (_e *StorageHandler_Expecter) Store(msg interface{}) *StorageHandler_Store_Call {
	return &StorageHandler_Store_Call{Call: _e.mock.On("Store", msg)}
}

func (_c *StorageHandler_Store_Call) Run(run func(msg string)) *StorageHandler_Store_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *StorageHandler_Store_Call) Return(_a0 error) *StorageHandler_Store_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *StorageHandler_Store_Call) RunAndReturn(run func(string) error) *StorageHandler_Store_Call {
	_c.Call.Return(run)
	return _c
}

// NewStorageHandler creates a new instance of StorageHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewStorageHandler(t interface {
	mock.TestingT
	Cleanup(func())
}) *StorageHandler {
	mock := &StorageHandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
