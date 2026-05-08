package executor

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

var _ Executor = (*Shell)(nil)

var (
	ErrRunCommand = fmt.Errorf("run command")
)

type Shell struct {
	dir string
}

func NewShell(dir string) *Shell {
	return &Shell{
		dir: dir,
	}
}

func (c *Shell) Execute(ctx context.Context, args []string, stdin io.Reader) (Result, error) {
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
