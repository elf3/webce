package permission

import (
	"gorm.io/gorm"
	"webce/internal/repositories/models/admins/permissions"
	"webce/pkg/library/databases"
)

type ApiPermission struct {
}

// GetMenus 获取权限菜单
func (m *ApiPermission) GetMenus() []*permissions.Permissions {
	allMenu := make([]*permissions.Permissions, 0)
	databases.DB.Model(&allMenu).Find(&allMenu)
	return allMenu
}

// GetByCount 获取条数 （分页使用）
func (m *ApiPermission) GetByCount(where string, values []interface{}) (count int64) {
	p := permissions.Permissions{}
	databases.DB.Model(&p).Where(where, values...).Count(&count)
	return
}

// Lists 获取权限列表
func (m *ApiPermission) Lists(where string, values []interface{}, offset, limit int) ([]permissions.Permissions, error) {
	list := make([]permissions.Permissions, limit)
	err := databases.DB.Model(&list).Select(permissions.AvailableQueryFields).
		Where(where, values...).Offset(offset).Limit(limit).Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

// Get 获取当个权限详情
func (m *ApiPermission) Get(where string, values []interface{}) (*permissions.Permissions, error) {
	p := permissions.Permissions{}
	err := databases.DB.Model(&p).Where(where, values...).First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// Update 更新权限
func (m *ApiPermission) Update(id uint64, p permissions.Permissions) (*permissions.Permissions, error) {
	err := databases.DB.Model(&p).Where("id = ?", id).Find(&p).Error
	if err != nil {
		return nil, err
	}

	err = databases.DB.Model(&p).Where("id = ?", id).Updates(p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// Create 更新权限
func (m *ApiPermission) Create(p permissions.Permissions) error {
	create := databases.DB.Model(&p).Create(&p)
	if create.Error != nil {
		return create.Error
	}
	return nil
}

// Delete 删除权限
func (m *ApiPermission) Delete(id uint64) error {
	p := permissions.Permissions{}
	find := databases.DB.Where("id = ?", id).Find(&p)
	if find.Error != nil || find.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	err := databases.DB.Model(&p).Where("id = ?", id).Delete(&p).Error
	if err != nil {
		return err
	}

	return nil
}
