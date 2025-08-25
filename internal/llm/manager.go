package llm

import (
	"context"
	"fmt"
	"strings"
	"time"

	"parrot/internal/config"
)

type LLMManager struct {
	config     *config.Config
	apiClient  *APIClient
	ollamaClient *OllamaClient
}

type Backend string

const (
	BackendAPI      Backend = "api"
	BackendLocal    Backend = "local"  
	BackendFallback Backend = "fallback"
)

func NewLLMManager(cfg *config.Config) *LLMManager {
	manager := &LLMManager{
		config: cfg,
	}
	
	// Initialize API client if enabled
	if cfg.API.Enabled && cfg.API.APIKey != "" {
		manager.apiClient = NewAPIClient(
			cfg.API.Endpoint,
			cfg.API.APIKey,
			cfg.API.Model,
			cfg.API.Timeout,
		)
	}
	
	// Initialize Ollama client if enabled
	if cfg.Local.Enabled {
		manager.ollamaClient = NewOllamaClient(
			cfg.Local.Endpoint,
			cfg.Local.Model,
		)
	}
	
	return manager
}

func (m *LLMManager) Generate(ctx context.Context, prompt string, commandType string) (string, Backend) {
	// If fallback mode is enabled, skip LLM backends
	if m.config.General.FallbackMode {
		return m.generateFallback(commandType), BackendFallback
	}
	
	// Try backends in priority order: API -> Local -> Fallback
	
	// 1. Try API first (if available)
	if m.apiClient != nil && m.config.API.Enabled {
		if m.config.General.Debug {
			fmt.Printf("🔍 Trying API backend...\n")
		}
		
		response, err := m.apiClient.Generate(ctx, prompt)
		if err == nil && response != "" {
			response = m.cleanResponse(response)
			if m.config.General.Debug {
				fmt.Printf("✅ API backend succeeded\n")
			}
			return response, BackendAPI
		}
		
		if m.config.General.Debug {
			fmt.Printf("❌ API backend failed: %v\n", err)
		}
	}
	
	// 2. Try local Ollama (if available)
	if m.ollamaClient != nil && m.config.Local.Enabled {
		if m.config.General.Debug {
			fmt.Printf("🔍 Trying local backend...\n")
		}
		
		// Create timeout context for local calls
		localCtx, cancel := context.WithTimeout(ctx, time.Duration(m.config.Local.Timeout)*time.Second)
		defer cancel()
		
		response, err := m.ollamaClient.Generate(localCtx, prompt)
		if err == nil && response != "" {
			response = m.cleanResponse(response)
			if m.config.General.Debug {
				fmt.Printf("✅ Local backend succeeded\n")
			}
			return response, BackendLocal
		}
		
		if m.config.General.Debug {
			fmt.Printf("❌ Local backend failed: %v\n", err)
		}
	}
	
	// 3. Fallback to hardcoded responses
	if m.config.General.Debug {
		fmt.Printf("🔄 Using fallback backend\n")
	}
	return m.generateFallback(commandType), BackendFallback
}

func (m *LLMManager) cleanResponse(response string) string {
	// Clean up the response
	response = strings.TrimSpace(response)
	
	// Remove tertiary output that starts with newline + "(Note:"
	if idx := strings.Index(response, "\n(Note:"); idx != -1 {
		response = response[:idx]
		response = strings.TrimSpace(response)
	}
	
	// Remove common prefixes from LLMs
	prefixes := []string{
		"Response:",
		"Parrot says:",
		"🦜",
	}
	
	for _, prefix := range prefixes {
		if strings.HasPrefix(response, prefix) {
			response = strings.TrimSpace(response[len(prefix):])
		}
	}
	
	// Remove quotes if the entire response is quoted
	if len(response) >= 2 && response[0] == '"' && response[len(response)-1] == '"' {
		response = response[1 : len(response)-1]
	}
	
	// Ensure response isn't too long (keep it snappy)
	if len(response) > 150 {
		// Try to cut at sentence boundary
		if idx := strings.LastIndex(response[:150], "."); idx > 50 {
			response = response[:idx+1]
		} else {
			response = response[:147] + "..."
		}
	}
	
	return strings.TrimSpace(response)
}

func (m *LLMManager) generateFallback(commandType string) string {
	fallbacks := map[string][]string{
		"git": {
			"Git good? More like git rekt!",
			"Did you forget to pull again? Classic amateur move.",
			"Another git genius strikes again!",
			"Your commits are as broken as your workflow.",
		},
		"nodejs": {
			"NPM install failed? Shocking! Nobody saw that coming.",
			"Your package.json is crying. Fix it.",
			"Node modules: where dependencies go to die.",
			"Even npm doesn't want to deal with your code.",
		},
		"docker": {
			"Docker container more like docker DISASTER!",
			"Even containers can't contain your incompetence.",
			"Your Dockerfile needs therapy.",
			"Container exit code: user error detected.",
		},
		"http": {
			"404: Competence not found.",
			"Even the internet doesn't want to talk to you.",
			"Connection refused? So is your logic.",
			"HTTP status: 500 Internal User Error.",
		},
		"generic": {
			"Wow, you managed to break something simple. Impressive!",
			"Maybe try reading the manual... oh wait, who am I kidding?",
			"Error code says it all: user error!",
			"Have you tried turning your brain on and off again?",
		},
	}
	
	responses, exists := fallbacks[commandType]
	if !exists {
		responses = fallbacks["generic"]
	}
	
	// Simple pseudo-random selection based on command type
	hash := 0
	for _, char := range commandType {
		hash = hash*31 + int(char)
	}
	if hash < 0 {
		hash = -hash
	}
	
	return responses[hash%len(responses)]
}

func (m *LLMManager) GetStatus() map[string]interface{} {
	status := map[string]interface{}{
		"fallback_mode": m.config.General.FallbackMode,
		"debug":         m.config.General.Debug,
		"personality":   m.config.General.Personality,
	}
	
	// Check API status
	if m.apiClient != nil && m.config.API.Enabled {
		status["api_enabled"] = true
		status["api_provider"] = m.config.API.Provider
		status["api_model"] = m.config.API.Model
		status["api_available"] = m.apiClient.IsAvailable()
	} else {
		status["api_enabled"] = false
		status["api_available"] = false
	}
	
	// Check local status  
	if m.ollamaClient != nil && m.config.Local.Enabled {
		status["local_enabled"] = true
		status["local_provider"] = m.config.Local.Provider
		status["local_model"] = m.config.Local.Model
		status["local_available"] = m.ollamaClient.IsAvailable()
	} else {
		status["local_enabled"] = false
		status["local_available"] = false
	}
	
	return status
}