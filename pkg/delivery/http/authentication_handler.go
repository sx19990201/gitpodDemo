package http

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/fire_boom/utils"
	"github.com/labstack/gommon/log"
	"net/http"
	"strconv"
	"time"

	"github.com/fire_boom/domain"
	"github.com/fire_boom/pkg/repository/sqllite"
	"github.com/fire_boom/pkg/usecase"
	"github.com/labstack/echo/v4"
)

const (
	openid        = "openid"
	github        = "github"
	google        = "google"
	DiscoveryURL  = ".well-known/openid-configuration"
	UserInfoPoint = "me"
	JwksURL       = ".well-known/jwks.json"
)

type AuthenticationHandler struct {
	AuthenticationUseCase domain.AuthenticationUseCase
	HookUseCase           domain.HooksUseCase
}

func InitAuthenticationUseCaseRouter(e *echo.Echo, db *sql.DB) {
	authenticationRepo := sqllite.NewAuthenticationRepository(db)
	timeoutContext := time.Duration(utils.GetTimeOut()) * time.Second
	au := usecase.NewAuthenticationUseCase(authenticationRepo, timeoutContext)
	hookUC := usecase.NewHooksUseCase(timeoutContext)
	NewAuthenticationHandler(e, au, hookUC)
}

func NewAuthenticationHandler(e *echo.Echo, aUseCase domain.AuthenticationUseCase, hookUseCase domain.HooksUseCase) {
	handler := &AuthenticationHandler{
		AuthenticationUseCase: aUseCase,
		HookUseCase:           hookUseCase,
	}
	v1 := e.Group("/api/v1")
	{
		authentication := v1.Group("/auth")
		{
			authentication.GET("", handler.FindAuthentication)
			authentication.POST("", handler.Store)
			authentication.PUT("", handler.Update)
			authentication.DELETE("/:id", handler.Delete)
			authentication.GET("/redirectUrl", handler.GetRedirectURL)
			authentication.POST("/redirectUrl", handler.UpdateRedirectURL)
			authentication.GET("/getDiscoveryURL", handler.GetDiscoveryURL)

			authentication.GET("/hooks", handler.GetHooks)
			authentication.POST("/hooks", handler.UpdateHooks)

		}
	}
}

// UpdateHooks @Title UpdateHooks
// @Description 修改hook
// @Accept  json
// @Tags  auth
// @Param fileName formData string true "路径"
// @Param content formData string true "内容"
// @Param hookSwitch formData bool true "开关"
// @Param hookName formData string true "preResolve, postResolve, customeResolve, mutatingPostResolve"
// @Success 200 "修改成功"
// @Failure 400	"修改失败"
// @Router /api/v1/auth/hooks  [POST]
func (d *AuthenticationHandler) UpdateHooks(c echo.Context) (err error) {
	var hooks domain.Hooks
	err = c.Bind(&hooks)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.ParamErr))
	}
	err = d.HookUseCase.UpdateHooksByPath(utils.GetAuthGlobalHookPathPrefix(), hooks)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.FileReadErr))
	}
	return c.JSON(http.StatusOK, SuccessResult())
}

// GetHooks @Title GetHooks
// @Description 获取hooks
// @Accept  json
// @Tags  auth
// @Param fileName formData string true "hook名称"
// @Success 200 {object} []domain.Hooks "获取成功"
// @Failure 400	"获取失败"
// @Router /api/v1/auth/hooks  [GET]
func (d *AuthenticationHandler) GetHooks(c echo.Context) (err error) {
	result, err := d.HookUseCase.FindHooksByPath(utils.GetAuthGlobalHookPathPrefix(), "", domain.AuthGlobalHook)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.ParamErr))
	}
	return c.JSON(http.StatusOK, SuccessWriteResult(result))
}

// GetDiscoveryURL @Title GetDiscoveryURL
// @Description 获取服务发现地址
// @Accept  json
// @Tags  auth
// @Param issuer formData string true "issuer地址"
// @Param id formData string true "供应商id"
// @Success 200 {object} []string "获取成功"
// @Failure 400	"获取失败"
// @Router /api/v1/auth/getDiscoveryURL  [GET]
func (d *AuthenticationHandler) GetDiscoveryURL(c echo.Context) (err error) {
	param := struct {
		Issuer string `json:"issuer"`
		ID     string `json:"id"`
	}{}
	err = c.Bind(&param)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.ParamErr))
	}
	url := fmt.Sprintf("%s/%s", param.Issuer, param.ID)
	return c.JSON(http.StatusOK, SuccessWriteResult(url))
}

