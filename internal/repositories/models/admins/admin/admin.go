package admin

import (
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"webce/internal/repositories/models"
	"webce/internal/repositories/models/admins/roles"
	"webce/pkg/library/databases"
	"webce/pkg/library/easycasbin"
	"webce/pkg/library/log"
)

// Admin 管理员
type Admin struct {
	models.BaseModel
	Username string `gorm:"type:char(50);unique; unique_index;not null;" validate:"min=6,max=32" form:"username" json:"username"`
	// 设置管理员账号 唯一并且不为空
	Password    string        `gorm:"size:255;not null;" form:"password" validate:"min=6,max=32"  json:"password" ` // 设置字段大小为255
	LastLoginIp int64         `gorm:"type:bigint(1);not null;" json:"last_login_ip"`                                // 上次登录IP
	IsSuper     int           `gorm:"type:tinyint(1);not null"  json:"is_super"`                                    // 是否超级管理员
	Roles       []roles.Roles `gorm:"many2many:admin_role;not null;" json:"roles"`                                  // 角色
}

// GetByCount 获取有多少条记录
func (u Admin) GetByCount(whereSql string, vals []interface{}) (count int64) {
	databases.DB.Model(u).Where(whereSql, vals...).Count(&count)
	return
}

// Lists 获取列表，按照 offest 和 limit参数进行分页
func (u Admin) Lists(whereSql string, vals []interface{}, offset, limit int) ([]Admin, error) {
	list := make([]Admin, limit)
	find := databases.DB.Preload("Roles").
		Model(&u).Select(AvailableQueryFields).Where(whereSql, vals...).Offset(offset).Limit(limit).Find(&list)
	if find.Error != nil && find.Error != gorm.ErrRecordNotFound {
		return nil, find.Error
	}
	return list, nil
}

// Get 获取单条记录
func (u Admin) Get(whereSql string, vals []interface{}) (Admin, error) {
	first := databases.DB.Preload("Roles").Model(&u).Where(whereSql, vals...).First(&u)
	if first.Error != nil {
		return u, first.Error
	}
	return u, nil
}

// GetById 通过主键ID
func (u Admin) GetById(id uint64) (Admin, error) {
	first := databases.DB.Preload("Roles").Model(&u).Where("id = ?", id).First(&u)
	if first.Error != nil {
		return u, first.Error
	}
	return u, nil
}

// Create 创建记录
func (u Admin) Create(roleIds []int64) (*Admin, error) {
	var role = make([]roles.Roles, 0)
	find := databases.DB.Where("id in (?)", roleIds).Find(&role)
	if find.Error != nil || find.RowsAffected == 0 {
		return nil, errors.New("角色未初始化")
	}
	create := databases.DB.Model(&u).Create(&u).Association("Roles").Append(role)
	if create != nil {
		return nil, errors.Wrap(create, "创建管理员失败，请重试")
	}
	return &u, nil
}

// Update 更新操作
func (u Admin) UpdateAdmin(id uint64, roleIds []int64) (*Admin, error) {
	var role = make([]roles.Roles, 10)
	if err := databases.DB.Where("id in (?)", roleIds).Find(&role).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("管理员没找到")
	}
	var admin Admin
	find := databases.DB.Model(&admin).Where("id = ?", id).Find(&admin)
	if find.Error != nil {
		return nil, find.Error
	}

	databases.DB.Model(&admin).Association("Roles").Replace(role)
	save := databases.DB.Model(&admin).Where("id = ?", id).Updates(&u)

	if save.Error != nil {
		return nil, save.Error
	}
	return &admin, nil
}

// Delete 删除操作
func (u Admin) Delete(id int) error {
	fmt.Println("id: ......... ", id)
	var data Admin
	databases.DB.Model(&data).Where("id = ?", id).Find(&data)

	err := databases.DB.Model(&data).Select("Roles").Delete(&data).Error
	if err != nil {
		log.Log.Error("delete admin role err", err)
		return err
	}
	_, err = easycasbin.GetEnforcer().DeleteUser(u.Username)
	if err != nil {
		log.Log.Error("delete casbin rule err", err)
		return err
	}
	return nil
}

// LoadPolicy 加载用户权限策略
func (u *Admin) LoadPolicy(id uint64) error {

	admin, err := u.GetById(id)
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
func (u Admin) GetUsersAll() ([]*Admin, error) {
	var admin []*Admin
	err := databases.DB.Model(&u).Preload("Roles").Find(&admin).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return admin, nil
}

// LoadAdminPolicy 加载用户策略
func (u *Admin) LoadAdminPolicy(id uint64) error {
	admin, err := u.GetById(id)
	if err != nil {
		return err
	}
	if len(admin.Roles) != 0 {
		err = u.LoadPolicy(admin.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

// LoadAllPolicy 加载所有的用户策略
func (u *Admin) LoadAllPolicy() error {
	admins, err := u.GetUsersAll()
	if err != nil {
		return err
	}
	for _, admin := range admins {
		if len(admin.Roles) != 0 {
			err = u.LoadPolicy(admin.ID)
			if err != nil {
				return err
			}
		}
	}
	//fmt.Println("角色权限关系", easycasbin.GetEnforcer().GetGroupingPolicy())
	return nil
}
