package agent

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cyradin/govnocode/internal/llm"
	"github.com/cyradin/govnocode/tools"
)

type toolProvider interface {
	Get(code string) (tools.Tool, error)
	All() []tools.Tool
}

type llmClient interface {
	Stream(ctx context.Context, messages []llm.ChatMessage) <-chan llm.ChatResult
}

type printer interface {
	PrintTaskText(task string) error
	PrintLLMResponse(messages <-chan PrinterLLMMessage) error
}

type CodingAgent struct {
	printer   printer
	llmClient llmClient
	tools     toolProvider
}

func NewCoding(
	printer printer,
	llmClient llmClient,
	toolProvider toolProvider,
) *CodingAgent {
	return &CodingAgent{
		printer:   printer,
		llmClient: llmClient,
		tools:     toolProvider,
	}
}

func (a *CodingAgent) Start(ctx context.Context, task string) error {
	systemPrompt, err := a.systemPrompt()
	if err != nil {
		return fmt.Errorf("make system prompt: %w", err)
	}

	if err := a.printer.PrintTaskText(task); err != nil {
		return fmt.Errorf("print task text: %w", err)
	}

	llmSession := llm.NewSession(a.llmClient, systemPrompt)

	_, err = a.writeMessage(ctx, task, llmSession)
	if err != nil {
		return err
	}

	return nil
}

func (a *CodingAgent) writeMessage(ctx context.Context, text string, llmSession *llm.Session) (llm.ChatMessage, error) {
	printerCh := make(chan PrinterLLMMessage)
	defer close(printerCh)

	go func() {
		_ = a.printer.PrintLLMResponse(printerCh)
	}()

	for part := range llmSession.WriteMessage(ctx, text) {
		if err := part.Err; err != nil {
			return llm.ChatMessage{}, fmt.Errorf("send message to llm :%w", err)
		}

		printerCh <- PrinterLLMMessage{
			Content:  part.Resp.Message.Content,
			Thinking: part.Resp.Message.Thinking,
		}
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

const codingAgentSystemPrompt = `
<instructions>
You are an autonomous Go coding agent operating in a sandboxed repository.
Your task is to solve software engineering problems using tools, code changes, and tests.
To perform your task, you have to use the provided list of available tools, following below.

- You can do only one thing per response:
	1) tool call JSON
	2) final answer
- If you use a tool, respond with valid JSON:
{
  "tool": "<tool_code>",
  "args": { }
}
</instructions>

<tools>
%s
</tools>
`
