package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type APIClient struct {
	Endpoint string
	APIKey   string
	Model    string
	client   *http.Client
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
	Temperature float64       `json:"temperature,omitempty"`
}

type ChatChoice struct {
	Message      ChatMessage `json:"message"`
	FinishReason string      `json:"finish_reason"`
}

type ChatResponse struct {
	Choices []ChatChoice `json:"choices"`
	Error   *APIError    `json:"error,omitempty"`
}

type APIError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Code    string `json:"code"`
}

func NewAPIClient(endpoint, apiKey, model string, timeout int) *APIClient {
	return &APIClient{
		Endpoint: endpoint,
		APIKey:   apiKey,
		Model:    model,
		client: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
	}
}

func (c *APIClient) Generate(ctx context.Context, prompt string) (string, error) {
	if c.APIKey == "" {
		return "", fmt.Errorf("API key not configured")
	}

	// Build chat request
	req := ChatRequest{
		Model: c.Model,
		Messages: []ChatMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		MaxTokens:   150, // Keep responses concise
		Temperature: 0.8, // Creative but focused
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	endpoint := c.Endpoint + "/chat/completions"
	httpReq, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.APIKey)

	// Send request
	resp, err := c.client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Parse response
	var chatResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	// Check for API errors
	if chatResp.Error != nil {
		return "", fmt.Errorf("API error: %s", chatResp.Error.Message)
	}

	// Extract response
	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no response choices returned")
	}

	response := chatResp.Choices[0].Message.Content
	if response == "" {
		return "", fmt.Errorf("empty response from API")
	}

	return response, nil
}

func (c *APIClient) IsAvailable() bool {
	if c.APIKey == "" {
		return false
	}
	
	// Simple check - try to create a request
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	// Create a minimal test request
	req := ChatRequest{
		Model: c.Model,
		Messages: []ChatMessage{
			{Role: "user", Content: "test"},
		},
		MaxTokens: 1,
	}
	
	reqBody, err := json.Marshal(req)
	if err != nil {
		return false
	}
	
	endpoint := c.Endpoint + "/chat/completions"
	httpReq, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(reqBody))
	if err != nil {
		return false
	}
	
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.APIKey)
	
	resp, err := c.client.Do(httpReq)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	
	// Consider 2xx status codes as available
	return resp.StatusCode >= 200 && resp.StatusCode < 300
}