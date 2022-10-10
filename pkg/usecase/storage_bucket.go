package usecase

import (
	"context"
	"github.com/fire_boom/domain"
	log "github.com/sirupsen/logrus"
	"time"
)

type StorageBucketUseCase struct {
	storageBucketRepo domain.StorageBucketRepository
	contextTimeout    time.Duration
}

func NewStorageBucketUseCase(s domain.StorageBucketRepository, timeout time.Duration) *StorageBucketUseCase {
	return &StorageBucketUseCase{
		storageBucketRepo: s,
		contextTimeout:    timeout,
	}
}

func (s *StorageBucketUseCase) Store(c context.Context, f *domain.FbStorageBucket) (result int64, err error) {
	ctx, cancel := context.WithTimeout(c, s.contextTimeout)
	defer cancel()

	existedStorageBucket, err := s.storageBucketRepo.GetByName(ctx, f.Name)
	if err != nil {
		log.Error("StorageBucketUseCase Store storageBucketRepo.GetByName err : ", err.Error())
		err = domain.DbCheckNameExistErr
		return
	}
	if existedStorageBucket != (domain.FbStorageBucket{}) {
		err = domain.DbNameExistErr
		return
	}
	result, err = s.storageBucketRepo.Store(ctx, f)
	if err != nil {
		err = domain.DbCreateErr
		return
	}
	return
}

func (s *StorageBucketUseCase) GetByID(c context.Context, id uint) (result domain.FbStorageBucket, err error) {
	ctx, cancel := context.WithTimeout(c, s.contextTimeout)
	defer cancel()

	result, err = s.storageBucketRepo.GetByID(ctx, id)
	if err != nil {
		err = domain.DbFindErr
		return
	}
	return
}

func (s *StorageBucketUseCase) Update(c context.Context, f *domain.FbStorageBucket) (affect int64, err error) {
	ctx, cancel := context.WithTimeout(c, s.contextTimeout)
	defer cancel()
	existedStorageBucket, err := s.storageBucketRepo.CheckExist(ctx, f)
	if err != nil {
		log.Error("StorageBucketUseCase Update storageBucketRepo.CheckExist err : ", err.Error())
		err = domain.DbCheckNameExistErr
		return
	}
	if existedStorageBucket != (domain.FbStorageBucket{}) {
		err = domain.DbNameExistErr
		return
	}
	affect, err = s.storageBucketRepo.Update(ctx, f)
	if err != nil {
		err = domain.DbUpdateErr
		return
	}
	return
}

func (s *StorageBucketUseCase) Delete(c context.Context, id uint) (affected int64, err error) {
	ctx, cancel := context.WithTimeout(c, s.contextTimeout)
	defer cancel()

	affected, err = s.storageBucketRepo.Delete(ctx, id)
	if err != nil {
		err = domain.DbDeleteErr
		return
	}
	return
}

func (s *StorageBucketUseCase) FindStorageBucket(c context.Context) (result []domain.FbStorageBucket, err error) {
	ctx, cancel := context.WithTimeout(c, s.contextTimeout)
	defer cancel()

	result, err = s.storageBucketRepo.FindStorageBucket(ctx)
	if err != nil {
		err = domain.DbFindErr
		return
	}

	return
}
