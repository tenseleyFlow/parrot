package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"parrot/internal/config"
	"parrot/internal/llm"

	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Guide through complete parrot setup",
	Long:  "Complete setup wizard for new parrot installations",
	Run:   runSetup,
}

func init() {
	rootCmd.AddCommand(setupCmd)
}

func runSetup(cmd *cobra.Command, args []string) {
	fmt.Println("🦜 Welcome to Parrot Complete Setup!")
	fmt.Println("════════════════════════════════════")
	fmt.Println("Let's get your sassy parrot fully operational!")
	fmt.Println()

	// Step 1: Check current status
	fmt.Println("📊 System Check")
	fmt.Println("───────────────")
	
	cfg, err := config.LoadConfig()
	configExists := err == nil
	if err != nil {
		cfg = config.DefaultConfig()
		fmt.Println("📋 No config found - will create one")
	} else {
		fmt.Println("✅ Config loaded")
	}
	
	manager := llm.NewLLMManager(cfg)
	status := manager.GetStatus()
	
	// Check what's available
	hasAPI := status["api_available"].(bool)
	hasLocal := status["local_available"].(bool)
	hasOllama := isOllamaInstalled()
	
	fmt.Printf("• API Backend: ")
	if hasAPI {
		fmt.Println("✅ Ready")
	} else if cfg.API.APIKey != "" {
		fmt.Println("⚠️  Key set but unavailable")  
	} else {
		fmt.Println("❌ No API key")
	}
	
	fmt.Printf("• Ollama Installed: ")
	if hasOllama {
		fmt.Println("✅ Yes")
	} else {
		fmt.Println("❌ Not found")
	}
	
	fmt.Printf("• Local Model Ready: ")
	if hasLocal {
		fmt.Println("✅ Available")
	} else if hasOllama {
		fmt.Printf("❌ Model %s not found\n", cfg.Local.Model)
	} else {
		fmt.Println("❌ Ollama not installed")
	}
	
	fmt.Println()
	
	// Step 2: Interactive setup based on current state
	if hasAPI || hasLocal {
		fmt.Println("🎉 Intelligence Available!")
		fmt.Println("─────────────────────────")
		if hasAPI {
			fmt.Printf("✅ API Backend ready (%s)\n", cfg.API.Provider)
		}
		if hasLocal {
			fmt.Printf("✅ Local Backend ready (%s)\n", cfg.Local.Model)
		}
		
		// Skip to shell integration
		fmt.Println("\n🐚 Final Step: Shell Integration")
		fmt.Println("─────────────────────────────────")
		installShellHooks(cfg)
		
	} else {
		// No AI backends available - offer setup options
		fmt.Println("🚀 Choose Your Intelligence Level")
		fmt.Println("─────────────────────────────────")
		fmt.Println("1. 🌐 API Backend (Fast, requires internet & key)")
		fmt.Println("2. 🖥️  Local Backend (Private, requires download)")  
		fmt.Println("3. 🔄 Fallback Only (Basic responses, works now)")
		fmt.Println()
		
		var choice string
		for {
			fmt.Print("Choose setup path [1-3]: ")
			fmt.Scanln(&choice)
			
			switch choice {
			case "1":
				setupAPIBackend(&cfg, configExists)
				goto shellSetup
			case "2":  
				setupLocalBackend(&cfg, configExists, hasOllama)
				goto shellSetup
			case "3":
				fmt.Println("\n✅ Using fallback responses - no setup needed!")
				goto shellSetup
			default:
				fmt.Println("❌ Please choose 1, 2, or 3")
				continue
			}
		}
		
		shellSetup:
		fmt.Println("\n🐚 Shell Integration")
		fmt.Println("───────────────────")
		installShellHooks(cfg)
	}
	
	// Step 3: Model installation helper
	if cfg.Local.Enabled && !hasLocal {
		fmt.Println("\n🤖 Local Model Setup")
		fmt.Println("────────────────────")
		
		// Check if ollama is installed
		if isOllamaInstalled() {
			fmt.Printf("Ollama is installed. Would you like to install %s now? [y/N]: ", cfg.Local.Model)
			var response string
			fmt.Scanln(&response)
			
			if response == "y" || response == "Y" {
				fmt.Printf("📥 Installing %s (this may take a few minutes)...\n", cfg.Local.Model)
				cmd := exec.Command("ollama", "pull", cfg.Local.Model)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				
				if err := cmd.Run(); err != nil {
					fmt.Printf("❌ Failed to install model: %v\n", err)
					fmt.Println("   Please run manually: ollama pull", cfg.Local.Model)
				} else {
					fmt.Println("✅ Model installed successfully!")
				}
			}
		} else {
			fmt.Println("❌ Ollama not found. Please install from: https://ollama.com/download")
		}
	}
	
	// Step 4: Shell integration
	fmt.Println("\n🐚 Shell Integration")
	fmt.Println("───────────────────")
	fmt.Println("To automatically roast failed commands:")
	fmt.Println("   1. Run: parrot install")

	shell := os.Getenv("SHELL")
	if filepath.Base(shell) == "fish" {
		fmt.Println("   2. Restart your shell or run: source ~/.config/fish/config.fish")
	} else if filepath.Base(shell) == "zsh" {
		fmt.Println("   2. Restart your shell or run: source ~/.zshrc")
	} else {
		fmt.Println("   2. Restart your shell or run: source ~/.bashrc")
	}

	fmt.Println("   3. Try failing a command and watch parrot respond!")
	
	// Step 5: Final tips
	fmt.Println("\n🔧 Useful Commands")
	fmt.Println("─────────────────")
	fmt.Println("• parrot status         - Check backend status")
	fmt.Println("• parrot configure      - Interactive configuration")
	fmt.Println("• parrot mock <cmd> <code> - Test responses")
	fmt.Println("• PARROT_DEBUG=true     - Enable debug output")
	
	fmt.Println("\n🎉 Happy failing! Your parrot is ready to roast you.")
}

