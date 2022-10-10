package http

import (
	"fmt"
	"github.com/fire_boom/domain"
	"github.com/fire_boom/utils"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
)

const (
	postAuthentication         = "postAuthentication"
	mutatingPostAuthentication = "mutatingPostAuthentication"
	onRequest                  = "onRequest"
	onResponse                 = "onResponse"
)

type ScriptHandler struct {
}

func InitScriptHandler(e *echo.Echo) {
	//timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second
	NewScriptHandler(e)
}

func NewScriptHandler(e *echo.Echo) {
	handler := &ScriptHandler{}
	v1 := e.Group("/api/v1")
	{
		authentication := v1.Group("/script")
		{
			authentication.POST("", handler.Store)
			authentication.GET("", handler.GetScriptByName)
		}
	}
}

// Store @Title Store
// @Description 添加脚本
// @Accept  json
// @Tags  script
// @Param name formData string true "名称"
// @Param content formData string true "内容"
// @Success 200 "添加成功"
// @Failure 400	"添加失败"
// @Router /api/v1/script  [POST]
func (s *ScriptHandler) Store(c echo.Context) (err error) {
	script := struct {
		Name    string `json:"name"`
		Content string `json:"content"`
	}{}
	err = c.Bind(&script)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.ScriptCode, domain.ParamErr))
	}
	if utils.Empty(script.Name) {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.ScriptCode, domain.ParamNameEmptyErr))
	}
	path := getPath(script.Name)

	err = utils.WriteFile(path, script.Content)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.ScriptCode, domain.FileWriteErr))
	}
	return c.JSON(http.StatusOK, SuccessResult())
}

// GetScriptByName @Title GetScriptByName
// @Description 根据脚本名获取脚本
// @Accept  json
// @Tags  script
// @Param scriptName formData string true "名称"
// @Success 200 "添加成功"
// @Failure 400	"添加失败"
// @Router /api/v1/script  [GET]
func (s *ScriptHandler) GetScriptByName(c echo.Context) (err error) {
	param := struct {
		Name string `json:"name"`
	}{}

	err = c.Bind(&param)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.ScriptCode, domain.ParamErr))
	}
	if utils.Empty(param.Name) {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.ScriptCode, domain.ParamNameEmptyErr))
	}
	path := getPath(param.Name)

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.ScriptCode, domain.FileReadErr))
	}
	return c.JSON(http.StatusOK, SuccessWriteResult(string(content)))
}

func getPath(name string) string {
	switch name {
	case postAuthentication:
		return fmt.Sprintf("%s/%s", utils.GetAuthGlobalHookPathPrefix(), postAuthentication)
	case mutatingPostAuthentication:
		return fmt.Sprintf("%s/%s", utils.GetAuthGlobalHookPathPrefix(), mutatingPostAuthentication)
	case onRequest:
		return fmt.Sprintf("%s/%s", utils.GetGlobalHookPathPrefix(), onRequest)
	case onResponse:
		return fmt.Sprintf("%s/%s", utils.GetGlobalHookPathPrefix(), onResponse)
	default:
		return ""
	}
}
