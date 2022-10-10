package oidcprovider

import (
	"encoding/json"
	"fmt"
	"github.com/fire_boom/domain"
	"github.com/fire_boom/utils"
	"github.com/ory/fosite"
	"github.com/prisma/prisma-client-go/engine"
	log "github.com/sirupsen/logrus"
	"net/http"
	"reflect"
	"strings"
)

var jwksJson = `{
  "issuer": "http://localhost:3846/oauth2",
  "authorization_endpoint": "http://localhost:3846/oauth2/auth",
  "token_endpoint": "http://localhost:3846/oauth2/token",
  "jwks_uri": "http://localhost:3846/oauth2/.well-known/jwks.json",
  "end_session_endpoint": "http://localhost:3846/oauth2/session/end",
  "userinfo_endpoint": "http://localhost:3846/oauth2/me",
  "introspection_endpoint": "http://localhost:3846/oauth2/introspection",
  "revocation_endpoint": "http://localhost:3846/oauth2/revocation",
  "claims_parameter_supported": false,
  "claims_supported": [
    "sub",
    "username",
    "phone_number",
    "phone_number_verified",
    "email",
    "email_verified",
    "address",
    "birthdate",
    "family_name",
    "gender",
    "given_name",
    "locale",
    "middle_name",
    "name",
    "nickname",
    "picture",
    "preferred_username",
    "profile",
    "updated_at",
    "website",
    "zoneinfo",
    "role",
    "roles",
    "unionid",
    "external_id",
    "extended_fields",
    "tenant_id",
    "userpool_id",
    "sid",
    "auth_time",
    "iss"
  ],
  "code_challenge_methods_supported": [
    "plain",
    "S256"
  ],
  "grant_types_supported": [
    "authorization_code",
    "password",
    "refresh_token"
  ],
  "response_types_supported": [
    "code"
  ],
  "response_modes_supported": [
    "query",
    "fragment",
    "form_post",
    "web_message"
  ],
  "scopes_supported": [
    "openid",
    "offline_access",
    "username",
    "phone",
    "email",
    "address",
    "profile",
    "role",
    "roles",
    "unionid",
    "external_id",
    "extended_fields",
    "tenant_id",
    "userpool_id"
  ],
  "token_endpoint_auth_methods_supported": [
    "client_secret_post",
    "client_secret_basic",
    "none"
  ],
  "request_parameter_supported": false,
  "request_uri_parameter_supported": false,
  "userinfo_signing_alg_values_supported": [
    "HS256",
    "RS256"
  ],
  "introspection_endpoint_auth_methods_supported": [
    "client_secret_post",
    "client_secret_basic",
    "none"
  ],
  "revocation_endpoint_auth_methods_supported": [
    "client_secret_post",
    "client_secret_basic",
    "none"
  ],
  "id_token_encryption_alg_values_supported": [
    "A128KW",
    "A256KW",
    "ECDH-ES",
    "ECDH-ES+A128KW",
    "ECDH-ES+A256KW",
    "RSA-OAEP"
  ],
  "id_token_encryption_enc_values_supported": [
    "A128CBC-HS256",
    "A128GCM",
    "A256CBC-HS512",
    "A256GCM"
  ],
  "userinfo_encryption_alg_values_supported": [
    "A128KW",
    "A256KW",
    "ECDH-ES",
    "ECDH-ES+A128KW",
    "ECDH-ES+A256KW",
    "RSA-OAEP"
  ],
  "userinfo_encryption_enc_values_supported": [
    "A128CBC-HS256",
    "A128GCM",
    "A256CBC-HS512",
    "A256GCM"
  ],
  "claim_types_supported": [
    "normal"
  ],
  "subject_types_supported": [
    "public"
  ],
  "id_token_signing_alg_values_supported": [
    "HS256",
    "RS256"
  ]
}`

func RegisterHandlers() {
	// openid 授权模式
	http.HandleFunc("/oauth2/.well-known/openid-configuration", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, jwksJson)
	})

	//http://localhost:3846/oauth2/session/end
	http.HandleFunc("/oauth2/session/end", sessionEnd)
	http.HandleFunc("/oauth2/me", getUserInfo)

	// Set up oauth2 endpoints. You could also use gorilla/mux or any other router.
	http.HandleFunc("/oauth2/auth", authorizeHandlerFunc)
	http.HandleFunc("/oauth2/token", tokenHandlerFunc)

	// revoke tokens
	http.HandleFunc("/oauth2/revoke", revokeEndpoint)
	http.HandleFunc("/oauth2/introspect", introspectionEndpoint)

	http.HandleFunc("/oauth2/login/username", UserNameLogin)
	http.HandleFunc("/oauth2/login/sms", SMSLogin)
	http.HandleFunc("/oauth2/login/social", SocialLogin)
	http.HandleFunc("/oauth2/login/email", EmailLogin)
	http.HandleFunc("/oauth2/register/username", UserNameRegister)
	http.HandleFunc("/oauth2/register/sms", SMSRegister)
	http.HandleFunc("/oauth2/register/email", EmailRegister)
}

