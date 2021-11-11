package manager

import (
	"webce/cmd/web/handlers/admin"
	admin2 "webce/internal/repositories/models/admins/admin"
	"webce/internal/repositories/repo/adminrepo"
	"webce/pkg/library/log"
)

type HandlerManager struct {
	admin.admin
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
		g.Error(400, err.Error())
		return
	}
	err = g.Validate(ad)
	if err != nil {
		log.Log.Error("error to get params: ", err.Error())
		g.Error(400, "invalid request")
		return
	}

	result := g.Repo.AddAdmin(ad, []int64{1})
	g.Success(result)
}