package manager

import (
	"webce/cmd/web/handlers/admin/controller"
	"webce/internal/repositories/models/admins"
	"webce/internal/repositories/repo/adminrepo"
	"webce/pkg/library/log"
)

type HandlerManager struct {
	admin.BaseHandler
	Repo *adminrepo.AdminUserRepository
}

func NewManager() *HandlerManager {
	return &HandlerManager{
		Repo: adminrepo.NewAdminUserRepository(),
	}
}

func (g *HandlerManager) GetCreateAdmin() {
	ad := &admins.Admin{}
	err := g.Ctx.ReadForm(ad)
	if err != nil {
		g.ApiJson(400, err.Error(), nil)
		return
	}
	err = ad.Validate()
	if err != nil {
		log.Log.Error("error to get params: ", err.Error())
		g.ApiJson(400, "invalid request", nil)
		return
	}

	resp := g.Repo.AddAdmin(ad, []int64{1})
	g.Api(resp)
}
