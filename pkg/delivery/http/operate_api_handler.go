package http

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/fire_boom/domain"
	"github.com/fire_boom/pkg/garphql/validator"
	"github.com/fire_boom/pkg/repository/sqllite"
	"github.com/fire_boom/pkg/usecase"
	"github.com/fire_boom/utils"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

//var operateAPIPathPrefix = viper.GetString("operateAPIPath.prefix")
//var operateAPIPathSuffix = viper.GetString("operateAPIPath.suffix")
//var operateAPIPathSwitchDefaultState = viper.GetString("operateAPIPath.switchDefaultState")

//1-preResolve 2-postResolve 3-customeResolve 4-mutatingPostResolve 5-mutatingPreResolve
const (
	operationSettingFlat       = 1 // operation设置
	operationGlobalSettingFlat = 2 // operation全局设置
)

// dirTreeNode 文件目录树形结构节点
type dirTreeNode struct {
	domain.FbOperations
	Child []dirTreeNode `json:"children"`
}

type dirTreeNodeResp struct {
	domain.FbOperationsResult
	Child []dirTreeNodeResp `json:"children"`
}

type trees []dirTreeNodeResp

func (s trees) Len() int { return len(s) }
func (s trees) Less(i, j int) bool {
	sorti := 1
	if s[i].Legal == false {
		sorti = 2
	}
	sortj := 1
	if s[j].Legal == false {
		sortj = 2
	}
	return sorti > sortj
}
func (s trees) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

type OperateAPIHandler struct {
	hooksUC      domain.HooksUseCase
	operationsUC domain.OperationsUseCase
}

func InitOperateAPIHandler(e *echo.Echo, db *sql.DB) {
	timeoutContext := time.Duration(utils.GetTimeOut()) * time.Second
	hookUC := usecase.NewHooksUseCase(timeoutContext)
	or := sqllite.NewOperationsRepository(db)
	operationUC := usecase.NewOperationsUseCase(or, timeoutContext)
	NewAPIHandler(e, hookUC, operationUC)
}

func NewAPIHandler(e *echo.Echo, hookUseCase domain.HooksUseCase, operationsUC domain.OperationsUseCase) {
	handler := &OperateAPIHandler{
		hooksUC:      hookUseCase,
		operationsUC: operationsUC,
	}
	v1 := e.Group("/api/v1")
	{
		operateApi := v1.Group("/operateApi")
		{
			operateApi.GET("", handler.FindOperateAPIs)                     // 列表
			operateApi.POST("", handler.CreateFile)                         // 创建文件
			operateApi.PUT("/:id", handler.Update)                          // 更新
			operateApi.PUT("/content/:id", handler.UpdateContentById)       // 更新
			operateApi.GET("/:id", handler.GetDetailByID)                   // 详情
			operateApi.DELETE("/:id", handler.Remove)                       // 删除
			operateApi.POST("/dir", handler.CreateDir)                      // 创建文件夹
			operateApi.PUT("/dir", handler.ReNameDir)                       // 更新文件夹
			operateApi.DELETE("/dir", handler.RemoveDir)                    // 删除文件夹
			operateApi.GET("/hooks/:id", handler.GetHooksByName)            // 获取operation钩子
			operateApi.PUT("/hooks/:id", handler.UpdateHooksByID)           // 修改operation钩子
			operateApi.GET("/hooks", handler.GetGlobalHooks)                // 全局钩子
			operateApi.PUT("/hooks", handler.UpdateGlobalHooks)             // 全局钩子
			operateApi.GET("/getGenerateSchema", handler.GetGenerateSchema) // 获取wdg生成的schema
			operateApi.GET("/setting/:id", handler.GetOperationSetting)     // 获取operation设置
			operateApi.PUT("/setting/:id", handler.UpdateOperationSetting)  // 修改operation设置
			operateApi.GET("/setting", handler.GetGlobalSetting)            // 获取全局operation设置
			operateApi.PUT("/setting", handler.UpdateGlobalSetting)         // 修改全局operation设置

			operateApi.GET("/mock/:id", handler.GetMock)
			operateApi.PUT("/mock/:id", handler.UpdateMock)
			operateApi.GET("/sdk", handler.GetSDK)
			operateApi.GET("/json", handler.GetJson)
		}
	}
}

// GetJson @Title GetJson
// @Description 获取json
// @Accept  json
// @Tags  operateApi
// @Param id formData integer true "id"
// @Param mockSwitch formData string true "mock开关"
// @Param content formData string true "内容"
// @Success 200 "获取成功"
// @Failure 400	"获取失败"
// @Router /api/v1/operateApi/json  [GET]
func (a *OperateAPIHandler) GetJson(c echo.Context) (err error) {
	jsonFile, err := ioutil.ReadFile("wundergraph/.wundergraph/generated/wundergraph.postman.json")
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, domain.FileReadErr))
	}
	return c.Stream(http.StatusOK, echo.MIMEApplicationJSON, bytes.NewReader(jsonFile))
}

// GetSDK @Title GetSDK
// @Description 获取sdk
// @Accept  json
// @Tags  operateApi
// @Param id formData integer true "id"
// @Param mockSwitch formData string true "mock开关"
// @Param content formData string true "内容"
// @Success 200 "获取成功"
// @Failure 400	"获取失败"
// @Router /api/v1/operateApi/sdk  [GET]
func (a *OperateAPIHandler) GetSDK(c echo.Context) (err error) {
	var zipPath = utils.GetSdkSrcPath()
	// 目标文件，压缩后的文件
	var dst = fmt.Sprintf("%s/sdk.zip", utils.GetSdkDstPath())
	if err := utils.ZipFiles(dst, utils.GetDicAllChildFilePath(zipPath)); err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, domain.ZipErr))
	}
	zipFile, err := ioutil.ReadFile(dst)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, domain.FileReadErr))
	}
	return c.Blob(http.StatusOK, "application/zip", zipFile)
}

