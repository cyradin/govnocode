package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
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

func (c *OllamaClient) Generate(ctx context.Context, messages []ChatMessage) (ChatResponse, error) {
	req, err := c.encodeChatRequest(ctx, messages, false)
	if err != nil {
		return ChatResponse{}, fmt.Errorf("create request: %w", err)
	}

	resp, err := c.inner.Do(req)
	if err != nil {
		return ChatResponse{}, fmt.Errorf("perform request: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return ChatResponse{}, c.readHTTPError(resp)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ChatResponse{}, fmt.Errorf("read body: %w", err)
	}

	var ollamaResp OllamaChatResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return ChatResponse{}, fmt.Errorf("unmarshal response: %w", err)
	}

	if ollamaResp.DoneReason != OllamaDoneReasonStop {
		return ChatResponse{}, fmt.Errorf("request failed with done reason: %s", ollamaResp.DoneReason)
	}

	return c.transformChatResponse(ollamaResp), nil
}

func (c *OllamaClient) Stream(ctx context.Context, messages []ChatMessage) <-chan ChatResult {
	out := make(chan ChatResult)

	go func() {
		defer close(out)

		req, err := c.encodeChatRequest(ctx, messages, true)
		if err != nil {
			out <- ChatResult{Err: fmt.Errorf("create request: %w", err)}

			return
		}

		resp, err := c.inner.Do(req)
		if err != nil {
			out <- ChatResult{Err: fmt.Errorf("perform request: %w", err)}

			return
		}

		defer func() { _ = resp.Body.Close() }()

		if resp.StatusCode != http.StatusOK {
			out <- ChatResult{Err: c.readHTTPError(resp)}

			return
		}

		dec := json.NewDecoder(resp.Body)

		for {
			var chunk OllamaChatResponse

			if err := dec.Decode(&chunk); err != nil {
				if errors.Is(err, io.EOF) {
					return
				}

				out <- ChatResult{Err: fmt.Errorf("unmarshal response: %w", err)}

				return
			}

			if chunk.Done {
				if chunk.DoneReason != OllamaDoneReasonStop {
					out <- ChatResult{
						Resp: c.transformChatResponse(chunk),
						Err:  fmt.Errorf("request failed with done reason: %s", chunk.DoneReason),
					}

					return
				}

				out <- ChatResult{Resp: c.transformChatResponse(chunk)}

				break
			}

			out <- ChatResult{Resp: c.transformChatResponse(chunk)}
		}
	}()

	return out
}

func (c *OllamaClient) encodeChatRequest(ctx context.Context, messages []ChatMessage, stream bool) (*http.Request, error) {
	reqBody := OllamaChatRequest{
		Model:    c.model,
		Stream:   stream,
		Messages: c.transformMessages(messages),
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("json marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.makeURL("/api/chat"),
		bytes.NewReader(data),
	)
	if err != nil {
		return nil, fmt.Errorf("make http request: %w", err)
	}

	return req, nil
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

func (c *OllamaClient) transformResultMessage(message *OllamaChatResponseMessage) ChatMessage {
	if message == nil {
		return ChatMessage{}
	}

	return ChatMessage{
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
