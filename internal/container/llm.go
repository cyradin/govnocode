package container

import (
	"net/http"
	"os"

	"github.com/cyradin/govnocode/internal/agent"
	"github.com/cyradin/govnocode/internal/llm"
	"github.com/cyradin/govnocode/tools"
)

func (c *Container) Printer() *agent.Printer {
	if c.printer == nil {
		c.printer = agent.NewPrinter(os.Stdout)
	}

	return c.printer
}

func (c *Container) CodingAgent() *agent.CodingAgent {
	if c.codingAgent == nil {
		toolRegistry := c.NewToolRegistry().
			MustRegister(tools.Git()...).
			MustRegister(tools.Filesystem()...)

		c.codingAgent = agent.NewCoding(
			c.Printer(),
			c.OllamaClient(),
			toolRegistry,
		)
	}

	return c.codingAgent
}

func (c *Container) OllamaClient() *llm.OllamaClient {
	if c.ollamaClient == nil {
		c.ollamaClient = llm.NewOllamaClient(
			c.cfg.LLM.Ollama.BaseURL,
			c.cfg.LLM.Ollama.Model,
			&http.Client{
				Timeout: c.cfg.LLM.Ollama.HTTPTimeout,
			},
			llm.OllamaOptions{
				Temperature:   c.cfg.LLM.Ollama.Temperature,
				TopP:          c.cfg.LLM.Ollama.TopP,
				TopK:          c.cfg.LLM.Ollama.TopK,
				RepeatPenalty: c.cfg.LLM.Ollama.RepeatPenalty,
				NumPredict:    c.cfg.LLM.Ollama.NumPredict,
			},
		)
	}

	return c.ollamaClient
}
