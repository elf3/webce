package admin

import (
	"github.com/kataras/iris/v12"
	"webce/library/apgs"
)

type BaseHandler struct {
	Ctx iris.Context
}

func (b BaseHandler) ApiJson(code int, msg string, data interface{}) {
	_, err := b.Ctx.JSON(apgs.ApiReturn(code, msg, data))
	if err != nil {
		return
	}
}

func (b BaseHandler) Api(data interface{}) {
	_, err := b.Ctx.JSON(data)
	if err != nil {
		return
	}
}
