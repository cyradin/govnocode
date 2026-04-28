package llm

import (
	"context"
	"fmt"

	"github.com/cyradin/govnocode/tools"
)

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type client interface {
	Generate(messages []ChatMessage) (string, error)
}

type toolProvider interface {
	Get(code string) (tools.Tool, error)
	All() []tools.Tool
}

type Agent struct {
	client client
	tools  toolProvider
}

func NewAgent(client client, tools toolProvider) *Agent {
	return &Agent{
		client: client,
		tools:  tools,
	}
}

func (a *Agent) Start(ctx context.Context) error {
	return fmt.Errorf("error")
}
