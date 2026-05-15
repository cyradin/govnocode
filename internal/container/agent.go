package container

import (
	"log/slog"

	"github.com/cyradin/govnocode/internal/agent/goagent"
)

func (c *Container) GoAgent() *goagent.Builder {
	return goagent.NewBuilder(
		c.Printer(),
		c.OllamaClient(),
		c.logger.With(slog.String("name", "coding_agent")),
	)
}
