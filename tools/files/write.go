package files

import (
	"context"
	"fmt"
	"strings"

	"github.com/cyradin/govnocode/internal/command"
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
	executor command.Executor,
	raw []byte,
) (command.Result, error) {
	a, err := parseArgs[WriteArgs](raw)
	if err != nil {
		return command.Result{}, fmt.Errorf("parse args: %w", err)
	}

	return executor.Execute(
		ctx,
		[]string{
			"tee",
			a.Path,
		},
		strings.NewReader(a.Content),
	)
}
