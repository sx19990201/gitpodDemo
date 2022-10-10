package http

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fire_boom/pkg/repository/sqllite"
	"github.com/fire_boom/pkg/wdgfunc"
	"github.com/fire_boom/utils"
	"github.com/prisma/prisma-client-go/engine"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fire_boom/pkg/usecase"
	"github.com/tidwall/gjson"

	"github.com/fire_boom/domain"
	"github.com/labstack/echo/v4"
)

const (
	Success     = "success"
	SuccessCode = "0"

	dbTypePGSql   = "pgsql"
	dbTypeMysql   = "mysql"
	dbTypeMongodb = "mongodb"
	dbTypeSqlite  = "sqlite"

	Mysql       = "mysql"
	Planetscale = "planetscale"
	Postgres    = "postgres"
	Sqlite      = "sqlite"
	SqlServer   = "sqlServer"
	Mongodb     = "mongodb"
	OpenAPI     = "openApi"
	Graphql     = "graphql"
)

// ResponseError represent the response error struct
type ResponseError struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}

type DataSourceHandler struct {
	DataSourceUseCase domain.DataSourceUseCase
	FileUseCase       domain.FileUseCase
	//wdg               wundergraph.Wdg
}

//func InitDataSourceRouter(e *echo.Echo, dsr domain.DataSourceRepository, wdg wundergraph.Wdg) {
//	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second
//	au := usecase.NewDataSourceUseCase(dsr, timeoutContext)
//	NewDataSourceHandler(e, au, wdg)
//}

func InitDataSourceRouter(e *echo.Echo, db *sql.DB, dsr domain.DataSourceRepository) {
	timeoutContext := time.Duration(utils.GetTimeOut()) * time.Second
	fsr := sqllite.NewFileRepository(db)
	au := usecase.NewDataSourceUseCase(dsr, timeoutContext)
	fileUseCase := usecase.NewFileUseCase(fsr, timeoutContext)
	NewDataSourceHandler(e, au, fileUseCase)
}

func NewDataSourceHandler(e *echo.Echo, dUseCase domain.DataSourceUseCase, fUseCase domain.FileUseCase) {
	handler := &DataSourceHandler{
		DataSourceUseCase: dUseCase,
		FileUseCase:       fUseCase,
	}
	v1 := e.Group("/api/v1")
	{
		dataSource := v1.Group("/dataSource")
		{
			dataSource.GET("", handler.FindDataSources)
			dataSource.POST("", handler.Store)
			dataSource.PUT("", handler.Update)
			dataSource.PUT("/content/:id", handler.UpdateContent)
			dataSource.DELETE("/:id", handler.Delete)
			dataSource.POST("/import", handler.ImportOpenAPI)
			dataSource.POST("/removeFile", handler.RemoveFile)
			dataSource.POST("/CheckDBConn", handler.CheckDataSourceConn)
		}
	}
}