// UpdateMock @Title UpdateMock
// @Description 修改mock
// @Accept  json
// @Tags  operateApi
// @Param id formData integer true "id"
// @Param mockSwitch formData string true "mock开关"
// @Param content formData string true "内容"
// @Success 200 "获取成功"
// @Failure 400	"获取失败"
// @Router /api/v1/operateApi/mock/:id  [PUT]
func (a *OperateAPIHandler) UpdateMock(c echo.Context) (err error) {
	ctx := c.Request().Context()

	param := struct {
		MockSwitch bool   `json:"mockSwitch"`
		Content    string `json:"content"`
	}{}
	err = c.Bind(&param)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, domain.ParamErr))
	}

	id := c.Param("id")
	operationID, _ := strconv.Atoi(id)
	if operationID == 0 {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, domain.ParamErr))
	}
	operation, err := a.operationsUC.GetByID(ctx, int64(operationID))
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, err))
	}

	// 判断该mock是否存在
	offMockName := fmt.Sprintf("%s%s%s%s", utils.GetMockPath(), operation.Path, utils.GetHooksSuffix(), utils.GetSwitchState(utils.OFF))
	onMockName := fmt.Sprintf("%s%s%s%s", utils.GetMockPath(), operation.Path, utils.GetHooksSuffix(), utils.GetSwitchState(utils.ON))
	// 获取所有mock
	filePaths := utils.GetDicAllChildFilePath(utils.GetMockPath())
	path := ""
	for _, mockPath := range filePaths {
		mockPath = filepath.ToSlash(mockPath)
		offMockName = filepath.ToSlash(offMockName)
		onMockName = filepath.ToSlash(onMockName)
		// 不存在则直接返回空
		if !(mockPath == offMockName || mockPath == onMockName) {
			continue
		}
		// 存在则读取mock
		path = mockPath
	}
	// 如果不存在，则拼接路径
	if path == "" {
		path = fmt.Sprintf("%s%s%s%s", utils.GetMockPath(), operation.Path, utils.GetHooksSuffix(), utils.GetSwitchState(param.MockSwitch))
		err = utils.WriteFile(path, param.Content)
		if err != nil {
			return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, domain.FileWriteErr))
		}
		return c.JSON(http.StatusOK, SuccessResult())
	}
	// 如果存在,则先重名，再修改
	oldPath := path
	newPath := fmt.Sprintf("%s%s%s%s", utils.GetMockPath(), operation.Path, utils.GetHooksSuffix(), utils.GetSwitchState(param.MockSwitch))
	// 先重命名
	err = os.Rename(oldPath, newPath)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, domain.FileReNameErr))
	}
	err = utils.WriteFile(newPath, param.Content)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, domain.FileWriteErr))
	}
	return c.JSON(http.StatusOK, SuccessResult())
}

// GetMock @Title GetMock
// @Description 获取mock
// @Accept  json
// @Tags  operateApi
// @Param id formData string true "id"
// @Success 200 {object} domain.Mock "获取成功"
// @Failure 400	"获取失败"
// @Router /api/v1/operateApi/mock/:id  [GET]
func (a *OperateAPIHandler) GetMock(c echo.Context) (err error) {
	ctx := c.Request().Context()
	id := c.Param("id")
	operationID, _ := strconv.Atoi(id)
	if operationID == 0 {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, domain.ParamErr))
	}
	operation, err := a.operationsUC.GetByID(ctx, int64(operationID))
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, err))
	}

	var result domain.Mock
	// 开关默认为关
	result.MockSwitch = utils.OFF

	// 判断该mock是否存在
	offMockName := filepath.ToSlash(fmt.Sprintf("%s%s%s%s", utils.GetMockPath(), operation.Path, utils.GetHooksSuffix(), utils.GetSwitchState(utils.OFF)))
	onMockName := filepath.ToSlash(fmt.Sprintf("%s%s%s%s", utils.GetMockPath(), operation.Path, utils.GetHooksSuffix(), utils.GetSwitchState(utils.ON)))
	// 获取所有mock
	filePaths := utils.GetDicAllChildFilePath(utils.GetMockPath())
	path := ""
	for _, mockPath := range filePaths {
		// 不存在则直接返回空
		if !(mockPath == offMockName || mockPath == onMockName) {
			continue
		}
		// 存在则读取mock
		path = mockPath
	}
	// 如果不存在，则直接返回
	if path == "" {
		return c.JSON(http.StatusOK, SuccessWriteResult(result))
	}

	// 拿到开关,如果包含.off
	if !strings.Contains(path, utils.GetSwitchState(utils.OFF)) {
		result.MockSwitch = utils.ON
	}

	settingByte, err := ioutil.ReadFile(path)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, domain.FileReadErr))
	}
	result.Content = string(settingByte)
	return c.JSON(http.StatusOK, SuccessWriteResult(result))
}

