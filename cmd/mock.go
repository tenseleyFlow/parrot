package cmd

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"parrot/internal/config"
	"parrot/internal/llm"
	"parrot/internal/prompts"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "parrot",
	Short: "A sassy CLI that mocks your failed commands",
	Long:  "Parrot listens for failed commands and responds with intelligent insults and mockery.",
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
	
	// Generate a smart mock response
	response := generateSmartResponse(cmdType, failedCmd, exitCode)
	
	fmt.Printf("🦜 %s\n", response)
}

func detectCommandType(command string) string {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return "unknown"
	}
	
	switch parts[0] {
	case "git":
		return "git"
	case "npm", "yarn", "pnpm":
		return "nodejs"
	case "docker", "docker-compose":
		return "docker"
	case "curl", "wget":
		return "http"
	case "ssh":
		return "ssh"
	case "cd":
		return "navigation"
	default:
		return "generic"
	}
}

func generateSmartResponse(cmdType, command, exitCode string) string {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		// If config loading fails, use fallback
		return getFallbackResponse(cmdType)
	}
	
	// Initialize LLM manager
	manager := llm.NewLLMManager(cfg)
	
	// Build context-aware prompt
	prompt := prompts.BuildPrompt(cmdType, command, exitCode)
	
	// Generate response with timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.API.Timeout)*time.Second)
	defer cancel()
	
	response, backend := manager.Generate(ctx, prompt, cmdType)
	
	// Add backend indicator in debug mode
	if cfg.General.Debug {
		switch backend {
		case llm.BackendAPI:
			fmt.Printf("🌐 API backend used\n")
		case llm.BackendLocal:
			fmt.Printf("🖥️ Local backend used\n")
		case llm.BackendFallback:
			fmt.Printf("🔄 Fallback backend used\n")
		}
	}
	
	return response
}

func getFallbackResponse(cmdType string) string {
	fallbacks := map[string][]string{
		"git": {
			"Git good? More like git rekt!",
			"Did you forget to pull again? Classic amateur move.",
			"Another git genius strikes again!",
		},
		"nodejs": {
			"NPM install failed? Shocking! Nobody saw that coming.",
			"Your package.json is crying. Fix it.",
			"Node modules: where dependencies go to die.",
		},
		"docker": {
			"Docker container more like docker DISASTER!",
			"Even containers can't contain your incompetence.",
			"Your Dockerfile needs therapy.",
		},
		"http": {
			"404: Competence not found.",
			"Even the internet doesn't want to talk to you.",
			"Connection refused? So is your logic.",
		},
		"generic": {
			"Wow, you managed to break something simple. Impressive!",
			"Maybe try reading the manual... oh wait, who am I kidding?",
			"Error code says it all: user error!",
		},
	}
	
	responses, exists := fallbacks[cmdType]
	if !exists {
		responses = fallbacks["generic"]
	}
	
	return responses[rand.Intn(len(responses))]
}