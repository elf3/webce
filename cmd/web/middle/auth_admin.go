package middle

import (
	"github.com/kataras/iris/v12"
	"net/http"
	"strings"
	"webce/internal/repositories/models/admins"
	"webce/internal/repositories/repo/adminrepo"
	"webce/pkg/library/apgs"
	"webce/pkg/library/easycasbin"
	"webce/pkg/library/jwt"
)

// AuthAdmin 中间件
func AuthAdmin(nocheck ...easycasbin.DontCheckFunc) iris.Handler {

	return func(c iris.Context) {
		if easycasbin.DontCheck(c, nocheck...) {
			c.Next()
			return
		}
		token, err := jwt.ParseToken(c.GetHeader("token"))
		if err != nil {
			c.JSON(apgs.ApiReturn(500, "error token", ""))
			return
		}
		//username := token["username"].(string)
		id := token["id"].(int)
		admin, err := adminrepo.NewAdminUserRepository().GetAdminById(id)
		if err != nil {
			c.Redirect("/admin/login", 302)
			return
		}
		// 超级管理员不验证权限
		if admin.IsSuper == 1 {
			c.Next()
			return
		}

		if len(admin.Roles) <= 0 || admin.Roles == nil || admin.ID <= 0 {
			c.JSON(apgs.ApiReturn(500, "permission denied", ""))
			_, _ = c.Problem(nil)
			return
		}
		_ = admin.LoadAllPolicy()

		var role admins.Roles
		_ = role.LoadAllPolicy()

		for _, i2 := range admin.Roles {
			role := i2.Title
			p := strings.ToLower(c.Path())
			m := strings.ToLower(c.Method())
			var b bool
			var err error

			if b, err = easycasbin.GetEnforcer().Enforce(role, p, m); err != nil {
				c.JSON(apgs.ApiReturn(http.StatusForbidden, "permission denied", ""))
				return
			}

			if !b {
				c.JSON(apgs.ApiReturn(http.StatusUnauthorized, "permission denied", ""))
				return
			}
		}

		c.Next()
	}
}
