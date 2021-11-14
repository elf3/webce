package admin

import (
	"webce/internal/repositories/models"
	"webce/internal/repositories/models/admins/roles"
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
