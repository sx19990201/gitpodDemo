package http

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/fire_boom/utils"
	"github.com/prisma/prisma-client-go/engine"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fire_boom/domain"
	"github.com/fire_boom/pkg/repository/sqllite"
	"github.com/fire_boom/pkg/usecase"
	"github.com/labstack/echo/v4"
)

type PrismaHandler struct {
	PrismaUseCase     domain.PrismaUseCase
	DataSourceUseCase domain.DataSourceUseCase
}

func InitPrismaUseCaseRouter(e *echo.Echo, db *sql.DB) {
	prismaRepo := sqllite.NewPrismaRepository(db)
	timeoutContext := time.Duration(utils.GetTimeOut()) * time.Second
	prismaAU := usecase.NewPrismaUseCase(prismaRepo, timeoutContext)

	dbRepo := sqllite.NewDataSourceRepository(db)
	au := usecase.NewDataSourceUseCase(dbRepo, timeoutContext)
	NewPrismaHandler(e, prismaAU, au)
}

func NewPrismaHandler(e *echo.Echo, pUseCase domain.PrismaUseCase, au domain.DataSourceUseCase) {
	handler := &PrismaHandler{
		PrismaUseCase:     pUseCase,
		DataSourceUseCase: au,
	}
	v1 := e.Group("/api/v1")
	{
		prsimaHandler := v1.Group("/prisma")
		{
			//authentication.GET("", handler.Fetch)
			//authentication.GET("/:id", handler.GetByID)
			prsimaHandler.POST("", handler.Create)       // TODO 这三个暂时没啥用
			prsimaHandler.PUT("", handler.Update)        // TODO 这三个暂时没啥用
			prsimaHandler.DELETE("/:id", handler.Delete) // TODO 这三个暂时没啥用
			// 实时内省
			prsimaHandler.POST("/migrate/:id", handler.Migrate)            // 迁移
			prsimaHandler.GET("/dmf/:id", handler.GetDMF)                  // 获取dmf
			prsimaHandler.GET("/introspection/:id", handler.Introspection) // 内省
		}
	}
}

//Introspection @Title Introspection
//@Description 内省
//@Accept  json
//@Tags  prisma
//@Param id path integer true "数据源id"
//@Success 200  "内省成功"
//@Failure 400  "内省失败"
//@Router /api/v1/prisma/introspection/:id  [POST]
func (p *PrismaHandler) Introspection(c echo.Context) (err error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id == 0 {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.DraftCode, domain.ParamErr))
	}
	ctx := c.Request().Context()
	// 创建草稿文件
	content, err := p.DataSourceUseCase.GetPrismaSchema(ctx, uint(id))
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.DraftCode, domain.DbSchemaErr))
	}

	//path := filepath.ToSlash(fmt.Sprintf("%s/%s/%s", utils.GetRootPath(), utils.GetDraftPrefixPath(), id))
	_, err = engine.Pull(content)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.DraftCode, domain.DbIntrospectionErr))
	}
	return c.JSON(http.StatusOK, SuccessResult())
}

//Migrate @Title Migrate
//@Description 迁移
//@Accept  json
//@Tags  prisma
//@Param id path integer true "数据源id"
//@Param schema path string true "schema内容"
//@Success 200  "迁移成功"
//@Failure 400  "迁移失败"
//@Router /api/v1/prisma/migrate/:id  [POST]
func (p *PrismaHandler) Migrate(c echo.Context) (err error) {
	req := struct {
		ID     int64  `json:"id"`
		Schema string `json:"schema"`
	}{}
	err = c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.PrismaCode, domain.ParamErr))
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id == 0 {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.DraftCode, domain.ParamErr))
	}
	//path := filepath.ToSlash(fmt.Sprintf("%s/%s/%s", utils.GetRootPath(), utils.GetDraftPrefixPath(), id))
	//ctx := c.Request().Context()
	// 创建草稿文件
	//content, err := p.DataSourceUseCase.GetPrismaSchema(ctx, uint(id))
	//if err != nil {
	//	return c.JSON(http.StatusBadRequest, GetResponseErr(domain.DraftCode, domain.DbSchemaErr))
	//}
	//moduleSchema, _ := ioutil.ReadFile(utils.GetModuleSchemaPath())
	//schemaContent := fmt.Sprintf("%s \n %s", content, req.Schema)
	path := fmt.Sprintf("schema%v", rand.Int())
	err = utils.WriteFile(path, req.Schema)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.DraftCode, domain.DbSchemaErr))
	}
	defer os.Remove(path)
	err = engine.Push(path)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.DraftCode, domain.NewBizErr("-1", err.Error())))
	}
	return c.JSON(http.StatusOK, SuccessResult())
}

