package databases

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
)

var DB *gorm.DB

// 初始化DB
func InitDB() {
	// 数据库初始化
	//DbType := viper.GetString("db.type")
	host := viper.GetString("db.host")
	user := viper.GetString("db.user")
	pass := viper.GetString("db.pass")
	dbname := viper.GetString("db.dbname")
	charset := viper.GetString("db.charset")
	loc := viper.GetString("db.loc")
	native := viper.GetString("db.native")
	prefix := viper.GetString("db.prefix")
	debug := viper.GetInt("db.debug")
	var err error
	dabs := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=True&loc=%s&allowNativePasswords=%s",
		user, pass, host,
		dbname, charset, loc,
		native,
	)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second,            // 慢 SQL 阈值
			LogLevel:      logger.LogLevel(debug), // Log level
			Colorful:      false,                  // 禁用彩色打印
		},
	)

	// 设置数据库连接数
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dabs,  // DSN data source name
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   prefix, // 表名前缀，`User` 的表名应该是 `t_users`
			SingularTable: true,   // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `t_user`
		},
		//Logger: logger.Default.LogMode(logger.LogLevel(debug)),
		Logger:                 newLogger,
		SkipDefaultTransaction: true, // 禁用默认事务
		PrepareStmt:            true, // 执行任何 SQL 时都创建 prepared statement 并缓存，可以提高后续的调用速度

	})
	//db, err := gorm.Open(mysql.Open(dabs), &gorm.Config{})
	if err != nil {
		logrus.Fatal("Cannot Connect : " + err.Error())
	}
	if db == nil {
		logrus.Fatal("数据库初始化错误")
	}

	sqlDB, err := db.DB()
	if sqlDB == nil {
		logrus.Fatal("无法获取到底层链接")
	}
	if err != nil {
		logrus.Fatal("获取database/sql错误")
	}
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)
	// 全局禁用表名复数
	//db.SingularTable(true)
	//db.LogMode(debug)
	//db.DB().SetMaxIdleConns(100)
	//db.DB().SetMaxOpenConns(1000)
	//db.DB().SetConnMaxLifetime(5 * time.Minute)

	//gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
	//	return prefix + defaultTableName
	//}
	DB = db

}

func GetPrefix() string {
	return viper.GetString("db.prefix")
}
func GetDB() *gorm.DB {
	return DB
}

// 关闭数据库连接
func CloseDB() {
	//if err := DB.Di; nil != err {
	//	log.Fatal("Disconnect from database failed: " + err.Error())
	//}
}
