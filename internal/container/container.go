package container

import (
	"log/slog"

	"github.com/cyradin/govnocode/internal/agent"
	"github.com/cyradin/govnocode/internal/config"
	"github.com/cyradin/govnocode/internal/llm"
	"github.com/cyradin/govnocode/pkg/logger"
)

type Container struct {
	version string
	cfg     *config.Config
	logger  *slog.Logger

	printer     *agent.Printer
	codingAgent *agent.CodingAgent

	ollamaClient *llm.OllamaClient
}

func New(
	version string,
	cfg *config.Config,
) *Container {
	return &Container{
		version: version,
		cfg:     cfg,
	}
}

func (c *Container) Logger() *slog.Logger {
	if c.logger == nil {
		c.logger = logger.New(c.version, c.cfg.Log.Level)
	}

	return c.logger
}

func (c *Container) Cfg() config.Config {
	return *c.cfg
}