// CheckDataSourceConn @Title CheckDataSourceConn
// @Description 检查数据源链接
// @Accept  json
// @Tags  datasource
// @Param name formData string true "数据源名称"
// @Param sourceType formData integer true "数据源类型"
// @Param config formData string true "数据源配置"
// @Param switch formData integer true "数据源开关 0-开 1-关"
// @Success 200 "添加成功"
// @Failure 400	"添加失败"
// @Router /api/v1/CheckDBConn  [POST]
func (d *DataSourceHandler) CheckDataSourceConn(c echo.Context) (err error) {
	var dataSource domain.FbDataSourceResp
	err = c.Bind(&dataSource)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.DataSourceCode, domain.ParamErr))
	}

	ds := dataSource.TransformDataSource()
	err = CheckDataSourceConn(ds)
	if err != nil {
		return c.JSON(http.StatusOK, SuccessWriteResult(err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessWriteResult("连接成功"))
}

func CheckDataSourceConn(dataSource domain.FbDataSource) error {
	wdg := wdgfunc.GetWunderCtlClient()
	switch dataSource.SourceType {
	case domain.SourceTypeDB:
		config := domain.DbConfig{}
		err := json.Unmarshal([]byte(dataSource.Config), &config)
		if err != nil {
			return err
		}
		return CheckDBConn(config.DBType, config.DatabaseURL.Val)
	case domain.SourceTypeRest:
		config := domain.RestConfig{}
		err := json.Unmarshal([]byte(dataSource.Config), &config)
		if err != nil {
			return err
		}
		return utils.CheckOASFileContent(fmt.Sprintf("%s/%s", utils.GetOASFilePath(), config.OASFileID))
	case domain.SourceTypeGraphQL:
		config := domain.GraphqlConfig{}
		err := json.Unmarshal([]byte(dataSource.Config), &config)
		if err != nil {
			return err
		}
		return wdg.WunderCtlClient.CheckIntrospect(Graphql, config.URL)
	}
	return nil
}

func CheckDBConn(dataSourceType, url string) error {
	wdgClient := wdgfunc.GetWunderCtlClient()
	dbType := ""
	switch strings.ToLower(dataSourceType) {
	case dbTypePGSql:
		dbType = Postgres
	case dbTypeMysql:
		dbType = Mysql
	case dbTypeSqlite:
		dbType = Sqlite
	case dbTypeMongodb:
		dbType = Mongodb
	}
	err := wdgClient.WunderCtlClient.CheckIntrospect(dbType, url)
	if err != nil {
		return errors.New("连接失败")
	}
	return nil
}

// Store @Title Store
// @Description 存储数据源信息
// @Accept  json
// @Tags  datasource
// @Param name formData string true "数据源名称"
// @Param sourceType formData integer true "数据源类型"
// @Param config formData string true "数据源配置"
// @Param switch formData integer true "数据源开关 0-开 1-关"
// @Success 200 "添加成功"
// @Failure 400	"添加失败"
// @Router /api/v1/dataSource  [POST]
func (d *DataSourceHandler) Store(c echo.Context) (err error) {
	var dataSource domain.FbDataSourceResp
	err = c.Bind(&dataSource)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.DataSourceCode, domain.ParamErr))
	}

	if err = dataSourceStoreCheckParams(dataSource); err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.DataSourceCode, err))
	}

	ds := dataSource.TransformDataSource()
	// 校验数据源连接
	if err = CheckDataSourceConn(ds); err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.DataSourceCode, domain.DbConnErr))
	}

	ctx := c.Request().Context()
	// 默认开关为关闭
	ds.Switch = domain.SwitchOff
	// 如果是自定义的，需要创建hooks文件
	if ds.SourceType == domain.SourceTypeCustomize {
		config := domain.CustomizeConfig{}
		err = json.Unmarshal([]byte(ds.Config), &config)
		if err != nil {
			return c.JSON(http.StatusBadRequest, GetResponseErr(domain.DataSourceCode, domain.JsonUnMarshalErr))
		}
		// TODO 检查数据源

		path := fmt.Sprintf("%s/%s%s", utils.GetCustomizeHookPathPrefix(), config.ApiNamespace, utils.GetHooksSuffix())
		err = utils.WriteFile(path, config.Schema)
		if err != nil {
			return c.JSON(http.StatusBadRequest, GetResponseErr(domain.DataSourceCode, domain.FileWriteErr))
		}
	}
	lastID, err := d.DataSourceUseCase.Store(ctx, &ds)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.DataSourceCode, err))
	}

	return c.JSON(http.StatusOK, SuccessWriteResult(lastID))
}

func dataSourceStoreCheckParams(dataSource domain.FbDataSourceResp) (err error) {
	if ok := utils.Empty(dataSource.Name); ok {
		return domain.ParamNameEmptyErr
	}
	if ok := dataSourceCheckSourceType(dataSource.SourceType); ok {
		return domain.ParamSourceTypeEmptyErr
	}
	return
}

