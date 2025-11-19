package llm

import (
	"context"
	"fmt"
)

type Client interface {
	Complete(ctx context.Context, prompt string) (string, error)
}

type MockClient struct{}

func (m *MockClient) Complete(ctx context.Context, prompt string) (string, error) {
	// Mock response for testing
	return "Mock response", nil
}

// OpenAIClient is a simple wrapper for OpenAI API
// (Implementation omitted for brevity, but structure is here)
type OpenAIClient struct {
	APIKey string
}

func NewOpenAIClient(apiKey string) *OpenAIClient {
	return &OpenAIClient{APIKey: apiKey}
}

func (c *OpenAIClient) Complete(ctx context.Context, prompt string) (string, error) {
	// TODO: Implement actual HTTP call
	return "", fmt.Errorf("not implemented")
}
