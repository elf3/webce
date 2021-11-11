package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"time"
	"webce/pkg/lib"
)

var (
	Exp       int64
	SecretKey string
)

func InitJwtConf() {
	Exp = viper.GetInt64("jwt.expire")
	SecretKey = viper.GetString("jwt.secret")
}

// CreateToken 生成token
func CreateToken(userId, username string) (tokenString string, err error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Second * time.Duration(Exp)).Unix()
	claims["iat"] = time.Now().Unix()
	claims["username"] = username
	claims["userId"] = userId

	token.Claims = claims
	tokenString, err = token.SignedString(lib.StringBytes(SecretKey))
	return
}

// ParseToken 解析token
func ParseToken(tokenSrt string) (claims jwt.MapClaims, err error) {
	token, err := jwt.Parse(tokenSrt, func(*jwt.Token) (interface{}, error) {
		return lib.StringBytes(SecretKey), nil
	})
	if err != nil {
		return nil, err
	}
	claims = token.Claims.(jwt.MapClaims)
	return
}
