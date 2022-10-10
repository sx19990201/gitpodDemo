package wundergraph

import (
	"encoding/json"
	"fmt"
	"github.com/fire_boom/domain"
	"github.com/fire_boom/utils"
	"github.com/labstack/gommon/log"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
)

type DatabaseIntrospection struct {
	DatabaseURL           string                              `json:"databaseURL,omitempty,optional"`
	ApiNamespace          string                              `json:"apiNamespace,omitempty,optional"`
	SchemaExtension       string                              `json:"schemaExtension,omitempty,optional"`
	ReplaceJSONTypeFields []ReplaceJSONTypeFieldConfiguration `json:"replaceJSONTypeFields,omitempty,optional"`
}

type ReplaceJSONTypeFieldConfiguration struct {
	EntityName              string `json:"entityName,omitempty,optional"`
	FieldName               string `json:"fieldName,omitempty,optional"`
	InputTypeReplacement    string `json:"inputTypeReplacement,omitempty,optional"`
	ResponseTypeReplacement string `json:"responseTypeReplacement,omitempty,optional"`
}

type GraphQLUpstream struct {
	Url              string `json:"url,omitempty,optional"`
	SubscriptionsURL string `json:"subscriptionsURL,omitempty,optional"`
	HTTPUpstream
	GraphQLIntrospectionOptions
}

type GraphQLIntrospectionOptions struct {
	LoadSchemaFromString string   `json:"loadSchemaFromString,omitempty,optional"`
	CustomFloatScalars   []string `json:"customFloatScalars,omitempty,optional"`
	CustomIntScalars     []string `json:"customIntScalars,omitempty,optional"`
	Internal             bool     `json:"internal,omitempty,optional"`
	SkipRenameRootFields []string `json:"skipRenameRootFields,omitempty,optional"`
}

type OpenAPIIntrospection struct {
	Source OpenAPIIntrospectionFile `json:"source,omitempty,optional"`
	HTTPUpstream
	StatusCodeUnions bool `json:"statusCodeUnions,omitempty,optional"`
}

type OpenAPIIntrospectionFile struct {
	Kind     string `json:"kind" default:"file"` // 文件类型,固定为file
	FilePath string `json:"filePath"`            // 文件路径
}

type HTTPUpstreamAuthentication struct {
	Kind                        string `json:"kind,omitempty,optional"`
	Secret                      string `json:"secret,omitempty,optional"`
	SigningMethod               string `json:"signingMethod,omitempty,optional"`
	AccessTokenExchangeEndpoint string `json:"accessTokenExchangeEndpoint,omitempty,optional"`
}

type HTTPUpstream struct {
	ApiNamespace   string                      `json:"apiNamespace,omitempty,optional"` // 命名空间
	Authentication *HTTPUpstreamAuthentication `json:"authentication,omitempty,optional"`
	Headers        string                      `json:"headers,omitempty,optional"`
}

func GetDataSourceConfig(source domain.FbDataSource, envMap map[string]string) (result string) {
	str := ""
	switch source.SourceType {
	case domain.SourceTypeDB:
		str = getDatabaseIntrospection(source, envMap)
	case domain.SourceTypeRest:
		str = getOpenAPIIntrospection(source, envMap)
	case domain.SourceTypeGraphQL:
		str = getGraphQLUpstream(source, envMap)
	//case domain.SourceTypeCustomize:
	//	str = ""
	default:
		return
	}
	// 处理双引号
	reg := regexp.MustCompile("\"([a-zA-Z]\\w*)\":")
	str = reg.ReplaceAllString(str, `$1:`)

	result = strings.ReplaceAll(str, `headers:builder =\u003e builder`, `headers: builder => builder`)

	result = strings.ReplaceAll(result, `)",statusCodeUnions`, `),statusCodeUnions`)
	result = strings.ReplaceAll(result, `)",loadSchemaFromString`, `),loadSchemaFromString`)
	// addStaticHeader\('[\w\d]*','[\w\d]*'\)"
	reg = regexp.MustCompile("addStaticHeader\\('[\\w\\d]*','[\\w\\d]*'\\)\"")
	result = reg.ReplaceAllString(result, strings.ReplaceAll(reg.FindString(result), "\"", ""))
	return
}

