package git

import (
	"context"
	"fmt"

	"github.com/cyradin/govnocode/internal/command"
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
	executor command.Executor,
	raw []byte,
) (command.Result, error) {
	a, err := parseArgs[CommitArgs](raw)
	if err != nil {
		return command.Result{}, fmt.Errorf("parse args: %w", err)
	}

	return executor.Execute(ctx, []string{
		"git",
		"commit",
		"-m",
		a.Message,
	}, nil)
}
