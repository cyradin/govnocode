package git

import (
	"context"
	"fmt"

	"github.com/cyradin/govnocode/tools/executor"
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
	e executor.Executor,
	raw []byte,
) (executor.Result, error) {
	a, err := parseArgs[PushArgs](raw)
	if err != nil {
		return executor.Result{}, fmt.Errorf("parse args: %w", err)
	}

	return e.Execute(ctx, []string{
		"git",
		"push",
		"-u",
		a.Remote,
		a.Branch,
	}, nil)
}
