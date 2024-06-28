package config

import (
	"errors"
	"github.com/spf13/viper"
	"os"
	"path"
)

type Config struct {
	JiraURL     string `mapstructure:"jira_url"`
	Email       string `mapstructure:"email"`
	AccessToken string `mapstructure:"access_token"`
}

type ConfigService struct {
	v *viper.Viper
}

type ConfigProvider interface {
	Load() (*Config, error)
	Save(config Config) error
}

const CONFIG_DIR = ".config/lazyjira/"
const CONFIG_NAME = "config"
const CONFIG_TYPE = "toml"

func NewConfigService() ConfigProvider {
	return ConfigService{
		v: viper.New(),
	}
}

func (c ConfigService) getBasePath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "."
	}

	return path.Join(homeDir, CONFIG_DIR)
}

func (c ConfigService) getFullPath() string {
	return path.Join(c.getBasePath(), c.getFilename())
}

func (c ConfigService) getFilename() string {
	return CONFIG_NAME + "." + CONFIG_TYPE
}

func (c ConfigService) createIfNotExist() error {
	if err := os.MkdirAll(c.getBasePath(), 0755); err != nil {
		return err
	}

	_, err := os.Stat(c.getFullPath())

	if errors.Is(err, os.ErrNotExist) {
		file, err := os.OpenFile(c.getFullPath(), os.O_RDONLY|os.O_CREATE, 0644)
		if err != nil {
			return err
		}
		return file.Close()
	}

	return nil
}

func (c ConfigService) Load() (*Config, error) {
	c.v.AddConfigPath(c.getBasePath())
	c.v.SetConfigName(CONFIG_NAME)
	c.v.SetConfigType(CONFIG_TYPE)
	c.v.AutomaticEnv()

	if err := c.v.ReadInConfig(); err != nil {

		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found, run `lazyjira config init`")
		}

		return nil, err
	}

	var cfg Config
	if err := c.v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c ConfigService) Save(config Config) error {
	c.v.AddConfigPath(c.getBasePath())
	c.v.SetConfigName(CONFIG_NAME)
	c.v.SetConfigType(CONFIG_TYPE)

	c.v.SetDefault("jira_url", config.JiraURL)
	c.v.SetDefault("email", config.Email)
	c.v.SetDefault("access_token", config.AccessToken)

	err := c.createIfNotExist()
	if err != nil {
		return err
	}

	if err := c.v.WriteConfig(); err != nil {
		return err
	}

	return nil
}
