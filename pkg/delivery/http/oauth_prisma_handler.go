package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fire_boom/domain"
	"github.com/fire_boom/pkg/usecase"
	"github.com/fire_boom/utils"
	"github.com/gocarina/gocsv"
	"github.com/labstack/echo/v4"
	"github.com/pborman/uuid"
	"github.com/prisma/prisma-client-go/engine"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var tableArr = []string{"oauth_user"}

type CSVOauthUser struct {
	ID                 string `db:"id" json:"id" csv:"id"`                                      // id
	Name               string `db:"name" json:"name" csv:"名称"`                                  // 姓名
	NickName           string `db:"nick_name" json:"nick_name" csv:"昵称"`                        // 昵称
	UserName           string `db:"user_name" json:"user_name" csv:"用户名称"`                      // 用户名
	EncryptionPassword string `db:"encryption_password" json:"encryption_password" csv:"加密后密码"` // 加密后密码
	Mobile             string `db:"mobile" json:"mobile" csv:"手机号"`                             // 手机号
	Email              string `db:"email" json:"email" csv:"邮箱"`                                // 邮箱
	LastLoginTime      string `db:"last_login_time" json:"-" csv:"最后一次登陆时间"`                    // 最后一次登陆时间
	Status             string `db:"status" json:"status" csv:"状态"`                              // 状态
	MateData           string `db:"mate_data" json:"mate_data" csv:"其他信息(json字符串保存)"`           // 其他信息(json字符串保存)
	CreateTime         string `db:"create_time" json:"create_time" csv:"创建时间"`                  // 创建时间
	UpdateTime         string `db:"update_time" json:"update_time" csv:"修改时间"`                  // 修改时间
}

type OauthUserLimitResp struct {
	CurrPage  int64                  `json:"currPage"`  // 当前页
	TotalPage int64                  `json:"totalPage"` // 总页数
	PageSize  int64                  `json:"pageSize"`  // 每页大小
	UserList  []domain.OauthUserResp `json:"userList"`  // 用户列表
}

type DbList struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type OauthUserDetailResp struct {
	ID            string `json:"id"`            // id
	Name          string `json:"name"`          // 姓名
	NickName      string `json:"nickName"`      // 昵称
	UserName      string `json:"userName"`      // 用户名
	Mobile        string `json:"mobile"`        // 手机号
	Email         string `json:"email"`         // 邮箱
	Gender        string `json:"gender"`        // 性别
	Birthday      string `json:"birthday"`      // 生日
	CountryCode   string `json:"countryCode"`   // 国家代码
	Address       string `json:"address"`       // 所在地
	Country       string `json:"country"`       // 国家
	City          string `json:"city"`          // 城市
	Province      string `json:"province"`      // 省/直辖市(手填)
	StreetAddress string `json:"streetAddress"` // 街道地址
	ExternalId    string `json:"externalId"`    // 原系统id
	PostalCode    string `json:"postalCode"`    // 邮政编码
	JsonString    string `json:"jsonString"`    // 原始json
}

type OauthUserImportData struct {
	ID                 string          `json:"id"`                  // id
	Name               string          `json:"name"`                // 姓名
	NickName           string          `json:"nick_name"`           // 昵称
	UserName           string          `json:"user_name"`           // 用户名
	EncryptionPassword string          `json:"encryption_password"` // 加密后密码
	Mobile             string          `json:"mobile"`              // 手机号
	Email              string          `json:"email"`               // 邮箱
	LastLoginTime      string          `json:"last_login_time"`     // 最后一次登陆时间
	Status             int64           `json:"status"`              // 状态
	MateData           domain.MateData `json:"mate_data"`           // 其他信息(json字符串保存)
	CreateTime         string          `json:"create_time"`         // 创建时间
	UpdateTime         string          `json:"update_time"`         // 修改时间
	IsDel              int64           `json:"is_del"`              // 是否删除
}

// OauthPrismaHandler 用户身份验证模块
type OauthPrismaHandler struct {
	DataSourceUseCase domain.DataSourceUseCase
}

func InitOauthPrismaUseCaseRouter(e *echo.Echo, dsr domain.DataSourceRepository) {
	timeoutContext := time.Duration(utils.GetTimeOut()) * time.Second
	au := usecase.NewDataSourceUseCase(dsr, timeoutContext)
	NewOauthPrismaHandler(e, au)
}

