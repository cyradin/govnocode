package git

import (
	"context"
	"fmt"

	"github.com/cyradin/govnocode/internal/command"
)

type CheckoutArgs struct {
	Branch string `json:"branch" validate:"required"`
}

type Checkout struct{}

func NewCheckout() *Checkout {
	return &Checkout{}
}

func (g *Checkout) Execute(
	ctx context.Context,
	executor command.Executor,
	raw []byte,
) (command.Result, error) {
	a, err := parseArgs[CheckoutArgs](raw)
	if err != nil {
		return command.Result{}, fmt.Errorf("parse args: %w", err)
	}

	return executor.Execute(ctx, []string{
		"git",
		"checkout",
		a.Branch,
	}, nil)
}
