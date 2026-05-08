package command

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os/exec"
)

type Executor interface {
	Execute(ctx context.Context, args []string, stdin io.Reader) (Result, error)
}

var _ Executor = (*ShellExecutor)(nil)

var (
	ErrRunCommand = fmt.Errorf("run command")
)

type ShellExecutor struct {
	dir string
}

func NewShellExecutor(dir string) *ShellExecutor {
	return &ShellExecutor{
		dir: dir,
	}
}

func (c *ShellExecutor) Execute(ctx context.Context, args []string, stdin io.Reader) (Result, error) {
	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
	)

	cmd := exec.Command(args[0], args[1:]...) //nolint:gosec
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Dir = c.dir
	cmd.Stdin = stdin

	err := cmd.Run()
	result := Result{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}

	if err != nil {
		return result, fmt.Errorf("%w: %w", ErrRunCommand, err)
	}

	return result, nil
}
