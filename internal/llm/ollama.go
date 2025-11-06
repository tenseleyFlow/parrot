package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type OllamaClient struct {
	BaseURL string
	Model   string
	client  *http.Client
}

type GenerateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type GenerateResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

func NewOllamaClient(baseURL, model string) *OllamaClient {
	if baseURL == "" {
		baseURL = "http://127.0.0.1:11434" // Use IPv4 explicitly to avoid IPv6 issues
	}
	if model == "" {
		model = "llama3.2:3b"
	}
	
	return &OllamaClient{
		BaseURL: baseURL,
		Model:   model,
		client: &http.Client{
			Timeout: 60 * time.Second, // Maximum timeout; actual timeout controlled by context
		},
	}
}

func (c *OllamaClient) Generate(ctx context.Context, prompt string) (string, error) {
	u, err := url.JoinPath(c.BaseURL, "/api/generate")
	if err != nil {
		return "", fmt.Errorf("invalid base URL: %w", err)
	}

	req := GenerateRequest{
		Model:  c.Model,
		Prompt: prompt,
		Stream: false,
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", u, bytes.NewReader(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("failed to send request to %s: %w", c.BaseURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Read the actual error response body
		bodyBytes := make([]byte, 512) // Read first 512 bytes for error message
		n, _ := resp.Body.Read(bodyBytes)
		bodyStr := string(bodyBytes[:n])
		return "", fmt.Errorf("ollama API returned status %d: %s", resp.StatusCode, bodyStr)
	}

	var genResp GenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&genResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return genResp.Response, nil
}

func (c *OllamaClient) IsAvailable() bool {
	u, err := url.JoinPath(c.BaseURL, "/api/version")
	if err != nil {
		return false
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
	if err != nil {
		return false
	}
	
	resp, err := c.client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	
	return resp.StatusCode == http.StatusOK
}

// WarmupModel preloads the model to avoid cold start delays
func (c *OllamaClient) WarmupModel() error {
	u, err := url.JoinPath(c.BaseURL, "/api/generate")
	if err != nil {
		return fmt.Errorf("invalid base URL: %w", err)
	}

	req := GenerateRequest{
		Model:  c.Model,
		Prompt: "test", // Minimal prompt to load model
		Stream: false,
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	// Use a longer timeout for initial model loading
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	httpReq, err := http.NewRequestWithContext(ctx, "POST", u, bytes.NewReader(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to warmup model: %w", err)
	}
	defer resp.Body.Close()

	return nil
}