package middle

import (
	"github.com/kataras/iris/v12"
	"strconv"
	"strings"
	api "webce/api/admin"
	"webce/api/role"
	"webce/internal/repositories/repo/adminrepo"
	"webce/pkg/library/easycasbin"
	"webce/pkg/library/jwt"
	"webce/pkg/library/log"
	"webce/pkg/library/resp"
)

func getToken(token string) string {
	authHeaderParts := strings.Split(token, " ")
	if len(authHeaderParts) != 2 {
		return ""
	}
	return authHeaderParts[1]
}

// AuthAdmin 中间件
func AuthAdmin(nocheck ...easycasbin.DontCheckFunc) iris.Handler {

	return func(c iris.Context) {
		if easycasbin.DontCheck(c, nocheck...) {
			c.Next()
			return
		}
		authToken := getToken(c.GetHeader("Authorization"))
		if authToken == "" {
			log.Log.Error("not auth: ")
			resp.Error(c, 500, "not auth")
			return
		}
		token, err := jwt.ParseToken(authToken)
		if err != nil {
			log.Log.Error("parse token err: ", err)
			resp.Error(c, 500, "error token")
			return
		}
		//username := token["username"].(string)
		id, ok := token["userId"]
		if !ok {
			log.Log.Error("parse token userId err: ", ok)
			resp.Error(c, 500, "error token")
			return
		}

		adminIdStr := id.(string)
		adminId, err := strconv.ParseUint(adminIdStr, 10, 64)
		if err != nil {
			log.Log.Error("parse token userId err: ", ok)
			resp.Error(c, 500, "error token")
			return
		}
		apiAdmin := api.ApiAdmin{}
		admin, err := adminrepo.NewAdminUserRepository().GetAdminById(adminId)
		if err != nil {
			resp.Error(c, 500, "error user")
			return
		}
		// 超级管理员不验证权限
		if admin.IsSuper == 1 {
			c.Next()
			return
		}

		if len(admin.Roles) <= 0 || admin.Roles == nil || admin.ID <= 0 {
			resp.Error(c, 500, "permission denied")
			return
		}
		err = apiAdmin.LoadAdminPolicy(admin.ID)
		if err != nil {
			log.Log.Error("load permission error : ", err)
			resp.Error(c, 500, "load permission denied")
			return
		}
		var role role.ApiRole
		_ = role.LoadAllPolicy()

		for _, i2 := range admin.Roles {
			role := i2.Title
			p := strings.ToLower(c.Path())
			m := strings.ToLower(c.Method())
			var b bool
			var err error

			if b, err = easycasbin.GetEnforcer().Enforce(role, p, m); err != nil {
				resp.Error(c, 500, "permission denied")
				return
			}

			if !b {
				resp.Error(c, 500, "permission denied")
				return
			}
		}

		c.Next()
	}
}
