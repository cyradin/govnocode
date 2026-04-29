package llm

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestOllamaClient_Generate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		serverHandler  http.HandlerFunc
		wantErr        bool
		expectedOutput ChatResponse
	}{
		{
			name: "success",
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, http.MethodPost, r.Method)
				require.Equal(t, "/api/chat", r.URL.Path)

				body, err := io.ReadAll(r.Body)
				require.NoError(t, err)

				require.Contains(t, string(body), `"model":"test-model"`)

				require.Contains(t, string(body), `"role":"user"`)
				require.Contains(t, string(body), `"content":"hello"`)

				_, _ = w.Write([]byte(`{
					"model": "batiai/gemma4-26b:iq4",
					"created_at": "2026-04-29T11:26:02.135006047Z",
					"message": {
						"role": "assistant",
						"content": "message",
						"thinking": "thinking"
					},
					"done": true,
					"done_reason": "stop",
					"total_duration": 3773606103,
					"load_duration": 108490601,
					"prompt_eval_count": 805,
					"prompt_eval_duration": 724666428,
					"eval_count": 74,
					"eval_duration": 2913794162
				}`))
			},
			wantErr: false,
			expectedOutput: ChatResponse{
				Message: ChatResponseMessage{
					Role:     "assistant",
					Content:  "message",
					Thinking: "thinking",
				},
				Metadata: ChatResponseMetadata{
					PromptProcessingTime:   time.Duration(724666428),
					ResponseProcessingTime: time.Duration(2913794162),
					PromptTokensUsed:       805,
					ResponseTokenUsed:      74,
				},
			},
		},
		{
			name: "done reason not stop",
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte(`{
					"message": {
						"role": "assistant",
						"content": "partial"
					},
					"done": true,
					"done_reason": "length"
				}`))
			},
			wantErr: true,
		},
		{
			name: "http error",
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			wantErr: true,
		},
		{
			name: "invalid json",
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte(`invalid-json`))
			},
			wantErr: true,
		},
		{
			name: "read body error",
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				hijacker, ok := w.(http.Hijacker)
				require.True(t, ok)

				conn, _, err := hijacker.Hijack()
				require.NoError(t, err)

				_ = conn.Close()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(tt.serverHandler)
			defer server.Close()

			client := NewOllamaClient(server.URL, "test-model", server.Client())

			messages := []ChatMessage{
				{Role: "user", Content: "hello"},
			}

			out, err := client.Generate(messages)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.expectedOutput, out)
		})
	}
}

func TestOllamaClient_Stream_Success(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/api/chat", r.URL.Path)

		flusher, ok := w.(http.Flusher)
		require.True(t, ok)

		chunks := []string{
			`{"message":{"role":"assistant","content":"Hel"},"done":false}`,
			`{"message":{"role":"assistant","content":"lo"},"done":false}`,
			`{
				"message":{"role":"assistant","content":""},
				"done":true,
				"done_reason":"stop",
				"prompt_eval_count":10,
				"eval_count":2,
				"prompt_eval_duration":100,
				"eval_duration":200
			}`,
		}

		for _, chunk := range chunks {
			_, err := w.Write([]byte(chunk + "\n"))
			require.NoError(t, err)
			flusher.Flush()
		}
	}))
	defer server.Close()

	client := NewOllamaClient(server.URL, "test-model", server.Client())

	stream, errCh := client.Stream([]ChatMessage{
		{Role: "user", Content: "hello"},
	})

	var results []ChatResponse

	for msg := range stream {
		results = append(results, msg)
	}

	require.NoError(t, <-errCh)

	require.Len(t, results, 3)

	require.Equal(t, "Hel", results[0].Message.Content)
	require.Equal(t, "lo", results[1].Message.Content)

	require.Equal(t, 10, results[2].Metadata.PromptTokensUsed)
	require.Equal(t, 2, results[2].Metadata.ResponseTokenUsed)
}

func TestOllamaClient_Stream_HTTPError(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"bad request"}`))
	}))
	defer server.Close()

	client := NewOllamaClient(server.URL, "test-model", server.Client())

	stream, errCh := client.Stream([]ChatMessage{
		{Role: "user", Content: "hello"},
	})

	for range stream {
	}

	err := <-errCh
	require.Error(t, err)
	require.Contains(t, err.Error(), "unexpected status")
	require.Contains(t, err.Error(), "400")
}

func TestOllamaClient_Stream_HTTPError_LargeBody(t *testing.T) {
	t.Parallel()

	large := strings.Repeat("a", 10<<10) // 10KB

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(large))
	}))
	defer server.Close()

	client := NewOllamaClient(server.URL, "test-model", server.Client())

	stream, errCh := client.Stream([]ChatMessage{
		{Role: "user", Content: "hello"},
	})

	for range stream {
	}

	err := <-errCh
	require.Error(t, err)

	require.Less(t, len(err.Error()), maxErrorBody+100)
}

func TestOllamaClient_Stream_DoneReasonError(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
		require.True(t, ok)

		_, _ = w.Write([]byte(`{"done":true,"done_reason":"length"}` + "\n"))
		flusher.Flush()
	}))
	defer server.Close()

	client := NewOllamaClient(server.URL, "test-model", server.Client())

	stream, errCh := client.Stream([]ChatMessage{
		{Role: "user", Content: "hello"},
	})

	var results []ChatResponse
	for msg := range stream {
		results = append(results, msg)
	}

	require.Empty(t, results)

	err := <-errCh
	require.Error(t, err)
	require.Contains(t, err.Error(), "length")
}

func TestOllamaClient_Stream_InvalidJSON(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
		require.True(t, ok)

		_, _ = w.Write([]byte(`invalid-json`))
		flusher.Flush()
	}))
	defer server.Close()

	client := NewOllamaClient(server.URL, "test-model", server.Client())

	stream, errCh := client.Stream([]ChatMessage{
		{Role: "user", Content: "hello"},
	})

	for range stream {
		// nothing
	}

	require.Error(t, <-errCh)
}

func TestOllamaClient_Stream_PostError(t *testing.T) {
	t.Parallel()

	brokenClient := &http.Client{
		Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
			return nil, errors.New("network error")
		}),
	}

	client := NewOllamaClient("http://localhost:1234", "test-model", brokenClient)

	stream, errCh := client.Stream([]ChatMessage{
		{Role: "user", Content: "hello"},
	})

	for range stream {
	}

	require.Error(t, <-errCh)
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}
