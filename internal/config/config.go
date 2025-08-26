package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	// API Configuration (Primary backend)
	API APIConfig `toml:"api"`
	
	// Local LLM Configuration (Secondary backend)
	Local LocalConfig `toml:"local"`
	
	// General Settings
	General GeneralConfig `toml:"general"`
}

type APIConfig struct {
	Enabled  bool   `toml:"enabled"`
	Provider string `toml:"provider"` // "openai", "anthropic", "custom"
	Endpoint string `toml:"endpoint"` // Custom endpoint URL
	APIKey   string `toml:"api_key"`  // API key
	Model    string `toml:"model"`    // Model name
	Timeout  int    `toml:"timeout"`  // Request timeout in seconds
}

type LocalConfig struct {
	Enabled  bool   `toml:"enabled"`
	Provider string `toml:"provider"` // "ollama"
	Endpoint string `toml:"endpoint"` // Ollama endpoint
	Model    string `toml:"model"`    // Model name
	Timeout  int    `toml:"timeout"`  // Request timeout in seconds
}

type GeneralConfig struct {
	Personality  string `toml:"personality"`   // "savage", "sarcastic", "mild"
	FallbackMode bool   `toml:"fallback_mode"` // Use hardcoded responses only
	Debug        bool   `toml:"debug"`         // Debug logging
	Colors       bool   `toml:"colors"`        // Enable colored output
	Enhanced     bool   `toml:"enhanced"`      // Enhanced formatting with borders/emphasis
}

// Default configuration
func DefaultConfig() *Config {
	return &Config{
		API: APIConfig{
			Enabled:  true,
			Provider: "openai",
			Endpoint: "https://api.openai.com/v1",
			APIKey:   "", // Must be set by user
			Model:    "gpt-3.5-turbo",
			Timeout:  3,  // Reduced from 10 to 3 seconds for responsiveness
		},
		Local: LocalConfig{
			Enabled:  true,
			Provider: "ollama", 
			Endpoint: "http://localhost:11434",
			Model:    "phi3.5:3.8b",
			Timeout:  5,  // Reduced from 30 to 5 seconds for responsiveness
		},
		General: GeneralConfig{
			Personality:  "savage",
			FallbackMode: false,
			Debug:        false,
			Colors:       true,
			Enhanced:     false,
		},
	}
}

// Config file paths in order of preference
func GetConfigPaths() []string {
	var paths []string
	
	// 0. Environment-specified config path (highest priority)
	if configPath := os.Getenv("PARROT_CONFIG"); configPath != "" {
		paths = append(paths, configPath)
	}
	
	// 1. System-wide config (for RPM installs)
	paths = append(paths, "/etc/parrot/config.toml")
	
	// 2. User config directory
	if configDir, err := os.UserConfigDir(); err == nil {
		paths = append(paths, filepath.Join(configDir, "parrot", "config.toml"))
	}
	
	// 3. Home directory
	if homeDir, err := os.UserHomeDir(); err == nil {
		paths = append(paths, filepath.Join(homeDir, ".parrot.toml"))
	}
	
	// 4. Current directory (for development)
	paths = append(paths, "./parrot.toml")
	
	return paths
}

// Load configuration from first available config file
func LoadConfig() (*Config, error) {
	config := DefaultConfig()
	
	// Try to load from config files
	for _, path := range GetConfigPaths() {
		if _, err := os.Stat(path); err == nil {
			if err := loadFromFile(config, path); err != nil {
				return nil, fmt.Errorf("error loading config from %s: %w", path, err)
			}
			break
		}
	}
	
	// Override with environment variables
	loadFromEnv(config)
	
	return config, nil
}

func loadFromFile(config *Config, path string) error {
	_, err := toml.DecodeFile(path, config)
	return err
}

func loadFromEnv(config *Config) {
	// API configuration from environment
	if key := os.Getenv("PARROT_API_KEY"); key != "" {
		config.API.APIKey = key
	}
	if endpoint := os.Getenv("PARROT_API_ENDPOINT"); endpoint != "" {
		config.API.Endpoint = endpoint
	}
	if model := os.Getenv("PARROT_API_MODEL"); model != "" {
		config.API.Model = model
	}
	
	// Local configuration from environment
	if endpoint := os.Getenv("PARROT_OLLAMA_ENDPOINT"); endpoint != "" {
		config.Local.Endpoint = endpoint
	}
	if model := os.Getenv("PARROT_OLLAMA_MODEL"); model != "" {
		config.Local.Model = model
	}
	
	// General configuration
	if personality := os.Getenv("PARROT_PERSONALITY"); personality != "" {
		config.General.Personality = personality
	}
	if os.Getenv("PARROT_FALLBACK_ONLY") == "true" {
		config.General.FallbackMode = true
	}
	if os.Getenv("PARROT_DEBUG") == "true" {
		config.General.Debug = true
	}
	if os.Getenv("PARROT_NO_COLOR") == "true" || os.Getenv("NO_COLOR") != "" {
		config.General.Colors = false
	}
	if os.Getenv("PARROT_ENHANCED") == "true" {
		config.General.Enhanced = true
	}
}

// Create a sample config file
func CreateSampleConfig(path string) error {
	config := DefaultConfig()
	
	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}
	
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer file.Close()
	
	encoder := toml.NewEncoder(file)
	if err := encoder.Encode(config); err != nil {
		return fmt.Errorf("failed to encode config: %w", err)
	}
	
	return nil
}