package request

import (
	"github.com/kataras/iris/v12"
	"webce/internal/repositories/models/admins/admin"
)

type ReqAddAdmin struct {
	admin.Admin
	RoleIds []int64 `form:"roles_ids[]" validate:"required,min=1,max=9,dive,required"`
}

type ReqEditAdmin struct {
	admin.Admin
	ID      uint64  `form:"id" validate:"required"`
	RoleIds []int64 `form:"roles_ids[]" validate:"required,min=1,max=9,dive,required"`
}

type ReqSearchAdmin struct {
	Username string `form:"username"`
}

func GetAdminSearchMap(c iris.Context) map[string]interface{} {
	where := map[string]interface{}{}
	req := ReqSearchAdmin{}
	err := c.ReadQuery(&req)
	if err != nil {
		return where
	}
	if req.Username != "" {
		where["username like"] = req.Username + "%"
	}
	return where
}
