package command

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
)

var (
	ErrRunCommand = fmt.Errorf("run command")
)

type Executor struct {
	dir string
}

func NewExecutor(dir string) *Executor {
	return &Executor{
		dir: dir,
	}
}

func (c *Executor) Run(ctx context.Context, args []string) (Result, error) {
	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
	)

	cmd := exec.Command(args[0], args[1:]...) //nolint:gosec
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Dir = c.dir

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
