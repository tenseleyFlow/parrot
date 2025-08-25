package cmd

import (
	"fmt"
	"os"
	"os/exec"

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
	fmt.Println("🦜 Welcome to Parrot Setup!")
	fmt.Println("═══════════════════════════")
	fmt.Println("This wizard will guide you through setting up your sassy parrot.")
	fmt.Println()

	// Step 1: Check current status
	fmt.Println("📊 Current Status")
	fmt.Println("─────────────────")
	
	cfg, err := config.LoadConfig()
	if err != nil {
		cfg = config.DefaultConfig()
		fmt.Println("⚠️  No configuration found - using defaults")
	} else {
		fmt.Println("✅ Configuration loaded")
	}
	
	manager := llm.NewLLMManager(cfg)
	status := manager.GetStatus()
	
	// Check what's available
	hasAPI := status["api_available"].(bool)
	hasLocal := status["local_available"].(bool)
	
	fmt.Printf("• API Backend: ")
	if hasAPI {
		fmt.Println("✅ Available")
	} else {
		fmt.Println("❌ Not configured")
	}
	
	fmt.Printf("• Local Backend: ")
	if hasLocal {
		fmt.Println("✅ Available")
	} else {
		fmt.Println("❌ Not available")
	}
	
	// Step 2: Recommend setup path
	fmt.Println("\n🎯 Recommended Setup")
	fmt.Println("────────────────────")
	
	if !hasAPI && !hasLocal {
		fmt.Println("🚀 Quick Start Option 1: API-based (fastest setup)")
		fmt.Println("   • Get an OpenAI API key (https://platform.openai.com/api-keys)")
		fmt.Println("   • Run: parrot configure")
		fmt.Println("   • Enable API backend and enter your key")
		fmt.Println()
		
		fmt.Println("🖥️  Quick Start Option 2: Local AI (privacy-focused)")
		fmt.Println("   • Install Ollama: https://ollama.com/download")
		fmt.Println("   • Pull model: ollama pull phi3.5:3.8b")
		fmt.Println("   • Run: parrot configure")
		fmt.Println()
		
		fmt.Println("🔄 Quick Start Option 3: Hardcoded responses (no setup)")
		fmt.Println("   • Already working! Just install shell hooks.")
		fmt.Println("   • Run: parrot install")
		
	} else if hasAPI {
		fmt.Println("✅ You're all set with API backend!")
		fmt.Println("   • Install shell hooks: parrot install")
		fmt.Println("   • Test it: parrot mock \"git push\" \"1\"")
		
	} else if hasLocal {
		fmt.Println("✅ You're all set with local backend!")
		fmt.Println("   • Install shell hooks: parrot install")
		fmt.Println("   • Test it: parrot mock \"git push\" \"1\"")
		
		if !status["local_enabled"].(bool) || cfg.Local.Model != "phi3.5:3.8b" {
			fmt.Println("\n💡 Tip: For better quality responses:")
			fmt.Println("   • Upgrade to phi3.5: ollama pull phi3.5:3.8b")
			fmt.Println("   • Update config: parrot configure")
		}
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
	fmt.Println("   2. Restart your shell or run: source ~/.bashrc")
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