// UpdateContent @Title UpdateContent
// @Description 修改数据源内容
// @Accept  json
// @Tags  datasource
// @Param id formData integer true "数据源id"
// @Param content formData string true "数据源内容"
// @Success 200 "修改成功"
// @Failure 400	"修改失败"
// @Router /api/v1/dataSource/content/:id  [PUT]
func (d *DataSourceHandler) UpdateContent(c echo.Context) (err error) {
	content := struct {
		Content string `json:"content"`
	}{}
	err = c.Bind(&content)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.DataSourceCode, domain.ParamErr))
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id == 0 {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.DataSourceCode, domain.ParamErr))
	}
	ctx := c.Request().Context()
	ds, err := d.DataSourceUseCase.GetByID(ctx, uint(id))
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.DataSourceCode, err))
	}
	config := domain.CustomizeConfig{}
	err = json.Unmarshal([]byte(ds.Config), &config)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.DataSourceCode, domain.JsonUnMarshalErr))
	}
	path := fmt.Sprintf("%s/%s%s", utils.GetCustomizeHookPathPrefix(), config.ApiNamespace, utils.GetHooksSuffix())
	err = utils.WriteFile(path, content.Content)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.DataSourceCode, domain.FileWriteErr))
	}

	return c.JSON(http.StatusOK, SuccessResult())
}

// Update @Title Update
// @Description 修改数据源信息
// @Accept  json
// @Tags  datasource
// @Param name formData string true "数据源名称"
// @Param sourceType formData integer true "数据源类型"
// @Param config formData string true "数据源配置"
// @Success 200 "修改成功"
// @Failure 400	"修改失败"
// @Router /api/v1/dataSource  [PUT]
func (d *DataSourceHandler) Update(c echo.Context) (err error) {
	var dataSource domain.FbDataSourceResp
	err = c.Bind(&dataSource)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.DataSourceCode, domain.ParamErr))
	}
	if err = dataSourceUpdateCheckParams(dataSource); err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.DataSourceCode, err))
	}

	newDs := dataSource.TransformDataSource()
	// TODO 开关开启，则校验数据库连接
	//if newDs.Switch == domain.SwitchOn {
	//	if err = CheckDataSourceConn(newDs); err != nil {
	//		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.DataSourceCode, domain.DbConnErr))
	//	}
	//}

	ctx := c.Request().Context()
	oldDs, err := d.DataSourceUseCase.GetByID(ctx, dataSource.ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.DataSourceCode, err))
	}

	// 如果是自定义的，需要创建hooks文件
	if newDs.SourceType == domain.SourceTypeCustomize {
		config := domain.CustomizeConfig{}
		err = json.Unmarshal([]byte(newDs.Config), &config)
		if err != nil {
			return c.JSON(http.StatusBadRequest, GetResponseErr(domain.DataSourceCode, domain.JsonUnMarshalErr))
		}
		oldPath := fmt.Sprintf("%s/%s%s", utils.GetCustomizeHookPathPrefix(), oldDs.Name, utils.GetHooksSuffix())
		newPath := fmt.Sprintf("%s/%s%s", utils.GetCustomizeHookPathPrefix(), newDs.Name, utils.GetHooksSuffix())
		// 如果存在则重命名后写入
		if utils.FileExist(oldPath) {
			// 重命名
			err := os.Rename(oldPath, newPath)
			if err != nil {
				return c.JSON(http.StatusBadRequest, GetResponseErr(domain.DataSourceCode, domain.FileReNameErr))
			}
		}
	}
	affect, err := d.DataSourceUseCase.Update(ctx, &newDs)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.DataSourceCode, err))
	}

	// 如果修改的是当前的数据源ID，则需要重载schema查询引擎
	if int64(newDs.ID) == utils.GetIndexDBID() {
		content, err := d.DataSourceUseCase.GetPrismaSchema(ctx, uint(newDs.ID))
		if err != nil {
			return c.JSON(http.StatusBadRequest, GetResponseErr(domain.DbErrCode, domain.FileReadErr))
		}
		_, err = engine.ReloadQueryEngineOnce(content)
		if err != nil {
			return c.JSON(http.StatusBadRequest, GetResponseErr(domain.DbErrCode, domain.NewBizErr("-1", err.Error())))
		}
	}

	return c.JSON(http.StatusOK, SuccessWriteResult(fmt.Sprintf("Affected rows : %d", affect)))
}

func dataSourceUpdateCheckParams(dataSource domain.FbDataSourceResp) (err error) {
	if ok := utils.Empty(dataSource.ID); ok {
		return domain.ParamIdEmptyErr
	}
	if ok := utils.Empty(dataSource.Name); ok {
		return domain.ParamNameEmptyErr
	}
	if ok := dataSourceCheckSourceType(dataSource.SourceType); ok {
		return domain.ParamSourceTypeEmptyErr
	}
	return
}