func NewOauthPrismaHandler(e *echo.Echo, dUseCase domain.DataSourceUseCase) {
	handler := &OauthPrismaHandler{
		DataSourceUseCase: dUseCase,
	}
	v1 := e.Group("/api/v1")
	{
		authentication := v1.Group("/oauth")
		{
			authentication.GET("", handler.FindUsers)                       // 查询列表
			authentication.GET("/:id", handler.GetById)                     // 查询详情
			authentication.POST("", handler.CreateUser)                     // 创建用户
			authentication.PUT("", handler.UpdateUser)                      // 修改用户
			authentication.PUT("/password", handler.UpdatePassword)         // 修改用户密码
			authentication.DELETE("", handler.DeleteUsers)                  // 批量删除用户
			authentication.PUT("/status", handler.UpdateStatus)             // 批量更改用户状态
			authentication.GET("/export", handler.ExportUsers)              // 导出用户信息
			authentication.GET("/exportTemplate", handler.ExportTemplate)   // 导出模版
			authentication.POST("/import", handler.ImportUsers)             // 导入
			authentication.GET("/forcedOffline/:id", handler.ForcedOffline) // TODO 强制下线 mock待实现
			authentication.GET("/role/:id", handler.GetUserRole)            // 获取用户角色
			authentication.POST("/role/:id", handler.UpdateRole)            // 修改用户角色 (新增/撤销)
			authentication.GET("/tables/:id", handler.FindTables)           // 查询当前表结构
			authentication.POST("/tables/:id/import", handler.TablesImport) // 导入表结构›
			authentication.GET("/dblist", handler.DbList)                   // 数据源列表
			authentication.POST("/changeDB/:id", handler.ChangeDB)          // 切换数据源
		}
	}
}

// ChangeDB @Title 切换数据源
// @Description 切换数据源
// @Accept  json
// @Tags  oauth
// @Param id formData integer true "数据库id"
// @Success 200 {object} TablesResp "查询成功"
// @Failure 400	"查询失败"
// @Router /api/v1/oauth/changeDB/:id  [POST]
func (d *OauthPrismaHandler) ChangeDB(c echo.Context) (err error) {
	// 切换数据源，需要重新组织生成schema.prisma、prismaDB文件
	// 获取数据源id
	ctx := c.Request().Context()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id == 0 {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.ParamErr))
	}
	schemaContent := ""
	// 如果是默认数据源
	if id == -1 {
		schemaContent = utils.GetDefaultDBSchema()
	} else {
		// 组织该数据源的schema文件
		content, err := d.DataSourceUseCase.GetPrismaSchema(ctx, uint(id))
		if err != nil {
			return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.FileReadErr))
		}
		schemaContent = content
	}

	//schemaPath := utils.GetPrismaSchemaFilePath()
	moduleSchema, _ := ioutil.ReadFile(utils.GetModuleSchemaPath())
	schemaContent = fmt.Sprintf("%s \n %s", schemaContent, moduleSchema)
	// 写入schema文件
	//err = utils.WriteFile(schemaPath, schemaContent)
	utils.SetIndexDBID(int64(id))

	// 重新加载
	_, err = engine.ReloadQueryEngineOnce(schemaContent)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.NewBizErr("-1", err.Error())))
	}
	// 将当前数据源id写到文件中
	utils.SetIndexDBID(int64(id))

	return c.JSON(http.StatusOK, SuccessResult())
}

// DbList @Title 数据源列表
// @Description 数据源列表
// @Accept  json
// @Tags  oauth
// @Param id formData integer true "数据库id"
// @Success 200 {object} TablesResp "查询成功"
// @Failure 400	"查询失败"
// @Router /api/v1/oauth/dblist  [POST]
func (d *OauthPrismaHandler) DbList(c echo.Context) (err error) {
	// 切换数据源，需要重新组织生成schema.prisma、prismaDB文件
	// 获取数据源id
	ctx := c.Request().Context()
	dataSources, err := d.DataSourceUseCase.FindDataSources(ctx)
	if err != nil {
		return c.JSON(http.StatusBadRequest, SuccessWriteResult(err))
	}
	result := struct {
		DBs       []DbList `json:"dbs"`
		IndexDBID int64    `json:"indexDBID"`
	}{}
	dbs := make([]DbList, 0)
	// 默认数据库
	dbs = append(dbs, DbList{
		ID:   -1,
		Name: "default",
	})
	for _, row := range dataSources {
		dbs = append(dbs, DbList{
			ID:   int64(row.ID),
			Name: row.Name,
		})
	}
	result.DBs = dbs
	result.IndexDBID = utils.GetIndexDBID()
	return c.JSON(http.StatusOK, SuccessWriteResult(result))
}

