package usecase

import (
	"context"
	"github.com/fire_boom/domain"
	log "github.com/sirupsen/logrus"
	"time"
)

type EnvUseCase struct {
	envRepo        domain.EnvRepository
	contextTimeout time.Duration
}

func NewEnvUseCase(e domain.EnvRepository, timeout time.Duration) *EnvUseCase {
	return &EnvUseCase{
		envRepo:        e,
		contextTimeout: timeout,
	}
}

func (e *EnvUseCase) Store(c context.Context, f *domain.FbEnv) (result int64, err error) {
	ctx, cancel := context.WithTimeout(c, e.contextTimeout)
	defer cancel()

	existedEnv, err := e.envRepo.GetByKey(ctx, f.Key)
	if err != nil {
		log.Error("EnvUseCase Store envRepo.GetByKey err : ", err.Error())
		err = domain.DbCheckNameExistErr
		return
	}
	if existedEnv != (domain.FbEnv{}) {
		err = domain.DbNameExistErr
		return
	}
	result, err = e.envRepo.Store(ctx, f)
	if err != nil {
		err = domain.DbCreateErr
		return
	}
	return
}

func (e *EnvUseCase) GetByKey(c context.Context, key string) (result domain.FbEnv, err error) {
	ctx, cancel := context.WithTimeout(c, e.contextTimeout)
	defer cancel()
	result, err = e.envRepo.GetByKey(ctx, key)
	if err != nil {
		log.Error("EnvUseCase GetByKey envRepo.GetByKey err : ", err.Error())
		err = domain.DbCheckNameExistErr
		return
	}
	return
}

func (e *EnvUseCase) Update(c context.Context, f *domain.FbEnv) (result int64, err error) {
	ctx, cancel := context.WithTimeout(c, e.contextTimeout)
	defer cancel()
	existedEnv, err := e.envRepo.Exist(ctx, int64(f.ID), f.Key)
	if err != nil {
		log.Error("EnvUseCase Update envRepo.GetByKey err : ", err.Error())
		err = domain.DbCheckNameExistErr
		return
	}

	if existedEnv != (domain.FbEnv{}) {
		err = domain.DbNameExistErr
		return
	}

	result, err = e.envRepo.Update(ctx, f)
	if err != nil {
		err = domain.DbUpdateErr
		return
	}
	return
}

func (e *EnvUseCase) Delete(c context.Context, id uint) (result int64, err error) {
	ctx, cancel := context.WithTimeout(c, e.contextTimeout)
	defer cancel()

	result, err = e.envRepo.Delete(ctx, id)
	if err != nil {
		err = domain.DbDeleteErr
		return
	}
	return
}

func (e *EnvUseCase) FindEnvs(c context.Context, key string) (result []domain.FbEnv, err error) {
	ctx, cancel := context.WithTimeout(c, e.contextTimeout)
	defer cancel()

	result, err = e.envRepo.FindEnvs(ctx, key)
	if err != nil {
		err = domain.DbFindErr
		return
	}
	return
}
