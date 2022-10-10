package domain

type GlobalConfig struct {
	AuthRedirectURL                 []string                           `json:"authorizedRedirectUris"` // 身份鉴权-重定向url配置
	Hooks                           HooksConfiguration                 `json:"hooksConfiguration"`
	ConfigureWunderGraphApplication WunderGraphConfigApplicationConfig `json:"configureWunderGraphApplication"`
}

type HooksConfiguration struct {
	Authentication Authentication `json:"authentication"`
	RestApi        RestApi        `json:"restApi"`
}

type Authentication struct {
	PostAuthenticationSwitch         bool `json:"postAuthenticationSwitch"`
	MutatingPostAuthenticationSwitch bool `json:"mutatingPostAuthenticationSwitch"`
}

type RestApi struct {
	OnRequestSwitch  bool `json:"onRequestSwitch"`
	OnResponseSwitch bool `json:"onResponseSwitch"`
}

type WunderGraphConfigApplicationConfig struct {
	Security *SecurityConfig    `json:"security"`
	Cors     *CorsConfiguration `json:"cors"`
}

type SecurityConfig struct {
	EnableGraphQLEndpoint bool     `json:"enableGraphQLEndpoint"` // GraphQL端点 0-关 1-开
	AllowedHosts          []string `json:"allowedHosts"`          // 允许主机,多个域名
}

type CorsConfiguration struct {
	AllowedOrigins   []string `json:"allowedOrigins"`   // 允许域名
	AllowedMethods   []string `json:"allowedMethods"`   // 允许方法 0-* 1-GET 2-POST 3-PUT
	AllowedHeaders   []string `json:"allowedHeaders"`   // 请求头部
	ExposedHeaders   []string `json:"exposedHeaders"`   // 排除头部
	MaxAge           int64    `json:"maxAge"`           // 跨域时间(s)
	AllowCredentials bool     `json:"allowCredentials"` // 允许证书开关 0-开 1-关
}
