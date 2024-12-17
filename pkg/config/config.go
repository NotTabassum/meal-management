package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	DBUser     string `mapstructure:"DBUser"`
	DBPass     string `mapstructure:"DBPass"`
	DBIP       string `mapstructure:"DBIP"`
	DBName     string `mapstructure:"DBName"`
	DBPort     string `mapstructure:"DBPort"`
	ServerPort string `mapstructure:"ServerPort"`
}

func InitConfig() *Config {
	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error reading env file, %s", err)
	}
	var config *Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal("Error reading env file, %s", err)
	}
	return config
}

var LocalConfig *Config

func SetConfig() {
	LocalConfig = InitConfig()
}
