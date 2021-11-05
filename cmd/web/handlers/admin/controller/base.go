package admin

import (
	"github.com/kataras/iris/v12"
	"webce/pkg/library/apgs"
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

func (b BaseHandler) ApiError(code int, msg string) {
	_, err := b.Ctx.JSON(apgs.ApiReturn(code, msg, nil))
	if err != nil {
		return
	}
}

func (b BaseHandler) ErrorRequest() {
	_, err := b.Ctx.JSON(apgs.ApiReturn(500, "invalid request", nil))
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
