package usecase

import (
	"context"
	"fmt"
	"github.com/fire_boom/domain"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

type OperationsUseCase struct {
	operationsRepo domain.OperationsRepository
	contextTimeout time.Duration
}

func NewOperationsUseCase(or domain.OperationsRepository, timeout time.Duration) *OperationsUseCase {
	return &OperationsUseCase{
		operationsRepo: or,
		contextTimeout: timeout,
	}
}

func (o *OperationsUseCase) Store(c context.Context, f *domain.FbOperations) (result int64, err error) {
	ctx, cancel := context.WithTimeout(c, o.contextTimeout)
	defer cancel()

	existedOperations, err := o.operationsRepo.GetByPath(ctx, f.Path)
	if err != nil {
		log.Error("OperationsUseCase Store operationsRepo.GetByPath err : ", err.Error())
		err = domain.DbCheckNameExistErr
		return
	}
	if existedOperations != (domain.FbOperations{}) {
		err = domain.DbNameExistErr
		return
	}
	result, err = o.operationsRepo.Store(ctx, f)
	if err != nil {
		err = domain.DbCreateErr
		return
	}
	return
}

func (o *OperationsUseCase) GetByID(c context.Context, id int64) (result domain.FbOperations, err error) {
	ctx, cancel := context.WithTimeout(c, o.contextTimeout)
	defer cancel()
	result, err = o.operationsRepo.GetByID(ctx, id)
	if err != nil {
		err = domain.DbFindErr
		return
	}
	return
}

func (o *OperationsUseCase) GetByPath(c context.Context, path string) (result domain.FbOperations, err error) {
	ctx, cancel := context.WithTimeout(c, o.contextTimeout)
	defer cancel()
	result, err = o.operationsRepo.GetByPath(ctx, path)
	if err != nil {
		err = domain.DbFindErr
		return
	}
	return
}

func (o *OperationsUseCase) Update(c context.Context, f *domain.FbOperations) (result int64, err error) {
	ctx, cancel := context.WithTimeout(c, o.contextTimeout)
	defer cancel()
	existedOperations, err := o.operationsRepo.GetByID(ctx, f.ID)
	if err != nil {
		log.Error("OperationsUseCase Update operationsRepo.GetByID err : ", err.Error())
		err = domain.DbCheckNameExistErr
		return
	}
	if existedOperations == (domain.FbOperations{}) {
		err = domain.DbNameNotExistErr
		return
	}
	result, err = o.operationsRepo.Update(ctx, f)
	if err != nil {
		err = domain.DbUpdateErr
		return
	}
	return
}

func (o *OperationsUseCase) Delete(c context.Context, id int64) (result int64, err error) {
	ctx, cancel := context.WithTimeout(c, o.contextTimeout)
	defer cancel()

	result, err = o.operationsRepo.Delete(ctx, id)
	if err != nil {
		err = domain.DbDeleteErr
		return
	}
	return
}
func (o *OperationsUseCase) BatchDelete(c context.Context, opers []domain.FbOperations) (err error) {
	ctx, cancel := context.WithTimeout(c, o.contextTimeout)
	defer cancel()

	ids := ""
	for _, row := range opers {
		ids = fmt.Sprintf("%s, %v", ids, row.ID)
	}
	if ids == "" {
		return
	}
	err = o.operationsRepo.BatchDelete(ctx, ids)
	if err != nil {
		err = domain.DbDeleteErr
		return
	}
	return
}
func (o *OperationsUseCase) FindOperations(c context.Context) (result []domain.FbOperations, err error) {
	ctx, cancel := context.WithTimeout(c, o.contextTimeout)
	defer cancel()

	result, err = o.operationsRepo.FindOperations(ctx)
	if err != nil {
		err = domain.DbFindErr
		return
	}
	return
}

func (o *OperationsUseCase) ReName(c context.Context, oldName, newName string, enable int64) (err error) {
	ctx, cancel := context.WithTimeout(c, o.contextTimeout)
	defer cancel()

	err = o.operationsRepo.ReName(ctx, oldName, newName, enable)
	if err != nil {
		err = domain.DbUpdateErr
		return
	}
	return
}

func (o *OperationsUseCase) ReDicName(c context.Context, oldName, newName string) (err error) {
	ctx, cancel := context.WithTimeout(c, o.contextTimeout)
	defer cancel()

	opertions, err := o.operationsRepo.FindByPath(ctx, oldName)
	if err != nil {
		err = domain.DbUpdateErr
		return
	}

	for _, row := range opertions {
		newPath := strings.ReplaceAll(row.Path, oldName, newName)
		err = o.operationsRepo.ReName(ctx, row.Path, newPath, row.Enable)
		if err != nil {
			err = domain.DbUpdateErr
			return
		}
	}
	return
}

func (o *OperationsUseCase) FindByPath(c context.Context, path string) (result []domain.FbOperations, err error) {
	ctx, cancel := context.WithTimeout(c, o.contextTimeout)
	defer cancel()

	result, err = o.operationsRepo.FindByPath(ctx, path)
	if err != nil {
		err = domain.DbUpdateErr
		return
	}
	return
}
