package usecase

import (
	"context"
	"github.com/fire_boom/domain"
	log "github.com/sirupsen/logrus"
	"time"
)

type AuthenticationUseCase struct {
	authenticationRepo domain.AuthenticationRepository
	contextTimeout     time.Duration
}

func NewAuthenticationUseCase(a domain.AuthenticationRepository, timeout time.Duration) *AuthenticationUseCase {
	return &AuthenticationUseCase{
		authenticationRepo: a,
		contextTimeout:     timeout,
	}
}

func (a *AuthenticationUseCase) Store(c context.Context, auth *domain.FbAuthentication) (result int64, err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	existedAuthentication, err := a.authenticationRepo.GetByName(ctx, auth.Name)
	if err != nil {
		log.Error("AuthenticationUseCase Store authenticationRepo.GetByName err : ", err.Error())
		err = domain.DbCheckNameExistErr
		return
	}
	if existedAuthentication.ID != 0 {
		err = domain.DbNameExistErr
		return
	}
	result, err = a.authenticationRepo.Store(ctx, auth)
	if err != nil {
		err = domain.DbCreateErr
		return
	}
	return
}

func (a *AuthenticationUseCase) Update(c context.Context, auth *domain.FbAuthentication) (affected int64, err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	existedAuthentication, err := a.authenticationRepo.CheckExist(ctx, auth)
	if err != nil {
		log.Error("AuthenticationUseCase Update authenticationRepo.CheckExist err : ", err.Error())
		err = domain.DbCheckNameExistErr
		return
	}
	if existedAuthentication.ID != 0 {
		err = domain.DbNameExistErr
		return
	}
	affected, err = a.authenticationRepo.Update(ctx, auth)
	if err != nil {
		err = domain.DbUpdateErr
		return
	}
	return
}

func (a *AuthenticationUseCase) Delete(c context.Context, id uint) (affected int64, err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	affected, err = a.authenticationRepo.Delete(ctx, id)
	if err != nil {
		err = domain.DbDeleteErr
		return
	}
	return
}

func (a *AuthenticationUseCase) FindAuthentication(c context.Context) (result []domain.FbAuthentication, err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	result, err = a.authenticationRepo.FindAuthentication(ctx)
	if err != nil {
		err = domain.DbFindErr
		return
	}
	return
}
