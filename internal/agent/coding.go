package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/cyradin/govnocode/internal/command"
	"github.com/cyradin/govnocode/internal/llm"
	"github.com/cyradin/govnocode/tools"
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

	container, err := command.NewDockerBuilder().BuildFromEmbed(ctx, command.DockerfileGo)
	if err != nil {
		return fmt.Errorf("run docker container: %w", err)
	}

	a.logger.InfoContext(ctx, "docker container started")

	executor := command.NewDockerExecutor(container)

	llmSession := llm.NewSession(a.llmClient, systemPrompt)

	return a.startLoop(ctx, task, executor, llmSession)
}

func (a *CodingAgent) startLoop(ctx context.Context, task string, executor *command.DockerExecutor, llmSession *llm.Session) error {
	var (
		msg llm.ChatMessage
		err error
	)

	msg, err = a.writeMessage(ctx, task, llmSession)
	if err != nil {
		return err
	}

	for {
		toolCall, err := NewParser(msg.Content).GetTool()
		if err != nil {
			if err := a.writeErrorMessage(ctx, err, llmSession); err != nil {
				return err
			}

			continue
		}

		tool, err := a.tools.Get(toolCall.Tool)
		if err != nil {
			if err := a.writeErrorMessage(ctx, err, llmSession); err != nil {
				return err
			}

			continue
		}

		res, err := tool.Execute(ctx, executor, toolCall.Args)
		if err != nil {
			if err := a.writeErrorMessage(ctx, err, llmSession); err != nil {
				return err
			}

			continue
		}

		msg, err = a.writeMessage(ctx, a.toolCallResultPrompt(res), llmSession)
		if err != nil {
			return err
		}
	}
}

func (a *CodingAgent) writeErrorMessage(ctx context.Context, err error, llmSession *llm.Session) error {
	if err := a.printer.PrintError(err); err != nil {
		return fmt.Errorf("print error message: %w", err)
	}

	if _, err := a.writeMessage(ctx, a.errorPrompt(err), llmSession); err != nil {
		return fmt.Errorf("write error message: %w", err)
	}

	return nil
}

func (a *CodingAgent) writeMessage(ctx context.Context, text string, llmSession *llm.Session) (llm.ChatMessage, error) {
	printerCh := make(chan PrinterLLMMessage)

	if err := a.printer.PrintUserMessage(text); err != nil {
		return llm.ChatMessage{}, fmt.Errorf("print user message: %w", err)
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
		return llm.ChatMessage{}, err
	}

	return llmSession.LastMessage(), nil
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

func (a *CodingAgent) errorPrompt(err error) string {
	return fmt.Sprintf(codingAgentErrorPrompt, err.Error())
}

func (a *CodingAgent) toolCallResultPrompt(res command.Result) string {
	return fmt.Sprintf(ccodingAgentToolCallResultPrompt, res.Stdout)
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

const codingAgentErrorPrompt = `
Your previous response was invalid:

<error>
%s
</error>

Fix it.

<rules>
Rules:
- Output ONLY valid tool call JSON
- No explanations
- No thinking
- No markdown
- No extra text before or after JSON
- Must strictly follow format:
  {
    "tool": "<tool_code>",
    "args": { }
  }
</rules>
Return a corrected response.`

const ccodingAgentToolCallResultPrompt = `<instruction>
You are working in an autonomous coding agent loop.

You receive the STDOUT output from a previously executed tool.

Your task:
- Use this output to continue solving the task
- Decide the next action (tool call or final answer)

</instruction>

<stdout>
%s
</stdout>

<rules>
- Output ONLY one of:
  1) tool call JSON
  2) final answer JSON
- No explanations
- No thinking
- No markdown
- No extra text

Tool call format:
{
  "tool": "<tool_code>",
  "args": { }
}
</rules>
`