// UpdateGlobalSetting @Title UpdateGlobalSetting
// @Description 修改operation设置
// @Accept  json
// @Tags  operateApi
// @Param id formData integer true "id"
// @Param authenticationRequired          formData boolean true "需要授权"
// @Param cachingEnable                   formData boolean true "开启缓存"
// @Param cachingMaxAge                   formData integer true "最大时长"
// @Param cachingStaleWhileRevalidate     formData integer true "重校验时长"
// @Param liveQueryEnable                 formData boolean true "开启实时"
// @Param liveQueryPollingIntervalSeconds formData integer true "轮询间隔"
// @Success 200 "修改成功"
// @Failure 400	"修改失败"
// @Router /api/v1/operateApi/setting/:id  [PUT]
func (a *OperateAPIHandler) UpdateGlobalSetting(c echo.Context) (err error) {
	operationSetting := domain.OperationSetting{}
	err = c.Bind(&operationSetting)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, domain.ParamErr))
	}
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, err))
	}
	settingByte, err := json.Marshal(operationSetting)
	content := string(settingByte)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, domain.JsonMarshalErr))
	}
	path := fmt.Sprintf("%s", utils.GetOperationGlobalSettingPath())

	err = utils.WriteFile(path, content)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, domain.FileWriteErr))
	}

	return c.JSON(http.StatusOK, SuccessResult())
}

// UpdateOperationSetting @Title UpdateOperationSetting
// @Description 修改operation设置
// @Accept  json
// @Tags  operateApi
// @Param settingType formData integer true "设置类型 1operation单个设置 2operation全局设置"
// @Param id formData integer true "id"
// @Param authenticationRequired          formData boolean true "需要授权"
// @Param cachingEnable                   formData boolean true "开启缓存"
// @Param cachingMaxAge                   formData integer true "最大时长"
// @Param cachingStaleWhileRevalidate     formData integer true "重校验时长"
// @Param liveQueryEnable                 formData boolean true "开启实时"
// @Param liveQueryPollingIntervalSeconds formData integer true "轮询间隔"
// @Success 200 "修改成功"
// @Failure 400	"修改失败"
// @Router /api/v1/operateApi/setting/:id  [PUT]
func (a *OperateAPIHandler) UpdateOperationSetting(c echo.Context) (err error) {
	ctx := c.Request().Context()
	operationSetting := domain.OperationSetting{}
	err = c.Bind(&operationSetting)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, domain.ParamErr))
	}
	id := c.Param("id")
	operationID, _ := strconv.Atoi(id)
	if operationID == 0 {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, domain.ParamIdEmptyErr))
	}
	operation, err := a.operationsUC.GetByID(ctx, int64(operationID))
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, err))
	}
	settingByte, err := json.Marshal(operationSetting)
	content := string(settingByte)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, domain.JsonMarshalErr))
	}

	path := fmt.Sprintf("%s%s", utils.GetOperationSettingPath(), operation.Path)
	err = utils.WriteFile(path, content)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, domain.FileWriteErr))
	}

	return c.JSON(http.StatusOK, SuccessResult())
}

// GetOperationSetting @Title GetOperationSetting
// @Description 获取operation设置
// @Accept  json
// @Tags  operateApi
// @Param settingType formData integer true "设置类型 1 operation单个设置 2 operation全局设置"
// @Param id formData integer true "id"
// @Success 200 {object} domain.OperationSetting "获取成功"
// @Failure 400	"获取失败"
// @Router /api/v1/operateApi/setting/:id  [GET]
func (a *OperateAPIHandler) GetOperationSetting(c echo.Context) (err error) {
	ctx := c.Request().Context()
	var result domain.OperationSetting

	id := c.Param("id")
	operationID, _ := strconv.Atoi(id)
	if operationID == 0 {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, domain.ParamErr))
	}
	operation, err := a.operationsUC.GetByID(ctx, int64(operationID))
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, err))
	}

	path := fmt.Sprintf("%s%s", utils.GetOperationSettingPath(), operation.Path)
	// 判断path是否存在
	if !utils.FileExist(path) {
		// 如果不存在则创建
		content, _ := json.Marshal(result)
		err := utils.WriteFile(path, string(content))
		if err != nil {
			return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, domain.FileWriteErr))
		}
		return c.JSON(http.StatusOK, SuccessWriteResult(result))
	}

	// 读取内容
	settingByte, err := ioutil.ReadFile(path)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, domain.FileReadErr))
	}
	if len(settingByte) == 0 {
		return c.JSON(http.StatusOK, SuccessWriteResult(result))
	}
	err = json.Unmarshal(settingByte, &result)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, domain.JsonUnMarshalErr))
	}

	return c.JSON(http.StatusOK, SuccessWriteResult(result))
}

// GetGlobalSetting @Title GetGlobalSetting
// @Description 获取operation设置
// @Accept  json
// @Tags  operateApi
// @Success 200 {object} domain.OperationSetting "获取成功"
// @Failure 400	"获取失败"
// @Router /api/v1/operateApi/setting  [GET]
func (a *OperateAPIHandler) GetGlobalSetting(c echo.Context) (err error) {
	var result domain.OperationSetting

	// 读取内容
	settingByte, err := ioutil.ReadFile(utils.GetOperationGlobalSettingPath())
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, domain.FileReadErr))
	}
	if len(settingByte) == 0 {
		return c.JSON(http.StatusOK, SuccessWriteResult(result))
	}
	err = json.Unmarshal(settingByte, &result)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, domain.JsonUnMarshalErr))
	}

	return c.JSON(http.StatusOK, SuccessWriteResult(result))
}

// UpdateGlobalHooks @Title UpdateGlobalHooks
// @Description 修改全局hook
// @Accept  json
// @Tags  operateApi
// @Param HookName formData string true "钩子名称"
// @Param content formData string true "内容"
// @Param hookSwitch formData bool true "开关"
// @Param hooksType formData string true "onRequest , onResponse"
// @Success 200 "修改成功"
// @Failure 400	"修改失败"
// @Router /api/v1/operateApi/hooks  [PUT]
func (a *OperateAPIHandler) UpdateGlobalHooks(c echo.Context) (err error) {
	var hooks domain.Hooks
	err = c.Bind(&hooks)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, domain.ParamErr))
	}
	err = a.hooksUC.UpdateHooksByPath(utils.GetGlobalHookPathPrefix(), hooks)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, domain.FileReadErr))
	}
	return c.JSON(http.StatusOK, SuccessWriteResult(hooks.Content))
}

