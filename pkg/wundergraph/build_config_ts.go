package wundergraph

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/fire_boom/domain"
	"github.com/fire_boom/utils"
	log "github.com/sirupsen/logrus"
	"regexp"
	"strings"
)

type configData struct {
	IntrospectDetailList string
	IntrospectNameList   []string
	IntrospectApis       string
	AuthConfig           string
	AuthRoleConfig       string
	S3UploadConfig       string
	CorsConfig           string
	SecurityConfig       string
}

func (w *wdg) buildConfigData(ctx context.Context) (configData, error) {
	cd := configData{}
	// 获取环境变量
	envs, err := w.repository.er.FindEnvs(ctx, "")
	if err != nil {
		log.Error("buildConfigData get envs info fail, err : ", err)
	}
	envMap, err := utils.GetEnvVal(envs)
	if err != nil {
		return cd, err
	}
	// 获取数据源更新内容
	w.buildDataSourceContent(ctx, &cd, envMap)
	cd.IntrospectApis = strings.Trim(cd.IntrospectApis, "[]")

	// 获取身份验证更新内容
	w.buildAuthContent(ctx, &cd)
	// 获取角色信息更新内容
	w.buildRolesContent(ctx, &cd)
	// 获取oss配置信息更新内容
	w.buildS3UploadContent(ctx, &cd)
	// 获取跨域配置信息更新内容
	w.buildCorsContent(ctx, &cd)
	// 获取安全更新内容
	w.buildSecurityContent(ctx, &cd)
	return cd, nil
}

func (w *wdg) buildDataSourceContent(ctx context.Context, cd *configData, envMap map[string]string) {
	ds, err := w.repository.dsr.FindDataSources(ctx)
	if err != nil {
		log.Error("buildDataSourceContent get datasource info fail, err : ", err)
		return
	}

	for i, v := range ds {
		if v.SourceType == domain.SourceTypeCustomize {
			continue
		}
		if v.SwitchOn() {
			// 用户命名可能会输入数字，在前面加上db前缀
			name := fmt.Sprintf("db_%s_%v", v.GetNameSpace(), i)
			str := GetDataSourceConfig(v, envMap)
			// 处理最外层的大括号
			str = strings.Trim(str, "{}")

			dbConfigStr := fmt.Sprintf(`const %s = introspect.%s({
						%s 
					}); `, name, v.GetSourceTypeStr(), str)
			cd.IntrospectDetailList += fmt.Sprintf("%s\n", dbConfigStr)
			cd.IntrospectNameList = append(cd.IntrospectNameList, name)
		}
	}

	jsonByte, _ := json.Marshal(cd.IntrospectNameList)
	cd.IntrospectApis = strings.ReplaceAll(string(jsonByte), "\"", "")
	if len(cd.IntrospectNameList) == 0 {
		cd.IntrospectApis = ""
	}
}

func (w *wdg) buildAuthContent(ctx context.Context, cd *configData) {
	auths, err := w.repository.ar.FindAuthentication(ctx)
	if err != nil {
		log.Error("buildAuthContent get auth info fail, err : ", err)
		return
	}
	var tokenProviders, cookieProviders []string
	for _, v := range auths {
		// 序列化
		var config domain.WdgAuthConfig
		err := json.Unmarshal([]byte(v.Config), &config)
		if err != nil {
			log.Error("auth unmarshal WdgAuthConfig fail , err : ", err)
		}
		if strings.Contains(v.SwitchState, "cookieBase") {
			// 根据供应商类型选择具体的连接方式 google、github、openIdConnect
			cookieProviders = append(cookieProviders, GetCookieProvidersConfig(v.AuthSupplier, config))
		}
		if strings.Contains(v.SwitchState, "tokenBase") {
			tokenProviders = append(tokenProviders, GetTokenProvidersConfig(config))
		}
	}
	var cookieBased CookieBased
	var tokenBased TokenBased
	// 获取全局配置的身份验证跳转url
	globalConfig, err := utils.GetGlobalConfig()
	if err != nil {
		log.Error("get auth redirecturis fail , err : ", err)
	}
	wdgAuth := WdgAuth{}
	if globalConfig.AuthRedirectURL != nil {
		authRedirectURLBytes, _ := json.Marshal(globalConfig.AuthRedirectURL)
		cookieBased.AuthorizedRedirectUris = string(authRedirectURLBytes)
	}
	if cookieProviders != nil {
		cookieProviderBytes, _ := json.Marshal(cookieProviders)
		cookie := strings.ReplaceAll(string(cookieProviderBytes), "\"", "")
		cookieBased.Providers = cookie
		wdgAuth.CookieBased = &cookieBased
	}
	if tokenProviders != nil {
		tokenProviderBytes, _ := json.Marshal(tokenProviders)
		token := strings.ReplaceAll(string(tokenProviderBytes), "\"", "")
		tokenBased.Providers = token
		wdgAuth.TokenBased = &tokenBased

	}

	wdgAuthBytes, _ := json.Marshal(wdgAuth)
	reg := regexp.MustCompile("\"([a-zA-Z]\\w*)\":")
	cd.AuthConfig = reg.ReplaceAllString(string(wdgAuthBytes), `$1:`)
	cd.AuthConfig = strings.ReplaceAll(cd.AuthConfig, "\"", "")
	if cd.AuthConfig != "" {
		cd.AuthConfig = fmt.Sprintf("authentication:%s ,", cd.AuthConfig)
	}
}

