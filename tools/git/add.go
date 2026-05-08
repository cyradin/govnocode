package git

import (
	"context"
	"fmt"

	"github.com/cyradin/govnocode/tools/executor"
)

type AddArgs struct {
	Path string `json:"path" validate:"required"`
}

type Add struct{}

func NewAdd() *Add {
	return &Add{}
}

func (g *Add) Execute(
	ctx context.Context,
	e executor.Executor,
	raw []byte,
) (executor.Result, error) {
	a, err := parseArgs[AddArgs](raw)
	if err != nil {
		return executor.Result{}, fmt.Errorf("parse args: %w", err)
	}

	return e.Execute(ctx, []string{
		"git",
		"add",
		a.Path,
	}, nil)
}
