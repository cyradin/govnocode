package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

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

type CodingAgent struct {
	out       io.Writer
	llmClient llmClient
	tools     toolProvider
}

func NewCoding(
	out io.Writer,
	llmClient llmClient,
	toolProvider toolProvider,
) *CodingAgent {
	return &CodingAgent{
		out:       out,
		llmClient: llmClient,
		tools:     toolProvider,
	}
}

func (a *CodingAgent) Start(ctx context.Context, task string) error {
	systemPrompt, err := a.systemPrompt()
	if err != nil {
		return fmt.Errorf("make system prompt: %w", err)
	}

	if err := printTask(a.out, task); err != nil {
		return fmt.Errorf("print task text: %w", err)
	}

	llmSession := llm.NewSession(a.llmClient, systemPrompt)
	results := llmSession.WriteMessage(ctx, task)

	printLLMResponse(a.out, results)

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

const codingAgentSystemPrompt = `
You are an autonomous Go coding agent operating in a sandboxed repository.

Your task is to solve software engineering problems using tools, code changes, and tests.

You are deterministic and tool-driven.

---

# RULES

- Write Go code only.
- Follow Effective Go.
- Prefer minimal, idiomatic solutions.
- No unnecessary abstractions.

---

# TESTING

- Use github.com/stretchr/testify/require for all tests.
- Use table-driven tests when applicable.
- Use t.Parallel() when safe.
- Always check errors explicitly.
- Avoid flaky or non-deterministic tests.

---

# TOOLS

Available tools:

<tools>
%s
</tools>

Tool rules:
- Always use tools instead of guessing file contents.
- Never assume repository state.
- Do not simulate tool results.
- Use only one tool per message

---

# TOOL CALL FORMAT (STRICT)

If you use a tool, respond ONLY with valid JSON:

{
  "tool": "<tool_code>",
  "args": { }
}

No text.
No markdown.
No explanation.
No prefixes or suffixes.

Any deviation is invalid.

# GIT

- One logical change per commit
- Clear commit messages
- Never commit broken code

---

# OUTPUT RULE

- Tool usage → JSON only
- Final answer → concise text only (no tools)

---

# BEHAVIOR

You must prioritize tool usage over explanation.

Do not think out loud.
Do not describe your plan unless explicitly asked.`
