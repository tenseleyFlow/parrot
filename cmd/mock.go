package cmd

import (
	"context"
	"fmt"
	"math/rand"
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
	
	// Use a shorter overall timeout for shell responsiveness (max 2 seconds)
	maxTimeout := 2 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), maxTimeout)
	defer cancel()
	
	// Create a channel for the response
	responseChan := make(chan struct {
		response string
		backend  llm.Backend
	}, 1)
	
	// Start generation in a goroutine
	go func() {
		response, backend := manager.Generate(ctx, prompt, cmdType)
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