// FindUsers @Title 获取用户信息
// @Description 查询身份验证信息
// @Accept  json
// @Tags  oauth
// @Param pageSize formData integer true "每页大小"
// @Param currPage formData integer true "当前页"
// @Param totalPage formData integer true "总页数"
// @Param search formData string true "搜索框(用户名/邮箱/手机号)"
// @Success 200 {object} OauthUserLimitResp "查询成功"
// @Failure 400	"查询失败"
// @Router /api/v1/oauth  [GET]
func (d *OauthPrismaHandler) FindUsers(c echo.Context) (err error) {
	type RequestParam struct {
		CurrPage  int64  `json:"currPage"`  // 当前页
		TotalPage int64  `json:"totalPage"` // 总页数
		PageSize  int64  `json:"pageSize"`  // 每页大小
		Email     string `json:"email"`
		Mobile    string `json:"mobile"`
		UserName  string `json:"userName"`
	}
	var req RequestParam
	err = c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.ParamErr))
	}
	if req.CurrPage == 0 {
		req.CurrPage = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	var result OauthUserLimitResp
	// 获取查询schema
	param := map[string]interface{}{
		"skip":     req.PageSize * (req.CurrPage - 1), // 当前页第一行
		"take":     req.PageSize,
		"email":    req.Email,
		"userName": req.UserName,
		"mobile":   req.Mobile,
	}
	schema := utils.GetQuerySchema("user", "FindAllLimit", param)
	userInfos := make([]domain.OauthUser, 0)
	// 查询
	err = engine.QuerySchema(schema, &userInfos)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.DbFindErr))
	}

	totalPage := int64(len(userInfos)) / req.PageSize
	if int64(len(userInfos))%req.PageSize > 0 {
		totalPage++
	}
	userInfoResp := make([]domain.OauthUserResp, 0)
	for _, row := range userInfos {
		userInfoResp = append(userInfoResp, row.TransformResp())
	}
	result.UserList = userInfoResp
	result.TotalPage = int64(len(userInfos))
	result.PageSize = req.PageSize
	result.CurrPage = req.CurrPage

	return c.JSON(http.StatusOK, SuccessWriteResult(result))
}

// GetById @Title 获取用户详细信息
// @Description 获取用户详细信息
// @Accept  json
// @Tags  oauth
// @Success 200 {object} OauthUserDetailResp "查询成功"
// @Failure 400	"查询失败"
// @Router /api/v1/oauth/:id  [GET]
func (d *OauthPrismaHandler) GetById(c echo.Context) (err error) {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.ParamErr))
	}

	// 获取查询schema
	param := map[string]interface{}{
		"id": id,
	}
	schema := utils.GetQuerySchema("user", "GetByID", param)
	userInfo := domain.OauthUser{}
	// 查询
	err = engine.QuerySchema(schema, &userInfo)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.DbFindErr))
	}
	var mataData domain.MateData
	err = json.Unmarshal([]byte(userInfo.MateData), &mataData)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.JsonUnMarshalErr))
	}
	jsonStringByte, _ := json.Marshal(userInfo)
	result := OauthUserDetailResp{
		ID:            userInfo.ID,
		Name:          userInfo.Name,
		NickName:      userInfo.NickName,
		UserName:      userInfo.UserName,
		Mobile:        userInfo.Mobile,
		Email:         userInfo.Email,
		Gender:        mataData.Gender,
		Address:       mataData.Address,
		Birthday:      mataData.Birthdate,
		CountryCode:   mataData.CountryCode,
		Country:       mataData.Country,
		City:          mataData.City,
		Province:      mataData.Province,
		StreetAddress: mataData.StreetAddress,
		ExternalId:    mataData.ExternalId,
		PostalCode:    mataData.PostalCode,
		JsonString:    string(jsonStringByte),
	}
	return c.JSON(http.StatusOK, SuccessWriteResult(result))
}

