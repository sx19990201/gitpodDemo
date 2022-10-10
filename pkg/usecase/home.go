package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/fire_boom/domain"
	"github.com/fire_boom/utils"
	"github.com/labstack/gommon/log"
	"io/ioutil"
	"time"
)

type HomeUseCase struct {
	dataSourceRepo    domain.DataSourceRepository
	operationsRepo    domain.OperationsRepository
	storageBucketRepo domain.StorageBucketRepository
	authRepo          domain.AuthenticationRepository
	contextTimeout    time.Duration
}

func NewHomeUseCase(timeout time.Duration, dsr domain.DataSourceRepository, or domain.OperationsRepository, sbr domain.StorageBucketRepository, ar domain.AuthenticationRepository) *HomeUseCase {
	return &HomeUseCase{
		dataSourceRepo:    dsr,
		operationsRepo:    or,
		storageBucketRepo: sbr,
		authRepo:          ar,
		contextTimeout:    timeout,
	}
}

func (h *HomeUseCase) GetDateSourceData(c context.Context) (result domain.HomeDataSource, err error) {
	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	dataSources, err := h.dataSourceRepo.FindDataSources(ctx)
	if err != nil {
		err = domain.DbFindErr
		return
	}
	for _, row := range dataSources {
		switch row.SourceType {
		case domain.SourceTypeDB:
			result.DBTotal++
		case domain.SourceTypeRest:
			result.RestTotal++
		case domain.SourceTypeGraphQL:
			result.GraphqlTotal++
		case domain.SourceTypeCustomize:
			result.CustomerTotal++
		}
	}
	return
}

func (h *HomeUseCase) GetApiData(c context.Context) (result domain.HomeApi, err error) {
	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	operations, err := h.operationsRepo.FindOperations(ctx)
	if err != nil {
		err = domain.DbFindErr
		return
	}
	queryOperationArr := make([]domain.FbOperations, 0)
	for _, row := range operations {
		switch row.OperationType {
		case domain.Queries:
			result.QueryTotal++
			queryOperationArr = append(queryOperationArr, row)
		case domain.Mutations:
			result.MutationsTotal++
		case domain.Subscriptions:
			result.SubscriptionsTotal++
		}
	}
	// 查询
	for _, row := range queryOperationArr {
		// 获取该operation设置路径
		path := fmt.Sprintf("%s%s", utils.GetOperationSettingPath(), row.Path)
		// 判断路径是否存在
		if !utils.FileExist(path) {
			// 如果不存在则创建
			continue
		}
		// 读取内容
		settingByte, err := ioutil.ReadFile(path)
		if err != nil {
			log.Info("HomeUseCase GetApiData get query operations Content fail, err := ", err)
			continue
		}
		if len(settingByte) == 0 {
			continue
		}
		var setting domain.OperationSetting
		err = json.Unmarshal(settingByte, &setting)
		if err != nil {
			log.Info("HomeUseCase GetApiData get query operations Content json.Unmarshal fail, err := ", err)
		}
		if !setting.LiveQueryEnable {
			continue
		}
		result.LiveQueryTotal++
	}

	return
}

func (h *HomeUseCase) GetOssData(c context.Context) (result domain.HomeOss, err error) {
	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()
	storageBuckets, err := h.storageBucketRepo.FindStorageBucket(ctx)
	if err != nil {
		err = domain.DbFindErr
		return
	}
	for _, row := range storageBuckets {
		if row.SwitchOn() {
			result.OssTotal++
		}
	}
	return
}

func (h *HomeUseCase) GetAuthData(c context.Context) (result domain.HomeAuth, err error) {
	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()
	auths, err := h.authRepo.FindAuthentication(ctx)
	if err != nil {
		err = domain.DbFindErr
		return
	}
	for _, row := range auths {
		if row.SwitchState != "0" {
			result.AuthTotal++
		}
	}
	result.TotalUser = 153
	result.TodayInsertUser = 17
	return
}
