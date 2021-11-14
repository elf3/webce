package adminrepo

import (
	"github.com/pkg/errors"
	"webce/internal/repositories/models/admins/admin"
	"webce/pkg/library/databases"
	"webce/pkg/library/password"
	"webce/pkg/library/repository"
)

type AdminUserRepository struct {
	repository.Repository
	admin.Admin
}

func NewAdminUserRepository() *AdminUserRepository {
	var model = admin.Admin{}
	newRepository, _ := repository.NewRepository(
		databases.GetDB().Model(&model),
	)
	return &AdminUserRepository{
		newRepository,
		model,
	}
}

func (a AdminUserRepository) Login(username, pass string) (*admin.Admin, error) {

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
