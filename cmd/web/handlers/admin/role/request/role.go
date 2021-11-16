package request

import (
	"webce/internal/repositories/models/admins/roles"
)

// ReqAddRole 添加角色
type ReqAddRole struct {
	roles.Roles
	PermissionIds []uint64 `form:"permission_ids[]" validate:"required,dive,required"`
}

// ReqEditRole 编辑角色
type ReqEditRole struct {
	roles.Roles
	Id            uint64   `form:"id" validate:"required,gt=0"`
	PermissionIds []uint64 `form:"permission_ids[]" validate:"required,dive,required"`
}

type ReqSearchRole struct {
	ID    uint64 `gorm:"primary_key" json:"id" structs:"id"`
	Title string `gorm:"type:varchar(50);unique_index" json:"title" form:"title"` // 角色标题
}