//GetDMF @Title GetDMF
//@Description GetDMF
//@Accept  json
//@Tags  prisma
//@Param id path integer true "数据源id"
//@Success 200  "成功"
//@Failure 400  "失败"
//@Router /api/v1/prisma/dmf/:id  [GET]
func (p *PrismaHandler) GetDMF(c echo.Context) (err error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id == 0 {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.DraftCode, domain.ParamErr))
	}
	ctx := c.Request().Context()
	// 创建草稿文件
	content, err := p.DataSourceUseCase.GetPrismaSchema(ctx, uint(id))
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.DraftCode, domain.DbSchemaErr))
	}

	content, err = engine.Pull(content)
	if err != nil {
		if strings.Contains(err.Error(), "introspect error: The introspected database was empty") {
			return c.JSON(http.StatusOK, SuccessWriteResult(domain.DMF{
				Models: []domain.Models{},
				Enums:  []domain.Enums{},
			}))
		}
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.DraftCode, domain.DbIntrospectionErr))
	}

	queryEngine, err := engine.NewDMFQueryEngine(content)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.DraftCode, domain.NewBizErr("-1", err.Error())))
	}
	defer queryEngine.Disconnect()

	dmf, err := queryEngine.IntrospectDMMF(context.Background())
	if err != nil {
		log.Info("get dmmf fail err : ", err)
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.DraftCode, domain.DbSchemaErr))
	}
	// 先内省
	result := domain.GetDMF(*dmf)
	result.SchemaContent = content

	return c.JSON(http.StatusOK, SuccessWriteResult(result))
}

// Fetch @Title Fetch
// @Description 查询prisma信息
// @Accept  json
// @Tags  prisma
// @Success 200 {object} []domain.Prisma "查询成功"
// @Failure 400	"查询失败"
// @Router /api/v1/prisma  [GET]
//func (p *PrismaHandler) Fetch(c echo.Context) (err error) {
//	ctx := c.Request().Context()
//	prismaModels, err := p.PrismaUseCase.Fetch(ctx)
//	if err != nil {
//		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.PrismaCode, err))
//	}
//
//	resultData := make([]domain.Prisma, 0)
//	for _, row := range prismaModels {
//		path := fmt.Sprintf("%s/%s", row.File.Path, row.File.Name)
//		content, err := ioutil.ReadFile(path)
//		if err != nil {
//			return c.JSON(http.StatusBadRequest, GetResponseErr(domain.PrismaCode, domain.FileReadErr))
//		}
//		row.File.Content = content
//		resultData = append(resultData, row)
//	}
//
//	return c.JSON(http.StatusOK, SuccessWriteResult(resultData))
//}

// GetByID @Title GetByID
// @Description 根据id查询prisma信息
// @Accept  json
// @Tags  prisma
// @Param id formData integer true "id"
// @Success 200 {object} domain.Prisma "查询成功"
// @Failure 400	"查询失败"
// @Router /api/v1/prisma  [GET]
//func (p *PrismaHandler) GetByID(c echo.Context) (err error) {
//	id, err := strconv.Atoi(c.Param("id"))
//	if err != nil {
//		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.PrismaCode, domain.ParamErr))
//	}
//	if Empty(id) {
//		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.PrismaCode, domain.ParamIdEmptyErr))
//	}
//	ctx := c.Request().Context()
//	prismaModel, err := p.PrismaUseCase.GetByID(ctx, int64(id))
//	if err != nil {
//		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.PrismaCode, err))
//	}
//
//	path := fmt.Sprintf("%s/%s", prismaModel.File.Path, prismaModel.File.Name)
//	content, err := ioutil.ReadFile(path)
//	if err != nil {
//		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.PrismaCode, domain.FileReadErr))
//	}
//	prismaModel.File.Content = content
//	return c.JSON(http.StatusOK, SuccessWriteResult(prismaModel))
//
//}

