package files

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cyradin/govnocode/tools/executor"
	"github.com/go-playground/validator/v10"
)

type ReadArgs struct {
	Path string `json:"path" validate:"required"`
}

type Read struct{}

func NewRead() *Read {
	return &Read{}
}

func (f *Read) Execute(
	ctx context.Context,
	e executor.Executor,
	raw []byte,
) (executor.Result, error) {
	a, err := parseArgs[ReadArgs](raw)
	if err != nil {
		return executor.Result{}, fmt.Errorf("parse args: %w", err)
	}

	return e.Execute(ctx, []string{
		"cat",
		a.Path,
	}, nil)
}

var validate = validator.New()

func parseArgs[T any](raw []byte) (T, error) {
	var a T

	if err := json.Unmarshal(raw, &a); err != nil {
		return a, fmt.Errorf("invalid args json: %w", err)
	}

	if err := validate.Struct(a); err != nil {
		return a, fmt.Errorf("validate args: %w", err)
	}

	return a, nil
}
