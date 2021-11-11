package migrate

import (
	"github.com/sirupsen/logrus"
	"webce/internal/repositories/models/admins/admin"
	"webce/internal/repositories/models/admins/permissions"
	"webce/internal/repositories/models/admins/roles"
	"webce/internal/repositories/models/business"
	"webce/pkg/library/databases"
)

var models = []interface{}{
	&admin.Admin{},
	&permissions.Permissions{},
	&roles.Roles{},
	&business.Domains{},
	&business.Node{},
}

// 数据自动填充
func AutoMigrate() {
	db := databases.GetDB()
	// 自动创建数据库
	if err := db.Set("gorm:table_options", "ENGINE=Innodb").AutoMigrate(models...); nil != err {
		logrus.Fatal("auto migrate tables failed: ", err)
	}
}
