package cmd

import (
	"fmt"
	"math/rand"
	"os"
	"strings"

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
	
	// Generate a mock response
	response := generateMockResponse(cmdType, failedCmd, exitCode)
	
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

func generateMockResponse(cmdType, command, exitCode string) string {
	responses := map[string][]string{
		"git": {
			"Oh look, another git genius who can't even commit properly!",
			"Did you forget to pull again? Classic amateur move.",
			"Git good? More like git rekt!",
		},
		"nodejs": {
			"Node modules strike again! Maybe try turning it off and on again?",
			"NPM install failed? Shocking! Nobody saw that coming.",
			"Your package.json is crying. Fix it.",
		},
		"docker": {
			"Docker container more like docker DISASTER!",
			"Even containers can't contain your incompetence.",
			"Did you try 'have you tried containerizing it differently'?",
		},
		"http": {
			"404: Competence not found.",
			"Looks like the internet doesn't want to talk to you.",
			"Even HTTP requests are rejecting you now.",
		},
		"generic": {
			"Wow, you managed to break something simple. Impressive!",
			"Error code? More like user error!",
			"Maybe try reading the manual... oh wait, who am I kidding?",
		},
	}
	
	cmdResponses, exists := responses[cmdType]
	if !exists {
		cmdResponses = responses["generic"]
	}
	
	response := cmdResponses[rand.Intn(len(cmdResponses))]
	return response
}