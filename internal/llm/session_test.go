package llm

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSession_WriteMessage_Success(t *testing.T) {
	t.Parallel()

	client := &mockStreamClient{
		chunks: []ChatResult{
			{Resp: ChatResponse{
				Message: ChatMessage{Content: "Hel"},
			}},
			{Resp: ChatResponse{
				Message: ChatMessage{Content: "lo"},
			}},
		},
	}

	s := NewSession(client, "system prompt")

	out := s.WriteMessage(context.Background(), "hello")

	var results []ChatResult
	for r := range out {
		results = append(results, r)
	}

	require.Len(t, results, 2)

	require.Equal(t, "Hel", results[0].Resp.Message.Content)
	require.Equal(t, "lo", results[1].Resp.Message.Content)
}

func TestSession_WriteMessage_ClientError(t *testing.T) {
	t.Parallel()

	client := &mockStreamClient{
		chunks: []ChatResult{
			{Resp: ChatResponse{Message: ChatMessage{Content: "Hel"}}},
			{Err: errors.New("boom")},
		},
	}

	s := NewSession(client, "system prompt")

	out := s.WriteMessage(context.Background(), "hello")

	var results []ChatResult
	for r := range out {
		results = append(results, r)
	}

	require.Len(t, results, 2)
	require.NoError(t, results[0].Err)
	require.Error(t, results[1].Err)
	require.Contains(t, results[1].Err.Error(), "boom")
}

func TestSession_WriteMessage_AppendsFinalMessage(t *testing.T) {
	t.Parallel()

	client := &mockStreamClient{
		chunks: []ChatResult{
			{Resp: ChatResponse{Message: ChatMessage{Content: "Hel"}}},
			{Resp: ChatResponse{Message: ChatMessage{Content: "lo"}}},
		},
	}

	s := NewSession(client, "system prompt")

	out := s.WriteMessage(context.Background(), "hello")

	for range out {
	}

	require.Len(t, s.messages, 2)
	require.Equal(t, RoleAssistant, s.messages[1].Role)
	require.Equal(t, "Hello", s.messages[1].Content)
}

type mockStreamClient struct {
	chunks []ChatResult
}

func (m *mockStreamClient) Stream(ctx context.Context, messages []ChatMessage) <-chan ChatResult {
	out := make(chan ChatResult)

	go func() {
		defer close(out)

		for _, c := range m.chunks {
			out <- c
		}
	}()

	return out
}
