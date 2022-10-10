package domain

import (
	"errors"
)

// 错误码前三位功能code
const (
	AuthCode          = "101" // 身份验证
	DataSourceCode    = "102" // 数据源
	DraftCode         = "103" // 草稿
	FileCode          = "104" // 文件
	OperateAPICode    = "105" // 操作api
	PrismaCode        = "106" // prisma
	RoleCode          = "107" // 角色
	S3UploadCode      = "108" // s3upload
	ScriptCode        = "109" // 脚本
	SettingCode       = "110" // 设置
	StorageBucketCode = "111" // 存储
	WunderGlobal      = "112" // 全局配置
)

// 中间两位类型码
const (
	ParamsErrCode  = "01"
	NetworkErrCode = "02"
	DbErrCode      = "03"
)

// 后三位错误编码
const (
	ParamsEmptyCode = "001"
	FieldIdErrCode  = "001"
)

// TODO
var (
	// ErrNotFound will throw if the requested item is not exists
	ErrNotFound = errors.New("your requested Item is not found")
)

type BizErr struct {
	code string `json:"code"`
	msg  string `json:"msg"`
}

func NewBizErr(code string, msg string) *BizErr {
	return &BizErr{code: code, msg: msg}
}

// 实现 Error 方法
func (be *BizErr) Error() string {
	return be.msg
}

func (be *BizErr) Code() string {
	return be.code
}

var (
	// OK 通用返回码
	OK = NewBizErr("10000000", "SUCCESS") // 成功

	// ParamErr 10-参数错误码
	ParamErr                 = NewBizErr("10000", "系统异常")
	ParamIdEmptyErr          = NewBizErr("10001", "id不能为空")
	ParamNameEmptyErr        = NewBizErr("10002", "名称不能为空")
	ParamNameExistErr        = NewBizErr("10003", "名称已经存在")
	ParamSupplierNotExistErr = NewBizErr("10004", "供应商不存在")
	ParamSourceTypeEmptyErr  = NewBizErr("10005", "数据源类型不能为空")
	ParamPathEmptyErr        = NewBizErr("10006", "路径不能为空")
	ParamCodeEmptyErr        = NewBizErr("10007", "code不能为空")
	ParamDateErr             = NewBizErr("10008", "时间参数不合法")
	ParamRePasswordErr       = NewBizErr("10009", "两次密码不一致")

	// NetworkErr 20-网络错误码
	NetworkErr = NewBizErr("20001", "网络异常")

	// DbErr 30-数据库错误码
	DbErr               = NewBizErr("30000", "内部异常")
	DbCreateErr         = NewBizErr("30001", "添加失败")
	DbUpdateErr         = NewBizErr("30002", "修改失败")
	DbDeleteErr         = NewBizErr("30003", "删除失败")
	DbFindErr           = NewBizErr("30004", "查询失败")
	DbCheckNameExistErr = NewBizErr("30005", "检查名称是否存在出现异常")
	DbNameExistErr      = NewBizErr("30005", "名称已经存在")
	DbPingErr           = NewBizErr("30006", "数据库连接失败，已将配置关闭，请检查配置后再开启")
	OasContentCheckErr  = NewBizErr("30007", "OAS文件不合法，已将配置关闭，请检查配置后再开启")
	DbNameNotExistErr   = NewBizErr("30008", "名称不存在")
	DbConnErr           = NewBizErr("30009", "数据源连接失败，请检查后重试")
	DbIntrospectionErr  = NewBizErr("30010", "数据源内省失败")
	DbMigrateErr        = NewBizErr("30011", "数据源迁移失败")
	DbSchemaErr         = NewBizErr("30012", "获取数据源schema失败")
	DbTableEmptyErr     = NewBizErr("30013", "数据源内表为空")

	// FileErr 40-文件错误码
	FileErr       = NewBizErr("40000", "文件异常")
	FileReadErr   = NewBizErr("40001", "文件读取错误")
	FileWriteErr  = NewBizErr("40002", "文件写入错误")
	FileDeleteErr = NewBizErr("40003", "文件删除错误")
	FileCreateErr = NewBizErr("40004", "文件创建错误")
	FileReNameErr = NewBizErr("40005", "文件重命名错误")
	FileOpenErr   = NewBizErr("40006", "文件打开错误")
	FileCopyErr   = NewBizErr("40007", "文件拷贝错误")
	FileTypeErr   = NewBizErr("40008", "该文件不是json或yaml文件")
	FileExportErr = NewBizErr("40009", "文件导出失败")
	FileImportErr = NewBizErr("40010", "文件导入失败")

	// JsonErr 50-json错误码
	JsonErr          = NewBizErr("50000", "系统异常")
	JsonMarshalErr   = NewBizErr("50001", "Json序列化错误")
	JsonUnMarshalErr = NewBizErr("50002", "Json反序列化错误")
	JsonSetErr       = NewBizErr("50003", "修改Json错误")

	// OssErr 60-oss错误码
	OssErr         = NewBizErr("60000", "系统异常")
	OssUploadErr   = NewBizErr("60001", "oss文件上传失败")
	OssDownloadErr = NewBizErr("60002", "oss文件下载失败")
	OssFindErr     = NewBizErr("60003", "oss文件查询失败")
	OssReNameErr   = NewBizErr("60004", "oss文件重命名失败")

	// ZipErr 70-zip压缩错误
	ZipErr = NewBizErr("70001", "压缩失败")

	// PrismaGeneratedErr 80-prisma错误吗
	PrismaGeneratedErr = NewBizErr("80001", "prisma生成错误")

	// TextErr 90-文本内容错误吗
	TextErr       = NewBizErr("90001", "文本内容错误")
	TextFormatErr = NewBizErr("90002", "文本内容格式错误")
)
