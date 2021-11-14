package roles

import (
	"webce/internal/repositories/models/admins/permissions"
)

// Roles 角色
type Roles struct {
	ID          uint64                    `gorm:"primary_key" json:"id" structs:"id"`
	Title       string                    `gorm:"type:varchar(50);unique_index" json:"title" form:"title"` // 角色标题
	Description string                    `gorm:"type:char(64);" json:"description" form:"description"`    // 角色注解
	Permissions []permissions.Permissions `gorm:"many2many:role_menu;" json:"permissions" `                // 关联权限
}
