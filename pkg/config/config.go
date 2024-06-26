package config

import (
	"errors"
	"github.com/spf13/viper"
)

type Config struct {
	JiraURL string `mapstructure:"jira_url"`
	Email   string `mapstructure:"email"`
	APIKey  string `mapstructure:"api_key"`
}

const CONFIG_PATH = "."
const CONFIG_NAME = "config"
const CONFIG_TYPE = "toml"

func Load() (*Config, error) {
	viper.AddConfigPath(CONFIG_PATH)
	viper.SetConfigName(CONFIG_NAME)
	viper.SetConfigType(CONFIG_TYPE)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {

		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found, run `lazyjira config init`")
		}

		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
