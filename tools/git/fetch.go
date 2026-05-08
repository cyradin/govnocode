package git

import (
	"context"

	"github.com/cyradin/govnocode/internal/command"
)

type Fetch struct{}

func NewFetch() *Fetch {
	return &Fetch{}
}

func (g *Fetch) Execute(
	ctx context.Context,
	executor command.Executor,
	raw []byte,
) (command.Result, error) {
	return executor.Execute(ctx, []string{
		"git",
		"fetch",
	}, nil)
}
