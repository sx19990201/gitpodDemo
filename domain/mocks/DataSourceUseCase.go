// Code generated by mockery v2.13.1. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/fire_boom/domain"
	mock "github.com/stretchr/testify/mock"
)

// DataSourceUseCase is an autogenerated mock type for the DataSourceUseCase type
type DataSourceUseCase struct {
	mock.Mock
}

// Delete provides a mock function with given fields: ctx, id
func (_m *DataSourceUseCase) Delete(ctx context.Context, id uint) (int64, error) {
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
func (_m *DataSourceUseCase) FindDataSources(ctx context.Context) ([]domain.FbDataSource, error) {
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
func (_m *DataSourceUseCase) GetByID(ctx context.Context, id uint) (domain.FbDataSource, error) {
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

// Store provides a mock function with given fields: ctx, f
func (_m *DataSourceUseCase) Store(ctx context.Context, f *domain.FbDataSource) (int64, error) {
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
func (_m *DataSourceUseCase) Update(ctx context.Context, f *domain.FbDataSource) (int64, error) {
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

type mockConstructorTestingTNewDataSourceUseCase interface {
	mock.TestingT
	Cleanup(func())
}

// NewDataSourceUseCase creates a new instance of DataSourceUseCase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewDataSourceUseCase(t mockConstructorTestingTNewDataSourceUseCase) *DataSourceUseCase {
	mock := &DataSourceUseCase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}