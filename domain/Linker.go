package domain

type AuthHome struct {
	TotalUser    int64    `json:"totalUser"`       // 总用户
	TodayUser    UserData `json:"todayInsertUser"` // 当日用户信息
	SevenDayUser UserData `json:"sevenDayUser"`    // 7日信息
}

type Active struct {
	ThirtyActive []DayActive `json:"thirtyActive"` // 近30天活跃数据
	WeekActive   UserData    `json:"weekActive"`   // 周活跃信息
	MonthActive  UserData    `json:"monthActive"`  // 月活跃信息
	DayActive    UserData    `json:"dayActive"`    // 日活跃信息
}

type DayActive struct {
	DataString string `json:"dataString"` // 日期
	Count      int64  `json:"count"`      // 人数
}

type UserData struct {
	Count   int64 `json:"count"`   // 新增数量
	UpCount int64 `json:"upCount"` // 对比前个时间段新增数量
}

type Linker struct {
	ID             string      `json:"id"`             // 链接器id 比如alipay-web, alipay-Native等等，该字段是通过target和platform拼接得到的
	Enabled        bool        `json:"enabled"`        // 是否创建 true-已创建 false-未创建
	Config         interface{} `json:"config"`         // 配置的数据结构,每个链接器配置的结构都不一样
	CreatedAt      int64       `json:"createdAt"`      // 创建时间，时间戳
	Target         string      `json:"target"`         // 标签,如:alipay、weChat、QQ等
	Types          string      `json:"types"`          // 类型，如:Social和SMS
	Platform       string      `json:"platform"`       // 平台类型，如Web、Native等
	Name           string      `json:"name"`           // 名称 如支付宝，微信等
	Logo           string      `json:"logo"`           // logo的url
	LogoDark       string      `json:"logoDark"`       // 深色模式下的logo
	Description    string      `json:"description"`    // 描述
	ConfigTemplate string      `json:"configTemplate"` // 配置模板,用于渲染配置页面 生产的配置对应config字段
}

type LoginBrandConfig struct {
	Color    Color    `json:"color"`
	Branding Branding `json:"branding"`
}
type Color struct {
	PrimaryColor      string `json:"primaryColor"`
	DarkPrimaryColor  string `json:"darkPrimaryColor"`
	IsDarkModeEnabled bool   `json:"isDarkModeEnabled"`
}
type Branding struct {
	Style       string `json:"style"`
	Slogan      string `json:"slogan"`
	LogoUrl     string `json:"logoUrl"`
	DarkLogoUrl string `json:"darkLogoUrl"`
}

type LoginConfig struct {
	SignInMethods                SignInMethods `json:"signInMethods"`
	SocialSignInConnectorTargets []string      `json:"socialSignInConnectorTargets"`
	Socials                      []Social      `json:"socials"`
	SignInMode                   string        `json:"signInMode,default=SignInAndRegister"`
}
type Social struct {
	Name     string `json:"name"`
	Logo     string `json:"logo"`
	Platform string `json:"platform"`
}
type SignInMethods struct {
	Sms      string `json:"sms"`
	Email    string `json:"email"`
	Social   string `json:"social"`
	Username string `json:"username"`
}

type LoginOtherConfig struct {
	Enabled    bool   `json:"enabled"`
	ContentUrl string `json:"contentUrl"`
}
