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
	"webce/pkg/library/resp"
)

// HandlerLogin 登陆Handler
type HandlerLogin struct {
	admin.BaseHandler
}

// NewHandlerLogin 实例化
func NewHandlerLogin() *HandlerLogin {
	return &HandlerLogin{}
}

// PostLogin  登陆获取token
func (h *HandlerLogin) PostLogin(ctx iris.Context) {
	req := request.ReqLogin{}
	err := ctx.ReadForm(&req)
	if err != nil {
		resp.Api(ctx, resp.ApiReturn(400, "invalid request", nil))
		return
	}
	valid := validator.New()
	err = valid.Struct(req)
	if err != nil {
		log.Log.Error("login invalid request data: ", err)
		resp.Api(ctx, resp.ApiReturn(500, "invalid request data", nil))
		return
	}

	model := auth.AdminAuth{}
	login, err := model.Login(req.Username, req.Password, ctx.RemoteAddr())
	if err != nil {
		resp.Api(ctx, resp.ApiReturn(500, "invalid login", nil))
		return
	}
	resp.Api(ctx, resp.ApiReturn(0, "", login))
}

// PostRefreshToken  刷新token
func (h *HandlerLogin) PostRefreshToken(ctx iris.Context) {
	token := ctx.GetHeader("token")
	if token == "" {
		resp.Api(ctx, resp.ApiReturn(400, "invalid request", nil))
		return
	}
	// 解析token
	claims, err := jwt.ParseToken(token)
	if err != nil {
		log.Log.Error("parse jwt token err:", err)
		resp.Api(ctx, resp.ApiReturn(500, "invalid request data", nil))
		return
	}
	// 判定token中是否存在userid
	if _, ok := claims["userId"]; !ok {
		log.Log.Error("parse jwt token exists userId")
		resp.Api(ctx, resp.ApiReturn(500, "invalid request data", nil))
		return
	}

	// todo 根据userid查询用户是否存在

	// 创建新token
	token, err = jwt.CreateToken(strconv.FormatUint(claims["userId"].(uint64), 10), fmt.Sprintf("%v", claims["username"]))
	if err != nil {
		log.Log.Error("create jwt token err:", err)
		resp.Api(ctx, resp.ApiReturn(500, "refreshToken err", nil))
	}

	resp.Api(ctx, resp.ApiReturn(0, "", token))
}
