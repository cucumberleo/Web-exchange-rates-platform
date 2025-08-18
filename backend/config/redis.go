package config

import (
	"exchangeapp/global"
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis"
)
var(
	redisHost = os.Getenv("REDIS_HOST")
	redisPort = os.Getenv("REDIS_PORT")
)
func InitRedis(){
	RedisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s",redisHost,redisPort),
		DB: 0,
		Password: "",
	})
	_,err :=RedisClient.Ping().Result()
	if err!=nil{
		log.Fatalf("Fail to connect Redis,got error: %v",err)
	}
	global.Redisdb = RedisClient
}