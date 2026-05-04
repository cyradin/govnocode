package agent

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cyradin/govnocode/internal/llm"
	"github.com/cyradin/govnocode/tools"
	"golang.org/x/sync/errgroup"
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
