package config

import (
	"fmt"
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
	"golang.org/x/net/context"
)

type Config struct {
	DBUser string `env:"DBUser"`

	DBPass string `env:"DBPass"`

	DBIP string `env:"DBIP"`

	DBName string `env:"DBName"`

	DBPort string `env:"DBPort"`

	ServerPort string `env:"ServerPort"`
}

var LocalConfig *Config

//func Get() *Secret {
//	if secret == nil {
//		loadSecret()
//		validateSecret()
//	}
//	return secret
//}

func SetConfig() {
	LocalConfig = &Config{}
	err := godotenv.Load()
	if err != nil {
		log.Error(context.Background(), fmt.Sprintf(" Error loading .env file. Error: %v. ", err))
	}
	if err := env.Parse(LocalConfig); err != nil {
		panic(fmt.Sprintf("Error reading the environment variables: %v", err))
	}
}

//func validateSecret() {
//	v := reflect.ValueOf(*Config)
//	t := reflect.TypeOf(*secret)
//
//	for i := 0; i < t.NumField(); i++ {
//		val := v.Field(i)
//		typ := t.Field(i)
//		if val.IsZero() {
//			panic("secret " + typ.Name + " not found")
//		}
//	}
//}
