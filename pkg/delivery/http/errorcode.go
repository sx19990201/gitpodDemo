package http

// 错误码前三位功能code
//const (
//	AuthCode          = "101"
//	DataSourceCode    = "102"
//	DraftCode         = "103"
//	FileCode          = "104"
//	OperateAPICode    = "105"
//	PrismaCode        = "106"
//	RoleCode          = "107"
//	S3UploadCode      = "108"
//	ScriptCode        = "109"
//	SettingCode       = "110"
//	StorageBucketCode = "111"
//	WunderGlobalCode  = "112"
//)

// 中间两位类型码
//const (
//	ParamsErrCode  = "01" // 参数错误
//	NetworkErrCode = "02" // 网络错误
//	DbErrCode      = "03" // 数据库错误
//	FileErrCode    = "04" // 文件操作错误
//	JsonErrCode    = "05" // json操作错误
//	OssErrCode     = "06" // oss操作错误
//)

// 参数后三位错误码
//const (
//	// BindErrCode 后三位错误码
//	BindErrCode  = "001"
//	EmptyErrCode = "002"
//)

// 数据库错误码
//const (
//	FindErrCode   = "001"
//	UpdateErrCode = "002"
//	CreateErrCode = "003"
//	DeleteErrCode = "004"
//)
//
//// 文件操作错误码
//const (
//	ReadErrCode       = "001" // 文件读取错误
//	WriteErrCode      = "002" // 文件写入错误
//	RemoveErrCode     = "003" // 文件删除错误
//	CreateFileErrCode = "004" // 文件创建错误
//	ReNameFileErrCode = "005" // 文件重命名错误
//	FileOpenErrCode   = "006" // 文件打开错误
//	FileCopyErrCode   = "007" // 文件拷贝错误
//)

// Json操作错误码
//const (
//	MarshalErrCode   = "001" // Json序列化错误
//	UnmarshalErrCode = "002" // Json反序列化错误
//	SetJsonErrCode   = "003" // 修改Json错误
//)
//
//// Oss操作错误码
//const (
//	UploadErrCode      = "001" // oss文件上传错误
//	DownloadErrCode    = "002" // oss文件下载错误
//	OssFindFileErrCode = "003" // oss文件查询错误
//	OssReNameErrCode   = "004" // oss文件重命名错误
//	OssDeleteErrCode   = "005" // oss文件重命名错误
//)
//
//func GetResponseCode(moduleCode, typeCode, errCode string) string {
//	return fmt.Sprintf("%s%s%s", moduleCode, typeCode, errCode)
//}
//
//var errTypeMap = map[string]string{
//	ParamsErrCode:  "参数异常",
//	NetworkErrCode: "网络异常",
//	DbErrCode:      "内部异常",
//	FileErrCode:    "系统错误",
//	JsonErrCode:    "系统错误",
//	OssErrCode:     "网络异常",
//}

//var errMap = map[string]map[string]string{
//	ParamsErrCode: {
//		BindErrCode:  "参数异常",
//		EmptyErrCode: "参数异常",
//	},
//	NetworkErrCode: {},
//	DbErrCode: {
//		FindErrCode:   "内部错误",
//		UpdateErrCode: "内部错误",
//		CreateErrCode: "内部错误",
//		DeleteErrCode: "内部错误",
//	},
//	FileErrCode: {
//		ReadErrCode:       "内部错误",
//		WriteErrCode:      "内部错误",
//		RemoveErrCode:     "内部错误",
//		CreateFileErrCode: "内部错误",
//		ReNameFileErrCode: "内部错误",
//		FileOpenErrCode:   "内部错误",
//		FileCopyErrCode:   "内部错误",
//	},
//	JsonErrCode: {
//		MarshalErrCode:   "内部错误",
//		UnmarshalErrCode: "内部错误",
//		SetJsonErrCode:   "内部错误",
//	},
//	OssErrCode: {
//		UploadErrCode:      "系统异常",
//		DownloadErrCode:    "系统异常",
//		OssFindFileErrCode: "系统异常",
//		OssReNameErrCode:   "系统异常",
//		OssDeleteErrCode:   "系统异常",
//	},
//}
//
//func GetErrMsg(typeCode, detailCode string) string {
//	return errMap[typeCode][detailCode]
//}
//
//func GetErrTypeMsg(code string) string {
//	return errTypeMap[code]
//}
//
//type Response struct {
//	Code int         `json:"code"`   // 错误码
//	Msg  string      `json:"msg"`    // 错误描述
//	Data interface{} `json:"result"` // 返回数据
//}

//// WithMsg 自定义响应信息
//func (res *Response) WithMsg(message string) Response {
//	return Response{
//		Code: res.Code,
//		Msg:  message,
//		Data: res.Data,
//	}
//}
//
//// WithData 追加响应数据
//func (res *Response) WithData(data interface{}) Response {
//	return Response{
//		Code: res.Code,
//		Msg:  res.Msg,
//		Data: data,
//	}
//}
//
//// ToString 返回 JSON 格式的错误详情
//func (res *Response) ToString() string {
//	err := &struct {
//		Code int         `json:"code"`
//		Msg  string      `json:"msg"`
//		Data interface{} `json:"data"`
//	}{
//		Code: res.Code,
//		Msg:  res.Msg,
//		Data: res.Data,
//	}
//	raw, _ := json.Marshal(err)
//	return string(raw)
//}
//
//// 构造函数
//func response(code int, msg string) *Response {
//	return &Response{
//		Code: code,
//		Msg:  msg,
//		Data: nil,
//	}
//}
