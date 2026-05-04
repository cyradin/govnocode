package agent

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParser_GetTool(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     string
		wantTool  string
		wantArgs  string
		expectErr bool
	}{
		{
			name:     "valid json direct",
			input:    `{"tool":"dir.list","args":{"path":"."}}`,
			wantTool: "dir.list",
			wantArgs: `{"path":"."}`,
		},
		{
			name:     "json in code block",
			input:    "some text\n```json\n{\"tool\":\"file.write\",\"args\":{\"path\":\"main.go\"}}\n```\nmore text",
			wantTool: "file.write",
			wantArgs: `{"path":"main.go"}`,
		},
		{
			name:     "code block without json tag",
			input:    "blah\n```\n{\"tool\":\"run\",\"args\":{\"cmd\":\"go test\"}}\n```",
			wantTool: "run",
			wantArgs: `{"cmd":"go test"}`,
		},
		{
			name:      "invalid json",
			input:     "```json\n{invalid json}\n```",
			expectErr: true,
		},
		{
			name:      "missing tool field",
			input:     `{"args":{"x":1}}`,
			expectErr: true,
		},
		{
			name:      "no json no code block",
			input:     "Thinking...\nI will do something later",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			p := NewParser(tt.input)

			call, err := p.GetTool()

			if tt.expectErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.wantTool, call.Tool)
			require.JSONEq(t, tt.wantArgs, string(call.Args))
		})
	}
}