// GetGlobalHooks @Title GetGlobalHooks
// @Description 获取全局hooks
// @Accept  json
// @Tags  operateApi
// @Success 200 {object} []domain.Hooks "获取成功"
// @Failure 400	"获取失败"
// @Router /api/v1/operateApi/hooks  [GET]
func (a *OperateAPIHandler) GetGlobalHooks(c echo.Context) (err error) {
	result, err := a.hooksUC.FindHooksByPath(utils.GetGlobalHookPathPrefix(), "", domain.GlobalHook)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.ParamErr))
	}
	return c.JSON(http.StatusOK, SuccessWriteResult(result))
}

// UpdateHooksByID @Title UpdateHooksByID
// @Description 修改hooks
// @Accept  json
// @Tags  operateApi
// @Param HookName formData string true "钩子名称"
// @Param content formData string true "内容"
// @Param hookSwitch formData bool true "开关"
// @Param hooksType formData string true "preResolve, postResolve, customeResolve, mutatingPostResolve"
// @Success 200 "修改成功"
// @Failure 400	"修改失败"
// @Router /api/v1/operateApi/hooks/:id  [PUT]
func (a *OperateAPIHandler) UpdateHooksByID(c echo.Context) (err error) {

	ctx := c.Request().Context()
	id := c.Param("id")
	operationID, _ := strconv.Atoi(id)
	operation, err := a.operationsUC.GetByID(ctx, int64(operationID))
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, err))
	}

	var hooks domain.Hooks
	err = c.Bind(&hooks)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.ParamErr))
	}
	hooks.FileName = operation.Path
	err = a.hooksUC.UpdateHooksByPath(utils.GetOperationsHookPathPrefix(), hooks)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, domain.FileReadErr))
	}
	return c.JSON(http.StatusOK, SuccessWriteResult(hooks.Content))
}

// GetHooksByName @Title GetHooksByName
// @Description 根据operations名称获取hooks
// @Accept  json
// @Tags  operateApi
// @Param id formData integer true "id"
// @Param hooksType formData integer true "hooks类型1-preResolve 2-postResolve 3-customeResolve 4-mutatingPostResolve"
// @Success 200 {object} []domain.Hooks "获取成功"
// @Failure 400	"获取失败"
// @Router /api/v1/operateApi/hooks/:id  [GET]
func (a *OperateAPIHandler) GetHooksByName(c echo.Context) (err error) {
	ctx := c.Request().Context()
	id := c.Param("id")
	operationID, _ := strconv.Atoi(id)
	operation, err := a.operationsUC.GetByID(ctx, int64(operationID))
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, err))
	}
	// 去除后缀  参数的路径是这样的：user/xxxx/xxxx.graphql 或者 user/xxxx/xxxx.graphql.off 去除后缀只保留user/xxxx/xxxx
	result, err := a.hooksUC.FindHooksByPath(utils.GetOperationsHookPathPrefix(), operation.Path, domain.OperationHook)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, domain.ParamErr))
	}
	return c.JSON(http.StatusOK, SuccessWriteResult(result))
}

// FindOperateAPIs @Title FindOperateAPIs
// @Description 查询api列表
// @Accept  json
// @Tags  operateApi
// @Success 200 {object} dirTreeNodeResp "查询成功"
// @Failure 400	"查询失败"
// @Router /api/v1/operateApi  [GET]
func (a *OperateAPIHandler) FindOperateAPIs(c echo.Context) (err error) {
	ctx := c.Request().Context()
	operations, err := a.operationsUC.FindOperations(ctx)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, err))
	}
	operationMap := make(map[string]domain.FbOperationsResult, 0)
	for _, row := range operations {
		operationMap[row.Path] = row.TransformToResult()
	}
	tree, err := getDirTreeResp(utils.GetApiPathPrefix(), operationMap)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, domain.FileReadErr))
	}
	result := make([]dirTreeNodeResp, 0)
	dirArr := make([]dirTreeNodeResp, 0)
	fileArr := make([]dirTreeNodeResp, 0)
	for _, row := range tree.Child {
		if row.IsDir {
			dirArr = append(dirArr, row)
			continue
		}
		fileArr = append(fileArr, row)
	}
	result = append(result, dirArr...)
	var sortFileArr trees
	sortFileArr = fileArr
	sort.Sort(sortFileArr)
	result = append(result, sortFileArr...)

	return c.JSON(http.StatusOK, SuccessWriteResult(result))
}

// GetGenerateSchema @Title GetGenerateSchema
// @Description 获取生成的schema
// @Accept  json
// @Tags  operateApi
// @Success 200 "查询成功"
// @Failure 400	"查询失败"
// @Router /api/v1/operateApi/getGenerateSchema  [GET]
func (a *OperateAPIHandler) GetGenerateSchema(c echo.Context) (err error) {
	tree, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", utils.GetWdgGeneratedPath(), domain.GenerateSchemaGraphql))
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, domain.FileReadErr))
	}
	return c.JSON(http.StatusOK, SuccessWriteResult(string(tree)))
}

