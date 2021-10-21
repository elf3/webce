package services

import (
	"github.com/kataras/iris/v12"
	"webce/apis/repositories/repo/adminrepo"
	"webce/library/apgs"
	"webce/library/log"
	"webce/pkg/lib"
)

type AdminAuth struct {
}

func (a *AdminAuth) Login(username, password, ip string) *apgs.Response {
	repo := adminrepo.NewAdminUserRepository()
	resp := repo.Login(username, password)
	if resp.Code == 0 {
		m := iris.Map{
			"last_login_ip": lib.Ip2long(ip),
		}
		where := iris.Map{
			"username": username,
		}
		_, err := repo.Update(where, m)
		if err != nil {
			log.Log.Error("login update last_login_ip error :", err)
			return apgs.ApiReturn(400, err.Error(), nil)
		}
	}

	return resp
}
