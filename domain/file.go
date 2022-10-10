package domain

import (
	"context"
	"database/sql"
)

type File struct {
	ID         string         `json:"id"`
	Name       string         `json:"name"`
	Path       string         `json:"path"` // 相对路径
	CreateTime sql.NullString `db:"create_time" json:"-"`
	UpdateTime sql.NullString `db:"update_time" json:"-"`
	IsDel      uint8          `db:"isDel" json:"-"`
}

type FileUseCase interface {
	Store(ctx context.Context, f *File) (string, error)
	Delete(ctx context.Context, fileUUID string) error
	GetById(ctx context.Context, fileUUID string) (File, error)
}

type FileRepository interface {
	// Store 上传文件到 static 目录，meta 信息存储到表
	Store(ctx context.Context, f *File) (string, error)
	Delete(ctx context.Context, fileUUID string) error
	GetByID(ctx context.Context, fileUUID uint) (File, error)
	//GetByName(ctx context.Context, name string) (File, error)
	//CheckExist(ctx context.Context, auth *File) (File, error)
	//FileWrite(ctx context.Context, fileParam *File) error
	//FileRead(ctx context.Context, path string) ([]byte, error)
}