// CreateUser @Title 创建用户
// @Description 创建用户
// @Accept  json
// @Tags  oauth
// @Param createType formData integer true "创建账号类型0-用户名密码 1-手机号 2-邮箱"
// @Param userName formData string true "用户名"
// @Param mobile formData string true "手机号"
// @Param email formData string true "邮箱"
// @Param password formData string true "密码"
// @Param rePassword formData string true "确认密码"
// @Param sendAddress formData boolean true "发送首次登陆地址"
// @Success 200 "创建成功"
// @Failure 400	"创建失败"
// @Router /api/v1/oauth  [POST]
func (d *OauthPrismaHandler) CreateUser(c echo.Context) (err error) {
	type RequestParam struct {
		CreateType  int64  `json:"createType"`  // 创建账号类型0-用户名密码 1-手机号 2-邮箱
		UserName    string `json:"userName"`    // 用户名
		Mobile      string `json:"mobile"`      // 手机号
		Email       string `json:"email"`       // 邮箱
		Password    string `json:"password"`    // 密码
		RePassword  string `json:"rePassword"`  // 重复密码
		SendAddress bool   `json:"sendAddress"` // 首次登陆发送地址
	}
	var req RequestParam
	err = c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.ParamErr))
	}
	param := make(map[string]interface{}, 0)
	param["id"] = uuid.New()
	param["userName"] = ""
	param["encryPswd"] = ""
	param["mobile"] = ""
	param["email"] = ""
	param["createTime"] = time.Now().Format(utils.DateTimeFormat)
	// 校验密码
	if req.CreateType == 0 {
		if req.Password != req.RePassword {
			return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.ParamRePasswordErr))
		}
		// 生成加密密码
		pswd, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 0)
		param["userName"] = req.UserName
		param["encryPswd"] = string(pswd)
	}
	// TODO 校验手机号验证码
	if req.CreateType == 1 {
		param["mobile"] = req.Mobile
	}
	// TODO 校验邮箱验证码
	if req.CreateType == 2 {
		param["email"] = req.Email
	}
	schema := utils.GetQuerySchema("user", "CreateOneUser", param)
	userInfos := domain.OauthUser{}
	// 查询
	err = engine.QuerySchema(schema, &userInfos)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.DbFindErr))
	}

	return c.JSON(http.StatusOK, SuccessResult())
}

// UpdateUser @Title 修改用户
// @Description 修改用户
// @Accept  json
// @Tags  oauth
// @Param name formData string true "姓名"
// @Param userName formData string true "用户名"
// @Param nickName formData string true "昵称"
// @Param gender formData string true "性别"
// @Param birthday formData string true "生日"
// @Param mobile formData string true "手机号"
// @Param email formData string true "邮箱"
// @Param password formData string true "密码"
// @Param countryCode    formData string true "国家代码"
// @Param address        formData string true "所在地"
// @Param company        formData string true "公司"
// @Param city           formData string true "城市"
// @Param province       formData string true "省/直辖市(手填)"
// @Param streetAddress  formData string true "街道地址"
// @Param externalId     formData string true "原系统id"
// @Param postalCode     formData string true "邮政编码"
// @Success 200 "修改成功"
// @Failure 400	"修改失败"
// @Router /api/v1/oauth/:id  [PUT]
func (d *OauthPrismaHandler) UpdateUser(c echo.Context) (err error) {
	type RequestParam struct {
		ID            string `json:"id"`            // id
		Name          string `json:"name"`          // 姓名
		NickName      string `json:"nickName"`      // 昵称
		UserName      string `json:"userName"`      // 用户名
		Mobile        string `json:"mobile"`        // 手机号
		Email         string `json:"email"`         // 邮箱
		Gender        string `json:"gender"`        // 性别
		Birthday      string `json:"birthday"`      // 生日
		CountryCode   string `json:"countryCode"`   // 国家代码
		Address       string `json:"address"`       // 所在地
		Country       string `json:"country"`       // 国家
		City          string `json:"city"`          // 城市
		Province      string `json:"province"`      // 省/直辖市(手填)
		StreetAddress string `json:"streetAddress"` // 街道地址
		ExternalId    string `json:"externalId"`    // 原系统id
		PostalCode    string `json:"postalCode"`    // 邮政编码
	}
	var req RequestParam
	err = c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.ParamErr))
	}
	mateData := domain.MateData{
		ProfileClaims: domain.ProfileClaims{
			Name:      req.Name,
			Gender:    req.Gender,
			Birthdate: req.Birthday,
			Nickname:  req.NickName,
		},
		PhoneClaims: domain.PhoneClaims{
			PhoneNumber: req.Mobile,
		},
		EmailClaims: domain.EmailClaims{
			Email: req.Email,
		},
		AddressClaims: domain.AddressClaims{
			Address: req.Address,
		},
		CountryCode:   req.CountryCode,
		Country:       req.Country,
		City:          req.City,
		Province:      req.Province,
		StreetAddress: req.StreetAddress,
		ExternalId:    req.ExternalId,
		PostalCode:    req.PostalCode,
	}
	mateDataBytes, _ := json.Marshal(mateData)
	param := make(map[string]interface{}, 0)
	param["id"] = req.ID
	param["name"] = req.Name
	param["nickName"] = req.NickName
	param["userName"] = req.UserName
	param["email"] = req.Email
	param["mobile"] = req.Mobile
	param["userName"] = req.UserName
	param["updateTime"] = time.Now().Format(utils.DateTimeFormat)
	param["mateData"] = strings.ReplaceAll(string(mateDataBytes), `"`, `\"`)

	schema := utils.GetQuerySchema("user", "UpdateByID", param)
	userInfo := domain.OauthUser{}
	// 查询
	err = engine.QuerySchema(schema, &userInfo)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.DbUpdateErr))
	}

	return c.JSON(http.StatusOK, SuccessResult())
}

