/**
domain 类似于 DDD 中领域服务, 这里主要存储领域模型、模型对应的方法、定义与 delivery、repositpry 层交互接口
*/
package domain

import (
	"context"
	"database/sql"
)

type Prisma struct {
	ID         int64          `json:"id"`
	Name       string         `json:"name"`
	File       File           `json:"file"`
	CreateTime sql.NullString `json:"-"`
	UpdateTime sql.NullString `json:"-"`
	IsDel      int8           `json:"-"`
}

type PrismaUseCase interface {
	Create(ctx context.Context, p *Prisma) error
	Update(ctx context.Context, p *Prisma) error
	Fetch(ctx context.Context) ([]Prisma, error)
	GetByID(ctx context.Context, id int64) (Prisma, error)
	Delete(ctx context.Context, id int64) error
}

type PrismaRepository interface {
	Create(ctx context.Context, p *Prisma) error
	Update(ctx context.Context, p *Prisma) error
	Fetch(ctx context.Context) ([]Prisma, error)
	GetByID(ctx context.Context, id int64) (Prisma, error)
	GetByName(ctx context.Context, name string) (Prisma, error)
	CheckExist(ctx context.Context, auth *Prisma) (Prisma, error)
	Delete(ctx context.Context, id int64) error
}