func sessionEnd(w http.ResponseWriter, r *http.Request) {
	username := r.Form.Get("username")
	session, _ := userSession.Get(r, username)

	// Revoke users authentication
	session.Values["authenticated"] = false
	session.Save(r, w)
}

// Claims 包含除oidc协议规定的Claims外，用户也可以自定义Claims

func GetProfileClaims(user domain.MateData) (result map[string]interface{}) {
	result = make(map[string]interface{}, 0)
	result["name"] = user.Name
	result["family_name"] = user.FamilyName
	result["given_name"] = user.GivenName
	result["middle_name"] = user.MiddleName
	result["nickname"] = user.Nickname
	result["preferred_username"] = user.PreferredUsername
	result["profile"] = user.Profile
	result["picture"] = user.Picture
	result["website"] = user.Website
	result["gender"] = user.Gender
	result["birthdate"] = user.Birthdate
	result["zoneinfo"] = user.ZoneInfo
	result["locale"] = user.Locale
	result["updated_at"] = user.UpdatedAt
	return
}
func GetEmailClaims(user domain.MateData) (result map[string]interface{}) {
	result = make(map[string]interface{}, 0)
	result["email"] = user.Email
	result["email_verified"] = user.EmailVerified
	return
}
func GetAddressClaims(user domain.MateData) (result map[string]interface{}) {
	result = make(map[string]interface{}, 0)
	result["address"] = user.Address
	return
}
func GetPhoneClaims(user domain.MateData) (result map[string]interface{}) {
	result = make(map[string]interface{}, 0)
	result["phone_number"] = user.PhoneNumber
	result["phone_number_verified"] = user.PhoneNumberVerified
	return
}

func getUserInfo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	_, ar, err := oauth2.IntrospectToken(ctx, fosite.AccessTokenFromRequest(r), fosite.AccessToken, new(fosite.DefaultSession))
	if err != nil {
		fmt.Fprintf(w, "系统异常")
		return
	}
	// 获取用户session
	session := ar.GetSession()
	// 获取用户id
	id := session.GetSubject()
	//engine := utils.GetPrismaEngine()
	schema := utils.GetQuerySchema("user", "GetByID", map[string]interface{}{"id": id})

	//userSchema, err := engine.Execute(schema)
	var userInfo domain.OauthUser
	err = engine.QuerySchema(schema, &userInfo)
	if err != nil {
		fmt.Fprintf(w, "该用户不存在")
		return
	}

	//user := gjson.Get(userSchema, "data.result").Raw
	//err = json.Unmarshal([]byte(user), &userInfo)
	//if err != nil {
	//	fmt.Fprintf(w, "用户信息获取失败")
	//	return
	//}
	mateData := domain.MateData{}
	err = json.Unmarshal([]byte(userInfo.MateData), mateData)
	if err != nil {
		fmt.Fprintf(w, "用户信息获取失败")
		return
	}
	scopes := ar.GetRequestedScopes()
	claimsMap := make(map[string]interface{})
	claimsMap["sub"] = userInfo.ID
	if scopes.Has(domain.ProfileScope) {
		for k, v := range GetProfileClaims(mateData) {
			claimsMap[k] = v
		}
	}
	if scopes.Has(domain.EmailScope) {
		for k, v := range GetEmailClaims(mateData) {
			claimsMap[k] = v
		}
	}
	if scopes.Has(domain.AddressScope) {
		for k, v := range GetAddressClaims(mateData) {
			claimsMap[k] = v
		}
	}
	if scopes.Has(domain.PhoneScope) {
		for k, v := range GetPhoneClaims(mateData) {
			claimsMap[k] = v
		}
	}
	result := ""
	for k, v := range claimsMap {
		val := reflect.TypeOf(v)
		if val.Kind() == reflect.String {
			result = fmt.Sprintf(`%s "%v":"%v",`, result, k, v)
			continue
		}
		result = fmt.Sprintf(`%s "%v":%v,`, result, k, v)
	}
	result = strings.TrimRight(result, ",")
	result = fmt.Sprintf(`{%s}`, result)
	fmt.Fprintf(w, result)
}