// UpdatePassword @Title 修改用户密码
// @Description 修改用户密码
// @Accept  json
// @Tags  oauth
// @Param id formData string true "id"
// @Param password formData string true "密码"
// @Success 200 "修改成功"
// @Failure 400	"修改失败"
// @Router /api/v1/oauth/password  [PUT]
func (d *OauthPrismaHandler) UpdatePassword(c echo.Context) (err error) {
	type RequestParam struct {
		ID       string `json:"id"`       // id
		Password string `json:"password"` // 邮政编码
	}
	var req RequestParam
	err = c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.ParamErr))
	}
	pswd, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 0)
	param := make(map[string]interface{}, 0)
	param["id"] = req.ID
	param["encryPassword"] = string(pswd)
	param["updateTime"] = time.Now().Format(utils.DateTimeFormat)

	schema := utils.GetQuerySchema("user", "UpdateUserPassword", param)
	userInfo := domain.OauthUser{}
	// 查询
	err = engine.QuerySchema(schema, &userInfo)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.DbUpdateErr))
	}

	return c.JSON(http.StatusOK, SuccessResult())
}

// DeleteUsers @Title 删除用户
// @Description 删除用户
// @Accept  json
// @Tags  oauth
// @Param ids formData []string true "id列表"
// @Success 200 "删除成功"
// @Failure 400	"删除失败"
// @Router /api/v1/oauth  [DELETE]
func (d *OauthPrismaHandler) DeleteUsers(c echo.Context) (err error) {
	req := struct {
		Ids []string `json:"ids"`
	}{}
	err = c.Bind(&req)
	if err != nil || len(req.Ids) == 0 {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.ParamErr))
	}
	param := make(map[string]interface{}, 0)
	param["ids"] = req.Ids
	param["updateTime"] = time.Now().Format(utils.DateTimeFormat)
	schema := utils.GetQuerySchema("user", "DeleteUsers", param)
	count := struct {
		Count int64 `json:"count"`
	}{}
	// 查询
	err = engine.QuerySchema(schema, &count)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.DbUpdateErr))
	}

	return c.JSON(http.StatusOK, SuccessResult())
}

