package ioc

import (
	"webook/internal/repository/dao"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	// 先连接到 MySQL 服务器（不指定数据库）以创建数据库
	// dsnWithoutDB := "root:123456@tcp(127.0.0.1:13316)/"
	// dbWithoutDB, err := gorm.Open(mysql.Open(dsnWithoutDB))
	// if err != nil {
	// 	panic(err)
	// }
	// // 创建数据库（如果不存在）
	// err = dbWithoutDB.Exec("CREATE DATABASE IF NOT EXISTS webook").Error
	// if err != nil {
	// 	panic(err)
	// }
	// // 关闭临时连接
	// sqlDB, _ := dbWithoutDB.DB()
	// sqlDB.Close()

	// 连接到 webook 数据库
	dsn := viper.GetString("db.mysql.DSN")

	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		//只会在初始化的时候panic
		//panic相当于整个goroutine结束
		//一旦初始化出错 应用就不要启动
		panic(err)
	}
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}
