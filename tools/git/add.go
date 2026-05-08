package git

import (
	"context"
	"fmt"

	"github.com/cyradin/govnocode/internal/command"
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
	executor command.Executor,
	raw []byte,
) (command.Result, error) {
	a, err := parseArgs[AddArgs](raw)
	if err != nil {
		return command.Result{}, fmt.Errorf("parse args: %w", err)
	}

	return executor.Execute(ctx, []string{
		"git",
		"add",
		a.Path,
	}, nil)
}
