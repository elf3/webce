package lib

import (
	"github.com/kataras/iris/v12"
	"webce/pkg/library/apgs"
)

func ErrJson(c iris.Context, code int, msg string) {
	_, err := c.JSON(apgs.ApiReturn(code, msg, nil))
	if err != nil {
		return
	}
}

func MJson(c iris.Context, code int, msg string, data interface{}) {
	_, err := c.JSON(apgs.ApiReturn(code, msg, data))
	if err != nil {
		return
	}
}

func Json(c iris.Context, data interface{}) {
	_, err := c.JSON(data)
	if err != nil {
		return
	}
}