// GetDetailByID @Title GetDetailByID
// @Description 查询api文件内容
// @Accept  json
// @Tags  operateApi
// @Param id formData integer true "id"
// @Success 200 {object} domain.FbOperations "查询成功"
// @Failure 400	"查询失败"
// @Router /api/v1/operateApi/:id  [GET]
func (a *OperateAPIHandler) GetDetailByID(c echo.Context) (err error) {
	ctx := c.Request().Context()
	id := c.Param("id")
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, domain.ParamErr))
	}

	operID, _ := strconv.Atoi(id)
	operation, err := a.operationsUC.GetByID(ctx, int64(operID))
	if err != nil || operation.ID == 0 {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, domain.FileDeleteErr))
	}

	readPath := fmt.Sprintf("%s%s%s%s", utils.GetApiPathPrefix(), operation.Path, utils.GetApiPathSuffix(), utils.GetSwitchState(operation.Enable == 0))
	content, err := ioutil.ReadFile(readPath)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, domain.FileReadErr))
	}

	result := operation.TransformToResult()

	result.Content = string(content)
	return c.JSON(http.StatusOK, SuccessWriteResult(result))
}

// Remove @Title Remove
// @Description 删除文件
// @Accept  json
// @Tags  operateApi
// @Param id formData integer true "id"
// @Success 200  "删除成功"
// @Failure 400	"删除失败"
// @Router /api/v1/operateApi/:id  [DELETE]
func (a *OperateAPIHandler) Remove(c echo.Context) (err error) {
	ctx := c.Request().Context()
	id := c.Param("id")

	operID, _ := strconv.Atoi(id)

	operation, err := a.operationsUC.GetByID(ctx, int64(operID))
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, domain.FileDeleteErr))
	}
	readPath := fmt.Sprintf("%s%s%s%s", utils.GetApiPathPrefix(), operation.Path, utils.GetApiPathSuffix(), utils.GetSwitchState(operation.Enable == 0))

	err = os.Remove(readPath)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, domain.FileDeleteErr))
	}

	_, err = a.operationsUC.Delete(ctx, operation.ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, domain.FileDeleteErr))
	}
	// 将hooks删掉
	a.hooksUC.RemoveOperationHooks(operation.Path)
	// 将setting删掉
	removeOperationSetting(operation.Path)
	// 将mock删掉
	removeMockSetting(operation.Path)
	return c.JSON(http.StatusOK, SuccessResult())
}

// RemoveDir @Title RemoveDir
// @Description 删除文件夹
// @Accept  json
// @Tags  operateApi
// @Param id formData integer true "id"
// @Success 200  "删除成功"
// @Failure 400	"删除失败"
// @Router /api/v1/operateApi/dir  [DELETE]hooksPath
func (a *OperateAPIHandler) RemoveDir(c echo.Context) (err error) {
	ctx := c.Request().Context()
	path := struct {
		Path string `json:"path"`
	}{}
	err = c.Bind(&path)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, domain.ParamErr))
	}
	if utils.Empty(path.Path) {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, domain.ParamNameEmptyErr))
	}
	// 查询该文件夹下所有文件
	operations, err := a.operationsUC.FindByPath(ctx, path.Path)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, domain.DbFindErr))
	}
	// 删除数据集库内文件的信息
	err = a.operationsUC.BatchDelete(ctx, operations)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, domain.DbDeleteErr))
	}
	operationsPath := fmt.Sprintf("%s%s", utils.GetApiPathPrefix(), path.Path)
	// 删除文件夹以及该文件夹下子文件
	err = os.Remove(operationsPath)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, domain.FileDeleteErr))
	}
	hookPath := fmt.Sprintf("%s%s", utils.GetOperationsHookPathPrefix(), path.Path)
	// 删除hooks文件夹
	err = os.Remove(hookPath)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, domain.FileDeleteErr))
	}

	// 删除设置文件夹
	settingPath := fmt.Sprintf("%s%s", utils.GetOperationSettingPath(), path.Path)
	err = os.Remove(settingPath)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, domain.FileDeleteErr))
	}

	// mock设置文件夹
	mockPath := fmt.Sprintf("%s%s", utils.GetMockPath(), path.Path)
	err = os.Remove(mockPath)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, domain.FileDeleteErr))
	}
	return c.JSON(http.StatusOK, SuccessResult())
}

func removeMockSetting(operationPath string) {
	for _, path := range utils.GetDicAllChildFilePath(utils.GetMockPath()) {
		path = filepath.ToSlash(path)
		operationPathOn := filepath.ToSlash(fmt.Sprintf("%s%s", utils.GetMockPath(), operationPath))
		operationPathOff := filepath.ToSlash(fmt.Sprintf("%s%s%s", utils.GetMockPath(), operationPath, utils.GetSwitchState(utils.OFF)))
		if path == operationPathOn || path == operationPathOff {
			os.Remove(path)
		}
	}
}

func removeOperationSetting(operationPath string) {
	for _, path := range utils.GetDicAllChildFilePath(utils.GetOperationSettingPath()) {
		path = filepath.ToSlash(path)
		operationPathOn := filepath.ToSlash(fmt.Sprintf("%s%s", utils.GetOperationSettingPath(), operationPath))
		operationPathOff := filepath.ToSlash(fmt.Sprintf("%s%s%s", utils.GetOperationSettingPath(), operationPath, utils.GetSwitchState(utils.OFF)))
		if path == operationPathOn || path == operationPathOff {
			os.Remove(path)
		}
	}
}

