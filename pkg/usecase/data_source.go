package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/fire_boom/domain"
)

type DataSourceUseCase struct {
	dataSourceRepo domain.DataSourceRepository
	contextTimeout time.Duration
}

func NewDataSourceUseCase(d domain.DataSourceRepository, timeout time.Duration) *DataSourceUseCase {
	return &DataSourceUseCase{
		dataSourceRepo: d,
		contextTimeout: timeout,
	}
}

func (d *DataSourceUseCase) Store(c context.Context, f *domain.FbDataSource) (result int64, err error) {
	ctx, cancel := context.WithTimeout(c, d.contextTimeout)
	defer cancel()

	existedDataSource, err := d.dataSourceRepo.GetByName(ctx, f.Name)
	if err != nil {
		log.Error("DataSourceUseCase Store dataSourceRepo.GetByName err : ", err.Error())
		err = domain.DbCheckNameExistErr
		return
	}
	if existedDataSource != (domain.FbDataSource{}) {
		err = domain.DbNameExistErr
		return
	}
	result, err = d.dataSourceRepo.Store(ctx, f)
	if err != nil {
		err = domain.DbCreateErr
		return
	}
	return
}

func (d *DataSourceUseCase) Update(c context.Context, f *domain.FbDataSource) (affect int64, err error) {
	ctx, cancel := context.WithTimeout(c, d.contextTimeout)
	defer cancel()
	existedDataSource, err := d.dataSourceRepo.CheckExist(ctx, f)
	if err != nil {
		log.Error("DataSourceUseCase Update dataSourceRepo.CheckExist err : ", err.Error())
		err = domain.DbCheckNameExistErr
		return
	}
	if existedDataSource != (domain.FbDataSource{}) {
		err = domain.DbNameExistErr
		return
	}
	affect, err = d.dataSourceRepo.Update(ctx, f)
	if err != nil {
		err = domain.DbUpdateErr
		return
	}
	return
}

func (d *DataSourceUseCase) Delete(c context.Context, id uint) (affected int64, err error) {
	ctx, cancel := context.WithTimeout(c, d.contextTimeout)
	defer cancel()

	affected, err = d.dataSourceRepo.Delete(ctx, id)
	if err != nil {
		err = domain.DbDeleteErr
		return
	}
	return
}

func (d *DataSourceUseCase) FindDataSources(c context.Context) (result []domain.FbDataSource, err error) {
	ctx, cancel := context.WithTimeout(c, d.contextTimeout)
	defer cancel()

	result, err = d.dataSourceRepo.FindDataSources(ctx)
	if err != nil {
		err = domain.DbFindErr
		return
	}
	return
}

func (d *DataSourceUseCase) GetByID(c context.Context, id uint) (result domain.FbDataSource, err error) {
	ctx, cancel := context.WithTimeout(c, d.contextTimeout)
	defer cancel()
	result, err = d.dataSourceRepo.GetByID(ctx, id)
	if err != nil {
		err = domain.DbFindErr
		return
	}
	return
}

func (d *DataSourceUseCase) GetPrismaSchema(c context.Context, id uint) (result string, err error) {
	ctx, cancel := context.WithTimeout(c, d.contextTimeout)
	defer cancel()
	dataSource, err := d.GetByID(ctx, id)
	if err != nil {
		return
	}
	dbConfig := dataSource.GetDbConfig()
	result = fmt.Sprintf(`datasource db {
    	// could be postgresql or mysql
    	provider = "%s"
    	url      = "%s"
	}
	generator db {
    	provider = "go run github.com/prisma/prisma-client-go"
	}
`, strings.ToLower(dbConfig.DBType), strings.ToLower(dbConfig.DBType)+"://"+dbConfig.DatabaseURL.Val)
	return
}
