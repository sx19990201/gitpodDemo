package wundergraph

import (
	"encoding/json"
	"fmt"
	"github.com/fire_boom/domain"
	"github.com/labstack/gommon/log"
	"regexp"
	"strings"
)

const (
	openIdConnect = "openIdConnect"
	github        = "github"
	google        = "google"
)

type WdgAuth struct {
	CookieBased *CookieBased `json:"cookieBased,omitempty,optional"`
	TokenBased  *TokenBased  `json:"tokenBased,omitempty,optional"`
}
type CookieBased struct {
	Providers                    string `json:"providers,omitempty,optional"`
	AuthorizedRedirectUris       string `json:"authorizedRedirectUris,omitempty,optional"`
	AuthorizedRedirectUriRegexes string `json:"authorizedRedirectUriRegexes,omitempty,optional"`
}

type TokenBased struct {
	Providers string `json:"providers,omitempty,optional"`
}

type TokenProviders struct {
	JwksJSON                string `json:"jwksJson,omitempty,optional"`
	JwksURL                 string `json:"jwksUrl,omitempty,optional"`
	UserInfoEndpoint        string `json:"userInfoEndpoint,omitempty,optional"`
	UserInfoCacheTtlSeconds int64  `json:"userInfoCacheTtlSeconds,omitempty,optional"`
}

type OpenIDConnectAuthProviderConfig struct {
	Id           string `json:"id,omitempty,optional"`
	Issuer       string `json:"issuer,omitempty,optional"`
	ClientId     string `json:"clientId,omitempty,optional"`
	ClientSecret string `json:"clientSecret,omitempty,optional"`
}

type GithubAuthProviderConfig struct {
	Id           string `json:"id,omitempty,optional"`
	ClientId     string `json:"clientId,omitempty,optional"`
	ClientSecret string `json:"clientSecret,omitempty,optional"`
}

type GoogleAuthProviderConfig struct {
	Id           string `json:"id,omitempty,optional"`
	ClientId     string `json:"clientId,omitempty,optional"`
	ClientSecret string `json:"clientSecret,omitempty,optional"`
}

func GetCookieProvidersConfig(authSupplier string, config domain.WdgAuthConfig) string {
	var jsonByte []byte
	var err error
	// 根据供应商类型选择结构体
	switch authSupplier {
	case domain.OpenID:
		openid := OpenIDConnectAuthProviderConfig{
			Id:           config.ID,
			Issuer:       config.Issuer,
			ClientId:     config.ClientId,
			ClientSecret: config.ClientSecret,
		}
		jsonByte, err = json.Marshal(openid)
	case domain.Github:
		github := GithubAuthProviderConfig{
			Id:           config.ID,
			ClientId:     config.ClientId,
			ClientSecret: config.ClientSecret,
		}
		jsonByte, err = json.Marshal(github)
	case domain.Google:
		google := GoogleAuthProviderConfig{
			Id:           config.ID,
			ClientId:     config.ClientId,
			ClientSecret: config.ClientSecret,
		}
		jsonByte, err = json.Marshal(google)
	}
	if err != nil {
		log.Error("authProviders marshal fail , err : ", err)
	}
	// 处理双引号
	reg := regexp.MustCompile("\"([a-zA-Z]\\w*)\":")
	str := reg.ReplaceAllString(string(jsonByte), `$1:`)

	result := fmt.Sprintf("authProviders.%s(%s),", getAuthProvidersName(authSupplier), str)
	result = strings.ReplaceAll(result, "\"", "'")
	return result
}

func GetTokenProvidersConfig(config domain.WdgAuthConfig) string {
	if (config.JwksURL == "" && config.JwksJSON == "") || config.UserInfoEndpoint == "" {
		return ""
	}
	token := TokenProviders{
		JwksJSON:                config.JwksJSON,
		JwksURL:                 config.JwksURL,
		UserInfoEndpoint:        config.UserInfoEndpoint,
		UserInfoCacheTtlSeconds: config.UserInfoCacheTtlSeconds,
	}
	jsonByte, err := json.Marshal(token)
	if err != nil {
		log.Error("TokenProviders marshal fail , err : ", err)
	}
	// 处理双引号
	reg := regexp.MustCompile("\"([a-zA-Z]\\w*)\":")
	result := reg.ReplaceAllString(string(jsonByte), `$1:`)
	result = strings.ReplaceAll(result, "\"", "'")
	return result
}

func getAuthProvidersName(authSupplier string) string {
	switch authSupplier {
	case domain.OpenID:
		return openIdConnect
	case domain.Github:
		return github
	case domain.Google:
		return google
	}
	return openIdConnect
}
