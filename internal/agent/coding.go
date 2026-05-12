package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log/slog"

	"github.com/cyradin/govnocode/internal/docker"
	"github.com/cyradin/govnocode/internal/llm"
	"github.com/cyradin/govnocode/tools"
	"github.com/cyradin/govnocode/tools/executor"
	"golang.org/x/sync/errgroup"
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
	PrintLLMResponse(messages <-chan PrinterLLMMessage) error
}

type CodingAgent struct {
	printer   printer
	llmClient llmClient
	tools     toolProvider
	logger    *slog.Logger
}

func NewCoding(
	printer printer,
	llmClient llmClient,
	toolProvider toolProvider,
	logger *slog.Logger,
) *CodingAgent {
	return &CodingAgent{
		printer:   printer,
		llmClient: llmClient,
		tools:     toolProvider,
		logger:    logger,
	}
}

func (a *CodingAgent) Start(ctx context.Context, dir string, task string) error {
	systemPrompt, err := a.systemPrompt()
	if err != nil {
		return fmt.Errorf("make system prompt: %w", err)
	}

	a.logger.InfoContext(ctx, "starting docker container...")

	container, err := docker.NewDockerBuilder().BuildFromEmbed(ctx, docker.DockerfileGo)
	if err != nil {
		return fmt.Errorf("run docker container: %w", err)
	}

	a.logger.InfoContext(ctx, "docker container started")

	executor := executor.NewDocker(container)

	llmSession := llm.NewSession(a.llmClient, systemPrompt)

	return a.startLoop(ctx, task, executor, llmSession)
}

func (a *CodingAgent) startLoop(ctx context.Context, task string, e *executor.Docker, llmSession *llm.Session) error {
	if err := a.writeMessage(ctx, task, llmSession); err != nil {
		return err
	}

	for {
		msg := llmSession.LastMessage()

		toolCall, err := NewParser(msg.Content).GetTool()
		if err != nil {
			if err := a.writeErrorMessage(ctx, executor.Result{}, err, llmSession); err != nil {
				return err
			}

			continue
		}

		tool, err := a.tools.Get(toolCall.Tool)
		if err != nil {
			if err := a.writeErrorMessage(ctx, executor.Result{}, err, llmSession); err != nil {
				return err
			}

			continue
		}

		res, err := tool.Execute(ctx, e, toolCall.Args)
		if err != nil {
			if err := a.writeErrorMessage(ctx, res, err, llmSession); err != nil {
				return err
			}

			continue
		}

		prompt, err := a.toolCallResultPrompt(res)
		if err != nil {
			return fmt.Errorf("create error prompt: %w", err)
		}

		if err := a.writeMessage(ctx, prompt, llmSession); err != nil {
			return err
		}
	}
}

func (a *CodingAgent) writeErrorMessage(ctx context.Context, result executor.Result, err error, llmSession *llm.Session) error {
	if err := a.printer.PrintError(err); err != nil {
		return fmt.Errorf("print error message: %w", err)
	}

	prompt, err := a.errorPrompt(result, err)
	if err != nil {
		return fmt.Errorf("create error prompt: %w", err)
	}

	if err := a.writeMessage(ctx, prompt, llmSession); err != nil {
		return fmt.Errorf("write error message: %w", err)
	}

	return nil
}

func (a *CodingAgent) writeMessage(ctx context.Context, text string, llmSession *llm.Session) error {
	printerCh := make(chan PrinterLLMMessage)

	if err := a.printer.PrintUserMessage(text); err != nil {
		return fmt.Errorf("print user message: %w", err)
	}

	eg := errgroup.Group{}

	eg.Go(func() error {
		return a.printer.PrintLLMResponse(printerCh)
	})

	eg.Go(func() error {
		defer close(printerCh)

		for part := range llmSession.WriteMessage(ctx, text) {
			if err := part.Err; err != nil {
				return fmt.Errorf("send message to llm :%w", err)
			}

			printerCh <- PrinterLLMMessage{
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

func (a *CodingAgent) systemPrompt() (string, error) {
	allTools := a.tools.All()

	specs := make([]tools.Spec, 0, len(allTools))
	for _, tool := range allTools {
		specs = append(specs, tool.Spec())
	}

	toolsRaw, err := json.Marshal(specs)
	if err != nil {
		return "", fmt.Errorf("json marshal: %w", err)
	}

	return fmt.Sprintf(codingAgentSystemPrompt, string(toolsRaw)), nil
}

func (a *CodingAgent) errorPrompt(res executor.Result, err error) (string, error) {
	var buf bytes.Buffer

	data := map[string]string{
		"status": "ERROR",
		"error":  err.Error(),
		"stdout": res.Stdout,
		"stderr": res.Stderr,
	}

	if err := codingAgentResultPrompt.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("execute template: %w", err)
	}

	return buf.String(), nil
}

func (a *CodingAgent) toolCallResultPrompt(res executor.Result) (string, error) {
	var buf bytes.Buffer

	data := map[string]string{
		"status": "SUCCESS",
		"error":  "",
		"stdout": res.Stdout,
		"stderr": res.Stderr,
	}

	if err := codingAgentResultPrompt.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("execute template: %w", err)
	}

	return buf.String(), nil
}

const codingAgentSystemPrompt = `
<instructions>
You are an autonomous Go coding agent.

STRICT RULES:

1. You MUST output EXACTLY ONE of:
   A) a tool call JSON
   B) a final answer

2. If you output a tool call:
   - Output ONLY valid JSON
   - DO NOT include any text before or after JSON
   - DO NOT include explanations
   - DO NOT include thinking
   - After the JSON, STOP immediately

3. DO NOT write "Thinking", plans, or reasoning.
4. DO NOT explain your actions.
5. DO NOT output code blocks.
6. DO NOT output multiple tool calls.

Tool call format:
{
  "tool": "<tool_code>",
  "args": {}
}

If the task is not finished → use a tool.
If the task is finished → output final answer.

</instructions>

<tools>
%s
</tools>
`

var codingAgentResultPrompt = template.Must(
	template.New("coding_agent_result").Parse(`
<result>
	{{- if .status }}
	<status>{{.status}}</status>
	{{- end }}

	{{- if .error }}
	<error>{{.error}}</error>
	{{- end }}

	{{- if .stdout }}
	<stdout>{{.stdout}}</stdout>
	{{- end }}

	{{- if .stderr }}
	<stderr>{{.stderr}}</stderr>
	{{- end }}
</result>
`),
)
