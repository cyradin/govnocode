package git

import (
	"context"

	"github.com/cyradin/govnocode/tools/executor"
)

type Fetch struct{}

func NewFetch() *Fetch {
	return &Fetch{}
}

func (g *Fetch) Execute(
	ctx context.Context,
	e executor.Executor,
	raw []byte,
) (executor.Result, error) {
	return e.Execute(ctx, []string{
		"git",
		"fetch",
	}, nil)
}
