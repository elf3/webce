package admin

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
	mod "webce/internal/repositories/models/admins/admin"
	"webce/internal/repositories/models/admins/roles"
	"webce/pkg/library/databases"
	"webce/pkg/library/easycasbin"
	"webce/pkg/library/log"
)

type ApiAdmin struct {
}

// GetByCount 获取有多少条记录
func (a ApiAdmin) GetByCount(whereSql string, vals []interface{}) (count int64) {
	admin := mod.Admin{}
	databases.DB.Model(&admin).Where(whereSql, vals...).Count(&count)
	return
}

// Lists 获取列表，按照 offest 和 limit参数进行分页
func (a ApiAdmin) Lists(whereSql string, vals []interface{}, offset, limit int) ([]mod.Admin, error) {
	list := make([]mod.Admin, limit)
	admin := mod.Admin{}
	find := databases.DB.Preload("Roles").
		Model(&admin).Select(mod.AvailableQueryFields).Where(whereSql, vals...).Offset(offset).Limit(limit).Find(&list)
	if find.Error != nil && find.Error != gorm.ErrRecordNotFound {
		return nil, find.Error
	}
	return list, nil
}

// Get 获取单条记录
func (a ApiAdmin) Get(whereSql string, vals []interface{}) (mod.Admin, error) {
	admin := mod.Admin{}
	first := databases.DB.Preload("Roles").Model(&admin).Where(whereSql, vals...).First(&admin)
	if first.Error != nil {
		return admin, first.Error
	}
	return admin, nil
}

// GetById 通过主键ID
func (a ApiAdmin) GetById(id uint64) (mod.Admin, error) {
	admin := mod.Admin{}
	first := databases.DB.Preload("Roles").Model(&admin).Where("id = ?", id).First(&admin)
	if first.Error != nil {
		return admin, first.Error
	}
	return admin, nil
}

// Create 创建记录
func (a ApiAdmin) Create(roleIds []int64, data *mod.Admin) (*mod.Admin, error) {
	var role = make([]roles.Roles, 0)
	admin := mod.Admin{}
	find := databases.DB.Where("id in (?)", roleIds).Find(&role)
	if find.Error != nil || find.RowsAffected == 0 {
		return nil, errors.New("角色未初始化")
	}
	create := databases.DB.Model(&admin).Create(&data).Association("Roles").Append(role)
	if create != nil {
		return nil, errors.Wrap(create, "创建管理员失败，请重试")
	}
	return &admin, nil
}

// Update 更新操作
func (a ApiAdmin) UpdateAdmin(id uint64, roleIds []int64, data mod.Admin) (*mod.Admin, error) {
	var role = make([]roles.Roles, 10)
	if err := databases.DB.Where("id in (?)", roleIds).Find(&role).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("管理员没找到")
	}
	var admin mod.Admin
	find := databases.DB.Model(&admin).Where("id = ?", id).Find(&admin)
	if find.Error != nil {
		return nil, find.Error
	}

	err := databases.DB.Model(&admin).Association("Roles").Replace(role)
	if err != nil {
		return nil, err
	}
	save := databases.DB.Model(&admin).Where("id = ?", id).Updates(&data)

	if save.Error != nil {
		return nil, save.Error
	}
	return &admin, nil
}

// Delete 删除操作
func (a ApiAdmin) Delete(id int) error {

	var data mod.Admin
	databases.DB.Model(&data).Where("id = ?", id).Find(&data)

	err := databases.DB.Model(&data).Select("Roles").Delete(&data).Error
	if err != nil {
		log.Log.Error("delete admin role err", err)
		return err
	}
	_, err = easycasbin.GetEnforcer().DeleteUser(data.Username)
	if err != nil {
		log.Log.Error("delete casbin rule err", err)
		return err
	}
	return nil
}

// LoadPolicy 加载用户权限策略
func (a ApiAdmin) LoadPolicy(id uint64) error {

	admin, err := a.GetById(id)
	if err != nil {
		return err
	}
	_, err = easycasbin.GetEnforcer().DeleteRolesForUser(admin.Username)
	if err != nil {
		log.Log.Error("LoadPolicy DeleteRolesForUser error : ", err)
		return err
	}

	for _, ro := range admin.Roles {
		_, err := easycasbin.GetEnforcer().AddRoleForUser(admin.Username, ro.Title)
		if err != nil {
			return err
		}
	}
	//fmt.Println("更新角色权限关系", easycasbin.GetEnforcer().GetGroupingPolicy())
	return nil
}

// GetUsersAll 获取所有管理员 - 包含角色
func (a ApiAdmin) GetUsersAll() ([]*mod.Admin, error) {
	var admin []*mod.Admin
	u := mod.Admin{}
	err := databases.DB.Model(&u).Preload("Roles").Find(&admin).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return admin, nil
}

// LoadAdminPolicy 加载用户策略
func (a ApiAdmin) LoadAdminPolicy(id uint64) error {
	admin, err := a.GetById(id)
	if err != nil {
		return err
	}
	if len(admin.Roles) != 0 {
		err = a.LoadPolicy(admin.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

// LoadAllPolicy 加载所有的用户策略
func (a ApiAdmin) LoadAllPolicy() error {
	admins, err := a.GetUsersAll()
	if err != nil {
		return err
	}
	for _, admin := range admins {
		if len(admin.Roles) != 0 {
			err = a.LoadPolicy(admin.ID)
			if err != nil {
				return err
			}
		}
	}
	//fmt.Println("角色权限关系", easycasbin.GetEnforcer().GetGroupingPolicy())
	return nil
}
