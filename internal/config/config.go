package config

import (
	"github.com/spf13/viper"
)

type DB struct {
	DSN     string `mapstructure:"dsn"`
	Timeout string `mapstructure:"timeout"`
}

type Config struct {
	DB *DB `mapstructure:"indexStorageDB"`
}

const (
	FileName = "application.yml"
)

func New(env string) (*Config, error) {
	viper.SetConfigType("yaml")
	viper.SetConfigFile(FileName)
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
