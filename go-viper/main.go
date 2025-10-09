package main

import (
	"github.com/spf13/viper"
)

type Config struct {
	Env    string `mapstructure:"ENV"`
	Port   int    `mapstructure:"APP_PORT"`
	DBHost string `mapstructure:"DB_HOST"`
	DBUser string `mapstructure:"DB_USER"`
	DBPass string `mapstructure:"DB_PASS"`
	// add more env variables here as your project grows
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
