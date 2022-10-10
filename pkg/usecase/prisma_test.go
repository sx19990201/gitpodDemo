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

func TestPrismaUseCase_Create(t *testing.T) {
	mockRepo := new(mocks.PrismaRepository)
	mockDS := domain.Prisma{
		Name:       "1",
		File:       domain.File{ID: "1"},
		CreateTime: sql.NullString{String: dataFormatStr},
		UpdateTime: sql.NullString{String: dataFormatStr},
		IsDel:      1,
	}
	t.Run("success", func(t *testing.T) {
		tempMockDS := mockDS
		tempMockDS.ID = 0
		mockRepo.On("GetByName", mock.Anything, mock.AnythingOfType("string")).Return(domain.Prisma{}, nil).Once()
		mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Prisma")).Return(nil).Once()

		u := NewPrismaUseCase(mockRepo, time.Second*3)
		err := u.Create(context.TODO(), &tempMockDS)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("name already exists", func(t *testing.T) {
		tempMockDS := mockDS
		tempMockDS.ID = 0
		errWant := errors.New("name already exists")
		mockRepo.On("GetByName", mock.Anything, mock.AnythingOfType("string")).Return(domain.Prisma{Name: "1"}, nil).Once()
		mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Prisma")).Return(errWant).Once()

		u := NewPrismaUseCase(mockRepo, time.Second*3)
		err := u.Create(context.TODO(), &tempMockDS)
		assert.Error(t, err, errWant)
	})
}

func TestPrismaUseCase_Update(t *testing.T) {
	mockRepo := new(mocks.PrismaRepository)
	mockDS := domain.Prisma{
		Name:       "1",
		File:       domain.File{ID: "1"},
		CreateTime: sql.NullString{String: dataFormatStr},
		UpdateTime: sql.NullString{String: dataFormatStr},
		IsDel:      1,
	}
	t.Run("success", func(t *testing.T) {
		want := mockDS
		want.ID = 1
		want.File.ID = "1"
		mockRepo.On("CheckExist", mock.Anything, mock.AnythingOfType("*domain.Prisma")).Return(domain.Prisma{}, nil).Once()
		mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Prisma")).Return(nil).Once()
		u := NewPrismaUseCase(mockRepo, time.Second*3)
		err := u.Update(context.TODO(), &want)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("fail", func(t *testing.T) {
		tempMock := mockDS
		tempMock.ID = 1
		wantErr := errors.New("name is not exists")
		mockRepo.On("CheckExist", mock.Anything, mock.AnythingOfType("*domain.Prisma")).Return(mockDS, nil).Once()
		mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Prisma")).Return(wantErr).Once()
		u := NewPrismaUseCase(mockRepo, time.Second*3)
		err := u.Update(context.TODO(), &tempMock)
		assert.Error(t, err, wantErr)
	})
}

func TestPrismaUseCase_Delete(t *testing.T) {
	mockRepo := new(mocks.PrismaRepository)
	t.Run("success", func(t *testing.T) {
		mockRepo.On("Delete", mock.Anything, mock.AnythingOfType("int64")).Return(nil).Once()
		u := NewPrismaUseCase(mockRepo, time.Second*3)
		err := u.Delete(context.TODO(), int64(1))
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestPrismaUseCase_Fetch(t *testing.T) {
	mockRepo := new(mocks.PrismaRepository)
	result := []domain.Prisma{
		{
			ID:   1,
			Name: "1",
			File: domain.File{ID: "1"},
		},
	}
	mockRepo.On("Fetch", mock.Anything).Return(result, nil).Once()
	u := NewPrismaUseCase(mockRepo, time.Second*3)
	rows, err := u.Fetch(context.TODO())
	assert.NoError(t, err)
	assert.Equal(t, rows, result)
	mockRepo.AssertExpectations(t)
}

func TestPrismaUseCase_GetByID(t *testing.T) {
	mockRepo := new(mocks.PrismaRepository)
	result := []domain.Prisma{
		{
			ID:   1,
			Name: "1",
			File: domain.File{ID: "1"},
		},
	}
	mockRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(result[0], nil).Once()
	u := NewPrismaUseCase(mockRepo, time.Second*3)
	rows, err := u.GetByID(context.TODO(), int64(1))
	assert.NoError(t, err)
	assert.Equal(t, rows, result[0])
	mockRepo.AssertExpectations(t)
}
