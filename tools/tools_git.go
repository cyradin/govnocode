package tools

import (
	"fmt"
	"os/exec"
)

const gitCommand = "git"

type GitInit struct {
	command string
	args    []string
}

func NewGitInit() *GitInit {
	return &GitInit{
		command: gitCommand,
		args:    []string{"init"},
	}
}

func (g *GitInit) Code() string {
	return gitCommand + ".init"
}

func (g *GitInit) Spec() Spec {
	return Spec{
		Code:        g.Code(),
		Description: "Initialize new git repository",
		Args:        nil,
	}
}

func (g *GitInit) Execute(dir string, args []byte) (Result, error) {
	cmd := exec.Command(g.command, g.args...) //nolint:gosec
	cmd.Dir = dir

	res, err := runCommand(cmd)
	if err != nil {
		return res, fmt.Errorf("git init: %w", err)
	}

	return res, nil
}

type gitCreateBranchArgs struct {
	Branch string `json:"branch" validate:"required"`
}

type GitCreateBranch struct {
	command string
	args    []string
}

func NewGetCreateBranch() *GitCreateBranch {
	return &GitCreateBranch{
		command: gitCommand,
		args:    []string{"checkout", "-b"},
	}
}

func (g *GitCreateBranch) Code() string {
	return gitCommand + ".create_branch"
}

func (g *GitCreateBranch) Spec() Spec {
	return Spec{
		Code:        g.Code(),
		Description: "Create new git branch",
		Args: gitCreateBranchArgs{
			Branch: "feature/login",
		},
	}
}

func (g *GitCreateBranch) Execute(dir string, args []byte) (Result, error) {
	parsedArgs, err := parseArgs[gitCreateBranchArgs](args)
	if err != nil {
		return Result{}, fmt.Errorf("parse args: %w", err)
	}

	cmd := exec.Command(
		g.command,
		append(g.args, parsedArgs.Branch)...,
	) //nolint:gosec
	cmd.Dir = dir

	res, err := runCommand(cmd)
	if err != nil {
		return res, fmt.Errorf("git init: %w", err)
	}

	return res, nil
}
