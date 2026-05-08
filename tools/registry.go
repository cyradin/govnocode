package tools

import (
	"context"
	"fmt"

	"github.com/cyradin/govnocode/tools/executor"
)

var (
	ErrNotFound          = fmt.Errorf("tool not found")
	ErrAlreadyRegistered = fmt.Errorf("already registered")
	ErrRunCommand        = fmt.Errorf("run command")
)

type Spec struct {
	Code        string `json:"code"`
	Description string `json:"description"`
	Args        any    `json:"args,omitempty"`
}

type toolExecutor interface {
	Execute(ctx context.Context, executor executor.Executor, args []byte) (executor.Result, error)
}

type Tool struct {
	spec  Spec
	inner toolExecutor
}

func NewTool(inner toolExecutor, spec Spec) *Tool {
	return &Tool{
		spec:  spec,
		inner: inner,
	}
}

func (t *Tool) Spec() Spec {
	return t.spec
}

func (t *Tool) Execute(ctx context.Context, executor executor.Executor, args []byte) (executor.Result, error) {
	return t.inner.Execute(ctx, executor, args)
}

type Registry struct {
	tools map[string]*Tool
}

func NewRegistry() *Registry {
	return &Registry{
		tools: make(map[string]*Tool),
	}
}

func (r *Registry) Get(code string) (*Tool, error) {
	tool, ok := r.tools[code]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrNotFound, code)
	}

	return tool, nil
}

func (r *Registry) All() []*Tool {
	result := make([]*Tool, 0, len(r.tools))

	for _, tool := range r.tools {
		result = append(result, tool)
	}

	return result
}

func (r *Registry) Register(tools ...*Tool) error {
	for _, tool := range tools {
		code := tool.Spec().Code

		_, ok := r.tools[code]
		if ok {
			return fmt.Errorf("%w: %s", ErrAlreadyRegistered, code)
		}

		r.tools[code] = tool
	}

	return nil
}

func (r *Registry) MustRegister(tools ...*Tool) *Registry {
	if err := r.Register(tools...); err != nil {
		panic(err)
	}

	return r
}
