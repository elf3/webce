package middle

import (
	"github.com/dgrijalva/jwt-go"
	jwtmiddleware "github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/v12"
	"webce/pkg/library/apgs"
	liberty "webce/pkg/library/jwt"
)

// OnError jwt验证失败返回内容
func OnError(ctx iris.Context, err error) {
	if err == nil {
		return
	}
	ctx.StopExecution()
	ctx.JSON(apgs.ApiReturn(999, err.Error(), nil))
}

// JwtHandler 验证Token
func JwtHandler() *jwtmiddleware.Middleware {
	return jwtmiddleware.New(jwtmiddleware.Config{
		ValidationKeyGetter: func(token *jwt.Token) (i interface{}, e error) {
			return liberty.SecretKey, nil
		},
		ErrorHandler:  OnError,
		SigningMethod: jwt.SigningMethodHS256,
	})
}
