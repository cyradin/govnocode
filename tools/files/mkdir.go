package files

import (
	"context"
	"fmt"

	"github.com/cyradin/govnocode/tools/executor"
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
	e executor.Executor,
	raw []byte,
) (executor.Result, error) {
	a, err := parseArgs[MkdirArgs](raw)
	if err != nil {
		return executor.Result{}, fmt.Errorf("parse args: %w", err)
	}

	return e.Execute(
		ctx,
		[]string{
			"mkdir",
			"-p",
			a.Path,
		},
		nil,
	)
}
