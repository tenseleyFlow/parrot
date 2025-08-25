package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"parrot/internal/config"

	"github.com/spf13/cobra"
)

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Interactively configure parrot backends and preferences",
	Long:  "Walk through interactive setup to configure API keys, models, and preferences",
	Run:   runConfigure,
}

func init() {
	rootCmd.AddCommand(configureCmd)
}

func runConfigure(cmd *cobra.Command, args []string) {
	fmt.Println("🦜 Parrot Configuration Wizard")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)
	
	// Load existing config or defaults
	cfg, err := config.LoadConfig()
	if err != nil {
		cfg = config.DefaultConfig()
	}

	// 1. Choose config location
	configPath := chooseConfigLocation(reader)
	
	// 2. Configure API backend
	fmt.Println("🌐 API Backend Configuration")
	fmt.Println("────────────────────────────")
	cfg.API.Enabled = askYesNo(reader, "Enable API backend? (recommended)", cfg.API.Enabled)
	
	if cfg.API.Enabled {
		cfg.API.Provider = askChoice(reader, "API Provider", []string{"openai", "anthropic", "custom"}, cfg.API.Provider)
		
		if cfg.API.Provider == "custom" {
			cfg.API.Endpoint = askString(reader, "API Endpoint URL", cfg.API.Endpoint)
		} else {
			// Set default endpoints
			switch cfg.API.Provider {
			case "openai":
				cfg.API.Endpoint = "https://api.openai.com/v1"
			case "anthropic":
				cfg.API.Endpoint = "https://api.anthropic.com/v1"
			}
		}
		
		cfg.API.APIKey = askString(reader, "API Key", cfg.API.APIKey)
		cfg.API.Model = askString(reader, "Model name", cfg.API.Model)
	}
	
	// 3. Configure Local backend
	fmt.Println("\n🖥️  Local Backend Configuration")
	fmt.Println("────────────────────────────")
	cfg.Local.Enabled = askYesNo(reader, "Enable local Ollama backend?", cfg.Local.Enabled)
	
	if cfg.Local.Enabled {
		cfg.Local.Endpoint = askString(reader, "Ollama endpoint", cfg.Local.Endpoint)
		availableModels := []string{"phi3.5:3.8b", "llama3.2:3b", "qwen2.5:0.5b", "custom"}
		selectedModel := askChoice(reader, "Local model", availableModels, cfg.Local.Model)
		
		if selectedModel == "custom" {
			cfg.Local.Model = askString(reader, "Custom model name", cfg.Local.Model)
		} else {
			cfg.Local.Model = selectedModel
		}
	}
	
	// 4. General preferences
	fmt.Println("\n⚙️  General Preferences")
	fmt.Println("────────────────────────")
	personalities := []string{"mild", "sarcastic", "savage"}
	cfg.General.Personality = askChoice(reader, "Personality level", personalities, cfg.General.Personality)
	cfg.General.Debug = askYesNo(reader, "Enable debug mode?", cfg.General.Debug)
	cfg.General.FallbackMode = askYesNo(reader, "Use only fallback responses? (disable AI)", cfg.General.FallbackMode)
	
	// 5. Save configuration
	fmt.Println("\n💾 Saving Configuration...")
	if err := config.CreateSampleConfig(configPath); err != nil {
		fmt.Printf("❌ Error creating config template: %v\n", err)
		return
	}
	
	// Load the template and update with user values
	if err := saveConfig(cfg, configPath); err != nil {
		fmt.Printf("❌ Error saving configuration: %v\n", err)
		return
	}
	
	fmt.Printf("✅ Configuration saved to: %s\n", configPath)
	
	// 6. Next steps
	fmt.Println("\n🎯 Next Steps:")
	if cfg.API.Enabled && cfg.API.APIKey != "" {
		fmt.Println("   • Test API backend: parrot status")
	}
	if cfg.Local.Enabled {
		fmt.Printf("   • Ensure model is available: ollama pull %s\n", cfg.Local.Model)
	}
	fmt.Println("   • Test parrot: parrot mock \"git push\" \"1\"")
	fmt.Println("   • Install shell hooks: parrot install")
}