// getDirTreeResp 递归遍历文件目录
func getDirTreeResp(pathName string, operationsMap map[string]domain.FbOperationsResult) (dirTreeNodeResp, error) {
	rd, err := ioutil.ReadDir(pathName)
	if err != nil {
		return dirTreeNodeResp{}, err
	}
	var tree, childNode dirTreeNodeResp
	childNode.IsDir = true
	tree.Path = strings.ReplaceAll(pathName, utils.GetApiPathPrefix(), "")
	var name, fullName string
	for _, fileDir := range rd {
		name = fileDir.Name()
		fullName = pathName + "/" + name
		if fileDir.IsDir() {
			childNode, err = getDirTreeResp(fullName, operationsMap)
			childNode.IsDir = true
			if err != nil {
				return dirTreeNodeResp{}, err
			}
		} else {
			key := tree.Path + "/" + strings.Split(name, ".")[0]
			if _, ok := operationsMap[key]; !ok {
				continue
			}
			childNode.ID = operationsMap[key].ID
			childNode.Method = operationsMap[key].Method
			childNode.OperationType = operationsMap[key].OperationType
			childNode.IsPublic = operationsMap[key].IsPublic
			childNode.Remark = operationsMap[key].Remark
			childNode.Legal = operationsMap[key].Legal
			childNode.Path = operationsMap[key].Path
			childNode.IsDir = false
			// 如果开关为关
			childNode.Enable = operationsMap[key].Enable
			childNode.CreateTime = operationsMap[key].CreateTime
			childNode.UpdateTime = operationsMap[key].UpdateTime
			childNode.Child = nil
		}
		tree.Child = append(tree.Child, childNode)
	}
	return tree, nil
}

// CreateFile @Title CreateFile
// @Description 创建api文件
// @Accept  json
// @Tags  operateApi
// @Param Method  formData string true "方法Get、POST等等"
// @Param OperationType  formData string true "类型  queries,mutations,subscriptions"
// @Param Status  formData integer true "状态 1共有 2私有"
// @Param Remark  formData string true "说明"
// @Param Legal  formData integer true "是否合法 1合法 2非法"
// @Param Path  formData string true "路径"
// @Param Content  formData string true "内容"
// @Param disable   formData bool true "开关"
// @Success 200 "添加成功"
// @Failure 400	"添加失败"
// @Router /api/v1/operateApi  [POST]
func (a *OperateAPIHandler) CreateFile(c echo.Context) (err error) {
	ctx := c.Request().Context()

	operateParam := domain.FbOperationsResult{}
	err = c.Bind(&operateParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, domain.ParamErr))
	}

	if utils.Empty(strings.TrimSpace(operateParam.Path)) {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, domain.ParamPathEmptyErr))
	}

	if strings.Split(operateParam.Path, "")[0] != "/" {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, domain.ParamErr))
	}

	exist, err := a.CheckFileExist(ctx, operateParam.Path)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, domain.DbFindErr))
	}
	// 判断数据库文件是否存在
	if exist {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, domain.ParamNameExistErr))
	}

	// 不存在则直接写入
	writePath := fmt.Sprintf("%s%s%s%s", utils.GetApiPathPrefix(), operateParam.Path, utils.GetApiPathSuffix(), utils.GetSwitchState(operateParam.Enable))
	err = utils.WriteFile(writePath, operateParam.Content)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, domain.FileWriteErr))
	}
	// 结构体转换
	operate := operateParam.Transform()

	// 设置该schema匹配的参数值
	v := validator.NewValidator()
	generateSchemaPath := fmt.Sprintf("%s/%s", utils.GetWdgGeneratedPath(), domain.GenerateSchemaGraphql)

	// TODO 校验schema
	schemaDocument, err := v.ValidateOperations(ctx, writePath, generateSchemaPath)
	// 默认给他合法
	operate.Legal = domain.Legitimate
	// 解析失败，设置格式为非法
	if err.Error() != "" {
		operate.Legal = domain.UnLegitimate
	}
	operate.SetField(schemaDocument)

	insertID, err := a.operationsUC.Store(ctx, &operate)
	if err != nil {
		// 如果报错添加不成功，则删除之前创建的文件
		os.Remove(writePath)
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, err))
	}

	result, err := a.operationsUC.GetByID(ctx, insertID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, err))
	}

	return c.JSON(http.StatusOK, SuccessWriteResult(result))
}

// CreateDir @Title CreateDir
// @Description 创建api文件夹
// @Accept  json
// @Tags  operateApi
// @Param path formData string true "路径"
// @Success 200 "添加成功"
// @Failure 400	"添加失败"
// @Router /api/v1/operateApi/dir  [POST]
func (a *OperateAPIHandler) CreateDir(c echo.Context) (err error) {
	operate := struct {
		DicPath string `json:"path"`
	}{}
	err = c.Bind(&operate)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, domain.ParamErr))
	}

	if utils.Empty(strings.TrimSpace(operate.DicPath)) {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, domain.ParamNameEmptyErr))
	}
	// 判断文件夹是否存在
	if utils.FileExist(operate.DicPath) {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, domain.ParamNameExistErr))
	}
	// 不存在则直接写入
	writePath := fmt.Sprintf("%s/%s", utils.GetApiPathPrefix(), operate.DicPath)
	err = os.Mkdir(writePath, 0777)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, domain.FileCreateErr))
	}
	return c.JSON(http.StatusOK, SuccessResult())
}

