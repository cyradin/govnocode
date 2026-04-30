package container

import (
	"net/http"
	"os"

	"github.com/cyradin/govnocode/internal/agent"
	"github.com/cyradin/govnocode/internal/llm"
)

func (c *Container) CodingAgent() *agent.CodingAgent {
	if c.agent == nil {
		c.agent = agent.NewCoding(
			os.Stdout,
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
