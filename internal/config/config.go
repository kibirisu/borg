package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	AppEnv      string `mapstructure:"appenv"`
	ListenHost  string `mapstructure:"listenhost"`
	ListenPort  string `mapstructure:"listenport"`
	Address     string `mapstructure:"address"`
	DatabaseURL string `mapstructure:"databaseurl"`
	JWTSecret   string `mapstructure:"jwtsecret"`
}

func GetConfig() *Config {
	viper.SetDefault("AppEnv", "prod")
	viper.SetDefault("ListenHost", "0.0.0.0")
	viper.SetDefault("ListenPort", "8080")
	viper.SetDefault("DatabaseUrl", "postgres://borg:borg@localhost:5432/borg")
	viper.SetDefault("Address", "localhost:8080")
	viper.SetDefault("JWTSecret", "changeme")
	viper.RegisterAlias("AppEnv", "app_env")
	viper.RegisterAlias("ListenHost", "listen_host")
	viper.RegisterAlias("ListenPort", "listen_port")
	viper.RegisterAlias("Address", "address")
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
