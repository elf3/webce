package router

import (
	"github.com/kataras/iris/v12"
	recover2 "github.com/kataras/iris/v12/middleware/recover"
	"github.com/spf13/viper"
	"webce/apis/migrate"
	"webce/cmd/web/boot/router/admin"
	"webce/cmd/web/boot/router/api"
	"webce/library/config"
	"webce/library/databases"
)

func InitRouter() *iris.Application {
	errConf := config.Init("./configs/app.yaml")
	if errConf != nil {
		panic(errConf)
	}
	databases.InitDB()

	app := iris.New()
	// 设置日志级别
	app.Logger().SetLevel(viper.GetString("runmode"))
	// 初始化DB
	// 重启
	app.Use(recover2.New())
	// API 路由
	api.InitRouter(app)
	// admin 路由
	admin.InitRouter(app)
	migrate.AutoMigrate()

	return app

}
