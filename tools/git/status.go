package git

import (
	"context"

	"github.com/cyradin/govnocode/tools/executor"
)

type Status struct{}

func NewStatus() *Status {
	return &Status{}
}

func (g *Status) Execute(
	ctx context.Context,
	e executor.Executor,
	raw []byte,
) (executor.Result, error) {
	return e.Execute(ctx, []string{
		"git",
		"status",
		"--porcelain",
	}, nil)
}
