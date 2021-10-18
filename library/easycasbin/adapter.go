package easycasbin

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	adapter "github.com/casbin/gorm-adapter/v3"
	"github.com/kataras/iris/v12"
	"github.com/spf13/viper"
	"webce/library/databases"
)

// Casbin is the auth services which contains the easycasbin enforcer.

var Enfocer *casbin.SyncedEnforcer

// GetEnforcer 获取已经实例化的对象
func GetEnforcer() *casbin.SyncedEnforcer {
	return Enfocer
}

// InitAdapter 初始化权限 数据库适配器
func InitAdapter() (*casbin.SyncedEnforcer, error) {
	//TODO 因为 前缀问题，此处要写死
	a, err := adapter.NewAdapterByDBUseTableName(databases.GetDB(), "wk", "casbin_rule")
	if err != nil {
		return nil, fmt.Errorf("can not Init: %v", err.Error())
	}
	e, err := casbin.NewSyncedEnforcer("./conf/rbac_model.conf", a)
	if err != nil {
		return nil, fmt.Errorf("can not Init: %v", err.Error())
	}
	// 开启AutoSave机制
	e.EnableAutoSave(true)
	_ = e.BuildRoleLinks()
	enableLog := viper.GetBool("casbin.debug")
	e.EnableLog(enableLog)
	// 10秒重新加载一次权限
	//e.StartAutoLoadPolicy(10 * time.Second)
	//e.EnableAutoBuildRoleLinks(true)
	// 因为开启了AutoSave机制，现在内存中的改变会同步回写到持久层中
	//e.AddPolicy("admin", "test", "test")
	Enfocer = e
	return e, err
}

type DontCheckFunc func(ctx iris.Context) bool

// NotCheck 指定路由不用检查
func NotCheck(prefixes ...string) DontCheckFunc {
	return func(c iris.Context) bool {
		path := c.GetCurrentRoute().Path()
		pathLen := len(path)
		for _, p := range prefixes {
			if pl := len(p); pathLen >= pl && path[:pl] == p {
				return true
			}
		}
		return false
	}
}

// Check 指定路由检查
func Check(prefixes ...string) DontCheckFunc {
	return func(c iris.Context) bool {
		// 获取路由地址
		path := c.GetCurrentRoute().Path()
		pathLen := len(path)
		for _, p := range prefixes {
			if pl := len(p); pathLen >= pl && path[:pl] == p {
				return false
			}
		}
		return true
	}
}

// DontCheck 不检查函数
func DontCheck(c iris.Context, skippers ...DontCheckFunc) bool {
	for _, skipper := range skippers {
		if skipper(c) {
			return true
		}
	}
	return false
}
