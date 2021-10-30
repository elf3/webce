package login

import (
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
	"webce/api/services"
	"webce/cmd/web/handlers/admin/controller/login/requests"
	"webce/pkg/library/apgs"
	"webce/pkg/library/log"
)

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
	login, err := auth.Login(req.Username, req.Password, addr)
	if err != nil {
		apgs.Api(ctx, apgs.ApiReturn(500, "invalid login", nil))
		return
	}
	apgs.Api(ctx, apgs.ApiReturn(0, "", login))
}
