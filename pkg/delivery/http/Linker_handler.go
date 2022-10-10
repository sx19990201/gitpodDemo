package http

import (
	"encoding/json"
	"github.com/fire_boom/domain"
	"github.com/fire_boom/utils"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
	"time"
)

type LinkerHandler struct {
}

func InitLinkerHandler(e *echo.Echo) {
	//timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second
	NewLinkerHandler(e)
}

func NewLinkerHandler(e *echo.Echo) {
	handler := &LinkerHandler{}
	v1 := e.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			home := auth.Group("/home")
			{
				home.GET("", handler.GetHome)             // 获取首页数据
				home.GET("/active", handler.GetDayActive) // 获取列表数据
			}
			linker := auth.Group("/linker")
			{
				linker.GET("", handler.FindLinkers)                // 获取连接器信息
				linker.GET("/default", handler.FindDefaultLinkers) // 更新连接器信息
				linker.POST("", handler.Update)                    // 更新连接器信息
			}
			brand := auth.Group("/brand")
			{
				brand.GET("", handler.GetBrandConfig)     // 获取品牌信息
				brand.POST("", handler.UpdateBrandConfig) // 修改品牌信息
			}
			loginConfig := auth.Group("/loginConfig")
			{
				loginConfig.GET("", handler.GetLoginConfig)     // 获取登陆配置
				loginConfig.POST("", handler.UpdateLoginConfig) // 修改登陆配置
			}
			otherConfig := auth.Group("/otherConfig")
			{
				otherConfig.GET("", handler.GetOtherConfig)     // 获取其他配置
				otherConfig.POST("", handler.UpdateOtherConfig) // 修改其他配置
			}
		}
	}
}

// GetHome @Title GetHome
// @Description 查询首页数据
// @Accept  json
// @Tags  auth/home
// @Success 200 {object} []domain.AuthHome "查询成功"
// @Failure 400	"查询失败"
// @Router /api/v1/auth/home  [GET]
func (l *LinkerHandler) GetHome(c echo.Context) error {
	//// 获取当前db的schema
	//content, _ := ioutil.ReadFile(utils.GetPrismaSchemaFilePath())
	//dbSchema := string(content)
	//// 获取查询schema
	//schema := ""
	//schema = utils.GetQuerySchema("user", "FindAll", nil)
	//userInfos := make([]domain.OauthUser, 0)
	//// 查询
	//err := engine.QuerySchema(dbSchema, schema, &userInfos)
	//if err != nil {
	//	return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.DbFindErr))
	//}
	//
	//local, _ := time.LoadLocation("Asia/Shanghai")
	//TodayUser := domain.UserData{}
	//toDayDate := time.Now().Format(utils.DateFormat)
	//
	//SevenDayUser := domain.UserData{}
	//sevenDayDate := time.Now().AddDate(0, 0, -7)
	//toDayBeforeCount := 0
	//for _, row := range userInfos {
	//	// 该用户注册的时间
	//	userCreateDate, err := time.ParseInLocation(utils.DateTimeFormat, row.CreateTime, local)
	//	if err != nil {
	//		log.Error("/api/v1/auth/home time.ParseInLocation err : ", err)
	//	}
	//	// 前一天的日期
	//	toDayBefore := time.Now().AddDate(0, 0, -1).Format(utils.DateFormat)
	//	// 当天增加人数
	//	if userCreateDate.Format(utils.DateFormat) == toDayDate {
	//		TodayUser.Count++
	//	}
	//	// 前一天的人数
	//	if userCreateDate.Format(utils.DateFormat) == toDayBefore {
	//		toDayBeforeCount++
	//	}
	//
	//	// 最近7天
	//
	//}
	//// 计算增长人数
	//TodayUser.UpCount = TodayUser.Count - int64(toDayBeforeCount)
	//// 负数置0
	//if TodayUser.UpCount <= 0 {
	//	TodayUser.UpCount = 0
	//}
	//SevenDayUser := domain.UserData{}

	result := domain.AuthHome{
		TotalUser: 1,
		TodayUser: domain.UserData{
			Count:   1,
			UpCount: 1,
		},
		SevenDayUser: domain.UserData{
			Count:   1,
			UpCount: 1,
		},
	}
	return c.JSON(http.StatusOK, SuccessWriteResult(result))
}

