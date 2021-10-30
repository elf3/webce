package login

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
	"webce/api/services"
	"webce/cmd/web/handlers/admin/controller/login/requests"
	"webce/library/apgs"
	"webce/library/log"
	"webce/library/session"
	"webce/pkg/lib"
)

func ShowLogin(ctx iris.Context) {
	err := ctx.View("login/login.html")
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
	login, err := auth.Login(req.Username, req.Password, addr)
	if err != nil {
		apgs.Api(ctx, apgs.ApiReturn(500, "invalid login", nil))
		return
	}
	marshal, _ := json.Marshal(login)
	session.SetSession(ctx, "userID", lib.BytesString(marshal))
	apgs.Api(ctx, apgs.ApiReturn(0, "", login))
}
