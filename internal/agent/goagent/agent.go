package goagent

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/cyradin/govnocode/internal/llm"
	"github.com/cyradin/govnocode/internal/output"
	"github.com/cyradin/govnocode/tools/executor"
	"golang.org/x/sync/errgroup"
)

type Agent struct {
	printer  printer
	tools    toolProvider
	executor executor.Executor
	session  *llm.Session

	logger *slog.Logger
}

func New(
	printer printer,
	session *llm.Session,
	toolProvider toolProvider,
	logger *slog.Logger,
	executor executor.Executor,
) *Agent {
	return &Agent{
		printer:  printer,
		session:  session,
		tools:    toolProvider,
		logger:   logger,
		executor: executor,
	}
}

func (a *Agent) Start(ctx context.Context, dir string, task string) error {
	system, err := makeSystemPrompt(a.tools.All())
	if err != nil {
		return fmt.Errorf("make system prompt: %w", err)
	}

	a.session.SetSystemPrompt(system)

	if err := a.writeMessage(ctx, task); err != nil {
		return err
	}

	for {
		msg := a.session.LastMessage()

		res, err := a.tryExecuteTool(ctx, msg)
		if err != nil {
			if err := a.writeErrorMessage(ctx, res, err); err != nil {
				return err
			}
		}

		prompt, err := makeResultPrompt(res)
		if err != nil {
			return fmt.Errorf("create error prompt: %w", err)
		}

		if err := a.writeMessage(ctx, prompt); err != nil {
			return err
		}
	}
}

func (a *Agent) tryExecuteTool(ctx context.Context, msg llm.ChatMessage) (executor.Result, error) {
	toolCall, err := llm.NewParser(msg.Content).GetTool()
	if err != nil {
		return executor.Result{}, err
	}

	tool, err := a.tools.Get(toolCall.Tool)
	if err != nil {
		return executor.Result{}, err
	}

	res, err := tool.Execute(ctx, a.executor, toolCall.Args)
	if err != nil {
		return executor.Result{}, err
	}

	return res, nil
}

func (a *Agent) writeErrorMessage(ctx context.Context, result executor.Result, err error) error {
	if err := a.printer.PrintError(err); err != nil {
		return fmt.Errorf("print error message: %w", err)
	}

	prompt, err := makeErrorPrompt(result, err)
	if err != nil {
		return fmt.Errorf("create error prompt: %w", err)
	}

	if err := a.writeMessage(ctx, prompt); err != nil {
		return fmt.Errorf("write error message: %w", err)
	}

	return nil
}

func (a *Agent) writeMessage(ctx context.Context, text string) error {
	printerCh := make(chan output.LLMMessage)

	if err := a.printer.PrintUserMessage(text); err != nil {
		return fmt.Errorf("print user message: %w", err)
	}

	eg := errgroup.Group{}

	eg.Go(func() error {
		return a.printer.PrintLLMResponse(printerCh)
	})

	eg.Go(func() error {
		defer close(printerCh)

		for part := range a.session.WriteMessage(ctx, text) {
			if err := part.Err; err != nil {
				return fmt.Errorf("send message to llm :%w", err)
			}

			printerCh <- output.LLMMessage{
				Content:  part.Resp.Message.Content,
				Thinking: part.Resp.Message.Thinking,
			}
		}

		return nil
	})

	if err := eg.Wait(); err != nil {
		return err
	}

	return nil
}