// GetDayActive @Title GetDayActive
// @Description 查询活跃信息
// @Accept  json
// @Tags  auth/home
// @Param endDate formData string true "结束日期"
// @Success 200 {object} domain.Active "查询成功"
// @Failure 400	"查询失败"
// @Router /api/v1/auth/home/active  [GET]
func (l *LinkerHandler) GetDayActive(c echo.Context) error {
	dataString := struct {
		EndDate string `json:"endDate"`
	}{}
	err := c.Bind(&dataString)
	if err != nil || dataString.EndDate == "" {
		dataString.EndDate = time.Now().Format(utils.DateFormat)
	}
	result := domain.Active{
		DayActive: domain.UserData{
			Count:   1,
			UpCount: 1,
		},
		WeekActive: domain.UserData{
			Count:   1,
			UpCount: 1,
		},
		MonthActive: domain.UserData{
			Count:   1,
			UpCount: 1,
		},
	}
	startDate, _ := time.Parse(utils.DateFormat, dataString.EndDate)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.ParamDateErr))
	}
	//endDate := startDate.AddDate(0, 0, -30)
	dayActive := make([]domain.DayActive, 0)
	for i := 0; i < 30; i++ {
		tempDate := startDate.AddDate(0, 0, -i).Format(utils.DateFormat)
		day := domain.DayActive{
			DataString: tempDate,
			Count:      30 - int64(i),
		}
		dayActive = append(dayActive, day)
	}
	result.ThirtyActive = dayActive
	return c.JSON(http.StatusOK, SuccessWriteResult(result))
}

// FindLinkers @Title FindLinkers
// @Description 查询链接器
// @Accept  json
// @Tags  auth/linker
// @Success 200 {object} []domain.Linker "查询成功"
// @Failure 400	"查询失败"
// @Router /api/v1/auth/linker  [GET]
func (l *LinkerHandler) FindLinkers(c echo.Context) error {
	linkerPath := utils.GetOauthLinkerPath()
	linkerBytes, err := ioutil.ReadFile(linkerPath)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.FileReadErr))
	}
	linker := make([]domain.Linker, 0)
	err = json.Unmarshal(linkerBytes, &linker)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.JsonUnMarshalErr))
	}

	return c.JSON(http.StatusOK, SuccessWriteResult(linker))
}

// FindDefaultLinkers @Title FindDefaultLinkers
// @Description 查询链接器
// @Accept  json
// @Tags  auth/linker
// @Param id formData string true "id"
// @Success 200 {object} []domain.Linker "查询成功"
// @Failure 400	"查询失败"
// @Router /api/v1/auth/linker  [GET]
func (l *LinkerHandler) FindDefaultLinkers(c echo.Context) error {
	req := struct {
		ID string `json:"id"`
	}{}
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.ParamErr))
	}
	if req.ID == "" {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.ParamIdEmptyErr))
	}
	linkerPath := utils.GetOauthDefaultLinkerPath()
	linkerBytes, err := ioutil.ReadFile(linkerPath)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.FileReadErr))
	}
	linker := make([]domain.Linker, 0)
	err = json.Unmarshal(linkerBytes, &linker)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.JsonUnMarshalErr))
	}
	result := domain.Linker{}
	for _, row := range linker {
		if row.ID == req.ID {
			result = row
		}
	}
	return c.JSON(http.StatusOK, SuccessWriteResult(result))
}

// Update @Title Update
// @Description 修改链接器
// @Accept  json
// @Tags  auth/linker
// @Param id formData string true "id"
// @Param enable formData boolean true "开关，true为创建链接器，false为删除链接器"
// @Param config formData string false "配置，创建和修改配置需要传入该字段"
// @Success 200 {object} domain.Linker "修改成功"
// @Failure 400	"修改失败"
// @Router /api/v1/auth/linker  [POST]
func (l *LinkerHandler) Update(c echo.Context) error {
	req := domain.Linker{}
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.ParamErr))
	}
	if req.ID == "" {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.ParamIdEmptyErr))
	}
	linkerPath := utils.GetOauthLinkerPath()
	linkerBytes, err := ioutil.ReadFile(linkerPath)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.FileReadErr))
	}
	linkers := make([]domain.Linker, 0)
	err = json.Unmarshal(linkerBytes, &linkers)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.JsonUnMarshalErr))
	}

	for i := 0; i < len(linkers); i++ {
		if linkers[i].ID == req.ID {
			linkers[i].Enabled = req.Enabled
			linkers[i].Config = req.Config
		}
	}
	content, err := json.Marshal(linkers)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.JsonMarshalErr))
	}
	err = utils.WriteFile(linkerPath, string(content))
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.FileWriteErr))
	}
	return c.JSON(http.StatusOK, SuccessWriteResult(linkers))
}

