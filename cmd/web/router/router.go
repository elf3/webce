package router

import (
	"github.com/kataras/iris/v12"
	recover2 "github.com/kataras/iris/v12/middleware/recover"
	"github.com/spf13/viper"
	"webce/cmd/web/middle"
	"webce/cmd/web/router/admin"
	"webce/cmd/web/router/api"
	"webce/internal/migrate"
	"webce/pkg/library/config"
	"webce/pkg/library/databases"
	"webce/pkg/library/easycasbin"
	"webce/pkg/library/jwt"
	"webce/pkg/library/log"
)

func InitRouter() *iris.Application {
	errConf := config.Init("./configs/app.yaml")
	if errConf != nil {
		panic(errConf)
	}
	log.InitLogger()
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
	// 重启
	app.Use(recover2.New())
	app.Use(middle.LoggerHandler)
	jwt.InitJwtConf()

	// 初始化DB
	databases.InitDB()
	// 初始化 casbin
	easycasbin.InitAdapter()

	// API 路由
	api.InitRouter(app)
	// admin 路由
	admin.InitRouter(app)
	// 数据填充
	migrate.AutoMigrate()

	return app

}
