package files

import (
	"context"
	"fmt"

	"github.com/cyradin/govnocode/internal/command"
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
	executor command.Executor,
	raw []byte,
) (command.Result, error) {
	a, err := parseArgs[DeleteArgs](raw)
	if err != nil {
		return command.Result{}, fmt.Errorf("parse args: %w", err)
	}

	return executor.Execute(
		ctx,
		[]string{
			"rm",
			"-rf",
			a.Path,
		},
		nil,
	)
}
