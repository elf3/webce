package admin

import (
	"github.com/kataras/iris/v12"
	"webce/apis/services"
)

func Main(ctx iris.Context) {

	err := ctx.View("login.html")
	if err != nil {
		return
	}
}

func Login(ctx iris.Context) {
	username := ctx.PostValueDefault("username", "123")
	password := ctx.PostValueDefault("password", "123")
	addr := ctx.RemoteAddr()
	auth := services.AdminAuth{}
	login := auth.Login(username, password, addr)
	_, err := ctx.JSON(login)
	if err != nil {
		return
	}
}
