package resp

import (
	"github.com/kataras/iris/v12"
)

func Error(c iris.Context, code int, msg string) {
	_, err := c.JSON(ApiReturn(code, msg, nil))
	if err != nil {
		return
	}
}
func Success(c iris.Context, data interface{}) {
	_, err := c.JSON(ApiReturn(200, "", data))
	if err != nil {
		return
	}
}

func ApiJson(c iris.Context, code int, msg string, data interface{}) {
	_, err := c.JSON(ApiReturn(code, msg, data))
	if err != nil {
		return
	}
}
func Page(c iris.Context, lists interface{}, page interface{}) {
	ApiJson(c, 200, "", iris.Map{
		"lists": lists,
		"page":  page,
	})
}

func Json(c iris.Context, data interface{}) {
	_, err := c.JSON(data)
	if err != nil {
		return
	}
}
