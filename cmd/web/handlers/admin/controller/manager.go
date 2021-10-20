package admin

import (
	"github.com/kataras/iris/v12"
	"webce/apis/repositories/repo/adminrepo"
	"webce/library/apgs"
)

type Manager struct {
	Ctx  iris.Context
	Repo *adminrepo.AdminUserRepository
}

func NewManager() *Manager {
	return &Manager{
		Repo: adminrepo.NewAdminUserRepository(),
	}
}

func (g *Manager) GetTest() {
	where := make(map[string]interface{})
	where["id"] = 1
	users1 := g.Repo.Select(where)
	_, err := g.Ctx.JSON(apgs.ApiReturn(0, "123123", users1))
	if err != nil {
		return
	}
	//users, _ := g.Repo.SelectById("select * from wk_admin_iser where id=?", 1)
	//fmt.Println(users)
	//fmt.Println(users1)
	//
	//_, _ = g.Ctx.JSON(apgs.ApiReturn(0, "123123", apgs.Map{
	//	"Users":  users,
	//	"Users1": users1,
	//}))
}

func (g *Manager) GetTest2() {
	repo := adminrepo.NewAdminUserRepository()
	adminresp := repo.AddAdmin()
	g.Ctx.JSON(adminresp)
}

func (g *Manager) Get() {
	_, err := g.Ctx.JSON(apgs.ApiReturn(0, "123123", nil))
	if err != nil {
		return
	}
}
