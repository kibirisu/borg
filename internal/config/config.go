package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	AppEnv      string `mapstructure:"appenv"`
	ListenPort  string `mapstructure:"listenport"`
	DatabaseURL string `mapstructure:"databaseurl"`
}

func GetConfig() *Config {
	viper.SetDefault("AppEnv", "prod")
	viper.SetDefault("ListenPort", "8080")
	// TODO make default url generic
	viper.SetDefault("DatabaseUrl", "postgres://postgres@localhost:5432/dev?sslmode=disable")
	viper.RegisterAlias("AppEnv", "app_env")
	viper.RegisterAlias("ListenPort", "listen_port")
	viper.RegisterAlias("DatabaseURL", "database_url")
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
