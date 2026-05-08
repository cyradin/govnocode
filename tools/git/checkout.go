package git

import (
	"context"
	"fmt"

	"github.com/cyradin/govnocode/tools/executor"
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
	e executor.Executor,
	raw []byte,
) (executor.Result, error) {
	a, err := parseArgs[CheckoutArgs](raw)
	if err != nil {
		return executor.Result{}, fmt.Errorf("parse args: %w", err)
	}

	return e.Execute(ctx, []string{
		"git",
		"checkout",
		a.Branch,
	}, nil)
}
