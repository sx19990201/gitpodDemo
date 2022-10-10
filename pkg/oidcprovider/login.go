package oidcprovider

import (
	"fmt"
	"github.com/fire_boom/domain"
	"github.com/fire_boom/utils"
	"github.com/gorilla/sessions"
	"github.com/ory/fosite"
	"github.com/prisma/prisma-client-go/engine"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	key         = []byte("super-secret-key")
	userSession = sessions.NewCookieStore(key)
)

func UserNameLogin(w http.ResponseWriter, r *http.Request) {
	con, _ := ioutil.ReadAll(r.Body) //获取post的数据
	username := strings.Split(string(con), "&")[0]
	password := strings.Split(string(con), "&")[1]
	username = strings.Split(username, "=")[1]
	password = strings.Split(password, "=")[1]
	//engine := utils.GetPrismaEngine()
	//schema := utils.GetQuerySchema("user", "FindByUserName", map[string]interface{}{"userName": username})
	//response, err := engine.Execute(schema)
	//
	//if strings.Contains(response, "error") || err != nil {
	//	fmt.Fprintf(w, "该用户不存在")
	//	return
	//}
	//var userInfo domain.OauthUser

	//user := gjson.Get(response, "data.result").Raw
	//err = json.Unmarshal([]byte(user), &userInfo)
	//if err != nil {
	//	fmt.Fprintf(w, "系统异常")
	//	return
	//}
	var userInfo domain.OauthUser
	schema := utils.GetQuerySchema("user", "FindByUserName", map[string]interface{}{"userName": username})

	err := engine.QuerySchema(schema, &userInfo)
	if err != nil {
		fmt.Fprintf(w, "系统异常")
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(userInfo.EncryptionPassword), []byte(password)); err != nil {
		fmt.Fprintf(w, "密码不正确")
		return
	}
	session, _ := userSession.Get(r, "user")
	session.Values["authenticated"] = true
	session.Values[userInfo.ID] = userInfo
	err = userSession.Save(r, w, session)
	// 创建授权码请求
	req, err := http.NewRequest("GET", r.Referer(), strings.NewReader(""))
	if err != nil {
		// todo
		fmt.Println(err)
	}
	ctx := r.Context()
	ar, err := oauth2.NewAuthorizeRequest(ctx, req)
	if err != nil {
		oauth2.WriteAuthorizeError(w, ar, err)
		return
	}
	// 现在用户授权已经通过,建立一个会话，用于验证/查找令牌，以及存储额外的信息
	mySessionData := &fosite.DefaultSession{
		Username: username,
		Subject:  userInfo.ID,
	}
	// 处理响应请求体
	authResponse, err := oauth2.NewAuthorizeResponse(ctx, ar, mySessionData)
	if err != nil {
		oauth2.WriteAuthorizeError(w, ar, err)
		return
	}
	// 重定向回客户端并向uri传递授权码
	oauth2.WriteAuthorizeResponse(w, ar, authResponse)
}

func UserNameRegister(w http.ResponseWriter, r *http.Request) {
	username := r.Form.Get("username")
	pswd := r.Form.Get("password")
	rePswd := r.Form.Get("rePassword")
	param := r.URL.RawQuery
	username = strings.Split(strings.Split(param, "&")[0], "=")[1]
	pswd = strings.Split(strings.Split(param, "&")[1], "=")[1]
	rePswd = strings.Split(strings.Split(param, "&")[2], "=")[1]
	if username == "" || pswd == "" || rePswd == "" {
		fmt.Fprintf(w, "用户名密码不能为空")
		return
	}
	// 包含空格，返回
	if utils.StrContainsSpace(username) || utils.StrContainsSpace(pswd) || utils.StrContainsSpace(rePswd) {
		fmt.Fprintf(w, "用户名密码不能包含空格")
		return
	}
	if pswd != rePswd {
		fmt.Fprintf(w, "密码与确认密码不一致")
		return
	}

	encryptionPswd, err := bcrypt.GenerateFromPassword([]byte(pswd), 0)
	if err != nil {
		fmt.Fprintf(w, "创建用户失败")
		return
	}
	//engine := utils.GetPrismaEngine()
	schema := utils.GetQuerySchema("user", "CreateOneUser", map[string]interface{}{"userName": username, "encryPswd": string(encryptionPswd)})
	//response, err := engine.Execute(schema)
	var userInfo domain.OauthUser
	err = engine.QuerySchema(schema, &userInfo)
	if err != nil {
		fmt.Fprintf(w, "创建用户失败")
		return
	}
	//if strings.Contains(response, "error") {
	//	fmt.Fprintf(w, "创建用户失败")
	//	return
	//}
	return
}