func (w *wdg) buildRolesContent(ctx context.Context, cd *configData) {
	roles, err := w.repository.rr.FindRoles(ctx)
	if err != nil {
		log.Error("buildAuthContent get auth info fail, err : ", err)
		return
	}
	rolesArr := make([]string, 0)
	for _, v := range roles {
		rolesArr = append(rolesArr, v.Code)
	}

	authorization := Authorization{}
	if len(rolesArr) != 0 {
		authorization.Roles = rolesArr
	}
	authorizationBytes, _ := json.Marshal(authorization)

	reg := regexp.MustCompile("\"([a-zA-Z]\\w*)\":")
	resultStr := reg.ReplaceAllString(string(authorizationBytes), `$1:`)
	if resultStr != "" {
		cd.AuthRoleConfig = fmt.Sprintf("authorization: %s ,", resultStr)
	}
}

func (w *wdg) buildS3UploadContent(ctx context.Context, cd *configData) {
	storageBuckets, err := w.repository.sbr.FindStorageBucket(ctx)
	if err != nil {
		log.Error("buildS3UploadContent get auth info fail, err : ", err)
		return
	}
	s3UploadProviderArr := make([]domain.S3UploadProvider, 0)
	for _, v := range storageBuckets {
		if v.SwitchOn() {
			var s3UploadProvider domain.S3UploadProvider
			err := json.Unmarshal([]byte(v.Config), &s3UploadProvider)
			if err != nil {
				log.Error("buildS3UploadContent fail ,err : ", err)
			}
			s3UploadProviderArr = append(s3UploadProviderArr, s3UploadProvider)
		}
	}
	configByte, _ := json.Marshal(s3UploadProviderArr)
	reg := regexp.MustCompile("\"([a-zA-Z]\\w*)\":")
	resultStr := reg.ReplaceAllString(string(configByte), `$1:`)

	cd.S3UploadConfig = fmt.Sprintf("s3UploadProvider:%s,", resultStr)
}

func (w *wdg) buildCorsContent(ctx context.Context, cd *configData) {
	// 读取配置文件
	globalConfig, err := utils.GetGlobalConfig()
	if err != nil {
		log.Error("buildCorsContent get setting config fail , err : ", err)
	}
	if globalConfig.ConfigureWunderGraphApplication.Cors == nil {
		return
	}
	corsConfig := WdgCors{
		AllowedOrigins:   globalConfig.ConfigureWunderGraphApplication.Cors.AllowedOrigins,
		AllowedMethods:   globalConfig.ConfigureWunderGraphApplication.Cors.AllowedMethods,
		AllowedHeaders:   globalConfig.ConfigureWunderGraphApplication.Cors.AllowedHeaders,
		ExposedHeaders:   globalConfig.ConfigureWunderGraphApplication.Cors.ExposedHeaders,
		MaxAge:           globalConfig.ConfigureWunderGraphApplication.Cors.MaxAge,
		AllowCredentials: globalConfig.ConfigureWunderGraphApplication.Cors.AllowCredentials,
	}

	if len(corsConfig.AllowedOrigins) == 0 {
		corsConfig.AllowedOrigins = append(corsConfig.AllowedOrigins, fmt.Sprintf("localhost:%s", utils.GetFireBoomPort()))
	}
	corsByte, _ := json.Marshal(corsConfig)
	reg := regexp.MustCompile("\"([a-zA-Z]\\w*)\":")
	corsStr := reg.ReplaceAllString(string(corsByte), `$1:`)
	corsStr = strings.Trim(corsStr, "{}")
	cd.CorsConfig = fmt.Sprintf("cors: {        ...cors.allowAll, %s },", corsStr)
}

func (w *wdg) buildSecurityContent(ctx context.Context, cd *configData) {
	// 读取配置文件
	globalConfig, err := utils.GetGlobalConfig()
	if err != nil {
		log.Error("buildCorsContent get setting config fail , err : ", err)
	}
	if globalConfig.ConfigureWunderGraphApplication.Security == nil {
		return
	}
	securityConfig := WdgSecurityConfig{
		EnableGraphQLEndpoint: globalConfig.ConfigureWunderGraphApplication.Security.EnableGraphQLEndpoint,
		AllowedHosts:          globalConfig.ConfigureWunderGraphApplication.Security.AllowedHosts,
	}

	securityByte, _ := json.Marshal(securityConfig)
	reg := regexp.MustCompile("\"([a-zA-Z]\\w*)\":")
	securityStr := reg.ReplaceAllString(string(securityByte), `$1:`)
	cd.SecurityConfig = fmt.Sprintf("security: %s ,", securityStr)
}