// Delete @Title Delete
// @Description 删除数据源信息
// @Accept  json
// @Tags  datasource
// @Param id formData integer true "数据源id"
// @Success 200 "删除成功"
// @Failure 400	"删除失败"
// @Router /api/v1/dataSource/:id  [DELETE]
func (d *DataSourceHandler) Delete(c echo.Context) (err error) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.DataSourceCode, domain.ParamErr))
	}
	if err = dataSourceDeleteCheckParams(uint(id)); err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.DataSourceCode, err))
	}
	ctx := c.Request().Context()

	ds, err := d.DataSourceUseCase.GetByID(ctx, uint(id))
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.DataSourceCode, err))
	}
	// 如果是自定义的，需要删除hooks文件
	if ds.SourceType == domain.SourceTypeCustomize {
		config := domain.CustomizeConfig{}
		err = json.Unmarshal([]byte(ds.Config), &config)
		if err != nil {
			return c.JSON(http.StatusBadRequest, GetResponseErr(domain.DataSourceCode, domain.JsonUnMarshalErr))
		}
		path := fmt.Sprintf("%s/%s%s", utils.GetCustomizeHookPathPrefix(), config.ApiNamespace, utils.GetHooksSuffix())

		err = os.Remove(path)
		if err != nil {
			return c.JSON(http.StatusBadRequest, GetResponseErr(domain.DataSourceCode, domain.FileWriteErr))
		}
	}
	affect, err := d.DataSourceUseCase.Delete(ctx, uint(id))
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.DataSourceCode, err))
	}

	//err = d.wdg.ReloadConfig(ctx)
	//if err != nil {
	//	log.Error("Store reload config meet err=", err)
	//}

	return c.JSON(http.StatusOK, SuccessWriteResult(fmt.Sprintf("Affected rows : %d", affect)))
}

// FindDataSources @Title FindDataSources
// @Description 查询数据源信息
// @Accept  json
// @Tags  datasource
// @Param datasourceType formData integer true "数据源类型"
// @Success 200 {object} []domain.FbDataSource "查询成功"
// @Failure 400	"查询失败"
// @Router /api/v1/dataSource  [GET]
func (d *DataSourceHandler) FindDataSources(c echo.Context) (err error) {
	ctx := c.Request().Context()
	dataSources, err := d.DataSourceUseCase.FindDataSources(ctx)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.DataSourceCode, err))
	}
	result := make([]domain.FbDataSourceResp, 0)
	for _, row := range dataSources {
		ds := row.GetFbDataSourceResp()
		// 自定义类型需要读取hooks
		if ds.SourceType == domain.SourceTypeCustomize {
			config := domain.CustomizeConfig{}
			err = json.Unmarshal([]byte(row.Config), &config)
			if err != nil {
				return c.JSON(http.StatusBadRequest, GetResponseErr(domain.DataSourceCode, domain.JsonUnMarshalErr))
			}
			path := fmt.Sprintf("%s/%s%s", utils.GetCustomizeHookPathPrefix(), ds.Name, utils.GetHooksSuffix())
			content, err := ioutil.ReadFile(path)
			if err != nil {
				return c.JSON(http.StatusBadRequest, GetResponseErr(domain.DataSourceCode, domain.FileReadErr))
			}
			config.Schema = string(content)
			//configResp, err := json.Marshal(config)
			//if err != nil {
			//	return c.JSON(http.StatusBadRequest, GetResponseErr(domain.DataSourceCode, domain.JsonMarshalErr))
			//}
			ds.Config = config
		}
		result = append(result, ds)
	}

	//
	param := struct {
		DatasourceType int64 `json:"datasourceType"`
	}{}
	c.Bind(&param)
	if param.DatasourceType != 0 {
		tempArr := make([]domain.FbDataSourceResp, 0)
		for _, row := range result {
			if row.SourceType != int(param.DatasourceType) {
				continue
			}
			tempArr = append(tempArr, row)
		}
		result = tempArr
	}

	return c.JSON(http.StatusOK, SuccessWriteResult(result))
}
func dataSourceDeleteCheckParams(id uint) (err error) {
	if ok := utils.Empty(id); ok {
		return domain.ParamIdEmptyErr
	}
	return
}
func dataSourceCheckSourceType(sourceType int) (ok bool) {
	switch sourceType {
	case domain.SourceTypeDB:
	case domain.SourceTypeRest:
	case domain.SourceTypeGraphQL:
	case domain.SourceTypeCustomize:
	default:
		return true
	}
	return false
}