// UpdateStatus @Title 修改状态
// @Description 修改状态
// @Accept  json
// @Tags  oauth
// @Param ids formData []string true "id列表"
// @Param status formData boolean true "false锁定 true 解锁"
// @Success 200 "修改成功"
// @Failure 400	"修改失败"
// @Router /api/v1/oauth/status  [PUT]
func (d *OauthPrismaHandler) UpdateStatus(c echo.Context) (err error) {
	req := struct {
		Ids    []string `json:"ids"`
		Status bool     `json:"status"`
	}{}
	err = c.Bind(&req)
	if err != nil || len(req.Ids) == 0 {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.ParamErr))
	}
	status := 0
	if req.Status == false {
		status = 1
	}
	param := make(map[string]interface{}, 0)

	param["ids"] = req.Ids
	param["status"] = status
	param["updateTime"] = time.Now().Format(utils.DateTimeFormat)

	schema := utils.GetQuerySchema("user", "UpdateUserStatus", param)
	count := struct {
		Count int64 `json:"count"`
	}{}
	// 查询
	err = engine.QuerySchema(schema, &count)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.DbUpdateErr))
	}

	return c.JSON(http.StatusOK, SuccessResult())
}

// ExportUsers @Title 导出
// @Description 导出
// @Accept  json
// @Tags  oauth
// @Param ids formData []string true "id列表"
// @Success 200 "修改成功"
// @Failure 400	"修改失败"
// @Router /api/v1/oauth/export  [GET]
func (d *OauthPrismaHandler) ExportUsers(c echo.Context) (err error) {
	req := struct {
		Ids []string `json:"ids"`
	}{}
	err = c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.ParamErr))
	}
	param := map[string]interface{}{"ids": req.Ids}

	// 获取当前db的schema
	// 获取查询schema
	schema := ""
	// 导出全部
	if len(req.Ids) == 0 {
		schema = utils.GetQuerySchema("user", "FindAll", param)
	} else {
		schema = utils.GetQuerySchema("user", "FindByIds", param)
	}
	userInfos := make([]domain.OauthUser, 0)
	// 查询
	err = engine.QuerySchema(schema, &userInfos)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.DbFindErr))
	}

	csvArr := make([]CSVOauthUser, 0)
	for i := 0; i < len(userInfos); i++ {
		row := userInfos[i]
		status := "正常"
		if row.Status == 1 {
			status = "锁定"
		}
		var lastLoginTime, createTime, updateTime time.Time
		if row.LastLoginTime != "" {
			lastLoginTime, _ = time.ParseInLocation(utils.DateTimeFormat, row.LastLoginTime, time.Local)
		}
		if row.CreateTime != "" {
			createTime, _ = time.ParseInLocation(utils.DateTimeFormat, row.CreateTime, time.Local)
		}
		if row.UpdateTime != "" {
			updateTime, _ = time.ParseInLocation(utils.DateTimeFormat, row.UpdateTime, time.Local)
		}
		csvArr = append(csvArr, CSVOauthUser{
			ID:                 row.ID,
			Name:               row.Name,
			NickName:           row.NickName,
			UserName:           row.UserName,
			EncryptionPassword: row.EncryptionPassword,
			Mobile:             row.Mobile,
			Email:              row.Email,
			LastLoginTime:      lastLoginTime.Format(utils.DateTimeFormat),
			Status:             status,
			MateData:           row.MateData,
			CreateTime:         createTime.Format(utils.DateTimeFormat),
			UpdateTime:         updateTime.Format(utils.DateTimeFormat),
		})

	}
	csvContent, err := gocsv.MarshalString(&csvArr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.FileExportErr))
	}
	return c.Stream(http.StatusOK, "text/csv; charset=UTF-8", bytes.NewReader([]byte(csvContent)))
}

// ExportTemplate @Title 导出模版
// @Description 导出
// @Accept  json
// @Tags  oauth
// @Success 200 "修改成功"
// @Failure 400	"修改失败"
// @Router /api/v1/oauth/export  [POST]
func (d *OauthPrismaHandler) ExportTemplate(c echo.Context) (err error) {
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.FileExportErr))
	}
	csvArr := make([]CSVOauthUser, 0)
	csvContent, err := gocsv.MarshalString(&csvArr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.FileExportErr))
	}
	return c.Stream(http.StatusOK, "text/csv; charset=UTF-8", bytes.NewReader([]byte(csvContent)))
}

