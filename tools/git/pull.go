package git

import (
	"context"

	"github.com/cyradin/govnocode/tools/executor"
)

type Pull struct{}

func NewPull() *Pull {
	return &Pull{}
}

func (g *Pull) Execute(
	ctx context.Context,
	e executor.Executor,
	raw []byte,
) (executor.Result, error) {
	return e.Execute(ctx, []string{
		"git",
		"pull",
	}, nil)
}
