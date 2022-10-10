package usecase

import (
	"context"
	"github.com/fire_boom/domain"
	log "github.com/sirupsen/logrus"
	"time"
)

type S3UploadClientUseCase struct {
	s3UploadClientRepository domain.S3UploadClientRepository
	contextTimeout           time.Duration
}

func NewS3UploadClientUseCase(a domain.S3UploadClientRepository, timeout time.Duration) *S3UploadClientUseCase {
	return &S3UploadClientUseCase{
		s3UploadClientRepository: a,
		contextTimeout:           timeout,
	}
}

func (s *S3UploadClientUseCase) S3UploadFiles(c context.Context, s3Upload domain.S3Upload, fileName string) ([]string, error) {
	ctx, cancel := context.WithTimeout(c, s.contextTimeout)
	defer cancel()

	return s.s3UploadClientRepository.S3UploadFiles(ctx, s3Upload, fileName)
}

func (s *S3UploadClientUseCase) ReName(c context.Context, s3Upload domain.S3Upload, oldName, newName string) (err error) {
	ctx, cancel := context.WithTimeout(c, s.contextTimeout)
	defer cancel()

	err = s.s3UploadClientRepository.Copy(ctx, s3Upload, oldName, newName)
	if err != nil {
		log.Error("S3UploadClientUseCase ReName s3UploadClientRepository.Copy err : ", err.Error())
		return
	}
	return
}

func (s *S3UploadClientUseCase) Delete(c context.Context, s3Upload domain.S3Upload, fileName string) (err error) {
	ctx, cancel := context.WithTimeout(c, s.contextTimeout)
	defer cancel()
	// 删除原文件
	err = s.s3UploadClientRepository.DeleteFiles(ctx, s3Upload, fileName)
	if err != nil {
		log.Error("S3UploadClientUseCase Delete s3UploadClientRepository.DeleteFiles err : ", err.Error())
		return
	}
	return
}

func (s *S3UploadClientUseCase) Download(c context.Context, s3Upload domain.S3Upload, fileName string) (result []string, err error) {
	ctx, cancel := context.WithTimeout(c, s.contextTimeout)
	defer cancel()
	result, err = s.s3UploadClientRepository.Download(ctx, s3Upload, fileName)
	if err != nil {
		err = domain.OssDownloadErr
		return
	}
	return
}

func (s *S3UploadClientUseCase) Detail(c context.Context, s3Upload domain.S3Upload, fileName string) (result domain.UploadedFile, err error) {
	ctx, cancel := context.WithTimeout(c, s.contextTimeout)
	defer cancel()
	result, err = s.s3UploadClientRepository.Detail(ctx, s3Upload, fileName)
	if err != nil {
		err = domain.OssFindErr
		return
	}
	return
}

func (s *S3UploadClientUseCase) Upload(c context.Context, s3Upload domain.S3Upload, ossPath, fileName string) (err error) {
	ctx, cancel := context.WithTimeout(c, s.contextTimeout)
	defer cancel()
	err = s.s3UploadClientRepository.CreateFile(ctx, s3Upload, ossPath, fileName)
	if err != nil {
		err = domain.OssUploadErr
		return
	}
	return
}
