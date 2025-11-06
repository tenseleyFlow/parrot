package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"parrot/internal/colors"
	"parrot/internal/config"
	"parrot/internal/llm"
	"parrot/internal/prompts"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "parrot",
	Version: "1.3.0",
	Short:   "A sassy CLI that mocks your failed commands",
	Long:    "Parrot listens for failed commands and responds with intelligent insults and mockery.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🦜 Parrot is watching... waiting for you to mess up!")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

var mockCmd = &cobra.Command{
	Use:   "mock [command] [exit_code]",
	Short: "Mock a failed command",
	Long:  "Called by shell hooks when a command fails",
	Args:  cobra.MinimumNArgs(2),
	Run:   mockCommand,
}

func init() {
	rootCmd.AddCommand(mockCmd)
}

func mockCommand(cmd *cobra.Command, args []string) {
	failedCmd := args[0]
	exitCode := args[1]
	
	// Basic command type detection
	cmdType := detectCommandType(failedCmd)
	
	// Show immediate feedback to user
	fmt.Print("🦜 ")
	
	// Generate a smart mock response
	response, cfg := generateSmartResponse(cmdType, failedCmd, exitCode)
	
	// Clear the loading indicator and show response
	fmt.Print("\r") // Clear current line
	
	// Format output with colors and personality
	if cfg.General.Colors {
		fmt.Println(colors.FormatParrotOutput(cfg.General.Personality, response, cfg.General.Enhanced))
	} else {
		fmt.Printf("🦜 %s\n", response)
	}
}

func detectCommandType(command string) string {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return "generic"
	}

	cmd := parts[0]

	// Check for common command patterns
	switch cmd {
	// Version control
	case "git":
		return "git"

	// Node.js ecosystem
	case "npm", "yarn", "pnpm", "node", "npx", "bun", "deno":
		return "nodejs"

	// Containers
	case "docker", "docker-compose", "podman", "kubectl", "k9s":
		return "docker"

	// HTTP/Network
	case "curl", "wget", "http", "https":
		return "http"

	// SSH/Remote
	case "ssh", "scp", "sftp", "rsync":
		return "ssh"

	// Navigation
	case "cd", "pushd", "popd":
		return "navigation"

	// Python
	case "python", "python3", "pip", "pip3", "poetry", "pipenv", "conda":
		return "python"

	// Rust
	case "cargo", "rustc", "rustup":
		return "rust"

	// Go
	case "go":
		return "go"

	// Java
	case "java", "javac", "mvn", "gradle":
		return "java"

	// C/C++
	case "gcc", "g++", "clang", "clang++", "cc", "c++":
		return "cpp"

	// Ruby
	case "ruby", "gem", "bundle", "rake", "rails":
		return "ruby"

	// PHP
	case "php", "composer":
		return "php"

	// Build systems (check for C files to categorize properly)
	case "make", "cmake", "ninja", "ant", "bazel":
		if strings.Contains(command, ".c ") || strings.HasSuffix(command, ".c") {
			return "c"
		}
		return "build"

	// Databases
	case "mysql", "psql", "postgres", "mongo", "mongosh", "redis-cli", "sqlite3":
		return "database"

	// Permission-related commands
	case "chmod", "chown", "chgrp", "sudo":
		return "permissions"

	default:
		return "generic"
	}
}

func generateSmartResponse(cmdType, command, exitCode string) (string, *config.Config) {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		// If config loading fails, use fallback with default config
		defaultCfg := config.DefaultConfig()
		return getFallbackResponse(cmdType), defaultCfg
	}
	
	// Initialize LLM manager
	manager := llm.NewLLMManager(cfg)
	
	// Build context-aware prompt with personality
	prompt := prompts.BuildPrompt(cmdType, command, exitCode, cfg.General.Personality)
	
	// Use a reasonable timeout for LLM responses (10 seconds max)
	maxTimeout := 10 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), maxTimeout)
	defer cancel()
	
	// Create a channel for the response
	responseChan := make(chan struct {
		response string
		backend  llm.Backend
	}, 1)
	
	// Start generation in a goroutine
	go func() {
		// Use GenerateWithContext for intelligent fallbacks
		response, backend := manager.GenerateWithContext(ctx, prompt, cmdType, command, exitCode)
		select {
		case responseChan <- struct {
			response string
			backend  llm.Backend
		}{response, backend}:
		case <-ctx.Done():
		}
	}()
	
	// Show progress indicator for anything longer than 500ms
	progressTimer := time.NewTimer(500 * time.Millisecond)
	defer progressTimer.Stop()
	
	select {
	case result := <-responseChan:
		progressTimer.Stop()
		// Add backend indicator in debug mode
		if cfg.General.Debug {
			switch result.backend {
			case llm.BackendAPI:
				fmt.Printf("🌐 API backend used\n")
			case llm.BackendLocal:
				fmt.Printf("🖥️ Local backend used\n")
			case llm.BackendFallback:
				fmt.Printf("🔄 Fallback backend used\n")
			}
		}
		return result.response, cfg
	case <-progressTimer.C:
		// Show thinking indicator after 500ms
		fmt.Print("💭")
		select {
		case result := <-responseChan:
			// Add backend indicator in debug mode
			if cfg.General.Debug {
				switch result.backend {
				case llm.BackendAPI:
					fmt.Printf("\n🌐 API backend used\n")
				case llm.BackendLocal:
					fmt.Printf("\n🖥️ Local backend used\n")
				case llm.BackendFallback:
					fmt.Printf("\n🔄 Fallback backend used\n")
				}
			}
			return result.response, cfg
		case <-ctx.Done():
			// Fallback to instant response if timeout reached
			return getFallbackResponse(cmdType), cfg
		}
	case <-ctx.Done():
		// Fallback to instant response if timeout reached
		return getFallbackResponse(cmdType), cfg
	}
}

func getFallbackResponse(cmdType string) string {
	// Use the expanded fallback database with hundreds of brutal insults
	// This is only called on config load failure, so we don't have full context
	return llm.GetExpandedFallback(cmdType, "")
}