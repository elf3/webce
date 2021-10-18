package session

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/sessions"
	"github.com/kataras/iris/v12/sessions/sessiondb/redis"
	"github.com/spf13/viper"
	"time"
)

// 初始化Session 方式 始终会返回一个数据存储方式，默认为 cookie
func NewSessionStore() iris.Handler {
	typeOf := viper.GetString("session.type")
	name := viper.GetString("session.name")
	// 判断
	sess := sessions.New(sessions.Config{
		Cookie:          name,
		Expires:         0, // defaults to 0: unlimited life. Another good value is: 45 * time.Minute,
		AllowReclaim:    true,
		CookieSecureTLS: true,
	})
	fmt.Println(typeOf)
	switch typeOf {
	case "redis":
		store := redisStore()
		sess.UseDatabase(store)
	}
	return sess.Handler()
}

// Redis
func redisStore() sessions.Database {
	host := viper.GetString("redis.host")
	port := viper.GetString("redis.port")
	pass := viper.GetString("redis.pass")
	db := redis.New(redis.Config{
		Network:   "tcp",
		Addr:      host + ":" + port,
		Timeout:   time.Duration(30) * time.Second,
		MaxActive: 10,
		Password:  pass,
		Database:  "",
		Prefix:    "session:",
		Driver:    redis.GoRedis(), // defaults.
	})
	return db
}

func GetSession(c iris.Context, key string) interface{} {
	session := sessions.Get(c)
	get := session.Get(key)
	return get
}
