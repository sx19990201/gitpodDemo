package domain

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/labstack/gommon/log"
	"strings"
)

const (
	SwitchOn  = 0
	SwitchOff = 1
)

const (
	SourceTypeDB = iota + 1
	SourceTypeRest
	SourceTypeGraphQL
	SourceTypeCustomize
)

const (
	IntrospectGraphql   = "graphql"
	IntrospectPG        = "postgresql"
	IntrospectMysql     = "mysql"
	IntrospectSqlLite   = "sqlite"
	IntrospectSqlServer = "sqlserver"
	IntrospectMongodb   = "mongodb"
	IntrospectOPenAPI   = "openApi"
	IntrospectUnSupport = "unSupprot"
)

const (
	Postgresql = "postgresql"
	Mysql      = "mysql"
	Sqlite     = "sqlite"
	Sqlserver  = "sqlserver"
	Mongodb    = "mongodb"
)

const (
	postgresql = iota + 1
	mysql
	planetscale
	sqlite
	sqlserver
	mongodb
)

// FbDataSource 数据源配置
type FbDataSource struct {
	ID         uint           `db:"id" json:"id"`
	Name       string         `db:"name" json:"name"`             // 数据源名称
	SourceType int            `db:"sourceType" json:"sourceType"` // 数据源类型: 1-db 2-rest 3-graphql 4-自定义
	Config     string         `db:"config" json:"config"`         // 数据源对应的配置项：命名空间、请求配置、连接配置、文件路径、是否配置为外部数据源等
	Switch     uint8          `db:"switch" json:"switch"`         // 开关 0-开 1-关
	CreateTime sql.NullString `db:"create_time" json:"-"`
	UpdateTime sql.NullString `db:"update_time" json:"-"`
	IsDel      uint8          `db:"isDel" json:"-"`
}

type FbDataSourceResp struct {
	ID         uint           `db:"id" json:"id"`
	Name       string         `db:"name" json:"name"`             // 数据源名称
	SourceType int            `db:"sourceType" json:"sourceType"` // 数据源类型: 1-db 2-rest 3-graphql 4-自定义
	Config     interface{}    `db:"config" json:"config"`         // 数据源对应的配置项：命名空间、请求配置、连接配置、文件路径、是否配置为外部数据源等
	Switch     uint8          `db:"switch" json:"switch"`         // 开关 0-开 1-关
	CreateTime sql.NullString `db:"create_time" json:"-"`
	UpdateTime sql.NullString `db:"update_time" json:"-"`
	IsDel      uint8          `db:"isDel" json:"-"`
}

func (f *FbDataSourceResp) TransformDataSource() (result FbDataSource) {
	configByte, _ := json.Marshal(f.Config)
	result.ID = f.ID
	result.Name = f.Name
	result.SourceType = f.SourceType
	result.Switch = f.Switch
	result.CreateTime = f.CreateTime
	result.UpdateTime = f.UpdateTime
	result.IsDel = f.IsDel
	result.Config = string(configByte)
	return
}

func (f *FbDataSource) GetFbDataSourceResp() (result FbDataSourceResp) {
	result.ID = f.ID
	result.Name = f.Name
	result.SourceType = f.SourceType
	result.Switch = f.Switch
	result.CreateTime = f.CreateTime
	result.UpdateTime = f.UpdateTime
	result.IsDel = f.IsDel
	if f.SourceType == SourceTypeDB {
		result.Config = f.GetDbConfig()
	}
	if f.SourceType == SourceTypeRest {
		result.Config = f.GetRestConfig()
	}
	if f.SourceType == SourceTypeGraphQL {
		result.Config = f.GetGraphqlConfig()
	}
	if f.SourceType == SourceTypeCustomize {
		result.Config = f.GetCustomizeConfig()
	}
	return
}

type RestConfig struct {
	ApiNameSpace     string  `json:"apiNameSpace"`     // 命名空间
	OASFileID        string  `json:"filePath"`         // oas文件id
	JWTType          string  `json:"jwtType"`          // jwt获取方式 0静态 1-动态
	Secret           Value   `json:"secret"`           // 秘钥
	SigningMethod    string  `json:"signingMethod"`    // 签名获取方法
	TokenPoint       string  `json:"tokenPoint"`       // token端点
	StatusCodeUnions bool    `json:"statusCodeUnions"` // 是否状态联合
	Headers          []Value `json:"headers"`          // 请求头
}

type GraphqlConfig struct {
	ApiNameSpace         string  `json:"apiNameSpace"`         // 命名空间
	URL                  string  `json:"url"`                  // graphql端点
	LoadSchemaFromString string  `json:"loadSchemaFromString"` // 指定schema
	Internal             bool    `json:"internal"`             // 是否内部
	CustomFloatScalars   string  `json:"customFloatScalars"`   // 自定义float标量
	CustomIntScalars     string  `json:"customIntScalars"`     // 自定义int标量
	SkipRenameRootFields string  `json:"skipRenameRootFields"` // 排除重命名根字段
	Headers              []Value `json:"headers"`              // 请求头
}

