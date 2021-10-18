package migrate

import (
	"github.com/sirupsen/logrus"
	"webce/apis/repositories/models/admins"
	"webce/library/databases"
)

var models = []interface{}{
	&admins.Admin{},
	&admins.Permissions{},
	&admins.Roles{},
}

// 数据自动填充
func AutoMigrate() {
	db := databases.GetDB()
	// 自动创建数据库
	if err := db.Set("gorm:table_options", "ENGINE=Innodb").AutoMigrate(models...); nil != err {
		logrus.Fatal("auto migrate tables failed: ", err)
	}
}
