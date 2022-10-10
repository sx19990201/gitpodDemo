package http

import (
	"github.com/fire_boom/domain"
	"github.com/fire_boom/pkg/usecase"
	"github.com/fire_boom/utils"
	"github.com/labstack/echo/v4"
	"github.com/prisma/prisma-client-go/engine"
	//"github.com/prisma/prisma-client-go/engine"

	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

type DraftHandler struct {
	DataSourceUseCase domain.DataSourceUseCase
}

func InitDraftHandler(e *echo.Echo, dsr domain.DataSourceRepository) {
	timeoutContext := time.Duration(utils.GetTimeOut()) * time.Second
	au := usecase.NewDataSourceUseCase(dsr, timeoutContext)
	NewDraftHandler(e, au)
}

func NewDraftHandler(e *echo.Echo, dUseCase domain.DataSourceUseCase) {
	handler := &DraftHandler{
		DataSourceUseCase: dUseCase,
	}
	v1 := e.Group("/api/v1")
	{
		authentication := v1.Group("/draft")
		{
			authentication.GET("/:id", handler.GetDraftById) // 获取草稿
			authentication.PUT("", handler.Update)           // 更新草稿
			authentication.DELETE("", handler.Delete)        // 删除草稿

		}
	}
}

// GetDraftById @Title GetDraftById
// @Description 根据数据源id获取草稿
// @Accept  json
// @Tags  draft
// @Param id path integer true "数据源id"
// @Success 200 {object} string  "查询成功"
// @Failure 400	"查询失败"
// @Router /api/v1/draft/:id  [GET]
func (d *DraftHandler) GetDraftById(c echo.Context) error {
	ctx := c.Request().Context()
	// 获取数据源id
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.DraftCode, domain.ParamErr))
	}
	if id == 0 {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.DraftCode, domain.ParamIdEmptyErr))
	}
	// 读取该数据源草稿信息
	//path := fmt.Sprintf("%s/%d", utils.GetDraftPrefixPath(), id)
	path := utils.GetPrismaSchemaFilePath()
	_, err = os.Stat(path)
	// 草稿不存在,创建草稿
	if os.IsNotExist(err) {
		// 创建草稿文件
		content, err := d.DataSourceUseCase.GetPrismaSchema(ctx, uint(id))
		if err != nil {
			return c.JSON(http.StatusBadRequest, GetResponseErr(domain.DraftCode, domain.DbIntrospectionErr))
		}
		err = utils.WriteFile(path, content)
		if err != nil {
			return c.JSON(http.StatusBadRequest, GetResponseErr(domain.DraftCode, domain.DbIntrospectionErr))
		}
		engine.Pull(path)
	}
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.DraftCode, domain.FileReadErr))
	}
	return c.JSON(http.StatusOK, SuccessWriteResult(string(content)))
}

// Update @Title Update
// @Description 更新草稿内容
// @Accept  json
// @Tags  draft
// @Param content formData string true "内容"
// @Param id formData integer true "内容"
// @Success 200 "更新成功"
// @Failure 400	"更新失败"
// @Router /api/v1/draft  [PUT]
func (d *DraftHandler) Update(c echo.Context) error {
	param := struct {
		Content string `json:"content"`
		ID      int64  `json:"id"`
	}{}
	err := c.Bind(&param)

	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.DraftCode, domain.ParamErr))
	}
	// path := fmt.Sprintf("%s/%d", utils.GetDraftPrefixPath(), param.ID)
	err = utils.WriteFile(utils.GetPrismaSchemaFilePath(), param.Content)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.DraftCode, domain.FileWriteErr))
	}
	return c.JSON(http.StatusOK, SuccessResult())
}

// Delete @Title Delete
// @Description 删除草稿
// @Accept  json
// @Tags  draft
// @Param id path integer true "数据源id"
// @Success 200 "删除成功"
// @Failure 400	"删除失败"
// @Router /api/v1/draft/:id  [DELETE]
func (d *DraftHandler) Delete(c echo.Context) error {
	//id, err := strconv.Atoi(c.Param("id"))
	//if err != nil {
	//	return c.JSON(http.StatusBadRequest, GetResponseErr(domain.DraftCode, domain.ParamErr))
	//}
	//path := fmt.Sprintf("%s/%d", utils.GetDraftPrefixPath(), id)
	err := os.Remove(utils.GetPrismaSchemaFilePath())
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.DraftCode, domain.FileDeleteErr))
	}

	return c.JSON(http.StatusOK, SuccessResult())
}
