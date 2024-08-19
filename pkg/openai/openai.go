package openai

import (
	"context"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
)

type Client struct {
	client   *openai.Client
	sessions map[int64]*Session
}

// New creates a new OpenAI API client.
func New(apiKey string) Client {
	client := openai.NewClient(apiKey)

	return Client{
		client:   client,
		sessions: make(map[int64]*Session),
	}
}

// ChatCompletion sends a prompt to the OpenAI API and returns the reponse.
func (c Client) ChatCompletion(ctx context.Context, key int64, prompt string) (string, error) {
	session := c.GetSession(key)

	// Set system role
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: session.GetRole(),
		},
	}

	// Add history context
	messages = append(messages, session.GetHistory()...)

	// Add user input
	input := openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: prompt,
	}
	messages = append(messages, input)

	// Send request
	resp, err := c.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
			Messages: messages,
		},
	)

	if err != nil {
		return "", fmt.Errorf("ChatCompletion error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("ChatCompletion response has no choices")
	}

	// Add input and response to history
	session.AddToHistory(input, resp.Choices[0].Message)

	return resp.Choices[0].Message.Content, nil
}

func (c Client) GetSession(key int64) *Session {
	if _, ok := c.sessions[key]; !ok {
		c.sessions[key] = NewSession()
	}

	return c.sessions[key]
}

func (c Client) GetHistory(key int64) []openai.ChatCompletionMessage {
	session := c.GetSession(key)
	return session.GetHistory()
}

func (c Client) GetRole(key int64) string {
	session := c.GetSession(key)
	return session.GetRole()
}

func (c Client) ClearHistory(key int64) {
	session := c.GetSession(key)
	session.ClearHistory()
}

func (c Client) ClearRole(key int64) {
	session := c.GetSession(key)
	session.ClearRole()
}

func (c Client) SetRole(key int64, role string) {
	session := c.GetSession(key)
	session.SetRole(role)
}
