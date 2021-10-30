package migrate

import (
	"github.com/sirupsen/logrus"
	admins2 "webce/internal/repositories/models/admins"
	business2 "webce/internal/repositories/models/business"
	"webce/pkg/library/databases"
)

var models = []interface{}{
	&admins2.Admin{},
	&admins2.Permissions{},
	&admins2.Roles{},
	&business2.Domains{},
	&business2.Node{},
}

// 数据自动填充
func AutoMigrate() {
	db := databases.GetDB()
	// 自动创建数据库
	if err := db.Set("gorm:table_options", "ENGINE=Innodb").AutoMigrate(models...); nil != err {
		logrus.Fatal("auto migrate tables failed: ", err)
	}
}
