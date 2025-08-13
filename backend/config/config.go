package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	App struct {
		Name string
		Port string
	}
	Database struct {
		Dsn          string
		MaxIdeConns  int
		MaxOpenConns int
	}
}

var Appconfig *Config

func InitConfig() {
	// 配置文件的名字是什么
	viper.SetConfigName("config")
	// 配置文件的type
	viper.SetConfigType("yml")
	// 路径是什么？
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file %v", err)
	}

	Appconfig = &Config{}

	if err := viper.Unmarshal(Appconfig); err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
	}
	InitDb()
	InitRedis()
}
