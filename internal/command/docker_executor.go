package command

import (
	"bytes"
	"context"
	"fmt"

	"github.com/moby/moby/api/pkg/stdcopy"
	"github.com/testcontainers/testcontainers-go"
)

type Result struct {
	Stdout string
	Stderr string
}

type DockerExecutor struct {
	container testcontainers.Container
}

func NewDockerExecutor(container testcontainers.Container) *DockerExecutor {
	return &DockerExecutor{
		container: container,
	}
}

func (c *DockerExecutor) Run(ctx context.Context, cmd []string) (Result, error) {
	exitCode, reader, err := c.container.Exec(ctx, cmd)
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
		return result, fmt.Errorf("exit code %d", exitCode)
	}

	return result, nil
}
