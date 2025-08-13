package config

import (
	"exchangeapp/global"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDb() {
	dsn := Appconfig.Database.Dsn
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Database fail to init,got error %v", err)
	}
	sqlDB, err := db.DB()
	// 设置连接池中最大的连接数量
	sqlDB.SetMaxIdleConns(Appconfig.Database.MaxIdeConns)
	// 设置最大打开连接数
	sqlDB.SetMaxOpenConns(Appconfig.Database.MaxOpenConns)
	// 最大连接时间
	sqlDB.SetConnMaxIdleTime(time.Hour)

	if err != nil {
		log.Fatalf("Fail to config database,got error %v", err)
	}
	// 全局实例
	global.Db = db
}
