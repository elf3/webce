package router

import (
	"github.com/kataras/iris/v12"
	recover2 "github.com/kataras/iris/v12/middleware/recover"
	"github.com/spf13/viper"
	"webce/apis/migrate"
	"webce/cmd/web/boot/router/admin"
	"webce/cmd/web/boot/router/api"
	"webce/cmd/web/middle"
	"webce/library/config"
	"webce/library/databases"
	"webce/library/log"
)

func InitRouter() *iris.Application {
	errConf := config.Init("./configs/app.yaml")
	if errConf != nil {
		panic(errConf)
	}
	log.InitLogger(log.LogConfig{
		Level:   "",
		Path:    "./",
		MaxSize: 1,
	})
	app := iris.New()
	app.Configure(iris.WithConfiguration(iris.Configuration{
		LogLevel: viper.GetString("runmode"), // 设置日志级别
		RemoteAddrHeaders: []string{
			"X-Real-Ip",
			"X-Forwarded-For",
			"CF-Connecting-IP",
			"True-Client-Ip",
		},
	}))
	app.Use(middle.LoggerHandler)

	// 初始化DB
	databases.InitDB()

	// 重启
	app.Use(recover2.New())
	// API 路由
	api.InitRouter(app)
	// admin 路由
	admin.InitRouter(app)
	migrate.AutoMigrate()

	return app

}
