package services

import (
	"github.com/kataras/iris/v12"
	"webce/apis/repositories/repo/adminrepo"
	"webce/library/apgs"
)

type AdminAuth struct {
}

func (a *AdminAuth) Login(username, password, ip string) *apgs.Response {
	//session.SaveSession(c, session.UserKey, admin)
	//
	//v := l.GetMap(1)
	//ip := c.ClientIP()
	//v["last_login_ip"] = ip
	//databases.DB.Model(&admin).Where("id = ?", admin.ID).Update(v)
	repo := adminrepo.NewAdminUserRepository()
	resp := repo.Login(username, password)
	if resp.Code == 0 {
		m := iris.Map{}
		m["last_login_ip"] = ip
		_, err := repo.Update(m, m)
		if err != nil {
			apgs.ApiReturn(400, err.Error(), nil)
		}
	}

	return resp
}
