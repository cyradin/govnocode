package git

import (
	"context"

	"github.com/cyradin/govnocode/internal/command"
)

type Pull struct{}

func NewPull() *Pull {
	return &Pull{}
}

func (g *Pull) Execute(
	ctx context.Context,
	executor command.Executor,
	raw []byte,
) (command.Result, error) {
	return executor.Execute(ctx, []string{
		"git",
		"pull",
	}, nil)
}
