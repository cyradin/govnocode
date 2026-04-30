package container

import (
	"log/slog"

	"github.com/cyradin/govnocode/internal/agent"
	"github.com/cyradin/govnocode/internal/config"
	"github.com/cyradin/govnocode/internal/llm"
	"github.com/cyradin/govnocode/pkg/logger"
	"github.com/cyradin/govnocode/tools"
)

type Container struct {
	version string
	cfg     *config.Config
	logger  *slog.Logger

	printer *agent.Printer
	agent   *agent.CodingAgent

	toolRegistry *tools.Registry
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
