package env

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Env             string
	DbConnectionUrl string
	Host            string
	Port            string
}

var Settings *Config

func ReadEnv() *Config {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("APP")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.SetConfigFile("config.env")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("failed to read config: %s", err)
	}

	viper.SetDefault("HTTP_PORT", "8080")
	viper.SetDefault("ENVIRONMENT", "dev")

	Settings = &Config{
		Host:            viper.GetString("HOST"),
		Port:            viper.GetString("PORT"),
		DbConnectionUrl: viper.GetString("DB_CONNECTION_URL"),
		Env:             viper.GetString("ENVIRONMENT"),
	}

	return Settings
}
