package git

import (
	"context"

	"github.com/cyradin/govnocode/tools/executor"
)

type BranchList struct{}

func NewBranchList() *BranchList {
	return &BranchList{}
}

func (g *BranchList) Execute(
	ctx context.Context,
	e executor.Executor,
	raw []byte,
) (executor.Result, error) {
	return e.Execute(ctx, []string{
		"git",
		"branch",
		"--all",
	}, nil)
}
