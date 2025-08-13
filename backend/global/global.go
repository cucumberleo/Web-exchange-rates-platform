package global

import (
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

// 全局变量
var (
	Db *gorm.DB
	Redisdb *redis.Client
)

