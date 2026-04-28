package container

import (
	"net/http"

	"github.com/cyradin/govnocode/internal/llm"
)

func (c *Container) Agent() *llm.Agent {
	if c.agent == nil {
		c.agent = llm.NewAgent(
			c.OllamaClient(),
			c.ToolRegistry(),
		)
	}

	return c.agent
}

func (c *Container) OllamaClient() *llm.OllamaClient {
	if c.ollamaClient == nil {
		c.ollamaClient = llm.NewOllamaClient(
			c.cfg.LLM.Ollama.BaseURL,
			c.cfg.LLM.Ollama.Model,
			&http.Client{
				Timeout: c.cfg.LLM.Ollama.HTTPTimeout,
			},
		)
	}

	return c.ollamaClient
}
