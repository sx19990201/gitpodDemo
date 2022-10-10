package domain

import (
	"context"
)

const (
	DevEnv         = "0"
	ProEnv         = "1"
	EnvVariable    = 1
	SystemVariable = 2
)

// FbEnv 环境变量
type FbEnv struct {
	ID         uint   `db:"id" json:"id"`
	Key        string `db:"key" json:"key"`          //环境变量key
	DevEnv     string `db:"dev_env" json:"devEnv"`   //开发环境
	ProEnv     string `db:"pro_env" json:"proEnv"`   //生产环境
	EnvType    uint8  `db:"env_type" json:"envType"` //环境类型: 1-环境变量 2-系统变量
	CreateTime string `db:"create_time" json:"createTime"`
	UpdateTime string `db:"update_time" json:"updateTime"`
	IsDel      uint8  `db:"is_del" json:"isDel"`
}

type FbEnvResp struct {
	System []FbEnv `json:"system"`
	Env    []FbEnv `json:"env"`
}

type EnvUseCase interface {
	Store(ctx context.Context, f *FbEnv) (int64, error)
	GetByKey(ctx context.Context, key string) (FbEnv, error)
	Update(ctx context.Context, f *FbEnv) (int64, error)
	Delete(ctx context.Context, id uint) (int64, error)
	FindEnvs(ctx context.Context, key string) ([]FbEnv, error)
}

type EnvRepository interface {
	Store(ctx context.Context, f *FbEnv) (int64, error)
	Update(ctx context.Context, f *FbEnv) (int64, error)
	GetByKey(ctx context.Context, key string) (FbEnv, error)
	Exist(ctx context.Context, id int64, key string) (FbEnv, error)
	Delete(ctx context.Context, id uint) (int64, error)
	FindEnvs(ctx context.Context, key string) ([]FbEnv, error)
}
