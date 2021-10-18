package middle

import (
	"encoding/json"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/sessions"
	"net/http"
	"strconv"
	"strings"
	"webce/apis/repositories/models/admins"
	"webce/library/easycasbin"
	"webce/library/session"
)

const (
	UserKey = "userID"
)

// 登录Session 中间件
func AuthSessionMiddle() iris.Handler {

	return func(c iris.Context) {
		sess := sessions.Get(c)
		userId, err := sess.GetInt64(UserKey)
		if userId <= 0 {
			c.Redirect("/admin/login", 302)
			_, _ = c.Problem("还没有登录")
			return
		}
		if err != nil {
			c.Redirect("/admin/login", 302)
			_, _ = c.Problem("发生错误")
			return
		}
		c.SetCookieKV(UserKey, strconv.FormatInt(userId, 10))
		c.Next()
	}
}

// AuthAdmin 中间件
func AuthAdmin(nocheck ...easycasbin.DontCheckFunc) iris.Handler {
	// Casbin
	enforcer, err := easycasbin.InitAdapter()
	if err != nil {
		panic(err)
	}
	return func(c iris.Context) {
		if easycasbin.DontCheck(c, nocheck...) {
			c.Next()
			return
		}
		// Session 判断权限
		get := session.GetSession(c, UserKey)

		if get == nil {
			c.Redirect("/admin/login", 302)
			return
		}

		var admin admins.Admin
		canGet := json.Unmarshal(get.([]byte), &admin)
		// 超级管理员不验证权限
		if admin.IsSuper == 1 {
			c.Next()
			return
		}

		if len(admin.Roles) <= 0 || admin.Roles == nil || admin.ID <= 0 || canGet != nil {
			//ginview.HTML(c, http.StatusUnauthorized, "err/401", apgs.NewApiRedirect(200, "登录异常", "/admin/login"))
			//c.Abort()
			c.Redirect("/admin/login", 302)
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

			if b, err = enforcer.Enforce(role, p, m); err != nil {
				// TODO 判断是是否为调试模式
				// TODO 调试模式下 判断 异步，同步 返回 JSON HTML
				//c.JSON(403, helpers.NewApiReturn(401, err.Error(), b))
				//c.AbortWithStatus(403)
				//ginview.HTML(c, http.StatusForbidden, "err/403", apgs.NewApiReturn(403, err.Error(), nil))
				//c.Abort()
				c.Redirect("/admin/login", http.StatusForbidden)
				_, _ = c.Problem(nil)
				break

			}

			if !b {
				//c.JSON(401, helpers.NewApiReturn(401, "权限验证失败", b))
				//c.Abort()
				//fmt.Println("Check:" + strconv.FormatBool(b))
				//c.Redirect(302, "/admin/login")
				//ginview.HTML(c, http.StatusUnauthorized, "err/401", apgs.NewApiRedirect(200, "无权限访问该内容", "/admin/login"))
				//c.Abort()
				c.Redirect("/admin/login", http.StatusUnauthorized)
				_, _ = c.Problem(nil)
				break
			}
		}

		c.Next()
	}
}
