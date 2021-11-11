package request

import (
	"webce/internal/repositories/models/admins/roles"
)

type ReqAddRole struct {
	roles.Roles
	PermissionIds []uint64 `form:"permission_ids[]" validate:"required,dive,required"`
}

type ReqEditRole struct {
	roles.Roles
	Id            uint64   `form:"id" validate:"required,gt=0"`
	PermissionIds []uint64 `form:"permission_ids[]" validate:"required,dive,required"`
}
