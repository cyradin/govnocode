package files

import (
	"context"
	"fmt"

	"github.com/cyradin/govnocode/internal/command"
)

type MoveArgs struct {
	From string `json:"from" validate:"required"`
	To   string `json:"to" validate:"required"`
}

type Move struct{}

func NewMove() *Move {
	return &Move{}
}

func (m *Move) Execute(
	ctx context.Context,
	executor command.Executor,
	raw []byte,
) (command.Result, error) {
	a, err := parseArgs[MoveArgs](raw)
	if err != nil {
		return command.Result{}, fmt.Errorf("parse args: %w", err)
	}

	return executor.Execute(
		ctx,
		[]string{
			"mv",
			a.From,
			a.To,
		},
		nil,
	)
}
