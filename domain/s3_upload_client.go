package domain

import (
	"context"
	"time"
)

const MaxUploadSize = 20 * 1024 * 1024 // 20MB
const MaxS3CreationTimeout = time.Duration(time.Second * 30)

type S3Upload struct {
	Name            string `json:"name"`
	EndPoint        string `json:"endPoint"`
	BucketName      string `json:"bucketName"`
	BucketLocation  string `json:"bucketLocation"`
	AccessKeyID     string `json:"accessKeyID"`
	SecretAccessKey string `json:"secretAccessKey"`
	UseSSL          bool   `json:"useSSL"`
}

type UploadResponse struct {
	Files []UploadedFile `json:"Files"`
}

type UploadedFile struct {
	URL        string `json:"url"`
	Name       string `json:"name"`
	MimeTypes  string `json:"mime"`
	Size       int64  `json:"size"`
	CreateTime string `json:"createTime"`
	UpdateTime string `json:"updateTime"`
	IsDir      bool   `json:"isDir"`
}

type S3UploadClientUseCase interface {
	S3UploadFiles(ctx context.Context, s3Upload S3Upload, fileName string) ([]string, error)
	ReName(ctx context.Context, s3Upload S3Upload, oldName, newName string) error
	Delete(ctx context.Context, s3Upload S3Upload, fileName string) error
	Upload(ctx context.Context, s3Upload S3Upload, ossPath, fileName string) error
	Download(ctx context.Context, s3Upload S3Upload, fileName string) ([]string, error)
	Detail(ctx context.Context, s3Upload S3Upload, fileName string) (UploadedFile, error)
}

type S3UploadClientRepository interface {
	S3UploadFiles(ctx context.Context, s3Upload S3Upload, fileName string) ([]string, error)
	CreateFile(ctx context.Context, s3Upload S3Upload, ossPath, fileName string) error
	DeleteFiles(ctx context.Context, s3Upload S3Upload, fileName string) error
	Copy(ctx context.Context, s3Upload S3Upload, oldName, newName string) error
	Download(ctx context.Context, s3Upload S3Upload, fileName string) ([]string, error)
	Detail(ctx context.Context, s3Upload S3Upload, fileName string) (UploadedFile, error)
}