// GetBrandConfig @Title GetBrandConfig
// @Description 查询品牌配置
// @Accept  json
// @Tags  auth/brand
// @Success 200 {object} domain.LoginBrandConfig "查询成功"
// @Failure 400	"查询失败"
// @Router /api/v1/auth/brand  [GET]
func (l *LinkerHandler) GetBrandConfig(c echo.Context) error {
	brandConfigBytes, err := ioutil.ReadFile(utils.GetOauthLoginBrandPath())
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.FileReadErr))
	}

	result := domain.LoginBrandConfig{}
	err = json.Unmarshal(brandConfigBytes, &result)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.JsonUnMarshalErr))
	}

	return c.JSON(http.StatusOK, SuccessWriteResult(result))
}

// UpdateBrandConfig @Title UpdateBrandConfig
// @Description 修改品牌配置链接器
// @Accept  json
// @Tags  auth/brand
// @Param color formData domain.Color true "颜色配置"
// @Param branding formData domain.Branding true "品牌配置"
// @Success 200 {object} domain.LoginBrandConfig "修改成功"
// @Failure 400	"修改失败"
// @Router /api/v1/auth/brand  [POST]
func (l *LinkerHandler) UpdateBrandConfig(c echo.Context) error {
	req := domain.LoginBrandConfig{}
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.ParamErr))
	}
	content, err := json.Marshal(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.JsonMarshalErr))
	}
	err = utils.WriteFile(utils.GetOauthLoginBrandPath(), string(content))
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.FileWriteErr))
	}

	return c.JSON(http.StatusOK, SuccessWriteResult(req))
}

// GetLoginConfig @Title GetLoginConfig
// @Description 查询登录配置
// @Accept  json
// @Tags  auth/loginConfig
// @Success 200 {object} domain.LoginConfig "查询成功"
// @Failure 400	"查询失败"
// @Router /api/v1/auth/loginConfig  [GET]
func (l *LinkerHandler) GetLoginConfig(c echo.Context) error {
	configBytes, err := ioutil.ReadFile(utils.GetOauthLoginConfigPath())
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.FileReadErr))
	}
	result := domain.LoginConfig{}
	err = json.Unmarshal(configBytes, &result)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.JsonUnMarshalErr))
	}
	return c.JSON(http.StatusOK, SuccessWriteResult(result))
}

// UpdateLoginConfig @Title UpdateLoginConfig
// @Description 修改登录配置
// @Accept  json
// @Tags  auth/loginConfig
// @Param signInMethods formData domain.LoginConfig true "登录方式配置"
// @Param socialSignInConnectorTargets formData string true "已开启的登录方式"
// @Success 200 {object} domain.LoginConfig "修改成功"
// @Failure 400	"修改失败"
// @Router /api/v1/auth/loginConfig  [POST]
func (l *LinkerHandler) UpdateLoginConfig(c echo.Context) error {
	req := domain.LoginConfig{}
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.ParamErr))
	}

	content, err := json.Marshal(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.FileReadErr))
	}

	err = utils.WriteFile(utils.GetOauthLoginConfigPath(), string(content))
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.FileWriteErr))
	}
	return c.JSON(http.StatusOK, SuccessWriteResult(req))
}

// GetOtherConfig @Title GetOtherConfig
// @Description 查询其他配置
// @Accept  json
// @Tags  auth/otherConfig
// @Success 200 {object} domain.LoginOtherConfig "查询成功"
// @Failure 400	"查询失败"
// @Router /api/v1/auth/otherConfig  [GET]
func (l *LinkerHandler) GetOtherConfig(c echo.Context) error {
	content, err := ioutil.ReadFile(utils.GetOauthConfigPath())
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.FileReadErr))
	}
	result := domain.LoginOtherConfig{}
	err = json.Unmarshal(content, &result)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.JsonUnMarshalErr))
	}
	return c.JSON(http.StatusOK, SuccessWriteResult(result))
}

// UpdateOtherConfig @Title UpdateOtherConfig
// @Description 查询其他配置
// @Accept  json
// @Tags  auth/otherConfig
// @Param Enabled    formData boolean true "开关 true开 false关"
// @Param ContentUrl formData string true "条款链接"
// @Success 200 {object} domain.LoginOtherConfig "查询成功"
// @Failure 400	"查询失败"
// @Router /api/v1/auth/otherConfig  [POST]
func (l *LinkerHandler) UpdateOtherConfig(c echo.Context) error {
	req := domain.LoginOtherConfig{}
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.ParamErr))
	}

	content, err := json.Marshal(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.JsonMarshalErr))
	}

	err = utils.WriteFile(utils.GetOauthConfigPath(), string(content))
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.FileWriteErr))
	}

	return c.JSON(http.StatusOK, SuccessWriteResult(req))
}
