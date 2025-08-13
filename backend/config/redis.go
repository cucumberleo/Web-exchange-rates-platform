package config

import (
	"exchangeapp/global"
	"log"

	"github.com/go-redis/redis"
)

func InitRedis(){
	RedisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB: 0,
		Password: "",
	})
	_,err :=RedisClient.Ping().Result()
	if err!=nil{
		log.Fatalf("Fail to connect Redis,got error: %v",err)
	}
	global.Redisdb = RedisClient
}