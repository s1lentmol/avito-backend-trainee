package config

import (
	"bytes"
	_ "embed"
	"fmt"
	"strings"

	"github.com/go-faster/errors"
	"github.com/spf13/viper"
)

//go:embed config.yaml
var defaultYAML []byte

type Config struct {
	DB struct {
		Username string
		Password string
		Host     string
		Port     string
		Name     string
	}
	App struct {
		Name string
		Port string
	}
}

func Load() (*Config, error) {
	var appConfig Config

	viper.SetConfigType("yaml")

	if err := viper.ReadConfig(bytes.NewBuffer(defaultYAML)); err != nil {
		return nil, errors.Wrap(err, "config")
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("env")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := viper.Unmarshal(&appConfig)
	if err != nil {
		return nil, errors.Wrap(err, "config")
	}

	return &appConfig, nil
}

func (cfg *Config) GetDSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		cfg.DB.Username,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Name,
	)
}
