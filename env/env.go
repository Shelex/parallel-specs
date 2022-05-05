package env

import (
	"github.com/spf13/viper"
)

type Config struct {
	Env             string
	DbConnectionUrl string
	HttpPort        string
}

var Settings *Config

func ReadEnv() *Config {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("APP")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.SetConfigFile("config.env")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	viper.SetDefault("HTTP_PORT", "3000")
	viper.SetDefault("ENVIRONMENT", "dev")

	Settings = &Config{
		HttpPort:        viper.GetString("HTTP_PORT"),
		DbConnectionUrl: viper.GetString("DB_CONNECTION_URL"),
		Env:             viper.GetString("ENVIRONMENT"),
	}

	return Settings
}
