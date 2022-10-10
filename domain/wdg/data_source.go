package wdg

import (
	"encoding/json"
	"fmt"
	"github.com/fire_boom/domain"
	"github.com/fire_boom/utils"
	"github.com/labstack/gommon/log"
)

type WgDataSourceConfig struct {
	IntrospectionType string
	ConfigString      string
}

func GetDataSourceConfig(source domain.FbDataSource) (result string) {

	str := ""
	switch source.SourceType {
	case domain.SourceTypeDB:
		str = getDatabaseIntrospection(source)
	case domain.SourceTypeRest:
		str = getOpenAPIIntrospection(source)
	case domain.SourceTypeGraphQL:
		str = getGraphQLUpstream(source)
	case domain.SourceTypeCustomize:
		str = ""
	default:
		return
	}
	result = fmt.Sprintf(`const %s = introspect.%s({
						%s
					});`, source.Name, source.GetSourceTypeStr(), str)
	return
}

func getGraphQLUpstream(databaseConfig domain.FbDataSource) string {
	var config domain.GraphqlConfig
	err := json.Unmarshal([]byte(databaseConfig.Config), &config)
	if err != nil {
		log.Error("getGraphQLUpstream err : ", err)
		return ""
	}
	var result GraphQLUpstream
	result.Url = config.URL
	result.ApiNamespace = config.ApiNameSpace

	headers := GetHeaders(config.Headers)
	if headers != "" {
		result.Headers = headers
	}

	result.CustomFloatScalars, err = utils.ArrStrFormatStrArr(config.CustomFloatScalars)
	result.CustomIntScalars, err = utils.ArrStrFormatStrArr(config.CustomIntScalars)
	result.SkipRenameRootFields, err = utils.ArrStrFormatStrArr(config.SkipRenameRootFields)

	if err != nil {
		log.Error("GetGraphQLUpstream err : ", err)
	}
	result.LoadSchemaFromString = config.LoadSchemaFromString
	result.Internal = config.Internal

	jsonByte, err := json.Marshal(result)
	return string(jsonByte)
}

func getOpenAPIIntrospection(databaseConfig domain.FbDataSource) string {
	var config domain.RestConfig
	err := json.Unmarshal([]byte(databaseConfig.Config), &config)
	if err != nil {
		log.Error("getGraphQLUpstream err : ", err)
		return ""
	}
	var result OpenAPIIntrospection
	// TODO 获取文件路径
	source := OpenAPIIntrospectionFile{
		FilePath: config.OASFileID,
	}
	auth := HTTPUpstreamAuthentication{
		Secret:                      config.Secret.Val,
		SigningMethod:               config.SigningMethod,
		AccessTokenExchangeEndpoint: config.TokenPoint,
	}

	result.Source = source
	result.ApiNamespace = config.ApiNameSpace
	result.StatusCodeUnions = config.StatusCodeUnions
	result.Authentication = auth
	headers := GetHeaders(config.Headers)
	if headers != "" {
		result.Headers = headers
	}

	jsonByte, err := json.Marshal(result)
	return string(jsonByte)
}

func getDatabaseIntrospection(databaseConfig domain.FbDataSource) string {
	var config domain.DbConfig
	err := json.Unmarshal([]byte(databaseConfig.Config), &config)
	if err != nil {
		log.Error("getGraphQLUpstream err : ", err)
		return ""
	}
	var result DatabaseIntrospection
	result.DatabaseURL = config.DatabaseURL.Val
	result.ApiNamespace = config.ApiNamespace

	jsonByte, err := json.Marshal(result)
	return string(jsonByte)
}

// GetHeaders 请求头转换
func GetHeaders(headerArr []domain.Value) (result string) {
	/*
		headers: builder => builder
		        .addStaticHeader("AuthToken", "staticToken")
		        .addStaticHeader("xxx", "xxx")
		        .addClientRequestHeader("Authorization", "Authorization")
	*/
	for _, row := range headerArr {
		// kind为2说明是转发客户端，该kv可为空，且只能有一个
		if row.Kind == "2" {
			// 设置后不用走下去直接下一次循环
			addClientRequestHeaderStr := fmt.Sprintf("addClientRequestHeader('%s', '%s')", row.Kay, row.Val)
			result = fmt.Sprintf("%s %s", result, addClientRequestHeaderStr)
			continue
		}
		// kind为1则说明是环境变量
		if row.Kind == "1" {
			// TODO 通过val读取环境变量
			row.Val = "读取的环境变量"
		}
		addStaticHeader := fmt.Sprintf("addStaticHeader('%s', '%s')", row.Kay, row.Val)
		result = fmt.Sprintf("%s %s", result, addStaticHeader)
	}
	if result == "" {
		return ""
	}
	return fmt.Sprintf("builder => builder %s", result)
}

type DatabaseIntrospection struct {
	DatabaseURL           string                              `json:"databaseURL"`
	ApiNamespace          string                              `json:"apiNamespace"`
	SchemaExtension       string                              `json:"schemaExtension"`
	ReplaceJSONTypeFields []ReplaceJSONTypeFieldConfiguration `json:"replaceJSONTypeFields"`
}

type ReplaceJSONTypeFieldConfiguration struct {
	EntityName              string `json:"entityName"`
	FieldName               string `json:"fieldName"`
	InputTypeReplacement    string `json:"inputTypeReplacement"`
	ResponseTypeReplacement string `json:"responseTypeReplacement"`
}

type GraphQLUpstream struct {
	Url              string `json:"url,omitempty"`
	SubscriptionsURL string `json:"subscriptionsURL,omitempty"`
	HTTPUpstream
	GraphQLIntrospectionOptions
}

type GraphQLIntrospectionOptions struct {
	LoadSchemaFromString string   `json:"loadSchemaFromString"`
	CustomFloatScalars   []string `json:"customFloatScalars"`
	CustomIntScalars     []string `json:"customIntScalars"`
	Internal             bool     `json:"internal"`
	SkipRenameRootFields []string `json:"skipRenameRootFields"`
}

type OpenAPIIntrospection struct {
	Source           OpenAPIIntrospectionFile `json:"source,omitempty"`
	StatusCodeUnions bool                     `json:"statusCodeUnions,omitempty"`
	HTTPUpstream
}

type OpenAPIIntrospectionFile struct {
	Kind     string `json:"kind,default=file,omitempty"` // 文件类型,固定为file
	FilePath string `json:"filePath,omitempty"`          // 文件路径
}

type HTTPUpstreamAuthentication struct {
	Secret                      string `json:"secret,omitempty"`
	SigningMethod               string `json:"signingMethod,default=HS256,omitempty"`
	AccessTokenExchangeEndpoint string `json:"accessTokenExchangeEndpoint,omitempty"`
}

type HTTPUpstream struct {
	ApiNamespace   string                     `json:"apiNamespace,omitempty"` // 命名空间
	Authentication HTTPUpstreamAuthentication `json:"authentication,omitempty"`
	Headers        string                     `json:"headers,omitempty"`
}
