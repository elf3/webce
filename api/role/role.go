package role

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"webce/internal/repositories/models/admins/permissions"
	"webce/internal/repositories/models/admins/roles"
	"webce/pkg/library/databases"
	"webce/pkg/library/easycasbin"
)

type ApiRole struct {
}

func (a ApiRole) Get(whereSql string, vals []interface{}) (roles.Roles, error) {
	r := roles.Roles{}
	first := databases.DB.Preload("Permissions").Model(&r).Where(whereSql, vals...).First(&r)
	if first.Error != nil {
		return r, first.Error
	}
	return r, nil
}

// GetPerRoleIds 获取权限绑定的角色ID列表
func (a ApiRole) GetPerRoleIds(id int) []int {
	var permission permissions.Permissions
	var role []roles.Roles

	databases.DB.Model(&permission).Where("id = ? ", id, 0)
	pf := databases.GetPrefix()
	joins := " left join " + pf + "role_menu b on " + pf + "roles.id=b.role_id left join " + pf + "permissions c on c.id=b.permissions_id"
	databases.DB.Joins(joins).Where("c.id = ?", id).Find(&role)

	var roleList []int
	for _, v := range role {
		roleList = append(roleList, int(v.ID))
	}
	return roleList
}

// FindByID 按照ID查找
func (a ApiRole) FindByID(id int) (bool, error) {
	var role roles.Roles
	err := databases.DB.Select("id").Where("id = ? ", id).First(&role).Error
	if err != nil {
		return false, err
	}
	if role.ID > 0 {
		return true, nil
	}
	return false, nil
}

// GetByCount 依据传入的条件查找条数
func (a ApiRole) GetByCount(whereSql string, vals []interface{}) (int64, error) {
	var count int64
	r := roles.Roles{}
	if err := databases.DB.Model(&r).Where(whereSql, vals...).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// GetRolesPage 获取角色列表
func (a ApiRole) GetRolesPage(whereSql string, vals []interface{}, offset, limit int) ([]*roles.Roles, error) {
	var role []*roles.Roles
	err := databases.DB.Where(whereSql, vals...).Offset(offset).Limit(limit).Find(&role).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return role, nil
}

// GetRoleByID 按照ID  获取角色
func (a ApiRole) GetRoleByID(id int) (*roles.Roles, error) {
	var role roles.Roles
	err := databases.DB.Preload("Permissions").Where("id = ?", id).First(&role).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &role, nil
}

// CheckRoleName 确认角色名称是否已存在
func (a ApiRole) CheckRoleName(name string) (bool, error) {
	var role roles.Roles
	err := databases.DB.Where("title=?", name).First(&role).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, nil
	}
	if role.ID > 0 {
		return true, nil
	}
	return false, nil
}

// EditRole 编辑角色
func (a ApiRole) EditRole(id uint64, permissionsIds []uint64, r roles.Roles) (*roles.Roles, error) {
	var persist = make([]permissions.Permissions, 10)
	if err := databases.DB.Where("id in (?)", permissionsIds).Find(&persist).Error; err != nil {
		return nil, errors.New("无法找到该权限，请刷新后重试")
	}
	var findRole roles.Roles
	err := databases.DB.Model(&findRole).Where("id = ?", id).Find(&findRole).Error
	if err != nil {
		return nil, err
	}
	err = databases.DB.Model(&findRole).Association("Permissions").Replace(persist)
	if err != nil {
		return nil, errors.New("无法更新权限")
	}

	if update := databases.DB.Model(&findRole).Updates(r).Error; update != nil {
		return nil, update
	}
	r.ID = id
	return &r, nil

}

// AddRole 添加角色
func (a ApiRole) AddRole(permissionIds []uint64, r roles.Roles) (*roles.Roles, error) {
	var per []permissions.Permissions
	err := databases.DB.Where("id in (?)", permissionIds).Find(&per).Error
	if err != nil {
		return nil, err
	}
	err = databases.DB.Create(&r).Association("Permissions").Append(&per)
	if err != nil {
		return nil, errors.New("不晓得什么鬼")
	}
	return &r, nil
}

// DeleteRole 删除角色
func (a ApiRole) DeleteRole(id int) error {
	r := roles.Roles{}
	databases.DB.Model(&r).Where("id = ?", id).First(&r)
	err := databases.DB.Model(&r).Select("Permissions").Delete(&r).Error
	if err != nil {
		return err
	}
	return nil
}

// CleanRole 删除所有角色
func (a ApiRole) CleanRole() error {
	//Unscoped 方法可以物理删除记录
	if err := databases.DB.Unscoped().Where("deleted_on != ? ", 0).Delete(&roles.Roles{}).Error; err != nil {
		return err
	}

	return nil
}

// GetRolesAll 获取所有角色
func (a ApiRole) GetRolesAll() ([]*roles.Roles, error) {
	var role []*roles.Roles
	err := databases.DB.Model(&role).Find(&role).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return role, nil
}

// LoadAllPolicy 加载所有的角色策略
func (a ApiRole) LoadAllPolicy() error {
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
func (a *ApiRole) LoadPolicy(id int) error {

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
