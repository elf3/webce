package api

import (
	"fmt"
	"github.com/kataras/iris/v12"
)

func Main(ctx iris.Context) {

	err := ctx.View("login.html")
	if err != nil {
		return
	}
}

func Login(ctx iris.Context) {
	value := ctx.PostValue("1")

	fmt.Println(value)
}
