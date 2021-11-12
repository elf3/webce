package admin

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/core/router"
	"github.com/kataras/iris/v12/mvc"
	"webce/cmd/web/handlers/admin/admin"
	"webce/cmd/web/handlers/admin/login"
	"webce/cmd/web/handlers/admin/node"
	"webce/cmd/web/handlers/admin/permission"
	"webce/cmd/web/handlers/admin/role"
	"webce/cmd/web/middle"
	"webce/pkg/library/easycasbin"
)

// InitRouter Admin 路由
func InitRouter(app *iris.Application) {
	// 使用SESSION
	//app.Use(session.NewSessionStore())
	// 静态模板类型
	//app.HandleDir("/static", "./resources")
	// 使用VIEW模板
	//templatesFS := iris.PrefixDir("./views", http.FS(conf.EmbedRoot))
	//app.RegisterView(iris.HTML(templatesFS, ".html"))
	// 免登陆的路由
	app.PartyFunc("/admin", func(p router.Party) {
		// 登陆后台请求接口
		mvc.New(p.Party("/account")).Handle(login.NewHandlerLogin())
	})

	// 使用中间件认证
	ntc := app.Party("/admin")
	{
		ntc.Use(middle.AuthAdmin(easycasbin.NotCheck("/admin/login", "/admin/logout")))
		mvc.New(ntc.Party("/node")).Handle(node.NewNode())
		mvc.New(ntc.Party("/permission")).Handle(permission.NewPermissionHandler())
		mvc.New(ntc.Party("/admin")).Handle(admin.NewAdminHandler())
		mvc.New(ntc.Party("/role")).Handle(role.NewRoleHandler())
	}

}
