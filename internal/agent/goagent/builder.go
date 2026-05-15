package goagent

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/cyradin/govnocode/internal/docker"
	"github.com/cyradin/govnocode/internal/llm"
	"github.com/cyradin/govnocode/internal/output"
	"github.com/cyradin/govnocode/tools"
	"github.com/cyradin/govnocode/tools/executor"
)

type toolProvider interface {
	Get(code string) (*tools.Tool, error)
	All() []*tools.Tool
}

type llmClient interface {
	Stream(ctx context.Context, messages []llm.ChatMessage) <-chan llm.ChatResult
}

type printer interface {
	PrintUserMessage(msg string) error
	PrintError(err error) error
	PrintLLMResponse(messages <-chan output.LLMMessage) error
}

type Builder struct {
	printer   printer
	llmClient llmClient
	logger    *slog.Logger

	workDir string
}

func NewBuilder(
	printer printer,
	llmClient llmClient,
	logger *slog.Logger,
) *Builder {
	return &Builder{
		printer:   printer,
		llmClient: llmClient,
		logger:    logger,
	}
}

func (b *Builder) WithWorkdir(workDir string) *Builder {
	b.workDir = workDir

	return b
}

func (b *Builder) Build(ctx context.Context) (*Agent, error) {
	toolRegistry := tools.NewRegistry().
		MustRegister(tools.Git()...).
		MustRegister(tools.Files()...)

	containerBuilder := docker.NewDockerBuilder()

	if b.workDir != "" {
		containerBuilder.WithMount(b.workDir, "/app")
	}

	b.logger.InfoContext(ctx, "docker container started")

	container, err := containerBuilder.BuildFromEmbed(ctx, docker.DockerfileGo)
	if err != nil {
		return nil, fmt.Errorf("run docker container: %w", err)
	}

	b.logger.InfoContext(ctx, "docker container started")

	executor := executor.NewDocker(container)

	session := llm.NewSession(b.llmClient)

	return New(b.printer, session, toolRegistry, b.logger, executor), nil
}