// RemoveFile @Title RemoveFile
// @Description 删除文件
// @Accept  json
// @Tags  file
// @Param id formData string true "id"
// @Success 200 "删除成功"
// @Failure 400	"删除失败"
// @Router /api/v1/file  [Delete]
func (d *DataSourceHandler) RemoveFile(c echo.Context) error {
	param := struct {
		UUID string `json:"id"`
	}{}
	err := c.Bind(&param)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.FileCode, domain.ParamErr))
	}

	filePath := fmt.Sprintf("%s/%s", utils.GetOASFilePath(), param.UUID)
	err = os.Remove(filePath)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.FileCode, domain.FileDeleteErr))
	}
	ctx := c.Request().Context()
	err = d.FileUseCase.Delete(ctx, param.UUID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.FileCode, domain.FileDeleteErr))
	}

	return c.JSON(http.StatusOK, SuccessResult())
}

// ImportOpenAPI @Title ImportOpenAPI
// @Description 导入外部数据源信息
// @Accept  json
// @Tags  datasource
// @Param file formData string true "数据源文件"
// @Success 200 "导入成功"
// @Failure 400	"导入失败"
// @Router /api/v1/dataSource/import  [POST]
func (d *DataSourceHandler) ImportOpenAPI(c echo.Context) (err error) {
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.DataSourceCode, domain.ParamErr))
	}
	fileUUID := utils.GenerateUUID()
	filePath := fmt.Sprintf("%s/%s", utils.GetOASFilePath(), fileUUID)
	err = utils.UploadWriteFile(file, filePath)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.DataSourceCode, domain.FileWriteErr))
	}
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.DataSourceCode, domain.FileReadErr))
	}
	// 校验文件名json 或者
	if !strings.Contains(file.Filename, utils.GetJsonSuffix()) && !strings.Contains(file.Filename, utils.GetYamlSuffix()) {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.DataSourceCode, domain.FileTypeErr))
	}

	if !utils.CheckJsonStr(content) && !utils.CheckYamlStr(content) {
		os.Remove(filePath)
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.DataSourceCode, domain.FileTypeErr))
	}

	ctx := c.Request().Context()
	_, err = d.FileUseCase.Store(ctx, &domain.File{
		ID:   fileUUID,
		Name: file.Filename,
		Path: filePath,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.FileCode, err))
	}
	result := fmt.Sprintf("%s", fileUUID)

	return c.JSON(http.StatusOK, SuccessWriteResult(result))
}

func getOpenApiJsonOperationIds(openAPIJson string) []string {
	// 获得paths
	pathsJson := gjson.Get(openAPIJson, "paths")
	operationIds := make([]string, 0)
	pathsJson.ForEach(func(uriName, row gjson.Result) bool {
		row.ForEach(func(uriMethod, value gjson.Result) bool {
			operationId := gjson.Get(value.Raw, "operationId").Str
			if utils.Empty(operationId) {
				return true
			}
			operationIds = append(operationIds, operationId)
			return true
		})
		return true
	})
	return operationIds
}

func GetResponseErr(code string, err error) (result ResponseError) {
	respErr := domain.NewBizErr("", "")
	if errors.As(err, &respErr) {
		result.Code = fmt.Sprintf("%s%s", code, respErr.Code())
		result.Message = respErr.Error()
	}
	return
}

func SuccessWriteResult(data interface{}) ResponseError {
	return ResponseError{
		Code:    domain.OK.Code(),
		Message: domain.OK.Error(),
		Result:  data,
	}
}

func SuccessResult() ResponseError {
	return ResponseError{
		Code:    domain.OK.Code(),
		Message: domain.OK.Error(),
	}
}
