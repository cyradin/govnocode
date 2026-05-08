package files

import (
	"context"
	"fmt"

	"github.com/cyradin/govnocode/tools/executor"
)

type DeleteArgs struct {
	Path string `json:"path" validate:"required"`
}

type Delete struct{}

func NewDelete() *Delete {
	return &Delete{}
}

func (f *Delete) Execute(
	ctx context.Context,
	e executor.Executor,
	raw []byte,
) (executor.Result, error) {
	a, err := parseArgs[DeleteArgs](raw)
	if err != nil {
		return executor.Result{}, fmt.Errorf("parse args: %w", err)
	}

	return e.Execute(
		ctx,
		[]string{
			"rm",
			"-rf",
			a.Path,
		},
		nil,
	)
}
