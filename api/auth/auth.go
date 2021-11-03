package auth

import (
	"github.com/pkg/errors"
	"strconv"
	"webce/api/auth/response"
	"webce/internal/repositories/repo/adminrepo"
	"webce/pkg/lib"
	"webce/pkg/library/jwt"
	"webce/pkg/library/log"
)

type AdminAuth struct {
}

func (a *AdminAuth) Login(username, password, ip string) (*response.LoginResponse, error) {
	repo := adminrepo.NewAdminUserRepository()
	resp, err := repo.Login(username, password)
	if err != nil {
		log.Log.Error("login err: ", err)
		return nil, errors.New("登陆错误，请重试")
	}
	m := map[string]interface{}{
		"last_login_ip": lib.Ip2long(ip),
	}
	where := map[string]interface{}{
		"username": username,
	}
	_, err = repo.Update(where, m)
	if err != nil {
		log.Log.Error("login update last_login_ip error :", err)
		return nil, errors.New("ip错误")
	}

	token, err := jwt.CreateToken(strconv.FormatUint(resp.ID, 10), resp.Username)
	if err != nil {
		log.Log.Error("token create err :", err)
		return nil, errors.New("token create err")
	}
	return &response.LoginResponse{
		Admin: resp,
		Token: token,
	}, nil
}
