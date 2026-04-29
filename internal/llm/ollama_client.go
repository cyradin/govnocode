package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const maxErrorBody = 4 << 10 // 4 KB

type OllamaChatRequest struct {
	Model    string              `json:"model"`
	Messages []OllamaChatMessage `json:"messages"`
	Stream   bool                `json:"stream"`
}

type OllamaChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OllamaDoneReason = string

const (
	OllamaDoneReasonStop   = "stop"
	OllamaDoneReasonLength = "length"
	OllamaDoneReasonLoad   = "load"
	OllamaDoneReasonUnload = "unload"
)

type OllamaChatResponse struct {
	Model     string                     `json:"model,omitempty"`
	CreatedAt string                     `json:"created_at"`
	Message   *OllamaChatResponseMessage `json:"message,omitempty"`

	Done       bool             `json:"done,omitempty"`
	DoneReason OllamaDoneReason `json:"done_reason,omitempty"`

	TotalDuration      int64 `json:"total_duration,omitempty"`
	LoadDuration       int64 `json:"load_duration,omitempty"`
	PromptEvalCount    int   `json:"prompt_eval_count,omitempty"`
	PromptEvalDuration int64 `json:"prompt_eval_duration,omitempty"`
	EvalCount          int   `json:"eval_count,omitempty"`
	EvalDuration       int64 `json:"eval_duration,omitempty"`
}

type OllamaChatResponseMessage struct {
	Role     string `json:"role"`
	Content  string `json:"content"`
	Thinking string `json:"thinking"`
}

type OllamaClient struct {
	baseURL string
	model   string
	inner   *http.Client
}

func NewOllamaClient(baseURL string, model string, inner *http.Client) *OllamaClient {
	return &OllamaClient{
		baseURL: baseURL,
		model:   model,
		inner:   inner,
	}
}

func (c *OllamaClient) Generate(messages []ChatMessage) (ChatResponse, error) {
	reqBody := OllamaChatRequest{
		Model:    c.model,
		Stream:   false,
		Messages: c.transformMessages(messages),
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return ChatResponse{}, err
	}

	resp, err := c.inner.Post(
		c.makeURL("/api/chat"),
		"application/json",
		bytes.NewBuffer(data),
	)
	if err != nil {
		return ChatResponse{}, err
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return ChatResponse{}, c.readHTTPError(resp)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ChatResponse{}, err
	}

	var ollamaResp OllamaChatResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return ChatResponse{}, err
	}

	if ollamaResp.DoneReason != OllamaDoneReasonStop {
		return ChatResponse{}, fmt.Errorf("request failed with done reason: %s", ollamaResp.DoneReason)
	}

	return c.transformChatResponse(ollamaResp), nil
}

func (c *OllamaClient) Stream(messages []ChatMessage) (<-chan ChatResponse, <-chan error) {
	out := make(chan ChatResponse)
	errCh := make(chan error, 1)

	go func() {
		defer close(out)
		defer close(errCh)

		reqBody := OllamaChatRequest{
			Model:    c.model,
			Stream:   true,
			Messages: c.transformMessages(messages),
		}

		data, err := json.Marshal(reqBody)
		if err != nil {
			errCh <- err
			return
		}

		resp, err := c.inner.Post(
			strings.TrimRight(c.baseURL, "/")+"/api/chat",
			"application/json",
			bytes.NewReader(data),
		)
		if err != nil {
			errCh <- err
			return
		}

		defer func() { _ = resp.Body.Close() }()

		if resp.StatusCode != http.StatusOK {
			errCh <- c.readHTTPError(resp)

			return
		}

		dec := json.NewDecoder(resp.Body)

		for {
			var chunk OllamaChatResponse

			if err := dec.Decode(&chunk); err != nil {
				if err == io.EOF {
					return
				}

				errCh <- err

				return
			}

			if chunk.Done {
				if chunk.DoneReason != OllamaDoneReasonStop {
					errCh <- fmt.Errorf("request failed with done reason: %s", chunk.DoneReason)

					return
				}

				out <- c.transformChatResponse(chunk)

				break
			}

			out <- c.transformChatResponse(chunk)
		}
	}()

	return out, errCh
}

func (c *OllamaClient) transformMessages(messages []ChatMessage) []OllamaChatMessage {
	result := make([]OllamaChatMessage, 0, len(messages))

	for _, m := range messages {
		result = append(result, OllamaChatMessage{
			Role:    m.Role,
			Content: m.Content,
		})
	}

	return result
}

func (c *OllamaClient) transformChatResponse(resp OllamaChatResponse) ChatResponse {
	return ChatResponse{
		Message:  c.transformResultMessage(resp.Message),
		Metadata: c.transformMetadata(resp),
	}
}

func (c *OllamaClient) transformResultMessage(message *OllamaChatResponseMessage) ChatResponseMessage {
	if message == nil {
		return ChatResponseMessage{}
	}

	return ChatResponseMessage{
		Role:     message.Role,
		Content:  message.Content,
		Thinking: message.Thinking,
	}
}

func (c *OllamaClient) transformMetadata(resp OllamaChatResponse) ChatResponseMetadata {
	return ChatResponseMetadata{
		PromptProcessingTime:   time.Duration(resp.PromptEvalDuration),
		ResponseProcessingTime: time.Duration(resp.EvalDuration),
		PromptTokensUsed:       resp.PromptEvalCount,
		ResponseTokenUsed:      resp.EvalCount,
	}
}

func (c *OllamaClient) makeURL(path string) string {
	return strings.TrimRight(c.baseURL, "/") + path
}

func (c *OllamaClient) readHTTPError(resp *http.Response) error {
	body, _ := io.ReadAll(io.LimitReader(resp.Body, maxErrorBody))
	return fmt.Errorf(
		"unexpected status %d: %s",
		resp.StatusCode,
		strings.TrimSpace(string(body)),
	)
}