// ReNameDir @Title ReNameDir
// @Description 重命名文件夹
// @Accept  json
// @Tags  operateApi
// @Param oldPath formData string true "旧路径"
// @Param newPath formData string true "新路径"
// @Param disable formData bool true "开关"
// @Param id formData integer true "id"
// @Success 200 "修改成功"
// @Failure 400	"修改失败"
// @Router /api/v1/operateApi/dir  [PUT]
func (a *OperateAPIHandler) ReNameDir(c echo.Context) (err error) {
	ctx := c.Request().Context()
	fileName := struct {
		OldPath string `json:"oldPath"`
		NewPath string `json:"newPath"`
	}{}
	err = c.Bind(&fileName)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, domain.ParamErr))
	}
	if utils.Empty(strings.TrimSpace(fileName.NewPath)) {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, domain.ParamNameEmptyErr))
	}
	operationOldPath := fmt.Sprintf("%s%s", utils.GetApiPathPrefix(), fileName.OldPath)
	operationNewPath := fmt.Sprintf("%s%s", utils.GetApiPathPrefix(), fileName.NewPath)
	// 修改数据库内的路径
	err = a.operationsUC.ReDicName(ctx, fileName.OldPath+"/", fileName.NewPath+"/")
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, domain.FileReNameErr))
	}

	// 修改operation
	err = os.Rename(operationOldPath, operationNewPath)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, domain.FileReNameErr))
	}

	// 修改hooks文件夹
	hooksOldPath := fmt.Sprintf("%s%s", utils.GetOperationsHookPathPrefix(), fileName.OldPath)
	hooksNewPath := fmt.Sprintf("%s%s", utils.GetOperationsHookPathPrefix(), fileName.NewPath)
	err = ReNameDic(hooksOldPath, hooksNewPath)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, domain.FileReNameErr))
	}

	// 修改设置文件夹
	settingOldPath := fmt.Sprintf("%s%s", utils.GetOperationSettingPath(), fileName.OldPath)
	settingNewPath := fmt.Sprintf("%s%s", utils.GetOperationSettingPath(), fileName.NewPath)
	err = ReNameDic(settingOldPath, settingNewPath)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, domain.FileReNameErr))
	}

	// mock设置文件夹
	mockOldPath := fmt.Sprintf("%s%s", utils.GetMockPath(), fileName.OldPath)
	mockNewPath := fmt.Sprintf("%s%s", utils.GetMockPath(), fileName.NewPath)
	err = ReNameDic(mockOldPath, mockNewPath)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, domain.FileReNameErr))
	}
	return c.JSON(http.StatusOK, SuccessResult())
}

func ReNameDic(oldPath, newPath string) (err error) {
	// 文件不存在
	if !utils.FileExist(oldPath) {
		// 创建
		err = os.MkdirAll(oldPath, os.ModePerm)
		if err != nil {
			return
		}
	}
	err = os.Rename(oldPath, newPath)
	if err != nil {
		return
	}
	return
}

func operationSettingReName(oldName, newName string) (err error) {
	operationPaths := utils.GetDicAllChildFilePath(utils.GetOperationSettingPath())
	for _, row := range operationPaths {
		row = filepath.ToSlash(row)
		oldNameOn := filepath.ToSlash(fmt.Sprintf("%s%s%s", utils.GetOperationSettingPath(), oldName, utils.GetSwitchState(utils.ON)))
		oldNameOff := filepath.ToSlash(fmt.Sprintf("%s%s%s", utils.GetOperationSettingPath(), oldName, utils.GetSwitchState(utils.OFF)))
		if !(row == oldNameOn || row == oldNameOff) {
			continue
		}
		replaceNewName := fmt.Sprintf("%s%s", utils.GetOperationSettingPath(), newName)
		if row == oldNameOn {
			replaceNewName = fmt.Sprintf("%s%s", newName, utils.GetSwitchState(utils.ON))
		}
		if row == oldNameOff {
			replaceNewName = fmt.Sprintf("%s%s", newName, utils.GetSwitchState(utils.OFF))
		}
		replaceNewName = filepath.ToSlash(replaceNewName)
		err = os.Rename(row, replaceNewName)
		if err != nil {
			log.Error("operationSettingReName fail ,err : ", err)
		}
	}
	return err
}

func (a *OperateAPIHandler) CheckFileExist(ctx context.Context, path string) (bool, error) {
	operation, err := a.operationsUC.GetByPath(ctx, path)
	if err != nil {
		return false, err
	}
	if operation.ID == 0 {
		return false, nil
	}
	return true, nil
	//offName := filepath.ToSlash(fmt.Sprintf("%s/%s%s%s", utils.GetApiPathPrefix(), path, utils.GetApiPathSuffix(), utils.GetSwitchOff()))
	//onName := filepath.ToSlash(fmt.Sprintf("%s/%s%s", utils.GetApiPathPrefix(), path, utils.GetApiPathSuffix()))
	//offName = strings.ReplaceAll(offName, "//", "/")
	//onName = strings.ReplaceAll(onName, "//", "/")
	//
	//// 获取所有operation
	//filePaths := utils.GetDicAllChildFilePath(utils.GetApiPathPrefix())
	//for _, apiPath := range filePaths {
	//	apiPath = filepath.ToSlash(apiPath)
	//	// 存在则直接返回，已存在
	//	if apiPath == onName || apiPath == offName {
	//		// 存在
	//		return true
	//	}
	//}
	//// 不存在
	//return false
}

