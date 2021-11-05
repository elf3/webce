package adminrepo

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"webce/internal/repositories/models/admins/admin"
	"webce/pkg/library/apgs"
	"webce/pkg/library/databases"
	"webce/pkg/library/log"
	"webce/pkg/library/password"
	repository2 "webce/pkg/library/repository"
)

type AdminUserRepository struct {
	repository2.Repository
	admin.Admin
}

func NewAdminUserRepository() *AdminUserRepository {
	var model = admin.Admin{}
	newRepository, _ := repository2.NewRepository(
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

func (a AdminUserRepository) AddAdmin(add *admin.Admin, roleIds []int64) *apgs.Response {
	findAdmin := &admin.Admin{}
	adminData := databases.DB.Model(&findAdmin).Where("username=?", add.Username).First(&findAdmin)
	if adminData != nil && adminData.Error != gorm.ErrRecordNotFound {
		return apgs.ApiReturn(500, "账号已存在，请重试", nil)
	}
	passCode, err := password.Encrypt(add.Password)
	if err != nil {
		log.Log.Error("error encrypt password：", err)
		return apgs.ApiReturn(500, "error encrypt password", nil)
	}
	add.Password = passCode
	create, err := findAdmin.Create(roleIds)
	if err != nil {
		log.Log.Error("create admin err: ", err.Error())
		return apgs.ApiReturn(400, "create admin err", nil)
	}
	return apgs.ApiReturn(200, "", create)
}

func (a AdminUserRepository) GetAdminById(id uint64) (admin.Admin, error) {
	return a.Admin.GetById(id)
}
