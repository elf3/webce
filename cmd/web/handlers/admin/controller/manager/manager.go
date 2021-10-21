package manager

import (
	"webce/apis/repositories/models/admins"
	"webce/apis/repositories/repo/adminrepo"
	"webce/cmd/web/handlers/admin/controller"
	"webce/library/log"
)

type Manager struct {
	admin.BaseHandler
	Repo *adminrepo.AdminUserRepository
}

func NewManager() *Manager {
	return &Manager{
		Repo: adminrepo.NewAdminUserRepository(),
	}
}

func (g *Manager) GetCreateAdmin() {
	admin := &admins.Admin{}
	err := g.Ctx.ReadForm(admin)
	if err != nil {
		g.ApiJson(400, err.Error(), nil)
		return
	}
	err = admin.Validate()
	if err != nil {
		log.Log.Error("error to get params: ", err.Error())
		g.ApiJson(400, "invalid request", nil)
		return
	}

	resp := g.Repo.AddAdmin(admin, []int64{1})
	g.Api(resp)
}