// Create @Title Create
// @Description 创建prisma信息
// @Accept  json
// @Tags  prisma
// @Param id formData integer true "id"
// @Param name formData string true "名称"
// @Param file_id formData integer true "文件id"
// @Success 200 "创建成功"
// @Failure 400	"创建失败"
// @Router /api/v1/prisma  [POST]
func (p *PrismaHandler) Create(c echo.Context) (err error) {
	var prisma domain.Prisma
	err = c.Bind(&prisma)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.PrismaCode, domain.ParamErr))
	}
	if err = p.createCheckParams(prisma); err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.PrismaCode, err))
	}
	ctx := c.Request().Context()
	err = p.PrismaUseCase.Create(ctx, &prisma)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.PrismaCode, err))
	}
	return c.JSON(http.StatusOK, SuccessResult())
}

func (p *PrismaHandler) createCheckParams(prisma domain.Prisma) error {
	if ok := utils.Empty(prisma.Name); ok {
		return domain.ParamNameEmptyErr
	}
	if ok := utils.Empty(prisma.File.ID); ok {
		return domain.ParamIdEmptyErr
	}
	return nil
}

// Update @Title Update
// @Description 修改prisma信息
// @Accept  json
// @Tags  prisma
// @Param id formData integer true "id"
// @Param name formData string true "名称"
// @Param file_id formData integer true "文件id"
// @Success 200 "修改成功"
// @Failure 400	"修改失败"
// @Router /api/v1/prisma  [PUT]
func (p *PrismaHandler) Update(c echo.Context) (err error) {
	var prisma domain.Prisma
	err = c.Bind(&prisma)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.PrismaCode, domain.ParamErr))
	}
	if err = p.updateCheckParams(prisma); err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.PrismaCode, err))
	}

	ctx := c.Request().Context()
	err = p.PrismaUseCase.Update(ctx, &prisma)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.PrismaCode, err))
	}
	return c.JSON(http.StatusOK, SuccessResult())
}

// Delete @Title Delete
// @Description 删除prisma信息
// @Accept  json
// @Tags  prisma
// @Param id formData integer true "id"
// @Success 200 "删除成功"
// @Failure 400	"删除失败"
// @Router /api/v1/prisma/:id  [DELETE]
func (p *PrismaHandler) Delete(c echo.Context) (err error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.PrismaCode, domain.ParamErr))
	}

	if err = p.checkID(uint(id)); err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.PrismaCode, err))
	}

	ctx := c.Request().Context()
	err = p.PrismaUseCase.Delete(ctx, int64(id))
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.PrismaCode, err))
	}

	return c.JSON(http.StatusOK, SuccessResult())
}

func (p *PrismaHandler) createUpdateParams(prisma domain.Prisma) error {
	if ok := utils.Empty(prisma.ID); ok {
		return domain.ParamIdEmptyErr
	}
	if ok := utils.Empty(prisma.Name); ok {
		return domain.ParamNameEmptyErr
	}
	if ok := utils.Empty(prisma.File.ID); ok {
		return domain.ParamIdEmptyErr
	}
	return nil
}

func (p *PrismaHandler) updateCheckParams(prisma domain.Prisma) (result error) {
	if ok := utils.Empty(prisma.ID); ok {
		return domain.ParamIdEmptyErr
	}
	if ok := utils.Empty(strings.TrimSpace(prisma.Name)); ok {
		return domain.ParamNameEmptyErr
	}
	if ok := utils.Empty(prisma.File.ID); ok {
		return domain.ParamIdEmptyErr
	}
	return
}

func (p *PrismaHandler) checkID(id uint) (err error) {
	if ok := utils.Empty(id); ok {
		return domain.ParamIdEmptyErr
	}
	return
}
