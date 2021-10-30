package admins

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"webce/pkg/library/databases"
	"webce/pkg/library/easycasbin"
)

// 角色
type Roles struct {
	ID          uint64        `gorm:"primary_key" json:"id" structs:"id"`
	Title       string        `gorm:"type:varchar(50);unique_index" json:"title"` // 角色标题
	Description string        `gorm:"type:char(64);" json:"description"`          // 角色注解
	Permissions []Permissions `gorm:"many2many:role_menu;" json:"permissions" `
}

func (r Roles) Get(whereSql string, vals []interface{}) (Roles, error) {
	first := databases.DB.Preload("Permissions").Model(&r).Where(whereSql, vals).First(&r)
	if first.Error != nil {
		return r, first.Error
	}
	return r, nil
}

// 按照ID查找
func (r *Roles) FindByID(id int) (bool, error) {
	var role Roles
	err := databases.DB.Select("id").Where("id = ? ", id).First(&role).Error
	if err != nil {
		return false, err
	}
	if role.ID > 0 {
		return true, nil
	}
	return false, nil
}

// 依据传入的条件查找条数
func (r *Roles) GetCount(whereSql string, vals []interface{}) (int64, error) {
	var count int64
	if err := databases.DB.Model(&Roles{}).Where(whereSql, vals).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// 获取角色列表
func (r *Roles) GetRolesPage(whereSql string, vals []interface{}, offset, limit int) ([]*Roles, error) {
	var role []*Roles
	err := databases.DB.Where(whereSql, vals).Offset(offset).Limit(limit).Find(&role).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return role, nil
}

// 按照ID  获取角色
func (r *Roles) GetRoleByID(id int) (*Roles, error) {
	var role Roles
	err := databases.DB.Preload("Permissions").Where("id = ?", id).First(&role).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &role, nil
}

// 确认角色名称是否已存在
func (r *Roles) CheckRoleName(name string) (bool, error) {
	var role Roles
	err := databases.DB.Where("title=?", name).First(&role).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, nil
	}
	if role.ID > 0 {
		return true, nil
	}
	return false, nil
}

// 编辑角色
func (r Roles) EditRole(id int, data map[string]interface{}) error {

	var permsiss = make([]Permissions, 10)
	if err2 := databases.DB.Where("id in (?)", data["permissions_id"]).Find(&permsiss).Error; err2 != nil {
		return errors.New("无法找到该权限，请刷新后重试")
	}

	if err := databases.DB.Model(&r).Where("id = ?", id).Find(&r).Error; err != nil {
		return err
	}

	if err := databases.DB.Model(&r).Association("Permissions").Replace(permsiss).Error; err != nil {
		return errors.New("123")
	}

	if update := databases.DB.Model(&r).Updates(r).Error; update != nil {
		return update
	}
	return nil

}

// 添加角色
func (r *Roles) AddRole(data map[string]interface{}) (int, error) {
	role := Roles{
		Title:       data["title"].(string),
		Description: data["description"].(string),
	}
	var per []Permissions
	databases.DB.Where("id in (?)", data["permissions_id"]).Find(&per)
	err := databases.DB.Create(&role).Association("Permissions").Append(&per).Error
	if err != nil {
		return 0, errors.New("不晓得什么鬼")
	}
	return int(role.ID), nil
}

// 删除角色
func (r *Roles) DeleteRole(id int) error {

	databases.DB.Where("id = ?", id).First(&r)
	err := databases.DB.Model(&r).Association("Permissions").Delete()
	if err != nil {
		return err
	}

	err = databases.DB.Where("id = ?", id).Delete(&r).Error
	if err != nil {
		return err
	}
	return nil
}

// 删除所有角色
func (r *Roles) CleanRole() error {
	//Unscoped 方法可以物理删除记录
	if err := databases.DB.Unscoped().Where("deleted_on != ? ", 0).Delete(&Roles{}).Error; err != nil {
		return err
	}

	return nil
}

// 获取所有角色
func (r *Roles) GetRolesAll() ([]*Roles, error) {
	var role []*Roles
	err := databases.DB.Model(&role).Find(&role).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return role, nil
}

// LoadAllPolicy 加载所有的角色策略
func (a *Roles) LoadAllPolicy() error {
	roles, err := a.GetRolesAll()
	if err != nil {
		return err
	}

	for _, role := range roles {
		err = a.LoadPolicy(int(role.ID))
		if err != nil {
			return err
		}
	}
	return nil
}

// LoadPolicy 加载角色权限策略
func (a *Roles) LoadPolicy(id int) error {

	role, err := a.GetRoleByID(id)
	if err != nil {
		return err
	}
	_, err = easycasbin.GetEnforcer().DeleteRole(role.Title)
	if err != nil {
		return err
	}

	for _, menu := range role.Permissions {
		if menu.HttpPath == "" || menu.Method == "" {
			continue
		}
		_, err := easycasbin.GetEnforcer().AddPermissionForUser(role.Title, menu.HttpPath, menu.Method)
		if err != nil {
			return err
		}
	}
	return nil
}
