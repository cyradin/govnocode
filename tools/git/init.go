package git

import (
	"context"

	"github.com/cyradin/govnocode/internal/command"
)

type Init struct{}

func NewInit() *Init {
	return &Init{}
}

func (g *Init) Execute(
	ctx context.Context,
	executor command.Executor,
	args []byte,
) (command.Result, error) {
	return executor.Execute(ctx, []string{
		"git",
		"init",
	}, nil)
}
