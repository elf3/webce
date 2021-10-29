package admin

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/core/router"
	"github.com/kataras/iris/v12/mvc"
	"net/http"
	"webce/cmd/web/conf"
	"webce/cmd/web/handlers/admin/controller/login"
	admin "webce/cmd/web/handlers/admin/controller/manager"
	"webce/cmd/web/handlers/admin/controller/node"
	"webce/cmd/web/middle"
	"webce/library/easycasbin"
	"webce/library/session"
)

// InitRouter Admin 路由
func InitRouter(app *iris.Application) {
	// 使用SESSION
	app.Use(session.NewSessionStore())
	// 静态模板类型
	app.HandleDir("/static", "./resources")
	// 使用VIEW模板
	templatesFS := iris.PrefixDir("./views", http.FS(conf.EmbedRoot))
	app.RegisterView(iris.HTML(templatesFS, ".html"))
	// 免登陆的路由
	app.PartyFunc("/admin", func(p router.Party) {
		// 显示登陆页面
		p.Get("/", login.ShowLogin)
		// 登陆后台请求接口
		p.Post("/login", login.Login)
	})

	// 使用中间件认证
	ntc := app.Party("/admin")
	{
		ntc.Use(middle.AuthAdmin(easycasbin.NotCheck("/admin/login", "/admin/logout", "/admin/user/create/admin")))
		mvc.New(ntc.Party("/user")).Handle(admin.NewManager())
		mvc.New(ntc.Party("/node")).Handle(node.NewNode())
	}

}
