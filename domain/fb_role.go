package domain

import (
	"context"
	"database/sql"
)

// FbRole 角色配置
type FbRole struct {
	ID         uint           `db:"id" json:"id"`
	Code       string         `db:"code" json:"code"`     // 角色编码
	Remark     string         `db:"remark" json:"remark"` // 描述
	CreateTime sql.NullString `db:"create_time" json:"-"`
	UpdateTime sql.NullString `db:"update_time" json:"-"`
	IsDel      uint8          `db:"isDel" json:"-"`
}

type RoleUseCase interface {
	Store(ctx context.Context, role *FbRole) (int64, error)
	Update(ctx context.Context, role *FbRole) (int64, error)
	Delete(ctx context.Context, id uint) (int64, error)
	FindRoles(ctx context.Context) ([]FbRole, error)
}

type RoleRepository interface {
	Store(ctx context.Context, role *FbRole) (int64, error)
	Update(ctx context.Context, role *FbRole) (int64, error)
	GetByCode(ctx context.Context, code string) (FbRole, error)
	CheckExist(ctx context.Context, auth *FbRole) (FbRole, error)
	Delete(ctx context.Context, id uint) (int64, error)
	FindRoles(ctx context.Context) ([]FbRole, error)
}
