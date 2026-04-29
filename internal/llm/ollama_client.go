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

	return ChatResponse{
		Message:  c.transformResultMessage(ollamaResp.Message),
		Metadata: c.transformMetadata(ollamaResp),
	}, nil
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

func (c *OllamaClient) transformMetadata(resp OllamaChatResponse) CharResponseMetadata {
	return CharResponseMetadata{
		PromptProcessingTime:   time.Duration(resp.PromptEvalDuration),
		ResponseProcessingTime: time.Duration(resp.EvalDuration),
		PromptTokensUsed:       resp.PromptEvalCount,
		ResponseTokenUsed:      resp.EvalCount,
	}
}

func (c *OllamaClient) makeURL(path string) string {
	return strings.TrimRight(c.baseURL, "/") + path
}
