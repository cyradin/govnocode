package llm

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

type OllamaChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OllamaChatRequest struct {
	Model    string              `json:"model"`
	Messages []OllamaChatMessage `json:"messages"`
	Stream   bool                `json:"stream"`
}

type OllamaChatResponse struct {
	Message OllamaChatMessage `json:"message"`
}

type OllamaClient struct {
	baseURL string
	inner   *http.Client
}

func NewOllamaClient(baseURL string, inner *http.Client) *OllamaClient {
	return &OllamaClient{
		baseURL: baseURL,
		inner:   inner,
	}
}
func (c *OllamaClient) Generate(model string, messages []OllamaChatMessage) (string, error) {
	reqBody := OllamaChatRequest{
		Model:    model,
		Stream:   false,
		Messages: messages,
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	resp, err := c.inner.Post(
		strings.TrimRight(c.baseURL, "/")+"/api/chat",
		"application/json",
		bytes.NewBuffer(data),
	)
	if err != nil {
		return "", err
	}

	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var ollamaResp OllamaChatResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return "", err
	}

	return ollamaResp.Message.Content, nil
}
