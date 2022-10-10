package http

import (
	"database/sql"
	"github.com/fire_boom/domain"
	"github.com/fire_boom/pkg/repository/sqllite"
	"github.com/fire_boom/pkg/usecase"
	"github.com/fire_boom/utils"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"time"
)

type EnvHandler struct {
	EnvUseCase domain.EnvUseCase
}

func InitEnvRouter(e *echo.Echo, db *sql.DB) {
	timeoutContext := time.Duration(utils.GetTimeOut()) * time.Second
	er := sqllite.NewEnvRepository(db)
	eu := usecase.NewEnvUseCase(er, timeoutContext)
	NewEnvHandler(e, eu)
}

func NewEnvHandler(e *echo.Echo, eUseCase domain.EnvUseCase) {
	handler := &EnvHandler{
		EnvUseCase: eUseCase,
	}
	v1 := e.Group("/api/v1")
	{
		dataSource := v1.Group("/env")
		{
			dataSource.GET("", handler.FindEnvs)
			dataSource.POST("", handler.Store)
			dataSource.PUT("", handler.Update)
			dataSource.DELETE("/:id", handler.Delete)
			dataSource.GET("/keys", handler.FindKeys)
		}
	}
}

// FindKeys @Title FindKeys
// @Description 查询环境变量keys
// @Accept  json
// @Tags  env
// @Param key   formData string true "key"
// @Success 200 {object} []domain.FbEnv "查询成功"
// @Failure 400	"查询失败"
// @Router /api/v1/env/keys  [GET]
func (e *EnvHandler) FindKeys(c echo.Context) (err error) {
	var param domain.FbEnv
	err = c.Bind(&param)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.DataSourceCode, domain.ParamErr))
	}
	ctx := c.Request().Context()
	envs, err := e.EnvUseCase.FindEnvs(ctx, param.Key)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.DataSourceCode, err))
	}
	result := make([]string, 0)
	for _, env := range envs {
		result = append(result, env.Key)
	}
	return c.JSON(http.StatusOK, SuccessWriteResult(result))
}

// FindEnvs @Title FindEnvs
// @Description 查询环境变量信息
// @Accept  json
// @Tags  env
// @Success 200 {object} []domain.FbEnv "查询成功"
// @Failure 400	"查询失败"
// @Router /api/v1/env  [GET]
func (e *EnvHandler) FindEnvs(c echo.Context) (err error) {
	ctx := c.Request().Context()
	envs, err := e.EnvUseCase.FindEnvs(ctx, "")
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.DataSourceCode, err))
	}
	return c.JSON(http.StatusOK, SuccessWriteResult(envs))
}

// Update @Title Update
// @Description 修改环境变量
// @Accept  json
// @Tags  env
// @Param id      formData integer true "id"
// @Param key     formData string true "key"
// @Param devEnv  formData string true "开发环境"
// @Param proEnv  formData string true "生产环境"
// @Param envType formData integer true "环境类型 1-环境变量 2-系统变量"
// @Success 200 "修改成功"
// @Failure 400	"修改失败"
// @Router /api/v1/env  [PUT]
func (e *EnvHandler) Update(c echo.Context) (err error) {
	var param domain.FbEnv
	err = c.Bind(&param)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.DataSourceCode, domain.ParamErr))
	}
	ctx := c.Request().Context()
	result, err := e.EnvUseCase.Update(ctx, &param)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.DataSourceCode, err))
	}
	return c.JSON(http.StatusOK, SuccessWriteResult(result))
}

// Store @Title Store
// @Description 添加环境变量
// @Accept  json
// @Tags  env
// @Param key     formData string true "key"
// @Param devEnv  formData string true "开发环境"
// @Param proEnv  formData string true "生产环境"
// @Param envType formData integer true "环境类型 1-环境变量 2-系统变量"
// @Success 200 "添加成功"
// @Failure 400	"添加失败"
// @Router /api/v1/env  [POST]
func (e *EnvHandler) Store(c echo.Context) (err error) {
	var param domain.FbEnv
	err = c.Bind(&param)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.DataSourceCode, domain.ParamErr))
	}
	if param.Key == "" || (param.ProEnv == "" || param.DevEnv == "") {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.DataSourceCode, domain.ParamNameEmptyErr))
	}
	ctx := c.Request().Context()
	result, err := e.EnvUseCase.Store(ctx, &param)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.DataSourceCode, err))
	}
	return c.JSON(http.StatusOK, SuccessWriteResult(result))
}

// Delete @Title Delete
// @Description 删除环境变量
// @Accept  json
// @Tags  env
// @Param id formData integer true "id"
// @Success 200 "删除成功"
// @Failure 400	"删除失败"
// @Router /api/v1/env/:id  [DELETE]
func (e *EnvHandler) Delete(c echo.Context) (err error) {
	idParam := c.Param("id")
	id, _ := strconv.Atoi(idParam)
	if id == 0 {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.DataSourceCode, domain.ParamIdEmptyErr))
	}
	ctx := c.Request().Context()
	result, err := e.EnvUseCase.Delete(ctx, uint(id))
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.DataSourceCode, err))
	}
	return c.JSON(http.StatusOK, SuccessWriteResult(result))
}