func chooseConfigLocation(reader *bufio.Reader) string {
	fmt.Println("📁 Configuration Location")
	fmt.Println("─────────────────────────")
	
	paths := config.GetConfigPaths()
	fmt.Println("Choose where to save your configuration:")
	for i, path := range paths {
		fmt.Printf("%d. %s\n", i+1, path)
	}
	
	for {
		fmt.Print("Choice [1]: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		
		if input == "" {
			return paths[0] // Default to first option
		}
		
		choice := 0
		if _, err := fmt.Sscanf(input, "%d", &choice); err == nil && choice >= 1 && choice <= len(paths) {
			return paths[choice-1]
		}
		
		fmt.Println("❌ Invalid choice. Please enter a number between 1 and", len(paths))
	}
}

func askString(reader *bufio.Reader, prompt, defaultValue string) string {
	if defaultValue != "" {
		fmt.Printf("%s [%s]: ", prompt, defaultValue)
	} else {
		fmt.Printf("%s: ", prompt)
	}
	
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	
	if input == "" {
		return defaultValue
	}
	return input
}

func askYesNo(reader *bufio.Reader, prompt string, defaultValue bool) bool {
	defaultStr := "n"
	if defaultValue {
		defaultStr = "y"
	}
	
	fmt.Printf("%s [%s]: ", prompt, defaultStr)
	input, _ := reader.ReadString('\n')
	input = strings.ToLower(strings.TrimSpace(input))
	
	if input == "" {
		return defaultValue
	}
	
	return input == "y" || input == "yes"
}

func askChoice(reader *bufio.Reader, prompt string, choices []string, defaultValue string) string {
	fmt.Printf("%s:\n", prompt)
	for i, choice := range choices {
		marker := " "
		if choice == defaultValue {
			marker = "*"
		}
		fmt.Printf("%s %d. %s\n", marker, i+1, choice)
	}
	
	for {
		fmt.Print("Choice: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		
		if input == "" && defaultValue != "" {
			return defaultValue
		}
		
		choice := 0
		if _, err := fmt.Sscanf(input, "%d", &choice); err == nil && choice >= 1 && choice <= len(choices) {
			return choices[choice-1]
		}
		
		// Allow text input too
		for _, validChoice := range choices {
			if strings.EqualFold(input, validChoice) {
				return validChoice
			}
		}
		
		fmt.Printf("❌ Invalid choice. Please enter 1-%d or the option name.\n", len(choices))
	}
}

func saveConfig(cfg *config.Config, path string) error {
	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	
	content := generateConfigContent(cfg)
	return os.WriteFile(path, []byte(content), 0644)
}

func generateConfigContent(cfg *config.Config) string {
	var content strings.Builder
	
	content.WriteString("# Parrot Configuration File\n")
	content.WriteString("# Generated by: parrot configure\n\n")
	
	// General section
	content.WriteString("[general]\n")
	content.WriteString(fmt.Sprintf("personality = \"%s\"\n", cfg.General.Personality))
	content.WriteString(fmt.Sprintf("fallback_mode = %t\n", cfg.General.FallbackMode))
	content.WriteString(fmt.Sprintf("debug = %t\n", cfg.General.Debug))
	content.WriteString("\n")
	
	// API section
	content.WriteString("[api]\n")
	content.WriteString(fmt.Sprintf("enabled = %t\n", cfg.API.Enabled))
	content.WriteString(fmt.Sprintf("provider = \"%s\"\n", cfg.API.Provider))
	content.WriteString(fmt.Sprintf("endpoint = \"%s\"\n", cfg.API.Endpoint))
	if cfg.API.APIKey != "" {
		content.WriteString(fmt.Sprintf("api_key = \"%s\"\n", cfg.API.APIKey))
	} else {
		content.WriteString("# api_key = \"your-api-key-here\"\n")
	}
	content.WriteString(fmt.Sprintf("model = \"%s\"\n", cfg.API.Model))
	content.WriteString(fmt.Sprintf("timeout = %d\n", cfg.API.Timeout))
	content.WriteString("\n")
	
	// Local section
	content.WriteString("[local]\n")
	content.WriteString(fmt.Sprintf("enabled = %t\n", cfg.Local.Enabled))
	content.WriteString(fmt.Sprintf("provider = \"%s\"\n", cfg.Local.Provider))
	content.WriteString(fmt.Sprintf("endpoint = \"%s\"\n", cfg.Local.Endpoint))
	content.WriteString(fmt.Sprintf("model = \"%s\"\n", cfg.Local.Model))
	content.WriteString(fmt.Sprintf("timeout = %d\n", cfg.Local.Timeout))
	
	return content.String()
}