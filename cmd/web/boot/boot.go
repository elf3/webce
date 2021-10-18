package boot

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/spf13/viper"
	"time"
)

func Run(app *iris.Application) {
	// 是否开启TLS
	tls := viper.GetBool("tls")
	fmt.Println(tls)
	var err error
	if !tls {
		addr := viper.GetString("addr")
		err = app.Run(
			iris.Addr(addr),
			iris.WithoutServerError(iris.ErrServerClosed),
			iris.WithOptimizations,
			iris.WithTimeFormat(time.RFC3339),
		)
	} else {
		host := fmt.Sprintf("%s:%d", viper.GetString("host"), 443)
		err = app.Run(iris.TLS(host, viper.GetString("cert"), viper.GetString("key")))
	}
	if err != nil {
		fmt.Println(err)
	}
}
