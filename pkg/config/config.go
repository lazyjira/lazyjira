package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	JiraURL string `mapstructure:"jira_url"`
	Email   string `mapstructure:"email"`
	APIKey  string `mapstructure:"api_key"`
}

func Load() (*Config, error) {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
