package usecase

import (
	"context"
	"github.com/fire_boom/domain"
	"time"
)

type FileUseCase struct {
	fileRepo       domain.FileRepository
	contextTimeout time.Duration
}

func NewFileUseCase(f domain.FileRepository, timeout time.Duration) *FileUseCase {
	return &FileUseCase{
		fileRepo:       f,
		contextTimeout: timeout,
	}
}

func (f *FileUseCase) Store(c context.Context, file *domain.File) (result string, err error) {
	ctx, cancel := context.WithTimeout(c, f.contextTimeout)
	defer cancel()

	result, err = f.fileRepo.Store(ctx, file)
	if err != nil {
		err = domain.DbCreateErr
	}
	return
}

func (f *FileUseCase) Delete(c context.Context, fileUUID string) (err error) {
	ctx, cancel := context.WithTimeout(c, f.contextTimeout)
	defer cancel()

	err = f.fileRepo.Delete(ctx, fileUUID)
	if err != nil {
		err = domain.DbUpdateErr
	}
	return err
}

func (f *FileUseCase) GetById(c context.Context, fileUUID string) (file domain.File, err error) {
	ctx, cancel := context.WithTimeout(c, f.contextTimeout)
	defer cancel()

	err = f.fileRepo.Delete(ctx, fileUUID)
	if err != nil {
		err = domain.DbUpdateErr
	}
	return
}
