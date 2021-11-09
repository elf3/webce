package request

import "webce/internal/repositories/models/admins/admin"

type ReqAddAdmin struct {
	admin.Admin
	RoleIds []int64 `form:"roles_ids[]" validate:"required,min=1,max=9,dive,required"`
}

type ReqEditAdmin struct {
	admin.Admin
	ID      uint64  `form:"id" validate:"required"`
	RoleIds []int64 `form:"roles_ids[]" validate:"required,min=1,max=9,dive,required"`
}
