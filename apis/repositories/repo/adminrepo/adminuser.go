package adminrepo

import (
	"webce/apis/repositories/models/admins"
	"webce/library/apgs"
	"webce/library/databases"
	"webce/library/log"
	"webce/library/password"
	"webce/library/repository"
)

type AdminUserRepository struct {
	repository.Repository
	admins.Admin
}

func NewAdminUserRepository() *AdminUserRepository {
	var model = admins.Admin{}
	newRepository, _ := repository.NewRepository(
		databases.GetDB().Model(&model),
	)
	return &AdminUserRepository{
		newRepository,
		model,
	}
}

func (a AdminUserRepository) Login(username, pass string) *apgs.Response {

	adminData := databases.DB.Model(&a.Admin).Preload("Roles").Where("username=?", username).First(&a.Admin)
	if adminData.Error != nil {
		return apgs.ApiReturn(400, "账号或密码错误", nil)
	}
	passBool := password.Compare(a.Admin.Password, pass)
	if passBool != nil {
		return apgs.ApiReturn(400, "账号或密码错误", nil)
	}
	a.Admin.Password = ""
	return apgs.ApiReturn(0, "", a.Admin)
}

func (a AdminUserRepository) AddAdmin(admin *admins.Admin, roleIds []int64) *apgs.Response {
	create, err := admin.Create(roleIds)
	if err != nil {
		log.Log.Error("create admin err: ", err.Error())
		return apgs.ApiReturn(400, "create admin err", nil)
	}
	return apgs.ApiReturn(200, "", create)
}
