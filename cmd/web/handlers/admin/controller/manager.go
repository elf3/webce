package admin

import (
	"github.com/kataras/iris/v12"
	"webce/apis/repositories/models/admins"
	"webce/apis/repositories/repo/adminrepo"
	"webce/library/apgs"
	"webce/library/log"
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
	admin := &admins.Admin{}
	err := g.Ctx.ReadForm(admin)
	if err != nil {
		g.Ctx.JSON(apgs.ApiReturn(400, err.Error(), nil))
		return
	}
	err = admin.Validate()
	if err != nil {
		log.Log.Error("error to get params: ", err.Error())
		g.Ctx.JSON(apgs.ApiReturn(400, "invalid request", nil))

		return
	}

	adminresp := repo.AddAdmin(admin, []int64{1})
	g.Ctx.JSON(adminresp)
}

func (g *Manager) Get() {
	_, err := g.Ctx.JSON(apgs.ApiReturn(0, "123123", nil))
	if err != nil {
		return
	}
}
