package resp

import (
	"github.com/kataras/iris/v12"
)

// Response 返回数据结构体
type Response struct {
	Code int         `json:"code"` // 返回码
	Msg  string      `json:"msg"`  // 返回消息
	Data interface{} `json:"data"` // 返回数据
}

// Redirect 重定向Json结构体
type Redirect struct {
	Code int    `json:"code"` // 返回码
	Msg  string `json:"msg"`  //返回消息
	Url  string `json:"url"`  // 返回数据
}

// ApiReturn 格式返回接口信息
func ApiReturn(code int, msg string, data interface{}) *Response {
	return &Response{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}

func Api(ctx iris.Context, data interface{}) {
	_, err := ctx.JSON(data)
	if err != nil {
		return
	}
}

// ApiRedirect 格式返回状态 消息 跳转连接
func ApiRedirect(code int, msg string, redirectUrl string) *Redirect {
	return &Redirect{
		Code: code,
		Msg:  msg,
		Url:  redirectUrl,
	}
}
