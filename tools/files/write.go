package files

import (
	"context"
	"fmt"
	"strings"

	"github.com/cyradin/govnocode/tools/executor"
)

type WriteArgs struct {
	Path    string `json:"path" validate:"required"`
	Content string `json:"content"`
}

type Write struct{}

func NewWrite() *Write {
	return &Write{}
}

func (f *Write) Execute(
	ctx context.Context,
	e executor.Executor,
	raw []byte,
) (executor.Result, error) {
	a, err := parseArgs[WriteArgs](raw)
	if err != nil {
		return executor.Result{}, fmt.Errorf("parse args: %w", err)
	}

	return e.Execute(
		ctx,
		[]string{
			"tee",
			a.Path,
		},
		strings.NewReader(a.Content),
	)
}