type DbConfig struct {
	ApiNamespace                      string                              `json:"apiNamespace"` // 连接名
	DBType                            string                              `json:"dbType"`       // db类型 mysql等等
	AppendType                        string                              `json:"appendType"`   // 类型 0-连接url 1-连接参数
	DatabaseURL                       Value                               `json:"databaseUrl"`  // 连接url
	SchemaExtension                   string                              `json:"schemaExtension"`
	ReplaceJSONTypeFieldConfiguration []ReplaceJSONTypeFieldConfiguration `json:"replaceJSONTypeFieldConfiguration"`
	Host                              string                              `json:"host"`     // 主机
	DBName                            string                              `json:"dbName"`   // 数据库名称
	Port                              string                              `json:"port"`     // 端口
	UserName                          string                              `json:"userName"` // 用户名
	Password                          string                              `json:"password"` // 密码
}

type ReplaceJSONTypeFieldConfiguration struct {
	EntityName              string `json:"entityName"`
	FieldName               string `json:"fieldName"`
	InputTypeReplacement    string `json:"inputTypeReplacement"`
	ResponseTypeReplacement string `json:"responseTypeReplacement"`
}

type CustomizeConfig struct {
	ApiNamespace string `json:"apiNamespace"`
	ServerName   string `json:"serverName"`
	Schema       string `json:"schema"`
}

type Value struct {
	Kay  string `json:"key"`
	Kind string `json:"kind"` // 0-值 1-环境变量 2-转发值客户端
	Val  string `json:"val"`
}

func (f *FbDataSource) GetDBStr() string {
	var dbConfig DbConfig
	err := json.Unmarshal([]byte(f.Config), &dbConfig)
	if err != nil {
		log.Error("序列化错误")
	}
	return strings.ToLower(dbConfig.DBType)
	//switch strings.ToLower(dbConfig.DBType) {
	//case postgresql:
	//	return IntrospectPG
	//case mysql:
	//	return IntrospectMysql
	//case sqlite:
	//	return IntrospectSqlLite
	//case sqlserver:
	//	return IntrospectSqlServer
	//case mongodb:
	//	return IntrospectMongodb
	//default:
	//	return IntrospectUnSupport
	//}
}
func (f *FbDataSource) GetSourceTypeStr() string {

	switch f.SourceType {
	case SourceTypeDB:
		return f.GetDBStr()
	case SourceTypeRest:
		return IntrospectOPenAPI
	case SourceTypeCustomize:
		return IntrospectSqlLite
	case SourceTypeGraphQL:
		return IntrospectGraphql
	default:
		return IntrospectUnSupport
	}
}

func (f *FbDataSource) SwitchOn() bool {
	return f.Switch == SwitchOn
}

// GetNameSpace TODO  获取nameSpace
func (f *FbDataSource) GetNameSpace() string {
	switch f.SourceType {
	case SourceTypeDB:
		return f.GetDbConfig().ApiNamespace
	case SourceTypeRest:
		return f.GetRestConfig().ApiNameSpace
	case SourceTypeGraphQL:
		return f.GetGraphqlConfig().ApiNameSpace
	case SourceTypeCustomize:
		return f.GetCustomizeConfig().ApiNamespace
	default:
		return ""
	}
}

// TODO
func (f *FbDataSource) GetDatabaseURL() string {
	return f.Config
}

func (f *FbDataSource) GetDbConfig() (result DbConfig) {
	err := json.Unmarshal([]byte(f.Config), &result)
	if err != nil {
		log.Error("GetDbConfig json.Unmarshal err : ", err)
	}
	return
}
func (f *FbDataSource) GetRestConfig() (result RestConfig) {
	err := json.Unmarshal([]byte(f.Config), &result)
	if err != nil {
		log.Error("GetRestConfig json.Unmarshal err : ", err)
	}
	return
}
func (f *FbDataSource) GetGraphqlConfig() (result GraphqlConfig) {
	err := json.Unmarshal([]byte(f.Config), &result)
	if err != nil {
		log.Error("GetGraphqlConfig json.Unmarshal err : ", err)
	}
	return
}
func (f *FbDataSource) GetCustomizeConfig() (result CustomizeConfig) {
	err := json.Unmarshal([]byte(f.Config), &result)
	if err != nil {
		log.Error("GetCustomizeConfig json.Unmarshal err : ", err)
	}
	return
}

type DataSourceUseCase interface {
	Store(ctx context.Context, f *FbDataSource) (int64, error)
	GetByID(ctx context.Context, id uint) (FbDataSource, error)
	Update(ctx context.Context, f *FbDataSource) (int64, error)
	Delete(ctx context.Context, id uint) (int64, error)
	FindDataSources(ctx context.Context) ([]FbDataSource, error)
	GetPrismaSchema(ctx context.Context, id uint) (result string, err error)
}

type DataSourceRepository interface {
	Store(ctx context.Context, f *FbDataSource) (int64, error)
	Update(ctx context.Context, f *FbDataSource) (int64, error)
	GetByName(ctx context.Context, name string) (FbDataSource, error)
	GetByID(ctx context.Context, id uint) (FbDataSource, error)
	CheckExist(ctx context.Context, auth *FbDataSource) (FbDataSource, error)
	Delete(ctx context.Context, id uint) (int64, error)
	FindDataSources(ctx context.Context) ([]FbDataSource, error)
}
