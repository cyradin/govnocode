package git

import (
	"context"

	"github.com/cyradin/govnocode/internal/command"
)

type BranchList struct{}

func NewBranchList() *BranchList {
	return &BranchList{}
}

func (g *BranchList) Execute(
	ctx context.Context,
	executor command.Executor,
	raw []byte,
) (command.Result, error) {
	return executor.Execute(ctx, []string{
		"git",
		"branch",
		"--all",
	}, nil)
}
