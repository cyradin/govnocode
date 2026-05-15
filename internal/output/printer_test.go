package output

import (
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

type fakeWriter struct {
	data []byte
}

func TestPrinter_PrintUserMessage(t *testing.T) {
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
				"Message:",
				"hello task",
				"-------------------------------------------------",
			},
		},
		{
			name: "empty message",
			task: "",
			check: []string{
				"Message:",
				"-------------------------------------------------",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			w := &fakeWriter{}
			p := NewPrinter(w)

			err := p.PrintUserMessage(tt.task)
			require.NoError(t, err)

			out := w.String()

			for _, c := range tt.check {
				require.Contains(t, out, c)
			}
		})
	}
}

func TestPrinter_PrintError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		err     error
		wantErr bool
		check   []string
		writer  io.Writer
	}{
		{
			name: "simple error",
			err:  errors.New("something went wrong"),
			check: []string{
				"-------------------------------------------------",
				"Error:",
				"something went wrong",
			},
			writer: &fakeWriter{},
		},
		{
			name: "empty error message",
			err:  errors.New(""),
			check: []string{
				"Error:",
			},
			writer: &fakeWriter{},
		},
		{
			name:    "writer error",
			err:     errors.New("boom"),
			wantErr: true,
			writer:  &errorWriter{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			p := NewPrinter(tt.writer)

			err := p.PrintError(tt.err)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			ww, _ := tt.writer.(*fakeWriter)
			out := ww.String()

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
		messages []LLMMessage
		check    []string
	}{
		{
			name: "content only",
			messages: []LLMMessage{
				{Content: "hello"},
			},
			check: []string{
				"LLM Response:",
				"hello",
			},
		},
		{
			name: "thinking then content",
			messages: []LLMMessage{
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
			messages: []LLMMessage{
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
			messages: []LLMMessage{},
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

			ch := make(chan LLMMessage, len(tt.messages))

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

type errorWriter struct{}

func (e *errorWriter) Write(p []byte) (int, error) {
	return 0, errors.New("write failed")
}
