package auth

import (
	"github.com/pkg/errors"
	"webce/internal/repositories/models/admins"
	"webce/internal/repositories/repo/adminrepo"
	"webce/pkg/lib"
	"webce/pkg/library/log"
)

type AdminAuth struct {
}

func (a *AdminAuth) Login(username, password, ip string) (*admins.Admin, error) {
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

	return resp, nil
}
