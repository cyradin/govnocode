package git

import (
	"context"

	"github.com/cyradin/govnocode/tools/executor"
)

type Init struct{}

func NewInit() *Init {
	return &Init{}
}

func (g *Init) Execute(
	ctx context.Context,
	e executor.Executor,
	args []byte,
) (executor.Result, error) {
	return e.Execute(ctx, []string{
		"git",
		"init",
	}, nil)
}
