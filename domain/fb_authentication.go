package domain

import (
	"context"
	"database/sql"
)

const (
	OpenID = "openid"
	Github = "github"
	Google = "google"
)

// FbAuthentication 身份验证配置
type FbAuthentication struct {
	ID           uint           `db:"id" json:"id"`
	Name         string         `db:"name" json:"name"`                                          // 身份验证名称
	AuthSupplier string         `db:"authSupplier" json:"authSupplier,default=openid,omitempty"` // 验证供应商:openid、github、google
	SwitchState  string         `db:"switchState" json:"switchState"`                            // 开关状态: 0-全部关闭 1-cookie 2-token 3-全部cookie、token都开启
	Config       string         `db:"config" json:"config"`                                      // 身份验证配置对应的配置项：供应商id、appID、appSecret、服务发现地址、重定向url等
	CreateTime   sql.NullString `db:"create_time" json:"-"`
	UpdateTime   sql.NullString `db:"update_time" json:"-"`
	IsDel        uint8          `db:"isDel" json:"-"`
}

// FbAuthenticationResp 身份验证配置
type FbAuthenticationResp struct {
	ID           uint           `db:"id" json:"id"`
	Name         string         `db:"name" json:"name"`                 // 身份验证名称
	AuthSupplier string         `db:"authSupplier" json:"authSupplier"` // 验证供应商:openid、github、google
	SwitchState  []string       `db:"switchState" json:"switchState"`   // 开关状态: 0-全部关闭 1-cookie 2-token 3-全部cookie、token都开启
	Config       WdgAuthConfig  `db:"config" json:"config"`             // 身份验证配置对应的配置项：供应商id、appID、appSecret、服务发现地址、重定向url等
	CreateTime   sql.NullString `db:"create_time" json:"-"`
	UpdateTime   sql.NullString `db:"update_time" json:"-"`
	IsDel        uint8          `db:"isDel" json:"-"`
}

type WdgAuthConfig struct {
	ID                      string `json:"id"`                      // 供应商id
	ClientId                string `json:"clientId"`                // appid
	ClientSecret            string `json:"clientSecret"`            // appSecret
	Issuer                  string `json:"issuer"`                  // issuer端点
	DiscoveryURL            string `json:"discoveryURL"`            // 服务发现
	Jwks                    int64  `json:"jwks"`                    // 类型 0-url 1-json
	JwksJSON                string `json:"jwksJSON"`                // jwksjson
	JwksURL                 string `json:"jwksURL"`                 // jwksurl
	UserInfoEndpoint        string `json:"userInfoEndpoint"`        // 用户信息端点
	UserInfoCacheTtlSeconds int64  `json:"userInfoCacheTtlSeconds"` // 用户信息缓存时间
}

type AuthenticationUseCase interface {
	Store(ctx context.Context, auth *FbAuthentication) (int64, error)
	Update(ctx context.Context, auth *FbAuthentication) (int64, error)
	Delete(ctx context.Context, id uint) (int64, error)
	FindAuthentication(ctx context.Context) ([]FbAuthentication, error)
}

type AuthenticationRepository interface {
	Store(ctx context.Context, auth *FbAuthentication) (int64, error)
	Update(ctx context.Context, auth *FbAuthentication) (int64, error)
	GetByName(ctx context.Context, name string) (FbAuthentication, error)
	CheckExist(ctx context.Context, auth *FbAuthentication) (FbAuthentication, error)
	Delete(ctx context.Context, id uint) (int64, error)
	FindAuthentication(ctx context.Context) ([]FbAuthentication, error)
}