// authorizeHandlerFunc 授权端点
func authorizeHandlerFunc(rw http.ResponseWriter, req *http.Request) {
	// 贯穿所有的上下文，用来做类似终止数据库等类似的操作
	ctx := req.Context()
	// 创建一个authRequest对象，用来分析请求并提取重要信息，如范围、响应类型等
	ar, err := oauth2.NewAuthorizeRequest(ctx, req)
	if err != nil {
		oauth2.WriteAuthorizeError(rw, ar, err)
		return
	}
	// TODO 判断用户是否登陆
	userCookie, err := req.Cookie("user")
	if err != nil || userCookie == nil {
		// TODO 跳转到登录页面
		rw.Header().Set("Content-Type", "text/html;charset=UTF-8")
		rw.Write([]byte(`<h1>Login page</h1>`))
		rw.Write([]byte(`
			<form action="http://localhost:3846/oauth2/login/username" method="POST">
				<input type="text" name="username" /><br>
				<input type="text" name="password" /><br>
				<input type="text" name="client_id"  value="my-client"/><br>
				<input type="submit">
			</form>
		`))
		return
	}
	userName := userCookie.Name
	session, _ := userSession.Get(req, userName)
	session.Name()
	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		// 将信息存入session
		//cookieSession, _ := userSession.Get(req, "cookie")
		//cookieSession.Values["state"] = ar.GetState()
		//cookieSession.Values["redirect_uri"] = ar.GetRedirectURI().String()
		//cookieSession.Values["success_redirect_uri"] = req.Header.Get("Referer")
		//rw.Header().Set("session-id", cookieSession.ID)
		// TODO 跳转到登录页面
		rw.Header().Set("Content-Type", "text/html;charset=UTF-8")
		rw.Write([]byte(`<h1>Login page</h1>`))
		rw.Write([]byte(`
			<form action="http://localhost:3846/oauth2/login/username" method="POST">
				<input type="text" name="username" /><br>
				<input type="text" name="password" /><br>
				<input type="text" name="client_id"  value="my-client"/><br>
				<input type="submit">
			</form>
		`))
		return
	}

	// 现在用户授权已经通过,建立一个会话，用于验证/查找令牌，以及存储额外的信息
	mySessionData := &fosite.DefaultSession{
		Username: userName,
	}

	// It's also wise to check the requested scopes, e.g.:
	// if authorizeRequest.GetScopes().Has("admin") {
	//     http.Error(rw, "you're not allowed to do that", http.StatusForbidden)
	//     return
	// }
	// 处理响应请求体
	response, err := oauth2.NewAuthorizeResponse(ctx, ar, mySessionData)
	if err != nil {
		oauth2.WriteAuthorizeError(rw, ar, err)
		return
	}
	// 重定向回客户端并向uri传递授权码
	oauth2.WriteAuthorizeResponse(rw, ar, response)
}

// token 端点
func tokenHandlerFunc(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	mySessionData := new(fosite.DefaultSession)
	// 创建一个请求体，并遍历注册的TokenEndpointHandlers以验证请求
	accessRequest, err := oauth2.NewAccessRequest(ctx, req, mySessionData)
	if err != nil {
		oauth2.WriteAccessError(rw, accessRequest, err)
		return
	}

	// 如果是超级管理员
	if mySessionData.Username == "super-admin-guy" {
		// do something...
	}
	// 简历一个响应请求，并在 次响应体聚合结果
	response, err := oauth2.NewAccessResponse(ctx, accessRequest)
	response.SetTokenType("Bearer")
	//response.SetScopes(accessRequest.GetRequestedScopes())
	if err != nil {
		oauth2.WriteAccessError(rw, accessRequest, err)
		return
	}

	// 上面处理全部完成后,发送响应给客户端
	oauth2.WriteAccessResponse(rw, accessRequest, response)

	// The client has a valid access token now

}

// revokeEndpoint 撤销token
func revokeEndpoint(rw http.ResponseWriter, req *http.Request) {
	// This context will be passed to all methods.
	ctx := req.Context()

	// This will accept the token revocation request and validate various parameters.
	err := oauth2.NewRevocationRequest(ctx, req)

	// All done, send the response.
	oauth2.WriteRevocationResponse(rw, err)
}

// introspectionEndpoint 内省端点
func introspectionEndpoint(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	mySessionData := newSession("")
	ir, err := oauth2.NewIntrospectionRequest(ctx, req, mySessionData)
	if err != nil {
		log.Printf("Error occurred in NewIntrospectionRequest: %+v", err)
		oauth2.WriteIntrospectionError(rw, err)
		return
	}

	oauth2.WriteIntrospectionResponse(rw, ir)
}
