package http

import (
	"encoding/json"
	"fmt"
	"github.com/fire_boom/domain"
	"github.com/fire_boom/utils"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
)

type HookHandler struct {
}

func InitHookHandler(e *echo.Echo) {
	NewHookHandler(e)
}

func NewHookHandler(e *echo.Echo) {
	handler := &HookHandler{}
	v1 := e.Group("/api/v1")
	{
		hook := v1.Group("/hook")
		{
			hook.POST("/script", handler.HookScript) // TODO mock保存钩子内容
			hook.POST("/input", handler.HookInput)   // TODO mock保存钩子输入参数
			hook.POST("/depend", handler.HookDepend) // TODO mock保存钩子依赖
			hook.POST("/switch", handler.HookSwitch) // TODO mock修改钩子开关
			hook.POST("/run", handler.HookRun)       // TODO mock执行hooks
			hook.GET("", handler.Hook)               // TODO mock获取脚本信息
		}
	}
}

// HookScript @Title HookScript
// @Description 保存钩子内容
// @Accept  json
// @Tags  hook
// @Param path formData string true "钩子路径"
// @Param switch formData boolean true "钩子类型"
// @Success 200 "修改成功"
// @Failure 400	"修改失败"
// @Router /api/v1/hook/script  [POST]
func (d *HookHandler) HookScript(c echo.Context) (err error) {
	param := domain.HookStruct{}
	err = c.Bind(&param)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.ParamErr))
	}
	if param.Path == "" {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.ParamPathEmptyErr))
	}
	var hooks domain.HookStruct
	err = json.Unmarshal([]byte(domain.HookDemo), &hooks)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.JsonUnMarshalErr))
	}
	hooks.HookSwitch = param.HookSwitch
	content, _ := json.Marshal(hooks)
	domain.HookDemo = string(content)
	return c.JSON(http.StatusOK, SuccessResult())
}

// HookInput @Title HookInput
// @Description 保存钩子输入参数
// @Accept  json
// @Tags  hook
// @Param path formData string true "路径"
// @Param input formData object true "输入"
// @Success 200 "修改成功"
// @Failure 400	"修改失败"
// @Router /api/v1/hook/input  [POST]
func (d *HookHandler) HookInput(c echo.Context) (err error) {
	param := domain.HookStruct{}
	err = c.Bind(&param)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.ParamErr))
	}
	if param.Path == "" {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.ParamPathEmptyErr))
	}
	var hooks domain.HookStruct
	err = json.Unmarshal([]byte(domain.HookDemo), &hooks)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.JsonUnMarshalErr))
	}
	hooks.Input = param.Input
	content, _ := json.Marshal(hooks)
	domain.HookDemo = string(content)
	return c.JSON(http.StatusOK, SuccessResult())
}

// HookDepend @Title HookDepend
// @Description 保存钩子依赖
// @Accept  json
// @Tags  hook
// @Param path formData string true "路径"
// @Param depend1 formData []object true "依赖数组,swagger的问题 把depend1中的1给去掉哈"
// @Success 200 "修改成功"
// @Failure 400	"修改失败"
// @Router /api/v1/hook/depend  [POST]
func (d *HookHandler) HookDepend(c echo.Context) (err error) {
	param := domain.HookStruct{}
	err = c.Bind(&param)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.ParamErr))
	}
	if param.Path == "" {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.ParamPathEmptyErr))
	}
	var hooks domain.HookStruct
	err = json.Unmarshal([]byte(domain.HookDemo), &hooks)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.JsonUnMarshalErr))
	}
	hooks.Depends = param.Depends
	content, _ := json.Marshal(hooks)
	domain.HookDemo = string(content)
	return c.JSON(http.StatusOK, SuccessResult())
}

// HookSwitch @Title HookSwitch
// @Description 保存钩子依赖
// @Accept  json
// @Tags  hook
// @Param path formData string true "路径"
// @Param switch formData boolean true "开关 true开 false 关"
// @Success 200 "修改成功"
// @Failure 400	"修改失败"
// @Router /api/v1/hook/switch  [POST]
func (d *HookHandler) HookSwitch(c echo.Context) (err error) {
	param := domain.HookStruct{}
	err = c.Bind(&param)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.ParamErr))
	}
	if param.Path == "" {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.ParamPathEmptyErr))
	}
	var hooks domain.HookStruct
	err = json.Unmarshal([]byte(domain.HookDemo), &hooks)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.JsonUnMarshalErr))
	}
	hooks.HookSwitch = param.HookSwitch
	content, _ := json.Marshal(hooks)
	domain.HookDemo = string(content)
	return c.JSON(http.StatusOK, SuccessResult())
}

// HookRun @Title HookRun
// @Description 执行钩子
// @Accept  json
// @Tags  hook
// @Param path formData string true "路径"
// @Param depend1 formData []object true "依赖数组,swagger的问题 把depend1中的1给去掉哈"
// @Param input formData object true "输入"
// @Param script formData string true "钩子内容"
// @Param scriptType formData string true "钩子类型"
// @Success 200 "修改成功"
// @Failure 400	"修改失败"
// @Router /api/v1/auth/hooks  [POST]
func (d *HookHandler) HookRun(c echo.Context) (err error) {
	param := domain.HookStruct{}
	err = c.Bind(&param)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.ParamErr))
	}
	if param.Path == "" {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.ParamPathEmptyErr))
	}

	var hooks domain.HookStruct
	err = json.Unmarshal([]byte(domain.HookDemo), &hooks)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.JsonUnMarshalErr))
	}
	if param.Depends != nil {
		hooks.Depends = param.Depends
	}
	if param.Input != nil {
		hooks.Input = param.Input
	}
	if param.Script != "" {
		hooks.Script = param.Script
	}
	if param.ScriptType != "" {
		hooks.ScriptType = param.ScriptType
	}

	content, _ := json.Marshal(hooks)
	domain.HookDemo = string(content)
	return c.JSON(http.StatusOK, SuccessWriteResult("恭喜你运行成功了哈哈哈哈哈哈哈哈哈哈哈哈哈哈哈哈哈哈"))
}

// Hook @Title Hook
// @Description 获取钩子内容
// @Accept  json
// @Tags  hook
// @Param path formData string true "钩子路径"
// @Success 200 {object} []domain.HookStruct "查询成功"
// @Failure 400	"查询失败"
// @Router /api/v1/auth/hooks  [GET]
func (d *HookHandler) Hook(c echo.Context) (err error) {
	param := domain.HookStruct{}
	err = c.Bind(&param)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.ParamErr))
	}
	if param.Path == "" {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.ParamPathEmptyErr))
	}
	// 拼接hooks路径
	hookPath := fmt.Sprintf("%s%s", utils.GetNewHookPathPrefix(), param.Path)
	if !utils.FileExist(hookPath) {
		return c.JSON(http.StatusOK, SuccessWriteResult(""))
	}
	bytes, err := ioutil.ReadFile(hookPath)

	var hooks domain.HookStruct
	err = json.Unmarshal(bytes, &hooks)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.JsonUnMarshalErr))
	}

	return c.JSON(http.StatusOK, SuccessWriteResult(hooks))
}
