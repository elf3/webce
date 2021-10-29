package adminrepo

import (
	"github.com/pkg/errors"
	"webce/internal/repositories/models/admins"
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

func (a AdminUserRepository) Login(username, pass string) (*admins.Admin, error) {

	adminData := databases.DB.Model(&a.Admin).Preload("Roles").Where("username=?", username).First(&a.Admin)
	if adminData.Error != nil {
		return nil, errors.New("账号或密码错误")
	}
	passBool := password.Compare(a.Admin.Password, pass)
	if passBool != nil {
		return nil, errors.New("账号或密码错误")
	}
	a.Admin.Password = ""
	return &a.Admin, nil
}

func (a AdminUserRepository) AddAdmin(admin *admins.Admin, roleIds []int64) *apgs.Response {
	passCode, err := password.Encrypt(admin.Password)
	if err != nil {
		log.Log.Error("error encrypt password：", err)
		return apgs.ApiReturn(500, "error encrypt password", nil)
	}
	admin.Password = passCode
	create, err := admin.Create(roleIds)
	if err != nil {
		log.Log.Error("create admin err: ", err.Error())
		return apgs.ApiReturn(400, "create admin err", nil)
	}
	return apgs.ApiReturn(200, "", create)
}
