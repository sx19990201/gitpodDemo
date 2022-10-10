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

func TestRoleUseCase_Store(t *testing.T) {
	mockRepo := new(mocks.RoleRepository)
	mockDS := domain.FbRole{
		Code:       "1",
		Remark:     "1",
		CreateTime: sql.NullString{String: dataFormatStr},
		UpdateTime: sql.NullString{String: dataFormatStr},
		IsDel:      1,
	}
	t.Run("success", func(t *testing.T) {
		tempMockDS := mockDS
		tempMockDS.ID = 0
		mockRepo.On("GetByCode", mock.Anything, mock.AnythingOfType("string")).Return(domain.FbRole{}, nil).Once()
		mockRepo.On("Store", mock.Anything, mock.AnythingOfType("*domain.FbRole")).Return(int64(1), nil).Once()

		u := NewRoleUseCase(mockRepo, time.Second*3)
		_, err := u.Store(context.TODO(), &tempMockDS)
		assert.NoError(t, err)
		assert.Equal(t, mockDS.Code, tempMockDS.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("name already exists", func(t *testing.T) {
		tempMockDS := mockDS
		tempMockDS.ID = 0
		errWant := errors.New("name already exists")
		var affectWant int64 = 0
		mockRepo.On("GetByCode", mock.Anything, mock.AnythingOfType("string")).Return(domain.FbRole{Code: "1"}, nil).Once()
		mockRepo.On("Store", mock.Anything, mock.AnythingOfType("*domain.FbRole")).Return(affectWant, errWant).Once()

		u := NewRoleUseCase(mockRepo, time.Second*3)
		affect, err := u.Store(context.TODO(), &tempMockDS)
		assert.Error(t, err, errWant)
		assert.Equal(t, affect, affectWant)
	})
}

func TestRoleUseCase_Update(t *testing.T) {
	mockRepo := new(mocks.RoleRepository)
	mockDS := domain.FbRole{
		Code:       "1",
		Remark:     "1",
		CreateTime: sql.NullString{String: dataFormatStr},
		UpdateTime: sql.NullString{String: dataFormatStr},
		IsDel:      1,
	}
	t.Run("success", func(t *testing.T) {
		want := mockDS
		want.ID = 1
		want.Code = "123456"
		mockRepo.On("CheckExist", mock.Anything, mock.AnythingOfType("*domain.FbRole")).Return(domain.FbRole{}, nil).Once()
		mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.FbRole")).Return(int64(1), nil).Once()
		u := NewRoleUseCase(mockRepo, time.Second*3)
		affect, err := u.Update(context.TODO(), &want)
		assert.NoError(t, err)
		assert.Equal(t, affect, int64(1))
		mockRepo.AssertExpectations(t)
	})

	t.Run("fail", func(t *testing.T) {
		tempMock := mockDS
		tempMock.ID = 1
		wantErr := errors.New("name is not exists")
		mockRepo.On("CheckExist", mock.Anything, mock.AnythingOfType("*domain.FbRole")).Return(mockDS, nil).Once()
		mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.FbRole")).Return(int64(0), wantErr).Once()
		u := NewRoleUseCase(mockRepo, time.Second*3)
		affect, err := u.Update(context.TODO(), &tempMock)
		assert.Error(t, err, wantErr)
		assert.Equal(t, affect, int64(0))
	})
}

func TestRoleUseCase_Delete(t *testing.T) {
	mockRepo := new(mocks.RoleRepository)
	t.Run("success", func(t *testing.T) {
		mockRepo.On("Delete", mock.Anything, mock.AnythingOfType("uint")).Return(int64(1), nil).Once()
		u := NewRoleUseCase(mockRepo, time.Second*3)
		_, err := u.Delete(context.TODO(), uint(1))
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestRoleUseCase_FindRoles(t *testing.T) {
	mockRepo := new(mocks.RoleRepository)
	result := []domain.FbRole{
		{
			ID:     1,
			Code:   "1",
			Remark: "1",
		},
	}
	mockRepo.On("FindRoles", mock.Anything).Return(result, nil).Once()
	u := NewRoleUseCase(mockRepo, time.Second*3)
	rows, err := u.FindRoles(context.TODO())
	assert.NoError(t, err)
	assert.Equal(t, rows, result)
	mockRepo.AssertExpectations(t)
}
