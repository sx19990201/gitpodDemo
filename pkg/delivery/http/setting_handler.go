package http

import (
	"encoding/json"
	"fmt"
	"github.com/fire_boom/domain"
	"github.com/fire_boom/utils"
	"github.com/labstack/echo/v4"
	"github.com/tidwall/sjson"
	"io/ioutil"
	"net/http"
)

type SettingHandler struct {
}

func InitSettingHandler(e *echo.Echo) {
	//timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second
	NewSettingHandler(e)
}

func NewSettingHandler(e *echo.Echo) {
	handler := &SettingHandler{}
	v1 := e.Group("/api/v1")
	{
		authentication := v1.Group("/setting")
		{
			authentication.POST("", handler.Update)
			authentication.GET("/systemConfig", handler.GetSystemConfig)
			authentication.GET("/versionConfig", handler.GetVersionConfig)
			authentication.GET("/environmentConfig", handler.GetEnvironmentConfig)
			authentication.GET("/securityConfig", handler.GetSecurityConfig)
			authentication.GET("/corsConfiguration", handler.GetCorsConfiguration)
			//authentication.GET("/restart", handler.ReStart)
		}
	}
}

// ReStart @Title ReStart
// @Description 获取跨域设置
// @Accept  json
// @Tags  setting
// @Success 200 {object} domain.CorsConfiguration "获取成功"
// @Failure 400	"获取失败"
// @Router /api/v1/setting/restart  [GET]
//func (s *SettingHandler) ReStart(c echo.Context) (err error) {
//	// 获取启动命令
//	cmd := GetStartCMD()
//
//	if err != nil {
//		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.SettingCode, err))
//	}
//
//	return c.JSON(http.StatusOK, SuccessResult())
//}

func GetStartCMD() string {
	//globalConfig, err := utils.GetGlobalConfig()
	setting, _ := utils.GetSettingConfig()

	s := setting.System
	return fmt.Sprintf("wunderctl %s %s %s %s %s", s.GetDevSwitchCMD(), s.GetLogLevelCMD(), s.GetApiPortCMD(), s.GetMiddlewarePortCMD(), s.GetForcedJumpCMD())
}

// GetSystemConfig @Title GetSystemConfig
// @Description 获取系统设置
// @Accept  json
// @Tags  setting
// @Success 200 {object} domain.SystemConfig "获取成功"
// @Failure 400	"获取失败"
// @Router /api/v1/setting/systemConfig  [GET]
func (s *SettingHandler) GetSystemConfig(c echo.Context) (err error) {
	setting, err := utils.GetSettingConfig()
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.SettingCode, err))
	}
	return c.JSON(http.StatusOK, SuccessWriteResult(setting.System))
}

// GetVersionConfig @Title GetVersionConfig
// @Description 获取版本设置
// @Accept  json
// @Tags  setting
// @Success 200 {object} domain.VersionConfig "获取成功"
// @Failure 400	"获取失败"
// @Router /api/v1/setting/versionConfig  [GET]
func (s *SettingHandler) GetVersionConfig(c echo.Context) (err error) {
	setting, err := utils.GetSettingConfig()
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.SettingCode, err))
	}
	return c.JSON(http.StatusOK, SuccessWriteResult(setting.Version))
}

// GetEnvironmentConfig @Title GetEnvironmentConfig
// @Description 获取环境变量设置
// @Accept  json
// @Tags  setting
// @Success 200 {object} domain.EnvironmentConfig "获取成功"
// @Failure 400	"获取失败"
// @Router /api/v1/setting/environmentConfig  [GET]
func (s *SettingHandler) GetEnvironmentConfig(c echo.Context) (err error) {
	setting, err := utils.GetSettingConfig()
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.SettingCode, err))
	}
	return c.JSON(http.StatusOK, SuccessWriteResult(setting.Environment))
}

// GetSecurityConfig @Title GetSecurityConfig
// @Description 获取安全设置
// @Accept  json
// @Tags  setting
// @Success 200 {object} domain.SecurityConfig "获取成功"
// @Failure 400	"获取失败"
// @Router /api/v1/setting/securityConfig  [GET]
func (s *SettingHandler) GetSecurityConfig(c echo.Context) (err error) {
	globalConfig, err := utils.GetGlobalConfig()
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.SettingCode, err))
	}
	return c.JSON(http.StatusOK, SuccessWriteResult(globalConfig.ConfigureWunderGraphApplication.Security))
}

// GetCorsConfiguration @Title GetCorsConfiguration
// @Description 获取跨域设置
// @Accept  json
// @Tags  setting
// @Success 200 {object} domain.CorsConfiguration "获取成功"
// @Failure 400	"获取失败"
// @Router /api/v1/setting/corsConfiguration  [GET]
func (s *SettingHandler) GetCorsConfiguration(c echo.Context) (err error) {
	globalConfig, err := utils.GetGlobalConfig()
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.SettingCode, err))
	}
	return c.JSON(http.StatusOK, SuccessWriteResult(globalConfig.ConfigureWunderGraphApplication.Cors))
}

// Update @Title Update
// @Description 修改设置
// @Accept  json
// @Tags  setting
// @Param key formData string true "key=路径"
// @Param val formData integer true "val=值"
// @Success 200 {object} domain.CorsConfiguration "修改成功"
// @Failure 400	"修改失败"
// @Router /api/v1/setting  [POST]
func (s *SettingHandler) Update(c echo.Context) (err error) {
	params := struct {
		Key string      `json:"key"`
		Val interface{} `json:"val"`
	}{}
	err = c.Bind(&params)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.SettingCode, domain.ParamErr))
	}
	config, err := ioutil.ReadFile(utils.GetSettingPath())
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.SettingCode, domain.FileReadErr))
	}
	newConfigJson, err := sjson.Set(string(config), fmt.Sprintf("*.%s", params.Key), params.Val)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.SettingCode, domain.JsonSetErr))
	}
	var setting domain.Setting
	err = json.Unmarshal([]byte(newConfigJson), &setting)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.SettingCode, domain.JsonSetErr))
	}
	if setting.System.DevSwitch == true {
		setting.System.ForcedJumpSwitch = true
	}
	content, err := json.Marshal(setting)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.SettingCode, domain.JsonSetErr))
	}
	err = utils.WriteFile(utils.GetSettingPath(), string(content))
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.SettingCode, domain.FileWriteErr))
	}
	return c.JSON(http.StatusOK, SuccessResult())
}
