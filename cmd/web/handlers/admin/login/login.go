package login

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
	"strconv"
	"webce/api/auth"
	"webce/cmd/web/handlers/admin"
	"webce/cmd/web/handlers/admin/login/request"
	"webce/pkg/library/jwt"
	"webce/pkg/library/log"
)

// HandlerLogin 登陆Handler
type HandlerLogin struct {
	admin.BaseHandler
	api *auth.AdminAuth
}

// NewHandlerLogin 实例化
func NewHandlerLogin() *HandlerLogin {
	return &HandlerLogin{
		api: &auth.AdminAuth{},
	}
}

// PostLogin  登陆获取token
func (h *HandlerLogin) PostLogin(ctx iris.Context) {
	req := request.ReqLogin{}
	err := ctx.ReadForm(&req)
	if err != nil {
		h.Error(400, "invalid request")
		return
	}
	valid := validator.New()
	err = valid.Struct(req)
	if err != nil {
		log.Log.Error("login invalid request data: ", err)
		h.Error(500, "invalid request data")
		return
	}

	login, err := h.api.Login(req.Username, req.Password, ctx.RemoteAddr())
	if err != nil {
		h.Error(500, "invalid login")
		return
	}
	h.Success(login)
}

// PostRefreshToken  刷新token
func (h *HandlerLogin) PostRefreshToken(ctx iris.Context) {
	token := ctx.GetHeader("token")
	if token == "" {
		h.Error(500, "invalid request")
		return
	}
	// 解析token
	claims, err := jwt.ParseToken(token)
	if err != nil {
		log.Log.Error("parse jwt token err:", err)
		h.Error(500, "invalid request data")
		return
	}
	// 判定token中是否存在userid
	if _, ok := claims["userId"]; !ok {
		log.Log.Error("parse jwt token exists userId")
		h.Error(500, "invalid request data")
		return
	}

	// todo 根据userid查询用户是否存在

	// 创建新token
	token, err = jwt.CreateToken(strconv.FormatUint(claims["userId"].(uint64), 10), fmt.Sprintf("%v", claims["username"]))
	if err != nil {
		log.Log.Error("create jwt token err:", err)
		h.Error(500, "refreshToken err")
		return
	}

	h.Success(token)
}