// UpdateContentById @Title UpdateContentById
// @Description 更新api
// @Accept  json
// @Tags  operateApi
// @Param id  formData integer true "id"
// @Param Content  formData string true "内容"
// @Param operation_type  formData string true "请求类型 queries,mutations,subscriptions"
// @Success 200 "更改成功"
// @Failure 400	"更改失败"
// @Router /api/v1/operateApi/content/:id  [PUT]
func (a *OperateAPIHandler) UpdateContentById(c echo.Context) (err error) {
	ctx := c.Request().Context()
	operateParam := domain.FbOperationsResult{}
	err = c.Bind(&operateParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, domain.ParamErr))
	}

	id := c.Param("id")
	operationID, _ := strconv.Atoi(id)
	if operationID == 0 {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, domain.ParamIdEmptyErr))
	}
	operation, err := a.operationsUC.GetByID(ctx, int64(operationID))
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, domain.DbFindErr))
	}
	writePath := fmt.Sprintf("%s%s%s%s", utils.GetApiPathPrefix(), operation.Path, utils.GetApiPathSuffix(), utils.GetSwitchState(operation.Enable == 0))

	// 写入文件内容，因为相比其他信息，文件内容优先级最低，放到最后
	err = utils.WriteFile(writePath, operateParam.Content)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, domain.FileWriteErr))
	}
	// 设置该schema匹配的参数值
	v := validator.NewValidator()
	generateSchemaPath := fmt.Sprintf("%s/%s", utils.GetWdgGeneratedPath(), domain.GenerateSchemaGraphql)

	// TODO 校验schema
	schemaDocument, err := v.ValidateOperations(ctx, writePath, generateSchemaPath)
	// 默认给他合法
	operation.Legal = domain.Legitimate
	// 解析失败，设置格式为非法
	if err.Error() != "" {
		operation.Legal = domain.UnLegitimate
	}
	operation.SetField(schemaDocument)
	// 修改operation的查询类型
	_, err = a.operationsUC.Update(ctx, &operation)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, err))
	}
	return c.JSON(http.StatusOK, SuccessWriteResult(operateParam.Content))
}

// Update @Title Update
// @Description 更新api
// @Accept  json
// @Tags  operateApi
// @Param Method  formData string true "方法Get、POST等等"
// @Param OperationType  formData string true "类型  queries,mutations,subscriptions"
// @Param Status  formData integer true "状态 1共有 2私有"
// @Param Path  formData string true "路径"
// @Param Remark  formData string true "说明"
// @Param Legal  formData integer true "是否合法 1合法 2非法"
// @Param Content  formData string true "内容"
// @Param disable   formData bool true "开关"
// @Success 200 "更改成功"
// @Failure 400	"更改失败"
// @Router /api/v1/operateApi/:id  [PUT]
func (a *OperateAPIHandler) Update(c echo.Context) (err error) {
	ctx := c.Request().Context()
	operateParam := domain.FbOperationsResult{}
	err = c.Bind(&operateParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, domain.ParamErr))
	}

	id := c.Param("id")
	operationID, _ := strconv.Atoi(id)
	operateParam.ID = int64(operationID)
	newOperation := operateParam.Transform()
	// 先拿到旧operation的信息
	oldOperation, _ := a.operationsUC.GetByID(ctx, operateParam.ID)
	// 如果更改了路径或者开关，需要重新重命名operation文件和hooks、setting、mooks文件
	//if oldOperation.Path != newOperation.Path {
	//
	//	// 修改数据库,operationsUC.ReName 该接口会将
	//	err = a.operationsUC.ReName(ctx, fileName.OldPath, fileName.NewPath, enable)
	//	if err != nil {
	//		// 将之前修改的operation文件改回来
	//		os.Rename(fileName.NewPath, fileName.OldPath)
	//		// 将之前修改的hooks文件改回来
	//		a.hooksUC.OperationHooksReName(fileName.NewPath, fileName.OldPath)
	//		// 将之前修改的setting文件改回来
	//		operationSettingReName(fileName.NewPath, fileName.OldPath)
	//		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, domain.FileReNameErr))
	//	}
	//	//
	//	//if operID == 0 {
	//	//	return c.JSON(http.StatusOK, SuccessResult())
	//	//}
	//	//result, err := a.operationsUC.GetByID(ctx, int64(operID))
	//	//if err != nil {
	//	//	return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, domain.FileReNameErr))
	//	//}
	//
	//}
	//
	//operate.Path = oldOperation.Path
	//updateId, err := a.operationsUC.Update(ctx, &operate)
	//if err != nil {
	//	return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, err))
	//}
	// 不为空则说明改名了

	oldName := fmt.Sprintf("%s%s%s%s", utils.GetApiPathPrefix(), oldOperation.Path, utils.GetApiPathSuffix(), utils.GetSwitchState(oldOperation.Enable == 0))
	newName := fmt.Sprintf("%s%s%s%s", utils.GetApiPathPrefix(), newOperation.Path, utils.GetApiPathSuffix(), utils.GetSwitchState(newOperation.Enable == 0))

	if !utils.Empty(operateParam.Path) {
		// 修改operation文件名
		err = os.Rename(oldName, newName)
		if err != nil {
			return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, domain.FileReNameErr))
		}

		// 修改hooks文件
		err = a.hooksUC.OperationHooksReName(oldOperation.Path, newOperation.Path)
		if err != nil {
			// 如果hooks文件名修改失败，则将之前更改的operation文件改回来
			os.Rename(newName, oldName)
			return c.JSON(http.StatusOK, GetResponseErr(domain.OperateAPICode, domain.FileReNameErr))
		}
		// 修改setting文件
		err = operationSettingReName(oldOperation.Path, newOperation.Path)
		if err != nil {
			// 如果setting文件名修改失败，则将之前更改的operation文件和hooks改回来
			os.Rename(newName, oldName)
			a.hooksUC.OperationHooksReName(newOperation.Path, oldOperation.Path)
			return c.JSON(http.StatusOK, GetResponseErr(domain.OperateAPICode, domain.FileReNameErr))
		}
	}

	// 修改数据库信息
	_, err = a.operationsUC.Update(ctx, &newOperation)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, err))
	}

	result, err := a.operationsUC.GetByID(ctx, newOperation.ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.OperateAPICode, err))
	}

	return c.JSON(http.StatusOK, SuccessWriteResult(result))
}
