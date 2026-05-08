package files

import (
	"context"
	"fmt"

	"github.com/cyradin/govnocode/tools/executor"
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
	e executor.Executor,
	raw []byte,
) (executor.Result, error) {
	a, err := parseArgs[DirListArgs](raw)
	if err != nil {
		return executor.Result{}, fmt.Errorf("parse args: %w", err)
	}

	path := "."
	if a.Path != "" {
		path = a.Path
	}

	return e.Execute(ctx, []string{
		"ls",
		"-1",
		path,
	}, nil)
}
