package config

import (
	"exchangeapp/global"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)
var (
	dbHost = os.Getenv("DB_HOST")
	dbPort = os.Getenv("DB_PORT")
	dbUser = os.Getenv("DB_USER")
	dbPassword = os.Getenv("DB_PASSWORD")
	dbName = os.Getenv("DB_NAME")
	
)
func InitDb() {
	// dsn := Appconfig.Database.Dsn
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",dbUser,dbPassword,dbHost,dbPort,dbName)
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
