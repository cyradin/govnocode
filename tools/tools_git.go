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

	return runCommand(cmd, dir)
}

type gitCreateBranchArgs struct {
	Branch string `json:"branch" validate:"required"`
}

type GitCreateBranch struct {
	command string
	args    []string
}

func NewGitCreateBranch() *GitCreateBranch {
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

	//nolint:gosec
	cmd := exec.Command(
		g.command,
		append(g.args, parsedArgs.Branch)...,
	)

	return runCommand(cmd, dir)
}

type gitCheckoutArgs struct {
	Branch string `json:"branch" validate:"required"`
}

type GitCheckout struct {
	command string
	args    []string
}

func NewGitCheckout() *GitCheckout {
	return &GitCheckout{
		command: "git",
		args:    []string{"checkout"},
	}
}

func (g *GitCheckout) Code() string {
	return "git.checkout"
}

func (g *GitCheckout) Spec() Spec {
	return Spec{
		Code:        g.Code(),
		Description: "Switch to existing branch",
		Args: gitCheckoutArgs{
			Branch: "main",
		},
	}
}

func (g *GitCheckout) Execute(dir string, raw []byte) (Result, error) {
	a, err := parseArgs[gitCheckoutArgs](raw)
	if err != nil {
		return Result{}, fmt.Errorf("parse args: %w", err)
	}

	//nolint:gosec
	cmd := exec.Command(g.command, append(g.args, a.Branch)...)

	return runCommand(cmd, dir)
}

type GitStatus struct{}

func NewGitStatus() *GitStatus {
	return &GitStatus{}
}

func (g *GitStatus) Code() string {
	return "git.status"
}

func (g *GitStatus) Spec() Spec {
	return Spec{
		Code:        g.Code(),
		Description: "Show working tree status",
		Args:        nil,
	}
}

func (g *GitStatus) Execute(dir string, raw []byte) (Result, error) {
	cmd := exec.Command("git", "status", "--porcelain")

	return runCommand(cmd, dir)
}

type gitAddArgs struct {
	Path string `json:"path" validate:"required"`
}

type GitAdd struct{}

func NewGitAdd() *GitAdd {
	return &GitAdd{}
}

func (g *GitAdd) Code() string {
	return "git.add"
}

func (g *GitAdd) Spec() Spec {
	return Spec{
		Code:        g.Code(),
		Description: "Stage files",
		Args: gitAddArgs{
			Path: ".",
		},
	}
}

func (g *GitAdd) Execute(dir string, raw []byte) (Result, error) {
	a, err := parseArgs[gitAddArgs](raw)
	if err != nil {
		return Result{}, fmt.Errorf("parse args: %w", err)
	}

	//nolint:gosec
	cmd := exec.Command("git", "add", a.Path)

	return runCommand(cmd, dir)
}

type gitCommitArgs struct {
	Message string `json:"message" validate:"required"`
}

type GitCommit struct{}

func NewGitCommit() *GitCommit {
	return &GitCommit{}
}

func (g *GitCommit) Code() string {
	return "git.commit"
}

func (g *GitCommit) Spec() Spec {
	return Spec{
		Code:        g.Code(),
		Description: "Create commit",
		Args: gitCommitArgs{
			Message: "update code",
		},
	}
}

func (g *GitCommit) Execute(dir string, raw []byte) (Result, error) {
	a, err := parseArgs[gitCommitArgs](raw)
	if err != nil {
		return Result{}, fmt.Errorf("parse args: %w", err)
	}

	//nolint:gosec
	cmd := exec.Command("git", "commit", "-m", a.Message)

	return runCommand(cmd, dir)
}

type gitPushArgs struct {
	Remote string `json:"remote" validate:"required"`
	Branch string `json:"branch" validate:"required"`
}

type GitPush struct{}

func NewGitPush() *GitPush {
	return &GitPush{}
}

func (g *GitPush) Code() string {
	return "git.push"
}

func (g *GitPush) Spec() Spec {
	return Spec{
		Code:        g.Code(),
		Description: "Push branch to remote",
		Args: gitPushArgs{
			Remote: "origin",
			Branch: "master",
		},
	}
}

func (g *GitPush) Execute(dir string, raw []byte) (Result, error) {
	a, err := parseArgs[gitPushArgs](raw)
	if err != nil {
		return Result{}, fmt.Errorf("parse args: %w", err)
	}

	//nolint:gosec
	cmd := exec.Command("git", "push", "-u", a.Remote, a.Branch)

	return runCommand(cmd, dir)
}

type GitPull struct{}

func NewGitPull() *GitPull {
	return &GitPull{}
}

func (g *GitPull) Code() string {
	return "git.pull"
}

func (g *GitPull) Spec() Spec {
	return Spec{
		Code:        g.Code(),
		Description: "Pull changes from remote",
		Args:        nil,
	}
}

func (g *GitPull) Execute(dir string, raw []byte) (Result, error) {
	cmd := exec.Command("git", "pull")

	return runCommand(cmd, dir)
}

type GitFetch struct{}

func NewGitFetch() *GitFetch {
	return &GitFetch{}
}

func (g *GitFetch) Code() string {
	return "git.fetch"
}

func (g *GitFetch) Spec() Spec {
	return Spec{
		Code:        g.Code(),
		Description: "Fetch remote changes",
		Args:        nil,
	}
}

func (g *GitFetch) Execute(dir string, raw []byte) (Result, error) {
	cmd := exec.Command("git", "fetch")

	return runCommand(cmd, dir)
}

type GitBranchList struct{}

func NewGitBranchList() *GitBranchList {
	return &GitBranchList{}
}

func (g *GitBranchList) Code() string {
	return "git.branch_list"
}

func (g *GitBranchList) Spec() Spec {
	return Spec{
		Code:        g.Code(),
		Description: "List branches",
		Args:        nil,
	}
}

func (g *GitBranchList) Execute(dir string, raw []byte) (Result, error) {
	cmd := exec.Command("git", "branch", "--all")

	return runCommand(cmd, dir)
}
