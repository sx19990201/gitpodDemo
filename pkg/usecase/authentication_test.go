package usecase

import (
	"context"
	"database/sql"
	"errors"
	"github.com/fire_boom/domain"
	"github.com/fire_boom/domain/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestAuthenticationUseCase_Store(t *testing.T) {
	mockRepo := new(mocks.AuthenticationRepository)
	mockDS := domain.FbAuthentication{
		Name:         "1",
		AuthSupplier: "1",
		Config:       "1",
		CreateTime:   sql.NullString{String: dataFormatStr},
		UpdateTime:   sql.NullString{String: dataFormatStr},
		IsDel:        1,
	}
	t.Run("success", func(t *testing.T) {
		tempMockDS := mockDS
		tempMockDS.ID = 0
		mockRepo.On("GetByName", mock.Anything, mock.AnythingOfType("string")).Return(domain.FbAuthentication{}, nil).Once()
		mockRepo.On("Store", mock.Anything, mock.AnythingOfType("*domain.FbAuthentication")).Return(int64(1), nil).Once()

		u := NewAuthenticationUseCase(mockRepo, time.Second*3)
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
		mockRepo.On("GetByName", mock.Anything, mock.AnythingOfType("string")).Return(domain.FbAuthentication{Name: "1"}, nil).Once()
		mockRepo.On("Store", mock.Anything, mock.AnythingOfType("*domain.FbAuthentication")).Return(affectWant, errWant).Once()

		u := NewAuthenticationUseCase(mockRepo, time.Second*3)
		affect, err := u.Store(context.TODO(), &tempMockDS)
		assert.Error(t, err, errWant)
		assert.Equal(t, affect, affectWant)
	})
}

func TestAuthenticationUseCase_Update(t *testing.T) {
	mockRepo := new(mocks.AuthenticationRepository)
	mockDS := domain.FbAuthentication{
		Name:         "1",
		AuthSupplier: "1",
		Config:       "1",
		CreateTime:   sql.NullString{String: dataFormatStr},
		UpdateTime:   sql.NullString{String: dataFormatStr},
		IsDel:        1,
	}
	t.Run("success", func(t *testing.T) {
		want := mockDS
		want.ID = 1
		want.Name = "123456"
		mockRepo.On("CheckExist", mock.Anything, mock.AnythingOfType("*domain.FbAuthentication")).Return(domain.FbAuthentication{}, nil).Once()
		mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.FbAuthentication")).Return(int64(1), nil).Once()
		u := NewAuthenticationUseCase(mockRepo, time.Second*3)
		affect, err := u.Update(context.TODO(), &want)
		assert.NoError(t, err)
		assert.Equal(t, affect, int64(1))
		mockRepo.AssertExpectations(t)
	})

	t.Run("fail", func(t *testing.T) {
		tempMock := mockDS
		tempMock.ID = 1
		wantErr := errors.New("name is exists")
		mockRepo.On("CheckExist", mock.Anything, mock.AnythingOfType("*domain.FbAuthentication")).Return(mockDS, nil).Once()
		mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.FbAuthentication")).Return(int64(0), wantErr).Once()
		u := NewAuthenticationUseCase(mockRepo, time.Second*3)
		affect, err := u.Update(context.TODO(), &tempMock)
		assert.Error(t, err, wantErr)
		assert.Equal(t, affect, int64(0))
	})
}

func TestAuthenticationUseCase_Delete(t *testing.T) {
	mockRepo := new(mocks.AuthenticationRepository)
	t.Run("success", func(t *testing.T) {
		mockRepo.On("Delete", mock.Anything, mock.AnythingOfType("uint")).Return(int64(1), nil).Once()
		u := NewAuthenticationUseCase(mockRepo, time.Second*3)
		_, err := u.Delete(context.TODO(), uint(1))
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestAuthenticationUseCase_FindAuthentication(t *testing.T) {
	mockRepo := new(mocks.AuthenticationRepository)
	result := []domain.FbAuthentication{
		{
			ID:           1,
			Name:         "1",
			AuthSupplier: "1",
			Config:       "1",
		},
	}
	mockRepo.On("FindAuthentication", mock.Anything).Return(result, nil).Once()
	u := NewAuthenticationUseCase(mockRepo, time.Second*3)
	rows, err := u.FindAuthentication(context.TODO())
	assert.NoError(t, err)
	assert.Equal(t, rows, result)
	mockRepo.AssertExpectations(t)
}
