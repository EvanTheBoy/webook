package ioc

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"webook/config"
	"webook/internal/repository/dao"
)

func InitDB() *gorm.DB {
	// 初始化数据库操作需要的组件
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	if err != nil {
		// 结束goroutine
		// 一旦初始化过程中出错, 应用就不要启动
		panic(err)
	}
	// 建表
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}
