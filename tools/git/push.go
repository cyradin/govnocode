package git

import (
	"context"
	"fmt"

	"github.com/cyradin/govnocode/internal/command"
)

type PushArgs struct {
	Remote string `json:"remote" validate:"required"`
	Branch string `json:"branch" validate:"required"`
}

type Push struct{}

func NewPush() *Push {
	return &Push{}
}

func (g *Push) Execute(
	ctx context.Context,
	executor command.Executor,
	raw []byte,
) (command.Result, error) {
	a, err := parseArgs[PushArgs](raw)
	if err != nil {
		return command.Result{}, fmt.Errorf("parse args: %w", err)
	}

	return executor.Execute(ctx, []string{
		"git",
		"push",
		"-u",
		a.Remote,
		a.Branch,
	}, nil)
}
