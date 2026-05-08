package files

import (
	"context"
	"fmt"

	"github.com/cyradin/govnocode/internal/command"
)

type DirListArgs struct {
	Path string `json:"path"`
}

type DirList struct{}

func NewDirList() *DirList {
	return &DirList{}
}

func (d *DirList) Execute(
	ctx context.Context,
	executor command.Executor,
	raw []byte,
) (command.Result, error) {
	a, err := parseArgs[DirListArgs](raw)
	if err != nil {
		return command.Result{}, fmt.Errorf("parse args: %w", err)
	}

	path := "."
	if a.Path != "" {
		path = a.Path
	}

	return executor.Execute(ctx, []string{
		"ls",
		"-1",
		path,
	}, nil)
}
