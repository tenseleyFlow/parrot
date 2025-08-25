package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"parrot/internal/config"
	"parrot/internal/llm"

	"github.com/spf13/cobra"
)

var firstrunCmd = &cobra.Command{
	Use:    "firstrun",
	Short:  "Check if this is the first run and guide setup",
	Long:   "Detects if parrot needs initial setup and guides the user",
	Hidden: true, // Hide from help - internal command
	Run:    checkFirstRun,
}

func init() {
	rootCmd.AddCommand(firstrunCmd)
}

func checkFirstRun(cmd *cobra.Command, args []string) {
	if isFirstRun() {
		fmt.Println("🦜 Welcome to Parrot!")
		fmt.Println("====================")
		fmt.Println("It looks like this is your first time running parrot.")
		fmt.Println("Let's get you set up for maximum sass!")
		fmt.Println()
		
		fmt.Print("Would you like to run the setup wizard? [Y/n]: ")
		var response string
		fmt.Scanln(&response)
		
		if response == "" || response == "y" || response == "Y" {
			// Run setup
			runSetup(cmd, args)
		} else {
			fmt.Println("⏭️  Skipped setup.")
			fmt.Println("💡 Run 'parrot setup' anytime to configure backends.")
			fmt.Println("🔄 Parrot works with fallback responses right now!")
		}
		
		// Mark as not first run
		markSetupComplete()
	} else {
		fmt.Println("✅ Parrot already set up!")
		fmt.Println("💡 Use 'parrot status' to check current configuration.")
	}
}

func isFirstRun() bool {
	// Check multiple indicators of setup completion
	cfg, err := config.LoadConfig()
	if err != nil {
		// No config file found - likely first run
		return true
	}
	
	// Check if any intelligent backend is available
	manager := llm.NewLLMManager(cfg)
	status := manager.GetStatus()
	
	hasAPI := status["api_available"].(bool)
	hasLocal := status["local_available"].(bool)
	
	// If no AI backends are working, consider it first run
	if !hasAPI && !hasLocal {
		return true
	}
	
	// Check for setup marker file
	markerPath := getSetupMarkerPath()
	if _, err := os.Stat(markerPath); os.IsNotExist(err) {
		return true
	}
	
	return false
}

func markSetupComplete() {
	markerPath := getSetupMarkerPath()
	
	// Create marker directory if needed
	if err := os.MkdirAll(filepath.Dir(markerPath), 0755); err != nil {
		return // Fail silently
	}
	
	// Create marker file
	file, err := os.Create(markerPath)
	if err != nil {
		return // Fail silently
	}
	defer file.Close()
	
	file.WriteString("setup_complete\n")
}

func getSetupMarkerPath() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		// Fallback to home directory
		homeDir, _ := os.UserHomeDir()
		return filepath.Join(homeDir, ".parrot_setup_complete")
	}
	return filepath.Join(configDir, "parrot", ".setup_complete")
}

// IsParrotSetup is a helper function for other commands to check setup status
func IsParrotSetup() bool {
	return !isFirstRun()
}