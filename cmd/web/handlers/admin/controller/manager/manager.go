package manager

import (
	"webce/cmd/web/handlers/admin/controller"
	admin2 "webce/internal/repositories/models/admins/admin"
	"webce/internal/repositories/repo/adminrepo"
	"webce/pkg/lib"
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
	ad := &admin2.Admin{}
	err := g.Ctx.ReadForm(ad)
	if err != nil {
		lib.ErrJson(g.Ctx, 400, err.Error())
		return
	}
	err = ad.Validate()
	if err != nil {
		log.Log.Error("error to get params: ", err.Error())
		lib.ErrJson(g.Ctx, 400, "invalid request")
		return
	}

	resp := g.Repo.AddAdmin(ad, []int64{1})
	lib.Json(g.Ctx, resp)
}
