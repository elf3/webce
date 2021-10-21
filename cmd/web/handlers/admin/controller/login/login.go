package login

import (
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
	"webce/apis/services"
	"webce/cmd/web/handlers/admin/controller/login/requests"
	"webce/library/apgs"
	"webce/library/log"
)

func Main(ctx iris.Context) {
	err := ctx.View("login.html")
	if err != nil {
		return
	}
}

func Login(ctx iris.Context) {
	req := requests.ReqLogin{}
	err := ctx.ReadForm(&req)
	if err != nil {
		apgs.Api(ctx, apgs.ApiReturn(400, "invalid request", nil))
		return
	}
	valid := validator.New()
	err = valid.Struct(req)
	if err != nil {
		log.Log.Error("login invalid request data: ", err)
		apgs.Api(ctx, apgs.ApiReturn(500, "invalid request data", nil))
		return
	}

	addr := ctx.RemoteAddr()
	auth := services.AdminAuth{}
	login := auth.Login(req.Username, req.Password, addr)
	apgs.Api(ctx, login)
}
