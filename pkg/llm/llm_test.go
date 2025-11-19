package llm

import (
	"context"
	"testing"
)

func TestMockClient(t *testing.T) {
	c := &MockClient{}
	resp, err := c.Complete(context.Background(), "test")
	if err != nil {
		t.Fatalf("Complete failed: %v", err)
	}
	if resp != "Mock response" {
		t.Errorf("Expected 'Mock response', got %q", resp)
	}
}

func TestOpenAIClient(t *testing.T) {
	c := NewOpenAIClient("test-key")
	if c.APIKey != "test-key" {
		t.Errorf("Expected APIKey 'test-key', got %q", c.APIKey)
	}

	_, err := c.Complete(context.Background(), "test")
	if err == nil {
		t.Error("Expected error from unimplemented Complete")
	}
}