// ImportUsers @Title 导入
// @Description 导入
// @Accept  json
// @Tags  oauth
// @Success 200 "修改成功"
// @Failure 400	"修改失败"
// @Router /api/v1/oauth/import  [POST]
func (d *OauthPrismaHandler) ImportUsers(c echo.Context) (err error) {
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.FileCode, domain.ParamErr))
	}
	path := "tempCSV.csv"
	defer os.Remove(path)
	err = uploadWriteFile(file, path)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.FileCode, domain.FileImportErr))
	}

	clientsFile, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.FileCode, domain.FileImportErr))
	}
	defer clientsFile.Close()
	clients := []*CSVOauthUser{}
	if err := gocsv.UnmarshalFile(clientsFile, &clients); err != nil { // Load clients from file
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.FileCode, domain.FileImportErr))
	}

	for i := 0; i < len(clients); i++ {
		row := clients[i]
		status := 0
		if row.Status != "正常" {
			status = 1
		}
		data := domain.OauthUser{
			ID:                 uuid.New(),
			Name:               row.Name,
			NickName:           row.NickName,
			UserName:           row.UserName,
			EncryptionPassword: row.EncryptionPassword,
			Mobile:             row.Mobile,
			Email:              row.Email,
			Status:             int64(status),
			MateData:           row.MateData,
			CreateTime:         time.Now().Format(utils.DateTimeFormat),
		}
		schemaBytes, _ := json.Marshal(data)
		// 处理双引号
		reg := regexp.MustCompile("\"([a-zA-Z]\\w*)\":")
		querySchema := reg.ReplaceAllString(string(schemaBytes), `$1:`)
		param := map[string]interface{}{"data": querySchema}
		schema := utils.GetQuerySchema("user", "CreateOneUserAllField", param)
		// 处理由于序列化参数最外层多出来的双引号
		schema = strings.ReplaceAll(schema, `{"{id`, `{id`)
		schema = strings.ReplaceAll(schema, `}"}){id`, `}){id`)
		result := domain.OauthUser{}
		// 查询
		err = engine.QuerySchema(schema, &result)
		if err != nil {
			return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.DbUpdateErr))
		}
	}

	return c.JSON(http.StatusOK, SuccessResult())
}

// ForcedOffline @Title 强制下线
// @Description 强制下线
// @Accept  json
// @Tags  oauth
// @Success 200 "成功"
// @Failure 400	"失败"
// @Router /api/v1/oauth/ForcedOffline/:id  [GET]
func (d *OauthPrismaHandler) ForcedOffline(c echo.Context) (err error) {
	return c.JSON(http.StatusOK, SuccessResult())
}

// GetUserRole @Title 获取用户角色
// @Description 获取用户角色
// @Accept  json
// @Tags  oauth
// @Success 200 {object} domain.FbRole "查询成功"
// @Failure 400	"查询失败"
// @Router /api/v1/oauth/role/:id  [GET]
func (d *OauthPrismaHandler) GetUserRole(c echo.Context) (err error) {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.ParamErr))
	}

	// 获取查询schema
	param := map[string]interface{}{
		"id": id,
	}
	schema := utils.GetQuerySchema("user", "GetByID", param)
	userInfo := domain.OauthUser{}
	// 查询
	err = engine.QuerySchema(schema, &userInfo)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.DbFindErr))
	}
	var mataData domain.MateData
	err = json.Unmarshal([]byte(userInfo.MateData), &mataData)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.JsonUnMarshalErr))
	}
	if mataData.Role == nil {
		mataData.Role = make([]domain.FbRole, 0)
	}
	return c.JSON(http.StatusOK, SuccessWriteResult(mataData.Role))
}

