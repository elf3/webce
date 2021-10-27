package services

import (
	"webce/apis/repositories/repo/adminrepo"
	"webce/library/apgs"
	"webce/library/log"
	"webce/pkg/lib"
)

type AdminAuth struct {
}

func (a *AdminAuth) Login(username, password, ip string) *apgs.Response {
	repo := adminrepo.NewAdminUserRepository()
	resp, err := repo.Login(username, password)
	if err != nil {
		return apgs.ApiReturn(500, "登陆错误，请重试", nil)
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
		return apgs.ApiReturn(400, "ip错误", nil)
	}

	return apgs.ApiReturn(200, "success", resp)
}
