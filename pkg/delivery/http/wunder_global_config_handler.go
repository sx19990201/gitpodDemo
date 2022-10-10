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

type GlobalConfigHandler struct {
}

func InitGlobalConfigHandler(e *echo.Echo) {
	NewGlobalConfigHandler(e)
}

func NewGlobalConfigHandler(e *echo.Echo) {
	handler := &GlobalConfigHandler{}
	v1 := e.Group("/api/v1")
	{
		authentication := v1.Group("/global")
		{
			authentication.POST("", handler.Update)
		}
	}
}

// Update @Title Update
// @Description 修改全局配置
// @Accept  json
// @Tags  global
// @Param key formData string true "key=路径"
// @Param val formData integer true "val=值"
// @Success 200 "修改成功"
// @Failure 400	"修改失败"
// @Router /api/v1/global  [POST]
func (s *GlobalConfigHandler) Update(c echo.Context) (err error) {
	params := struct {
		Key string      `json:"key"`
		Val interface{} `json:"val"`
	}{}
	err = c.Bind(&params)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.StorageBucketCode, domain.ParamErr))
	}
	config, err := ioutil.ReadFile(utils.GetGlobalConfigPath())
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.StorageBucketCode, domain.FileReadErr))
	}
	// 判断是否是跨域或者安全
	if params.Key == "allowedOrigins" || params.Key == "allowedMethods" ||
		params.Key == "allowedHeaders" || params.Key == "exposedHeaders" ||
		params.Key == "maxAge" || params.Key == "allowCredentials" {
		params.Key = fmt.Sprintf("cors.%s", params.Key)
	}
	if params.Key == "allowedHosts" || params.Key == "enableGraphQLEndpoint" {
		params.Key = fmt.Sprintf("security.%s", params.Key)
	}
	newConfigJson := ""
	if params.Key == "cors.allowedMethods" {
		var configStruct domain.GlobalConfig
		err := json.Unmarshal(config, &configStruct)
		if err != nil {
			return c.JSON(http.StatusBadRequest, GetResponseErr(domain.StorageBucketCode, domain.FileReadErr))
		}

		arr := make([]string, 0)
		bytes, _ := json.Marshal(params.Val)
		err = json.Unmarshal(bytes, &arr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, GetResponseErr(domain.StorageBucketCode, domain.FileReadErr))
		}
		configStruct.ConfigureWunderGraphApplication.Cors.AllowedMethods = arr
		configBytes, _ := json.Marshal(configStruct)
		newConfigJson = string(configBytes)
	} else {
		newConfigJson, err = sjson.Set(string(config), fmt.Sprintf("*.%s", params.Key), params.Val)
		if err != nil {
			return c.JSON(http.StatusBadRequest, GetResponseErr(domain.StorageBucketCode, domain.JsonSetErr))
		}
	}

	err = utils.WriteFile(utils.GetGlobalConfigPath(), newConfigJson)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.StorageBucketCode, domain.FileWriteErr))
	}
	return c.JSON(http.StatusOK, SuccessResult())
}
