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
		
		// Warm up the model in the background for better performance
		if manager.ollamaClient.IsAvailable() {
			go func() {
				if err := manager.ollamaClient.WarmupModel(); err != nil && cfg.General.Debug {
					fmt.Printf("🔥 Model warmup failed: %v\n", err)
				} else if cfg.General.Debug {
					fmt.Printf("🔥 Model warmed up successfully\n")
				}
			}()
		}
	}
	
	return manager
}

func (m *LLMManager) Generate(ctx context.Context, prompt string, commandType string) (string, Backend) {
	// If fallback mode is enabled, skip LLM backends
	if m.config.General.FallbackMode {
		return m.generateFallback(commandType, "", ""), BackendFallback
	}
	
	// Try backends in priority order: API -> Local -> Fallback
	
	// 1. Try API first (if available)
	if m.apiClient != nil && m.config.API.Enabled {
		if m.config.General.Debug {
			fmt.Printf("🔍 Trying API backend...\n")
		}

		// Create timeout context for API calls
		timeoutDuration := time.Duration(m.config.API.Timeout) * time.Second
		apiCtx, cancel := context.WithTimeout(ctx, timeoutDuration)
		defer cancel()

		response, err := m.apiClient.Generate(apiCtx, prompt)
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
		timeoutDuration := time.Duration(m.config.Local.Timeout) * time.Second
		localCtx, cancel := context.WithTimeout(ctx, timeoutDuration)
		defer cancel()
		
		response, err := m.ollamaClient.Generate(localCtx, prompt)
		if m.config.General.Debug {
			fmt.Printf("🐛 Raw Ollama response: '%s', error: %v\n", response, err)
		}
		if err == nil && response != "" {
			response = m.cleanResponse(response)
			if m.config.General.Debug {
				fmt.Printf("✅ Local backend succeeded with: '%s'\n", response)
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
	return m.generateFallback(commandType, "", ""), BackendFallback
}

// GenerateWithContext generates a response with full context for intelligent fallbacks
func (m *LLMManager) GenerateWithContext(ctx context.Context, prompt string, commandType string, fullCommand string, exitCode string) (string, Backend) {
	// If fallback mode is enabled, skip LLM backends
	if m.config.General.FallbackMode {
		return m.generateFallback(commandType, fullCommand, exitCode), BackendFallback
	}

	// Try backends in priority order: API -> Local -> Fallback

	// 1. Try API first (if available)
	if m.apiClient != nil && m.config.API.Enabled {
		if m.config.General.Debug {
			fmt.Printf("🔍 Trying API backend...\n")
		}

		// Create timeout context for API calls
		timeoutDuration := time.Duration(m.config.API.Timeout) * time.Second
		apiCtx, cancel := context.WithTimeout(ctx, timeoutDuration)
		defer cancel()

		response, err := m.apiClient.Generate(apiCtx, prompt)
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
		timeoutDuration := time.Duration(m.config.Local.Timeout) * time.Second
		localCtx, cancel := context.WithTimeout(ctx, timeoutDuration)
		defer cancel()

		response, err := m.ollamaClient.Generate(localCtx, prompt)
		if m.config.General.Debug {
			fmt.Printf("🐛 Raw Ollama response: '%s', error: %v\n", response, err)
		}
		if err == nil && response != "" {
			response = m.cleanResponse(response)
			if m.config.General.Debug {
				fmt.Printf("✅ Local backend succeeded with: '%s'\n", response)
			}
			return response, BackendLocal
		}

		if m.config.General.Debug {
			fmt.Printf("❌ Local backend failed: %v\n", err)
		}
	}

	// 3. Fallback to smart context-aware responses
	if m.config.General.Debug {
		fmt.Printf("🔄 Using smart fallback backend\n")
	}
	return m.generateFallback(commandType, fullCommand, exitCode), BackendFallback
}

func (m *LLMManager) cleanResponse(response string) string {
	// Clean up the response
	response = strings.TrimSpace(response)
	
	// Split at newlines and only keep the first meaningful part
	lines := strings.Split(response, "\n")
	if len(lines) > 1 {
		// Keep only the first line, discard any commentary after newlines
		response = strings.TrimSpace(lines[0])
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
	
	// Remove character count annotations like "(97 characters)"
	if idx := strings.Index(response, " ("); idx != -1 {
		remaining := response[idx:]
		if strings.Contains(remaining, "character") && strings.Contains(remaining, ")") {
			response = strings.TrimSpace(response[:idx])
		}
	}
	
	// Remove "Note:" annotations and similar commentary
	if idx := strings.Index(response, "Note:"); idx != -1 {
		response = strings.TrimSpace(response[:idx])
	}
	if idx := strings.Index(response, " *"); idx != -1 {
		// Remove asterisk annotations like "* This is a note"
		remaining := response[idx:]
		if strings.HasPrefix(strings.TrimSpace(remaining), "* ") {
			response = strings.TrimSpace(response[:idx])
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

func (m *LLMManager) generateFallback(commandType string, fullCommand string, exitCode string) string {
	// Use smart context-aware fallback when we have context
	if fullCommand != "" || exitCode != "" {
		ctx := ParseCommandContext(fullCommand, commandType, exitCode)
		insult := GenerateSmartFallback(ctx)
		if insult != "" {
			return insult
		}
	}

	// Fall back to expanded database if no smart match
	return GetExpandedFallback(commandType, fullCommand)
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