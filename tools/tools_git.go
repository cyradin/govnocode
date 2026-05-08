package tools

import (
	"github.com/cyradin/govnocode/tools/git"
)

func Git() []*Tool {
	return []*Tool{
		NewGitInit(),
		NewGitCreateBranch(),
		NewGitCheckout(),
		NewGitStatus(),
		NewGitAdd(),
		NewGitCommit(),
		NewGitPush(),
		NewGitPull(),
		NewGitFetch(),
		NewGitBranchList(),
	}
}

func NewGitInit() *Tool {
	return NewTool(
		git.NewInit(),
		Spec{
			Code:        "git.init",
			Description: "Initialize new git repository",
			Args:        nil,
		},
	)
}

func NewGitCreateBranch() *Tool {
	return NewTool(
		git.NewCreateBranch(),
		Spec{
			Code:        "git.create_branch",
			Description: "Create new git branch",
			Args: git.CreateBranchArgs{
				Branch: "feature/login",
			},
		},
	)
}

func NewGitCheckout() *Tool {
	return NewTool(
		git.NewCheckout(),
		Spec{
			Code:        "git.checkout",
			Description: "Switch to existing branch",
			Args: git.CheckoutArgs{
				Branch: "main",
			},
		},
	)
}

func NewGitStatus() *Tool {
	return NewTool(
		git.NewStatus(),
		Spec{
			Code:        "git.status",
			Description: "Show working tree status",
			Args:        nil,
		},
	)
}

func NewGitAdd() *Tool {
	return NewTool(
		git.NewAdd(),
		Spec{
			Code:        "git.add",
			Description: "Stage files",
			Args: git.AddArgs{
				Path: ".",
			},
		},
	)
}

func NewGitCommit() *Tool {
	return NewTool(
		git.NewCommit(),
		Spec{
			Code:        "git.commit",
			Description: "Create commit",
			Args: git.CommitArgs{
				Message: "update code",
			},
		},
	)
}

func NewGitPush() *Tool {
	return NewTool(
		git.NewPush(),
		Spec{
			Code:        "git.push",
			Description: "Push branch to remote",
			Args: git.PushArgs{
				Remote: "origin",
				Branch: "master",
			},
		},
	)
}

func NewGitPull() *Tool {
	return NewTool(
		git.NewPull(),
		Spec{
			Code:        "git.pull",
			Description: "Pull changes from remote",
			Args:        nil,
		},
	)
}

func NewGitFetch() *Tool {
	return NewTool(
		git.NewFetch(),
		Spec{
			Code:        "git.fetch",
			Description: "Fetch remote changes",
			Args:        nil,
		},
	)
}

func NewGitBranchList() *Tool {
	return NewTool(
		git.NewBranchList(),
		Spec{
			Code:        "git.branch_list",
			Description: "List branches",
			Args:        nil,
		},
	)
}
