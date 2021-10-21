package requests

import "webce/apis/repositories/models/admins"

type ReqManage struct {
	admins.Admin
	RoleIds []int64 `json:"role_ids" form:"[]role_ids" validate:"min=0,max=9,dive,required"`
}
