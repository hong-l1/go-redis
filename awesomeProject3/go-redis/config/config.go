package config

import (
	"github.com/spf13/viper"
)

var RedisConfig *Config

type Config struct {
	Self  string
	Peers []string
}

func init() {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		panic(err)
	}
	RedisConfig = &config
}
