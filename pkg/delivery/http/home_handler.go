package http

import (
	"database/sql"
	"github.com/fire_boom/domain"
	"github.com/fire_boom/pkg/repository/sqllite"
	"github.com/fire_boom/pkg/usecase"
	"github.com/fire_boom/utils"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type HomeHandler struct {
	homeUC domain.HomeUseCase
}

func InitHomeHandler(e *echo.Echo, db *sql.DB) {
	timeoutContext := time.Duration(utils.GetTimeOut()) * time.Second
	dsr := sqllite.NewDataSourceRepository(db)
	or := sqllite.NewOperationsRepository(db)
	sbr := sqllite.NewStorageBucketRepository(db)
	ar := sqllite.NewAuthenticationRepository(db)
	huc := usecase.NewHomeUseCase(timeoutContext, dsr, or, sbr, ar)
	NewHomeHandler(e, huc)
}

func NewHomeHandler(e *echo.Echo, huc domain.HomeUseCase) {
	handler := &HomeHandler{
		homeUC: huc,
	}
	v1 := e.Group("/api/v1")
	{
		home := v1.Group("/home")
		{
			home.GET("", handler.GetHomeData)
			home.GET("/bulletin", handler.GetBulletin)
		}
	}
}

// GetHomeData @Title GetHomeData
// @Description 获取首页数据
// @Accept  json
// @Tags  home
// @Success 200 {object} domain.Home "查询成功"
// @Failure 400	"查询失败"
// @Router /api/v1/home  [GET]
func (h *HomeHandler) GetHomeData(c echo.Context) (err error) {
	var result domain.Home
	ctx := c.Request().Context()

	result.DataSource, err = h.homeUC.GetDateSourceData(ctx)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, err))
	}
	result.Api, err = h.homeUC.GetApiData(ctx)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, err))
	}
	result.Oss, err = h.homeUC.GetOssData(ctx)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, err))
	}
	result.Auth, err = h.homeUC.GetAuthData(ctx)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, err))
	}
	return c.JSON(http.StatusOK, SuccessWriteResult(result))
}

// GetBulletin @Title GetBulletin
// @Description 获取公告
// @Accept  json
// @Tags  home
// @Success 200 {object} []domain.HomeBulletin "查询成功"
// @Failure 400	"查询失败"
// @Router /api/v1/home/bulletin  [GET]
func (h *HomeHandler) GetBulletin(c echo.Context) (err error) {
	result := make([]domain.HomeBulletin, 0)
	for i := 0; i < 10; i++ {
		result = append(result, domain.HomeBulletin{
			BulletinType: 1,
			Title:        "震惊!!!!",
			Date:         "30分钟前",
		})
	}
	return c.JSON(http.StatusOK, SuccessWriteResult(result))
}
