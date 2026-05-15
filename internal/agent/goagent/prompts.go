package goagent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"

	"github.com/cyradin/govnocode/tools"
	"github.com/cyradin/govnocode/tools/executor"
)

const systemPrompt = `
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
   - ONLY ONE TOOL CALL per message
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

func makeSystemPrompt(allTools []*tools.Tool) (string, error) {
	specs := make([]tools.Spec, 0, len(allTools))
	for _, tool := range allTools {
		specs = append(specs, tool.Spec())
	}

	toolsRaw, err := json.Marshal(specs)
	if err != nil {
		return "", fmt.Errorf("json marshal: %w", err)
	}

	return fmt.Sprintf(systemPrompt, string(toolsRaw)), nil
}

var resultPrompt = template.Must(
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

func makeResultPrompt(res executor.Result) (string, error) {
	var buf bytes.Buffer

	data := map[string]string{
		"status": "SUCCESS",
		"error":  "",
		"stdout": res.Stdout,
		"stderr": res.Stderr,
	}

	if err := resultPrompt.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("execute template: %w", err)
	}

	return buf.String(), nil
}

func makeErrorPrompt(res executor.Result, err error) (string, error) {
	var buf bytes.Buffer

	data := map[string]string{
		"status": "ERROR",
		"error":  err.Error(),
		"stdout": res.Stdout,
		"stderr": res.Stderr,
	}

	if err := resultPrompt.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("execute template: %w", err)
	}

	return buf.String(), nil
}
