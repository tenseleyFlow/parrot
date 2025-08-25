package cmd

import (
	"fmt"

	"parrot/internal/colors"

	"github.com/spf13/cobra"
)

var demoCmd = &cobra.Command{
	Use:   "demo",
	Short: "Show parrot personality and color demos",
	Long:  "Display examples of different personality levels and color schemes",
	Run:   runDemo,
}

func init() {
	rootCmd.AddCommand(demoCmd)
}

func runDemo(cmd *cobra.Command, args []string) {
	fmt.Println("🦜 Parrot Personality & Color Demo")
	fmt.Println("═══════════════════════════════════")
	fmt.Println()

	// Sample commands for demo
	commands := []struct {
		cmd      string
		cmdType  string
		exitCode string
	}{
		{"git push origin main", "git", "1"},
		{"npm install express", "nodejs", "1"},
		{"docker run myapp", "docker", "125"},
		{"curl https://api.example.com", "http", "7"},
	}

	personalities := []string{"mild", "sarcastic", "savage"}
	
	for _, personality := range personalities {
		fmt.Printf("🎭 %s Personality\n", personality)
		fmt.Println("─────────────────────")
		
		for _, test := range commands {
			// Use simple hardcoded responses for demo
			var response string
			switch personality {
			case "mild":
				response = getDemoResponse(test.cmdType, "mild")
			case "sarcastic":
				response = getDemoResponse(test.cmdType, "sarcastic")  
			case "savage":
				response = getDemoResponse(test.cmdType, "savage")
			}
			
			fmt.Printf("Command: %s\n", test.cmd)
			fmt.Printf("  Simple:   %s\n", colors.FormatParrotOutput(personality, response, false))
			fmt.Printf("  Enhanced: %s\n", colors.FormatParrotOutput(personality, response, true))
			fmt.Println()
		}
		
		fmt.Println()
	}
	
	// Color info
	fmt.Printf("🎨 Colors enabled: %t\n", colors.ColorEnabled())
	if !colors.ColorEnabled() {
		fmt.Println("   💡 To enable colors, ensure you're in a terminal and NO_COLOR is not set")
	}
}

func getDemoResponse(cmdType, personality string) string {
	responses := map[string]map[string]string{
		"mild": {
			"git":     "Git command failed. Maybe check your remote branch?",
			"nodejs":  "NPM seems unhappy. Try clearing your cache?",
			"docker":  "Container seems upset. Check your Dockerfile?",
			"http":    "Request didn't go through. Check the URL?",
			"generic": "Command didn't work as expected. Check the syntax?",
		},
		"sarcastic": {
			"git":     "Another git genius who forgot to pull first. Classic.",
			"nodejs":  "NPM install failed? Shocking! Nobody saw that coming.",
			"docker":  "Docker container more like docker DISASTER!",
			"http":    "404: Competence not found.",
			"generic": "Wow, you managed to break something simple. Impressive!",
		},
		"savage": {
			"git":     "Git rejected your code harder than everyone rejects you.",
			"nodejs":  "NPM refuses to install anything for someone this incompetent.",
			"docker":  "Your containers crash faster than your career prospects.",
			"http":    "The internet collectively rejected you. Impressive.",
			"generic": "Your command failed harder than you failed at life.",
		},
	}
	
	if personalityMap, exists := responses[personality]; exists {
		if response, exists := personalityMap[cmdType]; exists {
			return response
		}
		return personalityMap["generic"]
	}
	return "Error: Unknown personality"
}