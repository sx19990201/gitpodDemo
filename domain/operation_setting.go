package domain

type BaseOperationConfiguration struct {
}

type OperationSetting struct {
	AuthenticationRequired          bool  `json:"authenticationRequired"`          // 需要授权
	CachingEnable                   bool  `json:"cachingEnable"`                   // 开启缓存
	CachingMaxAge                   int64 `json:"cachingMaxAge"`                   // 最大时长
	CachingStaleWhileRevalidate     int64 `json:"cachingStaleWhileRevalidate"`     // 重校验时长
	LiveQueryEnable                 bool  `json:"liveQueryEnable"`                 // 开启实时
	LiveQueryPollingIntervalSeconds int64 `json:"liveQueryPollingIntervalSeconds"` // 轮询间隔
}
