package usecase

import (
	"context"
	"github.com/fire_boom/domain"
	log "github.com/sirupsen/logrus"
	"time"
)

type PrismaUseCase struct {
	prismaRepo     domain.PrismaRepository
	contextTimeout time.Duration
}

func NewPrismaUseCase(p domain.PrismaRepository, timeout time.Duration) *PrismaUseCase {
	return &PrismaUseCase{
		prismaRepo:     p,
		contextTimeout: timeout,
	}
}

func (p *PrismaUseCase) Create(c context.Context, prisma *domain.Prisma) (err error) {
	ctx, cancel := context.WithTimeout(c, p.contextTimeout)
	defer cancel()

	existedFile, err := p.prismaRepo.GetByName(ctx, prisma.Name)
	if err != nil {
		log.Error("PrismaUseCase Create prismaRepo.GetByName err : ", err.Error())
		err = domain.DbCheckNameExistErr
		return
	}
	if existedFile.ID != 0 {
		err = domain.DbNameExistErr
		return
	}
	err = p.prismaRepo.Create(ctx, prisma)
	if err != nil {
		err = domain.DbCreateErr
		return
	}
	return
}

func (p *PrismaUseCase) Update(c context.Context, prisma *domain.Prisma) (err error) {
	ctx, cancel := context.WithTimeout(c, p.contextTimeout)
	defer cancel()

	existedFile, err := p.prismaRepo.CheckExist(ctx, prisma)
	if err != nil {
		log.Error("PrismaUseCase Update prismaRepo.CheckExist err : ", err.Error())
		err = domain.DbCheckNameExistErr
		return
	}
	if existedFile.ID != 0 {
		err = domain.DbNameExistErr
		return
	}
	err = p.prismaRepo.Update(ctx, prisma)
	if err != nil {
		err = domain.DbUpdateErr
		return
	}
	return
}

func (p *PrismaUseCase) Fetch(c context.Context) (result []domain.Prisma, err error) {
	ctx, cancel := context.WithTimeout(c, p.contextTimeout)
	defer cancel()

	result, err = p.prismaRepo.Fetch(ctx)
	if err != nil {
		err = domain.DbFindErr
		return
	}

	return
}

func (p *PrismaUseCase) GetByID(c context.Context, id int64) (result domain.Prisma, err error) {
	ctx, cancel := context.WithTimeout(c, p.contextTimeout)
	defer cancel()

	result, err = p.prismaRepo.GetByID(ctx, id)
	if err != nil {
		err = domain.DbFindErr
		return
	}
	return
}
func (p *PrismaUseCase) Delete(c context.Context, id int64) (err error) {
	ctx, cancel := context.WithTimeout(c, p.contextTimeout)
	defer cancel()

	err = p.prismaRepo.Delete(ctx, id)
	if err != nil {
		err = domain.DbDeleteErr
		return
	}
	return
}