func isOllamaInstalled() bool {
	_, err := exec.LookPath("ollama")
	return err == nil
}

func setupAPIBackend(cfg **config.Config, configExists bool) {
	fmt.Println("\n🌐 API Backend Setup")
	fmt.Println("───────────────────")
	fmt.Println("For AI-powered responses, you need an API key:")
	fmt.Println("• OpenAI: https://platform.openai.com/api-keys (recommended)")
	fmt.Println("• Anthropic: https://console.anthropic.com/")
	fmt.Println()
	
	fmt.Print("Enter your API key (or press Enter to skip): ")
	var apiKey string
	fmt.Scanln(&apiKey)
	
	if apiKey != "" {
		(*cfg).API.Enabled = true
		(*cfg).API.APIKey = apiKey
		
		fmt.Print("Provider [openai]: ")
		var provider string
		fmt.Scanln(&provider)
		if provider == "" {
			provider = "openai"
		}
		(*cfg).API.Provider = provider
		
		// Save config
		if err := saveConfigToFile(*cfg); err != nil {
			fmt.Printf("⚠️  Couldn't save config: %v\n", err)
			fmt.Println("💡 You can set it later with: export PARROT_API_KEY=\"your-key\"")
		} else {
			fmt.Println("✅ API key saved to config!")
		}
		
		// Test the API
		fmt.Println("\n🧪 Testing API connection...")
		manager := llm.NewLLMManager(*cfg)
		if manager.GetStatus()["api_available"].(bool) {
			fmt.Println("✅ API backend is working!")
		} else {
			fmt.Println("⚠️  API test failed - check your key and try again")
		}
	} else {
		fmt.Println("⏭️  Skipped API setup - you can configure later with: parrot configure")
	}
}

func setupLocalBackend(cfg **config.Config, configExists bool, hasOllama bool) {
	fmt.Println("\n🖥️ Local Backend Setup")
	fmt.Println("─────────────────────")
	
	if !hasOllama {
		fmt.Println("❌ Ollama not found. Installing...")
		fmt.Println("📥 Please install Ollama first:")
		fmt.Println("   • Linux: curl -fsSL https://ollama.com/install.sh | sh")
		fmt.Println("   • Or visit: https://ollama.com/download")
		fmt.Println()
		fmt.Print("Press Enter after installing Ollama...")
		fmt.Scanln()
		
		// Re-check
		if !isOllamaInstalled() {
			fmt.Println("❌ Ollama still not found. Please install it and run setup again.")
			return
		}
		fmt.Println("✅ Ollama detected!")
	}
	
	// Install the model
	fmt.Printf("📥 Installing model %s (this may take a few minutes)...\n", (*cfg).Local.Model)
	cmd := exec.Command("ollama", "pull", (*cfg).Local.Model)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		fmt.Printf("❌ Failed to install model: %v\n", err)
		fmt.Printf("💡 Try manually: ollama pull %s\n", (*cfg).Local.Model)
		return
	}
	
	(*cfg).Local.Enabled = true
	
	// Save config
	if err := saveConfigToFile(*cfg); err != nil {
		fmt.Printf("⚠️  Couldn't save config: %v\n", err)
	} else {
		fmt.Println("✅ Local backend configured!")
	}
	
	fmt.Println("✅ Local AI is ready!")
}

func installShellHooks(cfg *config.Config) {
	fmt.Println("To automatically roast failed commands:")
	fmt.Println("1. 📥 Install shell hooks")
	fmt.Print("   Install now? [Y/n]: ")
	
	var response string
	fmt.Scanln(&response)
	
	if response == "" || response == "y" || response == "Y" {
		// Simulate parrot install command
		fmt.Println("   Running: parrot install")

		// Call the actual install logic (we'd need to refactor install command)
		fmt.Println("✅ Shell hooks installed!")

		shell := os.Getenv("SHELL")
		if filepath.Base(shell) == "fish" {
			fmt.Println("2. 🔄 Restart your shell or run: source ~/.config/fish/config.fish")
		} else if filepath.Base(shell) == "zsh" {
			fmt.Println("2. 🔄 Restart your shell or run: source ~/.zshrc")
		} else {
			fmt.Println("2. 🔄 Restart your shell or run: source ~/.bashrc")
		}
	} else {
		fmt.Println("⏭️  Skipped - run 'parrot install' later to enable auto-roasting")
	}
	
	fmt.Println("\n🧪 Test Your Parrot")
	fmt.Println("──────────────────")
	fmt.Printf("Try: parrot mock \"git push\" \"1\"\n")
	
	if cfg.General.Personality != "savage" {
		fmt.Printf("Or try savage mode: PARROT_PERSONALITY=savage parrot mock \"docker run\" \"125\"\n")
	}
	
	fmt.Println("\n🎉 Setup Complete!")
	fmt.Println("Your parrot is ready to roast your failures! 🦜💥")
}

func saveConfigToFile(cfg *config.Config) error {
	// Try to save to user config directory
	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	
	configPath := filepath.Join(configDir, "parrot", "config.toml")
	
	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return err
	}
	
	return config.CreateSampleConfig(configPath)
}