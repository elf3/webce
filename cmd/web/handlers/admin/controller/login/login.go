package login

import (
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
	"webce/api/auth"
	"webce/cmd/web/handlers/admin/controller/login/request"
	"webce/pkg/library/log"
	"webce/pkg/library/resp"
)

func Login(ctx iris.Context) {
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

	addr := ctx.RemoteAddr()
	auth := auth.AdminAuth{}
	login, err := auth.Login(req.Username, req.Password, addr)
	if err != nil {
		resp.Api(ctx, resp.ApiReturn(500, "invalid login", nil))
		return
	}
	resp.Api(ctx, resp.ApiReturn(0, "", login))
}
