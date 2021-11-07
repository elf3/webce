package permissions

import (
	"webce/pkg/library/databases"
)

// Permissions 菜单权限
type Permissions struct {
	ID          uint64 `gorm:"primary_key"  structs:"id" form:"id" json:"id"`                               // 主键ID
	Title       string `gorm:"type:varchar(50);unique_index" form:"title" json:"title" validate:"required"` // 权限标题
	Description string `gorm:"type:char(64);" form:"description" json:"description"`                        // 注解
	Slug        string `gorm:"type:varchar(50);" form:"slug" json:"slug"`                                   // 权限名称
	HttpPath    string `gorm:"type:text" form:"http_path" json:"http_path" validate:"required"`             // URI路径
	Method      string `gorm:"type:char(10);" form:"method" json:"method" validate:"required"`              // 请求方法
}

// GetMenus 获取权限菜单
func (m *Permissions) GetMenus() []*Permissions {
	allMenu := make([]*Permissions, 0)
	databases.DB.Model(&allMenu).Find(&allMenu)
	return allMenu
}

// GetByCount 获取条数 （分页使用）
func (m *Permissions) GetByCount(where string, values []interface{}) (count int64) {
	databases.DB.Model(&m).Where(where, values...).Count(&count)
	return
}

// Lists 获取权限列表
func (m *Permissions) Lists(where string, values []interface{}, offset, limit int) ([]Permissions, error) {
	list := make([]Permissions, limit)
	err := databases.DB.Model(&list).Select(AvailableQueryFields).
		Where(where, values...).Offset(offset).Limit(limit).Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

// Get 获取当个权限详情
func (m *Permissions) Get(where string, values []interface{}) (*Permissions, error) {
	err := databases.DB.Model(&m).Where(where, values...).First(&m).Error
	if err != nil {
		return nil, err
	}
	return m, nil
}

// Update 更新权限
func (m *Permissions) Update(id uint64) (*Permissions, error) {
	err := databases.DB.Model(&m).Where("id = ?", id).Find(&m).Error
	if err != nil {
		return nil, err
	}

	err = databases.DB.Model(&m).Where("id = ?", id).Updates(m).Error
	if err != nil {
		return nil, err
	}
	return m, nil
}

// Create 更新权限
func (m *Permissions) Create() error {
	create := databases.DB.Model(&m).Create(&m)
	if create.Error != nil {
		return create.Error
	}
	return nil
}

// Delete 删除权限
func (m *Permissions) Delete(id uint64) error {
	if err := databases.DB.Where("id = ?", id).Find(&m).Error; err != nil {
		return err
	}

	db := databases.DB.Model(&m).Where("id = ?", id).Delete(&m)
	if db.Error != nil {
		return db.Error
	}

	return nil
}
