package openai

import openai "github.com/sashabaranov/go-openai"

var (
	DefaultSystemRole = "You are a helpful assistant"
)

type Session struct {
	Role    string
	History []openai.ChatCompletionMessage
}

func NewSession() *Session {
	return &Session{}
}

func (s *Session) AddToHistory(messages ...openai.ChatCompletionMessage) {
	s.History = append(s.History, messages...)
}

func (s *Session) ClearHistory() {
	s.History = nil
}

func (s *Session) GetHistory() []openai.ChatCompletionMessage {
	return s.History
}

func (s *Session) ClearRole() {
	s.Role = ""
}

func (s *Session) GetRole() string {
	if s.Role == "" {
		return DefaultSystemRole
	}
	return s.Role
}

func (s *Session) SetRole(role string) {
	s.Role = role
}
