package files

import (
	"context"
	"fmt"

	"github.com/cyradin/govnocode/internal/command"
)

type MkdirArgs struct {
	Path string `json:"path" validate:"required"`
}

type Mkdir struct{}

func NewMkdir() *Mkdir {
	return &Mkdir{}
}

func (m *Mkdir) Execute(
	ctx context.Context,
	executor command.Executor,
	raw []byte,
) (command.Result, error) {
	a, err := parseArgs[MkdirArgs](raw)
	if err != nil {
		return command.Result{}, fmt.Errorf("parse args: %w", err)
	}

	return executor.Execute(
		ctx,
		[]string{
			"mkdir",
			"-p",
			a.Path,
		},
		nil,
	)
}
