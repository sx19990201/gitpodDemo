package usecase

import (
	"context"
	"errors"
	"github.com/fire_boom/domain"
	"github.com/fire_boom/domain/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestStorageBucketUseCase_Store(t *testing.T) {
	mockRepo := new(mocks.StorageBucketRepository)
	mockDS := domain.FbStorageBucket{
		Name:       "1",
		Switch:     1,
		Config:     "1",
		CreateTime: dataFormatStr,
		UpdateTime: dataFormatStr,
		IsDel:      1,
	}
	t.Run("success", func(t *testing.T) {
		tempMockDS := mockDS
		tempMockDS.ID = 0
		mockRepo.On("GetByName", mock.Anything, mock.AnythingOfType("string")).Return(domain.FbStorageBucket{}, nil).Once()
		mockRepo.On("Store", mock.Anything, mock.AnythingOfType("*domain.FbStorageBucket")).Return(int64(1), nil).Once()

		u := NewStorageBucketUseCase(mockRepo, time.Second*3)
		_, err := u.Store(context.TODO(), &tempMockDS)
		assert.NoError(t, err)
		assert.Equal(t, mockDS.Name, tempMockDS.Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("name already exists", func(t *testing.T) {
		tempMockDS := mockDS
		tempMockDS.ID = 0
		errWant := errors.New("name already exists")
		var affectWant int64 = 0
		mockRepo.On("GetByName", mock.Anything, mock.AnythingOfType("string")).Return(domain.FbStorageBucket{Name: "1"}, nil).Once()
		mockRepo.On("Store", mock.Anything, mock.AnythingOfType("*domain.FbStorageBucket")).Return(affectWant, errWant).Once()

		u := NewStorageBucketUseCase(mockRepo, time.Second*3)
		affect, err := u.Store(context.TODO(), &tempMockDS)
		assert.Error(t, err, errWant)
		assert.Equal(t, affect, affectWant)
	})
}

func TestStorageBucketUseCase_Update(t *testing.T) {
	mockRepo := new(mocks.StorageBucketRepository)
	mockDS := domain.FbStorageBucket{
		Name:       "1",
		Switch:     1,
		Config:     "1",
		CreateTime: dataFormatStr,
		UpdateTime: dataFormatStr,
		IsDel:      1,
	}
	t.Run("success", func(t *testing.T) {
		want := mockDS
		want.ID = 1
		want.Name = "123456"
		mockRepo.On("CheckExist", mock.Anything, mock.AnythingOfType("*domain.FbStorageBucket")).Return(domain.FbStorageBucket{}, nil).Once()
		mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.FbStorageBucket")).Return(int64(1), nil).Once()
		u := NewStorageBucketUseCase(mockRepo, time.Second*3)
		affect, err := u.Update(context.TODO(), &want)
		assert.NoError(t, err)
		assert.Equal(t, affect, int64(1))
		mockRepo.AssertExpectations(t)
	})

	t.Run("fail", func(t *testing.T) {
		tempMock := mockDS
		tempMock.ID = 1
		wantErr := errors.New("name is exists")
		mockRepo.On("CheckExist", mock.Anything, mock.AnythingOfType("*domain.FbStorageBucket")).Return(mockDS, nil).Once()
		mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.FbStorageBucket")).Return(int64(0), wantErr).Once()
		u := NewStorageBucketUseCase(mockRepo, time.Second*3)
		affect, err := u.Update(context.TODO(), &tempMock)
		assert.Error(t, err, wantErr)
		assert.Equal(t, affect, int64(0))
	})
}

func TestStorageBucketUseCase_Delete(t *testing.T) {
	mockDataSourceRepo := new(mocks.StorageBucketRepository)
	t.Run("success", func(t *testing.T) {
		mockDataSourceRepo.On("Delete", mock.Anything, mock.AnythingOfType("uint")).Return(int64(1), nil).Once()
		u := NewStorageBucketUseCase(mockDataSourceRepo, time.Second*3)
		_, err := u.Delete(context.TODO(), uint(1))
		assert.NoError(t, err)
		mockDataSourceRepo.AssertExpectations(t)
	})
}

func TestStorageBucketUseCase_FindStorageBucket(t *testing.T) {
	mockDataSourceRepo := new(mocks.StorageBucketRepository)
	result := []domain.FbStorageBucket{
		{
			ID:     1,
			Name:   "1",
			Switch: 1,
			Config: "1",
		},
	}
	mockDataSourceRepo.On("FindStorageBucket", mock.Anything).Return(result, nil).Once()
	u := NewStorageBucketUseCase(mockDataSourceRepo, time.Second*3)
	rows, err := u.FindStorageBucket(context.TODO())
	assert.NoError(t, err)
	assert.Equal(t, rows, result)
	mockDataSourceRepo.AssertExpectations(t)
}

func TestStorageBucketUseCase_GetByID(t *testing.T) {
	mockDataSourceRepo := new(mocks.StorageBucketRepository)
	result := []domain.FbStorageBucket{
		{
			ID:     1,
			Name:   "1",
			Switch: 1,
			Config: "1",
		},
	}
	mockDataSourceRepo.On("GetByID", mock.Anything, mock.AnythingOfType("uint")).Return(result[0], nil).Once()
	u := NewStorageBucketUseCase(mockDataSourceRepo, time.Second*3)
	_, err := u.GetByID(context.TODO(), uint(1))
	assert.NoError(t, err)
	mockDataSourceRepo.AssertExpectations(t)
}