// UpdateRole @Title 修改用户角色
// @Description 修改用户角色
// @Accept  json
// @Tags  oauth
// @Param roleId formData integer true "角色id"
// @Param role formData []string true "权限数组"
// @Success 200 "添加成功"
// @Failure 400	"添加失败"
// @Router /api/v1/oauth/role/:id  [POST]
func (d *OauthPrismaHandler) UpdateRole(c echo.Context) (err error) {
	req := struct {
		Role []domain.FbRole `json:"role"`
	}{}
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.ParamErr))
	}

	// 获取查询schema
	param := map[string]interface{}{
		"id": id,
	}
	schema := utils.GetQuerySchema("user", "GetByID", param)
	userInfo := domain.OauthUser{}

	// 查询
	err = engine.QuerySchema(schema, &userInfo)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.DbFindErr))
	}
	var mataData domain.MateData
	err = json.Unmarshal([]byte(userInfo.MateData), &mataData)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.JsonUnMarshalErr))
	}
	mataData.Role = req.Role
	mataDataBytes, _ := json.Marshal(mataData)

	// 查询
	err = engine.QuerySchema(schema, &userInfo)

	updateParam := make(map[string]interface{}, 0)
	updateParam["id"] = userInfo.ID
	updateParam["name"] = userInfo.Name
	updateParam["nickName"] = userInfo.NickName
	updateParam["userName"] = userInfo.UserName
	updateParam["email"] = userInfo.Email
	updateParam["mobile"] = userInfo.Mobile
	updateParam["userName"] = userInfo.UserName
	updateParam["updateTime"] = time.Now().Format(utils.DateTimeFormat)
	updateParam["mateData"] = strings.ReplaceAll(string(mataDataBytes), `"`, `\"`)

	updateSchema := utils.GetQuerySchema("user", "UpdateByID", updateParam)
	// 查询
	err = engine.QuerySchema(updateSchema, &userInfo)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.DbUpdateErr))
	}

	return c.JSON(http.StatusOK, SuccessResult())
}

type TablesResp struct {
	Exist    []string `json:"exist"`
	NotExist []string `json:"notExist"`
}

// FindTables @Title 查询表列表
// @Description 查询表列表
// @Accept  json
// @Tags  oauth
// @Param id formData integer true "数据库id"
// @Success 200 {object} TablesResp "查询成功"
// @Failure 400	"查询失败"
// @Router /api/v1/oauth/tables/:id  [GET]
func (d *OauthPrismaHandler) FindTables(c echo.Context) (err error) {
	result := TablesResp{
		Exist:    make([]string, 0),
		NotExist: make([]string, 0),
	}
	content, _ := ioutil.ReadFile(utils.GetPrismaSchemaFilePath())
	dmmf, err := engine.QueryDMMF(string(content))
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.DbFindErr))
	}
	tableMap := make(map[string]string, len(tableArr))
	for _, row := range tableArr {
		tableMap[row] = row
	}
	for _, model := range dmmf.Datamodel.Models {
		tableName := model.Name.String()
		if _, ok := tableMap[tableName]; !ok {
			continue
		}
		result.Exist = append(result.Exist, tableName)
		delete(tableMap, tableName)
	}

	for k, _ := range tableMap {
		result.NotExist = append(result.NotExist, k)
	}
	return c.JSON(http.StatusOK, SuccessWriteResult(result))
}

// TablesImport @Title 导入
// @Description 导入
// @Accept  json
// @Tags  oauth
// @Param id formData integer true "数据库id"
// @Success 200 {object} TablesResp "查询成功"
// @Failure 400	"查询失败"
// @Router /api/v1/oauth/tables/:id/import  [POST]
func (d *OauthPrismaHandler) TablesImport(c echo.Context) (err error) {
	// 拼接schema语句
	ctx := c.Request().Context()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id == 0 {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.ParamErr))
	}
	schemaContent := ""
	// 如果是默认数据源
	// 组织该数据源的schema文件
	content, err := d.DataSourceUseCase.GetPrismaSchema(ctx, uint(id))
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.AuthCode, domain.FileReadErr))
	}
	schemaContent = content
	//schemaPath := utils.GetPrismaSchemaFilePath()
	moduleSchema, _ := ioutil.ReadFile(utils.GetModuleSchemaPath())
	schemaContent = fmt.Sprintf("%s \n %s", schemaContent, moduleSchema)
	// 写入schema文件
	//err = utils.WriteFile(schemaPath, schemaContent)

	// 合并完文件后再执行迁移
	err = engine.Push(schemaContent)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.DraftCode, domain.DbMigrateErr))
	}
	return c.JSON(http.StatusOK, SuccessResult())
}

// TablesExport @Title 导出
// @Description 导出
// @Accept  json
// @Tags  oauth
// @Param id formData integer true "数据库id"
// @Success 200 {object} TablesResp "查询成功"
// @Failure 400	"查询失败"
// @Router /api/v1/oauth/tables/:id/export  [POST]
func (d *OauthPrismaHandler) TablesExport(c echo.Context) (err error) {
	return c.JSON(http.StatusOK, SuccessResult())
}

func InitPrismaEngine() {

}
