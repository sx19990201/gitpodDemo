package wundergraph

type WdgCors struct {
	AllowedOrigins   []string `json:"allowedOrigins,omitempty,optional"`   // 允许域名
	AllowedMethods   []string `json:"allowedMethods,omitempty,optional"`   // 允许方法 0-* 1-GET 2-POST 3-PUT
	AllowedHeaders   []string `json:"allowedHeaders,omitempty,optional"`   // 请求头部
	ExposedHeaders   []string `json:"exposedHeaders,omitempty,optional"`   // 排除头部
	MaxAge           int64    `json:"maxAge,omitempty,optional"`           // 跨域时间(s)
	AllowCredentials bool     `json:"allowCredentials,omitempty,optional"` // 允许证书开关 0-开 1-关
}