func EmailLogin(w http.ResponseWriter, r *http.Request) {
	// 获取邮箱和验证码
	email := r.Form.Get("email")
	code := r.Form.Get("code")
	// TODO 校验email 逻辑先空着，先查库看是邮箱否存在
	if email == "" {
		fmt.Fprintf(w, "邮箱格式不正确")
		return
	}
	// TODO 校验验证码 逻辑先空着
	if code == "" {
		fmt.Fprintf(w, "验证码不正确")
		return
	}
	// TODO 查库判断邮箱是否存在
	if email == "不存在" {
		fmt.Fprintf(w, "邮箱未注册")
		return
	}
	// TODO 拿到用户存入session
	session, _ := userSession.Get(r, email)
	session.Values["authenticated"] = true
	session.Values[email] = &email
	userSession.Save(r, w, session)
	fmt.Fprintf(w, "登陆成功")
}

func EmailRegister(w http.ResponseWriter, r *http.Request) {
	// 获取邮箱和验证码
	email := r.Form.Get("email")
	code := r.Form.Get("code")
	// TODO 校验email 逻辑先空着，先查库看是邮箱否存在
	if email == "" {
		fmt.Fprintf(w, "邮箱格式不正确")
		return
	}
	// TODO 校验验证码，不正确则返回
	if code == "" {
		fmt.Fprintf(w, "验证码不正确")
		return
	}
	// TODO 查库，看邮箱是否存在
	if email == "存在" {
		fmt.Fprintf(w, "邮箱已经存在")
		return
	}
	// TODO 邮箱用户入库
	fmt.Fprintf(w, "注册成功")
	return
}

func SMSLogin(w http.ResponseWriter, r *http.Request) {
	// 获取邮箱和验证码
	mobile := r.Form.Get("mobile")
	code := r.Form.Get("code")
	// TODO 校验手机号 逻辑先空着，先查库看是邮箱否存在
	if mobile == "" {
		fmt.Fprintf(w, "手机格式不正确")
		return
	}
	// TODO 校验验证码 逻辑先空着
	if code == "" {
		fmt.Fprintf(w, "验证码不正确")
		return
	}
	// TODO 查库判断手机号是否存在
	if mobile == "不存在" {
		fmt.Fprintf(w, "邮箱未注册")
		return
	}
	// TODO 拿到用户存入session
	session, _ := userSession.Get(r, mobile)
	session.Values["authenticated"] = true
	session.Values[mobile] = &mobile
	userSession.Save(r, w, session)
	fmt.Fprintf(w, "登陆成功")
}

func SMSRegister(w http.ResponseWriter, r *http.Request) {
	// 获取邮箱和验证码
	mobile := r.Form.Get("mobile")
	code := r.Form.Get("code")
	// TODO 校验手机号 逻辑先空着，先查库看是邮箱否存在
	if mobile == "" {
		fmt.Fprintf(w, "手机号格式不正确")
		return
	}
	// TODO 校验验证码，不正确则返回
	if code == "" {
		fmt.Fprintf(w, "验证码不正确")
		return
	}
	// TODO 查库，看邮箱是否存在
	if mobile == "存在" {
		fmt.Fprintf(w, "手机号已经存在")
		return
	}
	// TODO 手机用户入库
	fmt.Fprintf(w, "注册成功")
	return
}

func SocialLogin(w http.ResponseWriter, r *http.Request) {

}
