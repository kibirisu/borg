package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	AppEnv      string `mapstructure:"appenv"`
	ListenPort  string `mapstructure:"listenport"`
	DatabaseURL string `mapstructure:"databaseurl"`
	JWTSecret   string `mapstructure:"jwtsecret"`
}

func GetConfig() *Config {
	viper.SetDefault("AppEnv", "prod")
	viper.SetDefault("ListenPort", "8080")
	viper.SetDefault("DatabaseUrl", "postgres://borg:borg@localhost:5432/borg")
	viper.SetDefault("JWTSecret", "changeme")
	viper.RegisterAlias("AppEnv", "app_env")
	viper.RegisterAlias("ListenPort", "listen_port")
	viper.RegisterAlias("DatabaseURL", "database_url")
	viper.RegisterAlias("JWTSecret", "jwt_secret")
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Println(err)
	}
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal(err)
	}
	return &config
}