// GetRedirectURL @Title GetRedirectURL
// @Description 获取重定向url
// @Accept  json
// @Tags  auth
// @Success 200 {object} []string "获取成功"
// @Failure 400	"获取失败"
// @Router /api/v1/auth/redirectUrl  [GET]
func (d *AuthenticationHandler) GetRedirectURL(c echo.Context) (err error) {
	globalConfig, err := utils.GetGlobalConfig()
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, err))
	}
	return c.JSON(http.StatusOK, SuccessWriteResult(globalConfig.AuthRedirectURL))
}

// UpdateRedirectURL @Title UpdateRedirectURL
// @Description 修改重定向url
// @Accept  json
// @Tags  auth
// @Param redirectURLs formData []string true "重定向url数组"
// @Success 200 "修改成功"
// @Failure 400	"获取失败"
// @Router /api/v1/auth/redirectUrl  [POST]
func (d *AuthenticationHandler) UpdateRedirectURL(c echo.Context) (err error) {
	params := struct {
		RedirectURLs []string `json:"redirectURLs"`
	}{}
	err = c.Bind(&params)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.ParamErr))
	}
	globalConfig, err := utils.GetGlobalConfig()
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, err))
	}
	globalConfig.AuthRedirectURL = params.RedirectURLs
	jsonContent, err := json.Marshal(globalConfig)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.JsonMarshalErr))
	}
	err = utils.WriteFile(utils.GetGlobalConfigPath(), string(jsonContent))
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.StorageBucketCode, domain.FileWriteErr))
	}
	return c.JSON(http.StatusOK, SuccessResult())
}

// Store @Title Store
// @Description 添加身份验证信息
// @Accept  json
// @Tags  auth
// @Param name formData string true "名称"
// @Param authSupplier formData integer true "验证供应商类型"
// @Param switchState formData []string true "开关"
// @Param config formData string true "验证配置"
// @Success 200 "添加成功"
// @Failure 400	"添加失败"
// @Router /api/v1/auth  [POST]
func (d *AuthenticationHandler) Store(c echo.Context) (err error) {
	var authentication domain.FbAuthenticationResp
	err = c.Bind(&authentication)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.ParamErr))
	}
	if err = authenticationStoreCheckParams(authentication); err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, err))
	}
	if authentication.SwitchState == nil {
		authentication.SwitchState = make([]string, 0)
	}
	// 服务发现地址
	authentication.Config.DiscoveryURL = fmt.Sprintf("%s/%s", authentication.Config.Issuer, DiscoveryURL)
	authentication.Config.UserInfoEndpoint = fmt.Sprintf("%s/%s", authentication.Config.Issuer, UserInfoPoint)
	if authentication.Config.Jwks == 0 {
		authentication.Config.JwksURL = fmt.Sprintf("%s/%s", authentication.Config.Issuer, JwksURL)
	}
	jsonByte, err := json.Marshal(authentication.SwitchState)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.JsonMarshalErr))
	}
	authConfigByte, _ := json.Marshal(authentication.Config)

	insertData := domain.FbAuthentication{
		ID:           authentication.ID,
		Name:         authentication.Name,
		AuthSupplier: authentication.AuthSupplier,
		SwitchState:  string(jsonByte),
		Config:       string(authConfigByte),
		CreateTime:   authentication.CreateTime,
		UpdateTime:   authentication.UpdateTime,
		IsDel:        authentication.IsDel,
	}

	ctx := c.Request().Context()
	lastID, err := d.AuthenticationUseCase.Store(ctx, &insertData)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, err))
	}

	return c.JSON(http.StatusOK, SuccessWriteResult(lastID))
}

func authenticationStoreCheckParams(auth domain.FbAuthenticationResp) (err error) {
	if ok := utils.Empty(auth.Name); ok {
		return domain.ParamNameExistErr
	}
	//if ok := checkAuthSupplier(auth.AuthSupplier); ok {
	//	return domain.ParamSupplierNotExistErr
	//}
	return
}

func checkAuthSupplier(authSupplier string) bool {
	switch authSupplier {
	case openid:
	case github:
	case google:
	default:
		return true
	}
	return false
}

