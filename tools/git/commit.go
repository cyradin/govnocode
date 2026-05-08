package git

import (
	"context"
	"fmt"

	"github.com/cyradin/govnocode/tools/executor"
)

type CommitArgs struct {
	Message string `json:"message" validate:"required"`
}

type Commit struct{}

func NewCommit() *Commit {
	return &Commit{}
}

func (g *Commit) Execute(
	ctx context.Context,
	e executor.Executor,
	raw []byte,
) (executor.Result, error) {
	a, err := parseArgs[CommitArgs](raw)
	if err != nil {
		return executor.Result{}, fmt.Errorf("parse args: %w", err)
	}

	return e.Execute(ctx, []string{
		"git",
		"commit",
		"-m",
		a.Message,
	}, nil)
}
