package http

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/fire_boom/domain"
	"github.com/fire_boom/pkg/repository/sqllite"
	"github.com/fire_boom/pkg/usecase"
	"github.com/fire_boom/utils"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"time"
)

type StorageBucketHandler struct {
	StorageBucketUseCase domain.StorageBucketUseCase
}

func InitStorageBucketRouter(e *echo.Echo, db *sql.DB) {
	storageBucketRepo := sqllite.NewStorageBucketRepository(db)
	timeoutContext := time.Duration(utils.GetTimeOut()) * time.Second
	au := usecase.NewStorageBucketUseCase(storageBucketRepo, timeoutContext)
	NewStorageBucketHandler(e, au)
}

func NewStorageBucketHandler(e *echo.Echo, sUseCase domain.StorageBucketUseCase) {
	handler := &StorageBucketHandler{
		StorageBucketUseCase: sUseCase,
	}
	v1 := e.Group("/api/v1")
	{
		storageBucket := v1.Group("/storageBucket")
		{
			storageBucket.GET("", handler.FindStorageBuckets)
			storageBucket.POST("", handler.Store)
			storageBucket.PUT("", handler.Update)
			storageBucket.DELETE("/:id", handler.Delete)
		}
	}
}

// Store @Title Store
// @Description 添加存储配置信息
// @Accept  json
// @Tags  storageBucket
// @Param name formData string true "存储名称"
// @Param switch formData integer true "开关"
// @Param config formData string true "存储配置"
// @Success 200 "添加成功"
// @Failure 400	"添加失败"
// @Router /api/v1/storageBucket  [POST]
func (s *StorageBucketHandler) Store(c echo.Context) (err error) {
	var storageBucketParam domain.StorageBucketResult
	err = c.Bind(&storageBucketParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.StorageBucketCode, domain.ParamErr))
	}

	if ok := utils.Empty(storageBucketParam.Name); ok {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.StorageBucketCode, domain.ParamNameEmptyErr))
	}
	config, err := json.Marshal(storageBucketParam.Config)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.StorageBucketCode, domain.JsonMarshalErr))
	}
	storageBucket := domain.FbStorageBucket{
		ID:         storageBucketParam.ID,
		Name:       storageBucketParam.Name,
		Switch:     storageBucketParam.Switch,
		Config:     string(config),
		CreateTime: storageBucketParam.CreateTime,
		UpdateTime: storageBucketParam.UpdateTime,
		IsDel:      storageBucketParam.IsDel,
	}

	ctx := c.Request().Context()
	lastID, err := s.StorageBucketUseCase.Store(ctx, &storageBucket)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.StorageBucketCode, err))
	}

	result, err := s.StorageBucketUseCase.GetByID(ctx, uint(lastID))
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.StorageBucketCode, err))
	}

	return c.JSON(http.StatusOK, SuccessWriteResult(result))
}

// Update @Title Update
// @Description 修改存储信息
// @Accept  json
// @Tags  storageBucket
// @Param id formData integer true "id"
// @Param name formData string true "存储名称"
// @Param switch formData integer true "开关"
// @Param config formData string true "存储配置"
// @Success 200 "修改成功"
// @Failure 400	"修改失败"
// @Router /api/v1/storageBucket  [PUT]
func (s *StorageBucketHandler) Update(c echo.Context) (err error) {
	var storageBucketParam domain.StorageBucketResult
	err = c.Bind(&storageBucketParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.StorageBucketCode, domain.ParamErr))
	}

	if err = storageBucketUpdateCheckParams(storageBucketParam); err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.StorageBucketCode, err))
	}

	config, err := json.Marshal(storageBucketParam.Config)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.StorageBucketCode, domain.JsonMarshalErr))
	}
	storageBucket := domain.FbStorageBucket{
		ID:         storageBucketParam.ID,
		Name:       storageBucketParam.Name,
		Switch:     storageBucketParam.Switch,
		Config:     string(config),
		CreateTime: storageBucketParam.CreateTime,
		UpdateTime: storageBucketParam.UpdateTime,
		IsDel:      storageBucketParam.IsDel,
	}
	ctx := c.Request().Context()
	_, err = s.StorageBucketUseCase.Update(ctx, &storageBucket)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.StorageBucketCode, err))
	}

	result, err := s.StorageBucketUseCase.GetByID(ctx, uint(storageBucketParam.ID))
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.StorageBucketCode, err))
	}

	return c.JSON(http.StatusOK, SuccessWriteResult(result))
}

func storageBucketUpdateCheckParams(storageBucket domain.StorageBucketResult) (err error) {
	if ok := utils.Empty(storageBucket.ID); ok {
		return domain.ParamIdEmptyErr
	}
	if ok := utils.Empty(storageBucket.Name); ok {
		return domain.ParamNameEmptyErr
	}
	return
}

// Delete @Title Delete
// @Description 删除存储信息
// @Accept  json
// @Tags  storageBucket
// @Param id formData integer true "存储id"
// @Success 200 "删除成功"
// @Failure 400	"删除失败"
// @Router /api/v1/storageBucket/:id  [DELETE]
func (s *StorageBucketHandler) Delete(c echo.Context) (err error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.StorageBucketCode, domain.ParamErr))
	}

	if ok := utils.Empty(id); ok {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.StorageBucketCode, domain.ParamIdEmptyErr))
	}

	ctx := c.Request().Context()
	affect, err := s.StorageBucketUseCase.Delete(ctx, uint(id))
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.StorageBucketCode, err))
	}

	return c.JSON(http.StatusOK, SuccessWriteResult(fmt.Sprintf("Affected rows : %d", affect)))
}

// FindStorageBuckets @Title FindStorageBuckets
// @Description 查询存储信息
// @Accept  json
// @Tags  storageBucket
// @Success 200 {object} []domain.StorageBucketResult "查询成功"
// @Failure 400	"查询失败"
// @Router /api/v1/storageBucket  [GET]
func (s *StorageBucketHandler) FindStorageBuckets(c echo.Context) (err error) {
	var result []domain.StorageBucketResult
	ctx := c.Request().Context()
	storageBuckets, err := s.StorageBucketUseCase.FindStorageBucket(ctx)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.StorageBucketCode, err))
	}
	for _, bucket := range storageBuckets {
		var config domain.S3UploadProvider
		err = json.Unmarshal([]byte(bucket.Config), &config)
		if err != nil {
			return c.JSON(http.StatusBadRequest, GetResponseErr(domain.StorageBucketCode, domain.JsonUnMarshalErr))
		}
		result = append(result, domain.StorageBucketResult{
			ID:         bucket.ID,
			Name:       bucket.Name,
			Switch:     bucket.Switch,
			Config:     config,
			CreateTime: bucket.CreateTime,
			UpdateTime: bucket.UpdateTime,
			IsDel:      bucket.IsDel,
		})
	}
	return c.JSON(http.StatusOK, SuccessWriteResult(result))
}
