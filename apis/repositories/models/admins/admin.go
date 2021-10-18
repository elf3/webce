package admins

import (
	"errors"
	"fmt"
	"gopkg.in/go-playground/validator.v9"
	"gorm.io/gorm"
	"webce/apis/repositories/models"
	"webce/library/databases"
	"webce/library/easycasbin"
)

// Admin 管理员
type Admin struct {
	models.BaseModel
	Username string `gorm:"type:char(50); unique_index;not null;"  validate:"min=6,max=32"`
	// 设置管理员账号 唯一并且不为空
	Password    string `gorm:"size:255;not null;"  `     // 设置字段大小为255
	LastLoginIp uint32 `gorm:"type:int(1);not null;"`    // 上次登录IP
	IsSuper     int    `gorm:"type:tinyint(1);not null"` // 是否超级管理员

	Roles []Roles `json:"roles" gorm:"many2many:admin_role;not null;"`
}

// Validate the fields.
func (u *Admin) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}

// GetByCount 获取有多少条记录
func (u Admin) GetByCount(whereSql string, vals []interface{}) (count int64) {
	databases.DB.Model(u).Where(whereSql, vals).Count(&count)
	return
}

// Lists 获取列表，按照 offest 和 limit参数进行分页
func (u Admin) Lists(fields string, whereSql string, vals []interface{}, offset, limit int) ([]Admin, error) {
	list := make([]Admin, limit)
	find := databases.DB.Preload("Roles").Model(&u).Select(fields).Where(whereSql, vals).Offset(offset).Limit(limit).Find(&list)
	if find.Error != nil && find.Error != gorm.ErrRecordNotFound {
		return nil, find.Error
	}
	return list, nil
}

// Get 获取单条记录
func (u Admin) Get(whereSql string, vals []interface{}) (Admin, error) {
	first := databases.DB.Preload("Roles").Model(&u).Where(whereSql, vals).First(&u)
	if first.Error != nil {
		return u, first.Error
	}
	return u, nil
}

// GetById 通过主键ID
func (u Admin) GetById(id int) (Admin, error) {
	first := databases.DB.Preload("Roles").Model(&u).Where("id = ?", id).First(&u)
	if first.Error != nil {
		return u, first.Error
	}
	return u, nil
}

// Create 创建记录
func (u Admin) Create(data map[string]interface{}) (*Admin, error) {
	var role = make([]Roles, 10)
	databases.DB.Where("id in (?)", data["role_id"]).Find(&role)
	create := databases.DB.Model(&u).Create(&u).Association("Roles").Append(role).Error()
	if create != "" {
		return nil, errors.New("测试")
	}
	return &u, nil
}

// 更新操作
//func (u Admin) Update(id int, data map[string]interface{}) error {
//	var role = make([]Roles, 10)
//	if err := databases.DB.Where("id in (?)", data["role_id"]).Find(&role).Error; errors.Is(err, gorm.ErrRecordNotFound) {
//		return errors.New("管理员没找到")
//	}
//
//	find := databases.DB.Model(&u).Where("id = ?", id).Find(&u)
//	if find.Error != nil {
//		return find.Error
//	}
//
//	databases.DB.Model(&u).Association("Roles").Replace(role)
//	save := databases.DB.Model(&u).Updates(data)
//
//	if save.Error != nil {
//		return save.Error
//	}
//	return nil
//}

// Delete 删除操作
func (u Admin) Delete(id int) (bool, error) {
	databases.DB.Where("id = ?", id).Find(&u)
	err := databases.DB.Model(&u).Association("Roles").Delete()
	if err != nil {
		return false, err
	}
	db := databases.DB.Model(&u).Where("id = ?", id).Delete(&u)
	if db.Error != nil {
		return false, db.Error
	}

	_, err = easycasbin.GetEnforcer().DeleteUser(u.Username)
	if err != nil {
		return false, err
	}
	return true, nil
}

// LoadPolicy 加载用户权限策略
func (u *Admin) LoadPolicy(id int) error {

	admin, err := u.GetById(id)
	if err != nil {
		return err
	}
	_, err = easycasbin.GetEnforcer().DeleteRolesForUser(admin.Username)
	if err != nil {
		return err
	}

	for _, ro := range admin.Roles {
		_, err := easycasbin.GetEnforcer().AddRoleForUser(admin.Username, ro.Title)
		if err != nil {
			return err
		}
	}
	fmt.Println("更新角色权限关系", easycasbin.GetEnforcer().GetGroupingPolicy())
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

// LoadAllPolicy 加载所有的用户策略
func (u *Admin) LoadAllPolicy() error {
	admins, err := u.GetUsersAll()
	if err != nil {
		return err
	}
	for _, admin := range admins {
		if len(admin.Roles) != 0 {
			err = u.LoadPolicy(int(admin.ID))
			if err != nil {
				return err
			}
		}
	}
	fmt.Println("角色权限关系", easycasbin.GetEnforcer().GetGroupingPolicy())
	return nil
}
