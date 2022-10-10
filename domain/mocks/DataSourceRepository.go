// Code generated by mockery v2.13.1. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/fire_boom/domain"
	mock "github.com/stretchr/testify/mock"
)

// DataSourceRepository is an autogenerated mock type for the DataSourceRepository type
type DataSourceRepository struct {
	mock.Mock
}

// CheckExist provides a mock function with given fields: ctx, auth
func (_m *DataSourceRepository) CheckExist(ctx context.Context, auth *domain.FbDataSource) (domain.FbDataSource, error) {
	ret := _m.Called(ctx, auth)

	var r0 domain.FbDataSource
	if rf, ok := ret.Get(0).(func(context.Context, *domain.FbDataSource) domain.FbDataSource); ok {
		r0 = rf(ctx, auth)
	} else {
		r0 = ret.Get(0).(domain.FbDataSource)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *domain.FbDataSource) error); ok {
		r1 = rf(ctx, auth)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: ctx, id
func (_m *DataSourceRepository) Delete(ctx context.Context, id uint) (int64, error) {
	ret := _m.Called(ctx, id)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context, uint) int64); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uint) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindDataSources provides a mock function with given fields: ctx
func (_m *DataSourceRepository) FindDataSources(ctx context.Context) ([]domain.FbDataSource, error) {
	ret := _m.Called(ctx)

	var r0 []domain.FbDataSource
	if rf, ok := ret.Get(0).(func(context.Context) []domain.FbDataSource); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.FbDataSource)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByID provides a mock function with given fields: ctx, id
func (_m *DataSourceRepository) GetByID(ctx context.Context, id uint) (domain.FbDataSource, error) {
	ret := _m.Called(ctx, id)

	var r0 domain.FbDataSource
	if rf, ok := ret.Get(0).(func(context.Context, uint) domain.FbDataSource); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(domain.FbDataSource)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uint) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByName provides a mock function with given fields: ctx, name
func (_m *DataSourceRepository) GetByName(ctx context.Context, name string) (domain.FbDataSource, error) {
	ret := _m.Called(ctx, name)

	var r0 domain.FbDataSource
	if rf, ok := ret.Get(0).(func(context.Context, string) domain.FbDataSource); ok {
		r0 = rf(ctx, name)
	} else {
		r0 = ret.Get(0).(domain.FbDataSource)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Store provides a mock function with given fields: ctx, f
func (_m *DataSourceRepository) Store(ctx context.Context, f *domain.FbDataSource) (int64, error) {
	ret := _m.Called(ctx, f)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context, *domain.FbDataSource) int64); ok {
		r0 = rf(ctx, f)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *domain.FbDataSource) error); ok {
		r1 = rf(ctx, f)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, f
func (_m *DataSourceRepository) Update(ctx context.Context, f *domain.FbDataSource) (int64, error) {
	ret := _m.Called(ctx, f)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context, *domain.FbDataSource) int64); ok {
		r0 = rf(ctx, f)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *domain.FbDataSource) error); ok {
		r1 = rf(ctx, f)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, f
func (_m *DataSourceRepository) CheckDbConn(ctx context.Context, f *domain.FbDataSource) (bool) {


	return false
}

type mockConstructorTestingTNewDataSourceRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewDataSourceRepository creates a new instance of DataSourceRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewDataSourceRepository(t mockConstructorTestingTNewDataSourceRepository) *DataSourceRepository {
	mock := &DataSourceRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}