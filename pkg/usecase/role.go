package usecase

import (
	"context"
	"github.com/fire_boom/domain"
	log "github.com/sirupsen/logrus"
	"time"
)

type RoleUseCase struct {
	roleRepo       domain.RoleRepository
	contextTimeout time.Duration
}

func NewRoleUseCase(a domain.RoleRepository, timeout time.Duration) *RoleUseCase {
	return &RoleUseCase{
		roleRepo:       a,
		contextTimeout: timeout,
	}
}

func (r *RoleUseCase) Store(c context.Context, role *domain.FbRole) (result int64, err error) {
	ctx, cancel := context.WithTimeout(c, r.contextTimeout)
	defer cancel()

	existedRole, err := r.roleRepo.GetByCode(ctx, role.Code)
	if err != nil {
		log.Error("RoleUseCase Store roleRepo.GetByCode err : ", err.Error())
		err = domain.DbCheckNameExistErr
		return
	}
	if existedRole != (domain.FbRole{}) {
		err = domain.DbNameExistErr
		return
	}
	result, err = r.roleRepo.Store(ctx, role)
	if err != nil {
		err = domain.DbFindErr
		return
	}
	return
}

func (r *RoleUseCase) Update(c context.Context, role *domain.FbRole) (affect int64, err error) {
	ctx, cancel := context.WithTimeout(c, r.contextTimeout)
	defer cancel()
	existedRole, err := r.roleRepo.CheckExist(ctx, role)
	if err != nil {
		log.Error("RoleUseCase Update roleRepo.CheckExist err : ", err.Error())
		err = domain.DbCheckNameExistErr
		return
	}
	if existedRole != (domain.FbRole{}) {
		err = domain.DbNameExistErr
		return
	}

	affect, err = r.roleRepo.Update(ctx, role)
	if err != nil {
		err = domain.DbUpdateErr
		return
	}
	return
}

func (r *RoleUseCase) Delete(c context.Context, id uint) (affected int64, err error) {
	ctx, cancel := context.WithTimeout(c, r.contextTimeout)
	defer cancel()

	affected, err = r.roleRepo.Delete(ctx, id)
	if err != nil {
		err = domain.DbDeleteErr
		return
	}
	return
}

func (r *RoleUseCase) FindRoles(c context.Context) (result []domain.FbRole, err error) {
	ctx, cancel := context.WithTimeout(c, r.contextTimeout)
	defer cancel()
	result, err = r.roleRepo.FindRoles(ctx)
	if err != nil {
		err = domain.DbFindErr
		return
	}
	return
}
