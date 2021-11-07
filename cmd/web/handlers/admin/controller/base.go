package admin

import (
	"github.com/kataras/iris/v12"
	"webce/pkg/library/resp"
)

type BaseHandler struct {
	Ctx iris.Context
}

func (b BaseHandler) Error(code int, msg string) {
	_, err := b.Ctx.JSON(resp.ApiReturn(code, msg, nil))
	if err != nil {
		return
	}
}

func (b BaseHandler) Success(data interface{}) {
	_, err := b.Ctx.JSON(resp.ApiReturn(200, "", data))
	if err != nil {
		return
	}
}

func (b BaseHandler) Page(lists interface{}, page interface{}) {
	_, err := b.Ctx.JSON(resp.ApiReturn(200, "", iris.Map{
		"lists": lists,
		"page":  page,
	}))
	if err != nil {
		return
	}
}
