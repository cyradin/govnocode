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
					"prompt_eval_count": 805,
					"prompt_eval_duration": 724666428,
					"eval_count": 74,
					"eval_duration": 2913794162
				}`))
			},
			wantErr: false,
			expectedOutput: ChatResponse{
				Message: ChatMessage{
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
				_, _ = w.Write([]byte(`{"error":"fail"}`))
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

			ctx := t.Context()

			out, err := client.Generate(ctx, []ChatMessage{
				{Role: "user", Content: "hello"},
			})

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

	stream := client.Stream(t.Context(), []ChatMessage{
		{Role: "user", Content: "hello"},
	})

	var results []ChatResult

	for r := range stream {
		require.NoError(t, r.Err)
		results = append(results, r)
	}

	require.Len(t, results, 3)

	require.Equal(t, "Hel", results[0].Resp.Message.Content)
	require.Equal(t, "lo", results[1].Resp.Message.Content)

	require.Equal(t, 10, results[2].Resp.Metadata.PromptTokensUsed)
	require.Equal(t, 2, results[2].Resp.Metadata.ResponseTokenUsed)
}

func TestOllamaClient_Stream_HTTPError(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"bad request"}`))
	}))
	defer server.Close()

	client := NewOllamaClient(server.URL, "test-model", server.Client())

	stream := client.Stream(t.Context(), []ChatMessage{
		{Role: "user", Content: "hello"},
	})

	var gotErr error
	for r := range stream {
		gotErr = r.Err
	}

	require.Error(t, gotErr)
	require.Contains(t, gotErr.Error(), "unexpected status")
	require.Contains(t, gotErr.Error(), "400")
}

func TestOllamaClient_Stream_HTTPError_LargeBody(t *testing.T) {
	t.Parallel()

	large := strings.Repeat("a", 10<<10)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(large))
	}))
	defer server.Close()

	client := NewOllamaClient(server.URL, "test-model", server.Client())

	stream := client.Stream(t.Context(), []ChatMessage{
		{Role: "user", Content: "hello"},
	})

	var gotErr error
	for r := range stream {
		gotErr = r.Err
	}

	require.Error(t, gotErr)
	require.Less(t, len(gotErr.Error()), maxErrorBody+200)
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

	stream := client.Stream(t.Context(), []ChatMessage{
		{Role: "user", Content: "hello"},
	})

	var results []ChatResult
	for r := range stream {
		results = append(results, r)
	}

	require.Len(t, results, 1)
	require.Error(t, results[0].Err)
	require.Contains(t, results[0].Err.Error(), "length")
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

	stream := client.Stream(t.Context(), []ChatMessage{
		{Role: "user", Content: "hello"},
	})

	var gotErr error
	for r := range stream {
		gotErr = r.Err
	}

	require.Error(t, gotErr)
}

func TestOllamaClient_Stream_PostError(t *testing.T) {
	t.Parallel()

	brokenClient := &http.Client{
		Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
			return nil, errors.New("network error")
		}),
	}

	client := NewOllamaClient("http://localhost:1234", "test-model", brokenClient)

	stream := client.Stream(t.Context(), []ChatMessage{
		{Role: "user", Content: "hello"},
	})

	var gotErr error
	for r := range stream {
		gotErr = r.Err
	}

	require.Error(t, gotErr)
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}
