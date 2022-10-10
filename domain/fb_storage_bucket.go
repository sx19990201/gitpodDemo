package domain

import (
	"context"
)

// FbStorageBucket  存储配置
type FbStorageBucket struct {
	ID         int64  `db:"id" json:"id"`
	Name       string `db:"name" json:"name"`     //  名称
	Switch     uint8  `db:"switch" json:"switch"` //  开关: 0-关 1-开
	Config     string `db:"config" json:"config"` //  配置
	CreateTime string `db:"create_time" json:"createTime"`
	UpdateTime string `db:"update_time" json:"updateTime"`
	IsDel      int64  `db:"isDel" json:"isDel"`
}

type StorageBucketResult struct {
	ID         int64            `db:"id" json:"id"`
	Name       string           `db:"name" json:"name"`     //  名称
	Switch     uint8            `db:"switch" json:"switch"` //  开关: 0-关 1-开
	Config     S3UploadProvider `db:"config" json:"config"` //  配置
	CreateTime string           `db:"create_time" json:"createTime"`
	UpdateTime string           `db:"update_time" json:"updateTime"`
	IsDel      int64            `db:"isDel" json:"isDel"`
}

func (f *FbStorageBucket) SwitchOn() bool {
	return f.Switch == SwitchOn
}

type StorageBucketUseCase interface {
	Store(ctx context.Context, f *FbStorageBucket) (int64, error)
	GetByID(ctx context.Context, id uint) (FbStorageBucket, error)
	Update(ctx context.Context, f *FbStorageBucket) (int64, error)
	Delete(ctx context.Context, id uint) (int64, error)
	FindStorageBucket(ctx context.Context) ([]FbStorageBucket, error)
}

type S3UploadProvider struct {
	Name            string `json:"name"`
	Endpoint        string `json:"endpoint"`
	AccessKeyID     string `json:"accessKeyID"`
	SecretAccessKey string `json:"secretAccessKey"`
	BucketName      string `json:"bucketName"`
	BucketLocation  string `json:"bucketLocation"`
	UseSSL          bool   `json:"useSSL"`
}

type StorageBucketRepository interface {
	Store(ctx context.Context, f *FbStorageBucket) (int64, error)
	Update(ctx context.Context, f *FbStorageBucket) (int64, error)
	GetByName(ctx context.Context, name string) (FbStorageBucket, error)
	GetByID(ctx context.Context, id uint) (FbStorageBucket, error)
	CheckExist(ctx context.Context, auth *FbStorageBucket) (FbStorageBucket, error)
	Delete(ctx context.Context, id uint) (int64, error)
	FindStorageBucket(ctx context.Context) ([]FbStorageBucket, error)
}
