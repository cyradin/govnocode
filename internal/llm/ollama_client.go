package llm

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
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

type OllamaChatResponse struct {
	Message OllamaChatResponseMessage `json:"message"`
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
		strings.TrimRight(c.baseURL, "/")+"/api/chat",
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

	return ChatResponse{
		Message: c.transformResultMessage(ollamaResp.Message),
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

func (c *OllamaClient) transformResultMessage(message OllamaChatResponseMessage) ChatResponseMessage {
	return ChatResponseMessage{
		Role:     message.Role,
		Content:  message.Content,
		Thinking: message.Thinking,
	}
}
