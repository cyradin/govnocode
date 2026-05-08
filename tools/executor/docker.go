package executor

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/moby/moby/api/pkg/stdcopy"
	"github.com/testcontainers/testcontainers-go"
	tcexec "github.com/testcontainers/testcontainers-go/exec"
)

var _ Executor = (*Docker)(nil)

type Result struct {
	Stdout string
	Stderr string
}

type Docker struct {
	container testcontainers.Container
}

func NewDocker(container testcontainers.Container) *Docker {
	return &Docker{
		container: container,
	}
}

func (e *Docker) Execute(
	ctx context.Context,
	cmd []string,
	stdin io.Reader,
) (Result, error) {
	opts := []tcexec.ProcessOption{}

	if stdin != nil {
		opts = append(opts, tcexec.ProcessOptionFunc(func(p *tcexec.ProcessOptions) {
			p.Reader = stdin
		}))
	}

	exitCode, reader, err := e.container.Exec(
		ctx,
		cmd,
		opts...,
	)
	if err != nil {
		return Result{}, fmt.Errorf("exec in container: %w", err)
	}

	var (
		outBuf bytes.Buffer
		errBuf bytes.Buffer
	)

	_, err = stdcopy.StdCopy(&outBuf, &errBuf, reader)
	if err != nil {
		return Result{}, fmt.Errorf("decode output: %w", err)
	}

	result := Result{
		Stdout: outBuf.String(),
		Stderr: errBuf.String(),
	}

	if exitCode != 0 {
		return result, fmt.Errorf("%w: exit code %d", ErrRunCommand, exitCode)
	}

	return result, nil
}
