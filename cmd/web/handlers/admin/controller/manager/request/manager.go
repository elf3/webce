package request

import (
	"webce/internal/repositories/models/admins/admin"
)

type ReqManage struct {
	admin.Admin
	RoleIds []int64 `json:"role_ids" form:"[]role_ids" validate:"min=0,max=9,dive,required"`
}
