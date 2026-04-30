package agent

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type fakeWriter struct {
	data []byte
}

func TestPrinter_PrintTaskText(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		task  string
		check []string
	}{
		{
			name: "simple task",
			task: "hello task",
			check: []string{
				"Task:",
				"hello task",
				"-------------------------------------------------",
			},
		},
		{
			name: "empty task",
			task: "",
			check: []string{
				"Task:",
				"-------------------------------------------------",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			w := &fakeWriter{}
			p := NewPrinter(w)

			err := p.PrintTaskText(tt.task)
			require.NoError(t, err)

			out := w.String()

			for _, c := range tt.check {
				require.Contains(t, out, c)
			}
		})
	}
}

func TestPrinter_PrintLLMResponse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		messages []PrinterLLMMessage
		check    []string
	}{
		{
			name: "content only",
			messages: []PrinterLLMMessage{
				{Content: "hello"},
			},
			check: []string{
				"LLM Response:",
				"hello",
			},
		},
		{
			name: "thinking then content",
			messages: []PrinterLLMMessage{
				{Thinking: "step1"},
				{Content: "final"},
			},
			check: []string{
				"Thinking...",
				"step1",
				"final",
			},
		},
		{
			name: "mixed stream",
			messages: []PrinterLLMMessage{
				{Thinking: "t1"},
				{Thinking: "t2"},
				{Content: "c1"},
				{Thinking: "t3"},
				{Content: "c2"},
			},
			check: []string{
				"Thinking...",
				"c1",
				"c2",
			},
		},
		{
			name:     "empty messages",
			messages: []PrinterLLMMessage{},
			check: []string{
				"LLM Response:",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			w := &fakeWriter{}
			p := NewPrinter(w)

			ch := make(chan PrinterLLMMessage, len(tt.messages))

			go func() {
				for _, m := range tt.messages {
					ch <- m
				}

				close(ch)
			}()

			err := p.PrintLLMResponse(ch)
			require.NoError(t, err)

			out := w.String()

			for _, c := range tt.check {
				require.Contains(t, out, c)
			}
		})
	}
}

func (w *fakeWriter) Write(p []byte) (int, error) {
	w.data = append(w.data, p...)
	return len(p), nil
}

func (w *fakeWriter) String() string {
	return string(w.data)
}