func getGraphQLUpstream(databaseConfig domain.FbDataSource, envMap map[string]string) string {
	config := databaseConfig.GetGraphqlConfig()
	var result GraphQLUpstream
	result.Url = config.URL
	result.ApiNamespace = config.ApiNameSpace
	headers := GetHeaders(config.Headers, envMap)
	if headers != "" {
		result.Headers = headers
	}
	strings.Split(config.CustomFloatScalars, ",")
	if config.CustomFloatScalars != "" {
		result.CustomFloatScalars = strings.Split(config.CustomFloatScalars, ",")
	}
	if config.CustomIntScalars != "" {
		result.CustomIntScalars = strings.Split(config.CustomIntScalars, ",")
	}
	if config.SkipRenameRootFields != "" {
		result.SkipRenameRootFields = strings.Split(config.SkipRenameRootFields, ",")
	}

	if config.LoadSchemaFromString != "" {
		// 此时config.LoadSchemaFromString内保存的是文件名，需要读取文件内容再更换
		schemaContent, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", utils.GetOASFilePath(), filepath.Base(config.LoadSchemaFromString)))
		if err != nil {
			log.Error("getGraphQLUpstream read schema content fail ,err : ", err)
			return ""
		}
		result.LoadSchemaFromString = string(schemaContent)
	}
	result.Internal = config.Internal
	// 以这个为header结束开始的结束符号
	result.Headers = fmt.Sprintf(`***@@@@%s***@@@@`, result.Headers)
	jsonByte, _ := json.Marshal(result)
	resultStr := strings.ReplaceAll(string(jsonByte), `"***@@@@`, "")
	resultStr = strings.ReplaceAll(resultStr, `***@@@@"`, "")
	return resultStr
}

func getOpenAPIIntrospection(databaseConfig domain.FbDataSource, envMap map[string]string) string {
	config := databaseConfig.GetRestConfig()
	var result OpenAPIIntrospection
	source := OpenAPIIntrospectionFile{
		Kind:     "file",
		FilePath: fmt.Sprintf("../../%s/%s", utils.GetOASFilePath(), filepath.Base(config.OASFileID)),
	}
	result.Authentication = nil
	if config.Secret.Val != "" && config.TokenPoint != "" {
		result.Authentication = &HTTPUpstreamAuthentication{
			Kind:                        "jwt_with_access_token_exchange",
			SigningMethod:               "HS256",
			Secret:                      config.Secret.Val,
			AccessTokenExchangeEndpoint: config.TokenPoint,
		}
	}
	result.Source = source
	result.ApiNamespace = config.ApiNameSpace
	result.StatusCodeUnions = config.StatusCodeUnions
	headers := GetHeaders(config.Headers, envMap)
	if headers != "" {
		result.Headers = headers
	}

	resultByte, _ := json.Marshal(result)
	resultStr := string(resultByte)
	return resultStr
}

func getDatabaseIntrospection(databaseConfig domain.FbDataSource, envMap map[string]string) string {
	config := databaseConfig.GetDbConfig()
	// 如果databaseurl是环境变量需要读取环境变量
	if config.DatabaseURL.Kind == "1" {
		//读取环境变量
		config.DatabaseURL.Val = envMap[config.DatabaseURL.Val]
	}
	// 如果是连接参数，组装连接url
	if config.AppendType == "1" {
		//"root:shaoxiong123456@tcp(8.142.115.204:3306)/gotrue_development?parseTime=true&multiStatements=true"
		config.DatabaseURL.Val = fmt.Sprintf("%s://%s:%s@%s:%s/%s", strings.ToLower(config.DBType), config.UserName, config.Password, config.Host, config.Port, config.DBName)
	} else {
		config.DatabaseURL.Val = fmt.Sprintf("%s://%s", strings.ToLower(config.DBType), config.DatabaseURL.Val)
	}
	var result DatabaseIntrospection
	result.ApiNamespace = config.ApiNamespace
	result.DatabaseURL = config.DatabaseURL.Val

	jsonByte, _ := json.Marshal(result)
	return string(jsonByte)
}

// GetHeaders 请求头转换
func GetHeaders(headerArr []domain.Value, envMap map[string]string) (result string) {
	/*
		headers: builder => builder
		        .addStaticHeader("AuthToken", "staticToken")
		        .addStaticHeader("xxx", "xxx")
		        .addClientRequestHeader("Authorization", "Authorization")
	*/
	for _, row := range headerArr {
		// kind为2说明是转发客户端，该kv可为空，如果有则只能有一个
		if row.Kind == "2" {
			// 设置后不用走下去直接下一次循环
			addClientRequestHeaderStr := fmt.Sprintf(".addClientRequestHeader('%s', '%s')", row.Kay, row.Val)
			result = fmt.Sprintf("%s %s", result, addClientRequestHeaderStr)
			continue
		}
		// kind为1则说明是环境变量
		if row.Kind == "1" {
			// TODO 通过val读取环境变量
			row.Val = envMap[row.Val]
		}
		addStaticHeader := fmt.Sprintf(".addStaticHeader('%s','%s')", row.Kay, row.Val)
		result = fmt.Sprintf("%s %s", result, addStaticHeader)
	}
	if result == "" {
		return ""
	}
	result = fmt.Sprintf("builder => builder %s", result)
	result = strings.ReplaceAll(result, "\"", "")
	return result
}
