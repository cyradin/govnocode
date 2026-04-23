package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Log LogConfig
}

type LogConfig struct {
	Level string `envconfig:"FIXIK_LOG_LEVEL" required:"true" default:"info"`
}

func New() (*Config, error) {
	var cfg Config

	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	return &cfg, nil
}
