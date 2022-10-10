package domain

import "context"

type Home struct {
	DataSource HomeDataSource `json:"homeDataSource"`
	Api        HomeApi        `json:"homeApi"`
	Oss        HomeOss        `json:"homeOss"`
	Auth       HomeAuth       `json:"homeAuth"`
}

type HomeBulletin struct {
	BulletinType int64  `json:"bulletinType"`
	Title        string `json:"title"`
	Date         string `json:"date"`
}

type HomeDataSource struct {
	DBTotal       int64 `json:"dbTotal"`       // db个数
	RestTotal     int64 `json:"RestTotal"`     // rest个数
	GraphqlTotal  int64 `json:"GraphqlTotal"`  // graphql个数
	CustomerTotal int64 `json:"CustomerTotal"` // 自定义个数
}

type HomeApi struct {
	QueryTotal         int64 `json:"queryTotal"`         // 查询个数
	LiveQueryTotal     int64 `json:"liveQueryTotal"`     // 实时查询个数
	MutationsTotal     int64 `json:"mutationsTotal"`     // 更改个数
	SubscriptionsTotal int64 `json:"subscriptionsTotal"` // 订阅个数
}

type HomeOss struct {
	OssTotal    int64  `json:"ossTotal"`    // oss存储个数
	TotalMemory string `json:"totalMemory"` // 总内存
	UseMemory   string `json:"useMemory"`   // 已使用内存
}

type HomeAuth struct {
	AuthTotal       int64 `json:"authTotal"`       // 累计身份验证商
	TotalUser       int64 `json:"totalUser"`       // 累计用户
	TodayInsertUser int64 `json:"todayInsertUser"` // 今日新增用户
}

type HomeUseCase interface {
	GetDateSourceData(c context.Context) (result HomeDataSource, err error)
	GetApiData(c context.Context) (result HomeApi, err error)
	GetOssData(c context.Context) (result HomeOss, err error)
	GetAuthData(c context.Context) (result HomeAuth, err error)
}
