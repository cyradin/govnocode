package llm

import (
	"context"
	"slices"
	"time"
)

type Role = string

const (
	RoleSystem    = "system"
	RoleUser      = "user"
	RoleAssistant = "assistant"
)

type ChatMessage struct {
	Role     string
	Content  string
	Thinking string
}

type ChatResponse struct {
	Message  ChatMessage
	Metadata ChatResponseMetadata
}

type ChatResult struct {
	Resp ChatResponse
	Err  error
}

type ChatResponseMetadata struct {
	PromptProcessingTime   time.Duration
	ResponseProcessingTime time.Duration

	PromptTokensUsed  int
	ResponseTokenUsed int
}

type client interface {
	Stream(ctx context.Context, messages []ChatMessage) <-chan ChatResult
}

type Session struct {
	messages []ChatMessage
	client   client
}

func NewSession(
	client client,
	systemPrompt string,
) *Session {
	return &Session{
		client: client,
		messages: []ChatMessage{
			{
				Content: systemPrompt,
				Role:    RoleSystem,
			},
		},
	}
}

func (s *Session) WriteMessage(ctx context.Context, text string) <-chan ChatResult {
	messages := slices.Clone(s.messages)
	messages = append(messages, ChatMessage{
		Role:    RoleUser,
		Content: text,
	})

	in := s.client.Stream(ctx, messages)
	out := make(chan ChatResult)

	go func() {
		defer close(out)

		message := ChatMessage{
			Role:    RoleAssistant,
			Content: "",
		}

		for r := range in {
			if err := r.Err; err != nil {
				out <- r

				return
			}

			message.Content += r.Resp.Message.Content
			message.Thinking += r.Resp.Message.Thinking

			out <- r
		}

		s.messages = append(s.messages, message)
	}()

	return out
}

func (s *Session) LastMessage() ChatMessage {
	if len(s.messages) == 0 {
		return ChatMessage{}
	}

	return s.messages[len(s.messages)-1]
}
