package llm

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOllamaClient_Generate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		serverHandler  http.HandlerFunc
		wantErr        bool
		expectedOutput string
	}{
		{
			name: "success",
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, http.MethodPost, r.Method)
				require.Equal(t, "/api/chat", r.URL.Path)

				body, err := io.ReadAll(r.Body)
				require.NoError(t, err)

				// проверяем model
				require.Contains(t, string(body), `"model":"test-model"`)

				// проверяем messages
				require.Contains(t, string(body), `"role":"user"`)
				require.Contains(t, string(body), `"content":"hello"`)

				_, _ = w.Write([]byte(`{
					"message": {
						"role": "assistant",
						"content": "echo hello"
					}
				}`))
			},
			wantErr:        false,
			expectedOutput: "echo hello",
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

			client := NewOllamaClient(server.URL, server.Client())

			messages := []OllamaChatMessage{
				{Role: "user", Content: "hello"},
			}

			out, err := client.Generate("test-model", messages)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.expectedOutput, out)
		})
	}
}

func TestOllamaClient_Generate_PostError(t *testing.T) {
	t.Parallel()

	brokenClient := &http.Client{
		Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
			return nil, errors.New("network error")
		}),
	}

	client := NewOllamaClient("http://localhost:1234", brokenClient)

	messages := []OllamaChatMessage{
		{Role: "user", Content: "p"},
	}

	out, err := client.Generate("m", messages)

	require.Error(t, err)
	require.Empty(t, out)
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}
