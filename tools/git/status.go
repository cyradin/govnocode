package git

import (
	"context"

	"github.com/cyradin/govnocode/internal/command"
)

type Status struct{}

func NewStatus() *Status {
	return &Status{}
}

func (g *Status) Execute(
	ctx context.Context,
	executor command.Executor,
	raw []byte,
) (command.Result, error) {
	return executor.Execute(ctx, []string{
		"git",
		"status",
		"--porcelain",
	}, nil)
}