// FindAuthentication @Title FindAuthentication
// @Description 查询身份验证信息
// @Accept  json
// @Tags  auth
// @Success 200 {object} []domain.FbAuthentication "查询成功"
// @Failure 400	"查询失败"
// @Router /api/v1/auth  [GET]
func (d *AuthenticationHandler) FindAuthentication(c echo.Context) (err error) {
	ctx := c.Request().Context()
	authentications, err := d.AuthenticationUseCase.FindAuthentication(ctx)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, err))
	}
	result := make([]domain.FbAuthenticationResp, 0)
	for _, row := range authentications {
		stateArr := make([]string, 0)
		err = json.Unmarshal([]byte(row.SwitchState), &stateArr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.JsonUnMarshalErr))
		}
		var authConfig domain.WdgAuthConfig
		err := json.Unmarshal([]byte(row.Config), &authConfig)
		if err != nil {
			log.Error(row.Name, "json.Unmarshal WdgAuthConfig fail ,err : ", err)
		}
		result = append(result, domain.FbAuthenticationResp{
			ID:           row.ID,
			Name:         row.Name,
			AuthSupplier: row.AuthSupplier,
			SwitchState:  stateArr,
			Config:       authConfig,
			CreateTime:   row.CreateTime,
			UpdateTime:   row.UpdateTime,
			IsDel:        row.IsDel,
		})
	}

	return c.JSON(http.StatusOK, SuccessWriteResult(result))
}

// Delete @Title Delete
// @Description 删除身份验证信息
// @Accept  json
// @Tags  auth
// @Param id formData integer true "id"
// @Success 200 "删除成功"
// @Failure 400	"删除失败"
// @Router /api/v1/auth/:id  [DELETE]
func (d *AuthenticationHandler) Delete(c echo.Context) (err error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.ParamErr))
	}

	if err = authenticationDeleteCheckParams(uint(id)); err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, err))
	}

	ctx := c.Request().Context()
	affect, err := d.AuthenticationUseCase.Delete(ctx, uint(id))
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, err))
	}

	return c.JSON(http.StatusOK, SuccessWriteResult(fmt.Sprintf("Affected rows : %d", affect)))
}

func authenticationDeleteCheckParams(id uint) (err error) {
	if ok := utils.Empty(id); ok {
		return domain.ParamIdEmptyErr
	}
	return
}

// Update @Title Update
// @Description 修改身份验证信息
// @Accept  json
// @Tags  auth
// @Param id formData integer true "id"
// @Param name formData string true "名称"
// @Param authSupplier formData integer true "验证供应商类型"
// @Param switchState formData string true "开关"
// @Param config formData string true "验证配置"
// @Success 200 "修改成功"
// @Failure 400	"修改失败"
// @Router /api/v1/auth  [PUT]
func (d *AuthenticationHandler) Update(c echo.Context) (err error) {
	var authenticationResp domain.FbAuthenticationResp
	err = c.Bind(&authenticationResp)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.ParamErr))
	}
	// 服务发现地址
	authenticationResp.Config.DiscoveryURL = fmt.Sprintf("%s/%s", authenticationResp.Config.Issuer, DiscoveryURL)
	authenticationResp.Config.UserInfoEndpoint = fmt.Sprintf("%s/%s", authenticationResp.Config.Issuer, UserInfoPoint)
	if authenticationResp.Config.Jwks == 0 {
		authenticationResp.Config.JwksURL = fmt.Sprintf("%s/%s", authenticationResp.Config.Issuer, JwksURL)
	}

	authConfigByte, err := json.Marshal(authenticationResp.Config)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.JsonUnMarshalErr))
	}
	switchStateJsonByte, err := json.Marshal(authenticationResp.SwitchState)
	authentication := domain.FbAuthentication{
		ID:           authenticationResp.ID,
		Name:         authenticationResp.Name,
		AuthSupplier: authenticationResp.AuthSupplier,
		SwitchState:  string(switchStateJsonByte),
		Config:       string(authConfigByte),
		CreateTime:   authenticationResp.CreateTime,
		UpdateTime:   authenticationResp.UpdateTime,
		IsDel:        authenticationResp.IsDel,
	}

	if err = authenticationUpdateCheckParams(authentication); err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, err))
	}
	ctx := c.Request().Context()
	resp, err := d.AuthenticationUseCase.Update(ctx, &authentication)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, err))
	}

	return c.JSON(http.StatusOK, SuccessWriteResult(fmt.Sprintf("Affected rows : %d", resp)))
}

func authenticationUpdateCheckParams(auth domain.FbAuthentication) (result error) {
	if ok := utils.Empty(auth.ID); ok {
		return domain.ParamIdEmptyErr
	}
	if ok := utils.Empty(auth.Name); ok {
		return domain.ParamNameEmptyErr
	}
	//if ok := checkAuthSupplier(auth.AuthSupplier); ok {
	//	return domain.ParamSupplierNotExistErr
	//}
	return
}
