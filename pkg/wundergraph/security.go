package wundergraph

type WdgSecurityConfig struct {
	EnableGraphQLEndpoint bool     `json:"enableGraphQLEndpoint,omitempty,optional"` // GraphQL端点 0-关 1-开
	AllowedHosts          []string `json:"allowedHosts,omitempty,optional"`          // 允许主机,多个域名
}
