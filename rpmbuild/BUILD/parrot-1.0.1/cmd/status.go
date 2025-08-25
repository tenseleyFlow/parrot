package cmd

import (
	"fmt"
	"os"

	"parrot/internal/config"
	"parrot/internal/llm"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show parrot configuration and backend status",
	Long:  "Display current configuration, available backends, and their status",
	Run:   showStatus,
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func showStatus(cmd *cobra.Command, args []string) {
	fmt.Println("🦜 Parrot Status Report")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━")
	
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("❌ Configuration Error: %v\n", err)
		return
	}
	
	// Show configuration source
	fmt.Println("\n📁 Configuration:")
	configFound := false
	for _, path := range config.GetConfigPaths() {
		if _, err := os.Stat(path); err == nil {
			fmt.Printf("   ✅ Loaded from: %s\n", path)
			configFound = true
			break
		}
	}
	if !configFound {
		fmt.Println("   ℹ️  Using default configuration (no config file found)")
	}
	
	// General settings
	fmt.Printf("   • Personality: %s\n", cfg.General.Personality)
	fmt.Printf("   • Debug mode: %t\n", cfg.General.Debug)
	fmt.Printf("   • Fallback only: %t\n", cfg.General.FallbackMode)
	
	// Initialize LLM manager to get status
	manager := llm.NewLLMManager(cfg)
	status := manager.GetStatus()
	
	// API Backend Status
	fmt.Println("\n🌐 API Backend:")
	if status["api_enabled"].(bool) {
		fmt.Printf("   • Enabled: ✅\n")
		fmt.Printf("   • Provider: %s\n", status["api_provider"])
		fmt.Printf("   • Model: %s\n", status["api_model"])
		if status["api_available"].(bool) {
			fmt.Printf("   • Status: ✅ Available\n")
		} else {
			fmt.Printf("   • Status: ❌ Unavailable (check API key/endpoint)\n")
		}
	} else {
		fmt.Printf("   • Enabled: ❌ (no API key configured)\n")
	}
	
	// Local Backend Status  
	fmt.Println("\n🖥️  Local Backend:")
	if status["local_enabled"].(bool) {
		fmt.Printf("   • Enabled: ✅\n")
		fmt.Printf("   • Provider: %s\n", status["local_provider"])
		fmt.Printf("   • Model: %s\n", status["local_model"])
		if status["local_available"].(bool) {
			fmt.Printf("   • Status: ✅ Available\n")
		} else {
			fmt.Printf("   • Status: ❌ Unavailable (check if Ollama is running)\n")
		}
	} else {
		fmt.Printf("   • Enabled: ❌\n")
	}
	
	// Fallback Status
	fmt.Println("\n🔄 Fallback Backend:")
	fmt.Printf("   • Status: ✅ Always available\n")
	
	// Show active backend priority
	fmt.Println("\n⚡ Backend Priority:")
	if cfg.General.FallbackMode {
		fmt.Println("   1. 🔄 Fallback (forced)")
	} else {
		priority := 1
		if status["api_enabled"].(bool) {
			if status["api_available"].(bool) {
				fmt.Printf("   %d. 🌐 API (ready)\n", priority)
			} else {
				fmt.Printf("   %d. 🌐 API (unavailable)\n", priority)
			}
			priority++
		}
		if status["local_enabled"].(bool) {
			if status["local_available"].(bool) {
				fmt.Printf("   %d. 🖥️  Local (ready)\n", priority)
			} else {
				fmt.Printf("   %d. 🖥️  Local (unavailable)\n", priority)
			}
			priority++
		}
		fmt.Printf("   %d. 🔄 Fallback (always)\n", priority)
	}
	
	// Configuration hints
	fmt.Println("\n💡 Quick Setup:")
	if !status["api_enabled"].(bool) {
		fmt.Println("   • Set API key: export PARROT_API_KEY=\"your-key-here\"")
	}
	if !status["local_available"].(bool) && status["local_enabled"].(bool) {
		fmt.Printf("   • Install model: ollama pull %s\n", status["local_model"])
	}
	
	fmt.Println("\n   📖 Use 'parrot config' to create a configuration file")
}