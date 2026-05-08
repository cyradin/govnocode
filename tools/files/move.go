package files

import (
	"context"
	"fmt"

	"github.com/cyradin/govnocode/tools/executor"
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
	e executor.Executor,
	raw []byte,
) (executor.Result, error) {
	a, err := parseArgs[MoveArgs](raw)
	if err != nil {
		return executor.Result{}, fmt.Errorf("parse args: %w", err)
	}

	return e.Execute(
		ctx,
		[]string{
			"mv",
			a.From,
			a.To,
		},
		nil,
	)
}
