package admin

import (
	"github.com/kataras/iris/v12"
	"webce/pkg/lib"
	"webce/pkg/library/page"
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
func (b BaseHandler) Msg(message string) {
	_, err := b.Ctx.JSON(resp.ApiReturn(200, message, nil))
	if err != nil {
		return
	}
}

func (b BaseHandler) Page(lists interface{}, page *page.PageResponse) {
	_, err := b.Ctx.JSON(resp.ApiReturn(200, "", iris.Map{
		"lists": lists,
		"page":  page,
	}))
	if err != nil {
		return
	}
}

func (b BaseHandler) Validate(u interface{}) error {
	return lib.Validate(u)
}
