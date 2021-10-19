package admin

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/core/router"
	"github.com/kataras/iris/v12/mvc"
	"net/http"
	"webce/cmd/web/conf"
	admin "webce/cmd/web/handlers/admin/controller"
	"webce/cmd/web/middle"
	"webce/library/easycasbin"
	"webce/library/session"
)

// InitRouter Admin 路由
func InitRouter(app *iris.Application) {
	// 使用SESSION
	app.Use(session.NewSessionStore())
	// 使用VIEW模板
	fmt.Println(http.FS(conf.EmbedRoot))
	templatesFS := iris.PrefixDir("./views", http.FS(conf.EmbedRoot))
	app.RegisterView(iris.HTML(templatesFS, ".html"))
	fmt.Println("初始化路由")
	// 免登陆的路由
	app.PartyFunc("/admin", func(p router.Party) {
		p.Get("/", admin.Main)
		p.Get("/login", admin.Login)
		p.Post("/login", admin.Login)
	})

	// 使用中间件认证
	ntc := app.Party("/admin")
	{
		ntc.Use(middle.AuthAdmin(easycasbin.NotCheck("/admin/login", "/admin/logout")))
		mvc.New(ntc.Party("/user")).Handle(admin.NewManager())
	}

}
