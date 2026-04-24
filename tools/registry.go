package tools

import (
	"bytes"
	"fmt"
	"os/exec"
)

var (
	ErrNotFound          = fmt.Errorf("tool not found")
	ErrAlreadyRegistered = fmt.Errorf("already registered")
	ErrRunCommand        = fmt.Errorf("run command")
)

type Result struct {
	Stdout string
	StdErr string
}

type Tool interface {
	Code() string
	Execute(dir string, args []byte) (Result, error)
}

type Registry struct {
	tools map[string]Tool
}

func NewRegistry() *Registry {
	return &Registry{
		tools: make(map[string]Tool),
	}
}

func (r *Registry) Get(code string) (Tool, error) {
	tool, ok := r.tools[code]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrNotFound, code)
	}

	return tool, nil
}

func (r *Registry) Register(tool Tool) error {
	_, ok := r.tools[tool.Code()]
	if ok {
		return fmt.Errorf("%w: %s", ErrAlreadyRegistered, tool.Code())
	}

	r.tools[tool.Code()] = tool

	return nil
}

func runCommand(cmd *exec.Cmd) (Result, error) {
	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
	)

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	result := Result{
		Stdout: stdout.String(),
		StdErr: stderr.String(),
	}

	if err != nil {
		return result, fmt.Errorf("%w: %w", ErrRunCommand, err)
	}

	return result, nil
}
