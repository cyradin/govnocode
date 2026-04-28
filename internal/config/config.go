package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Log LogConfig

	LLM LLMConfig
}

type LogConfig struct {
	Level string `envconfig:"GOVNOCODE_LOG_LEVEL" required:"true" default:"info"`
}

type LLMConfig struct {
	Ollama struct {
		BaseURL string `envconfig:"GOVNOCODE_OLLAMA_BASE_URL" required:"true" default:"localhost:11434"`
		Modell  string `envconfig:"GOVNOCODE_OLLAMA_MODEL" required:"true"`
	}
}

func New() (*Config, error) {
	var cfg Config

	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	return &cfg, nil
}
