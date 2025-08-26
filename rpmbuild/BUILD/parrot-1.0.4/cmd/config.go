package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"parrot/internal/config"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage parrot configuration",
	Long:  "Create and manage parrot configuration files",
}

var configInitCmd = &cobra.Command{
	Use:   "init [path]",
	Short: "Create a sample configuration file",
	Long:  "Create a sample configuration file with default values",
	Args:  cobra.MaximumNArgs(1),
	Run:   initConfig,
}

func init() {
	configCmd.AddCommand(configInitCmd)
	rootCmd.AddCommand(configCmd)
}

func initConfig(cmd *cobra.Command, args []string) {
	var configPath string
	
	if len(args) > 0 {
		// User specified path
		configPath = args[0]
	} else {
		// Use default user config path
		configDir, err := os.UserConfigDir()
		if err != nil {
			fmt.Printf("❌ Error getting config directory: %v\n", err)
			return
		}
		configPath = filepath.Join(configDir, "parrot", "config.toml")
	}
	
	// Check if config already exists
	if _, err := os.Stat(configPath); err == nil {
		fmt.Printf("❌ Configuration file already exists at: %s\n", configPath)
		fmt.Println("   Remove it first if you want to recreate it.")
		return
	}
	
	// Create the config file
	if err := config.CreateSampleConfig(configPath); err != nil {
		fmt.Printf("❌ Error creating config file: %v\n", err)
		return
	}
	
	fmt.Printf("✅ Created configuration file at: %s\n", configPath)
	fmt.Println("\n💡 Next steps:")
	fmt.Println("   1. Edit the config file to add your API key:")
	fmt.Printf("      api_key = \"your-openai-api-key-here\"\n")
	fmt.Println("   2. Test with: parrot status")
	fmt.Println("   3. Try: parrot mock \"git push\" \"1\"")
	
	fmt.Println("\n🔧 Configuration options:")
	fmt.Println("   • API providers: openai, anthropic, custom")
	fmt.Println("   • Personalities: mild, sarcastic, savage") 
	fmt.Println("   • Local models: phi3.5:3.8b, llama3.2:3b")
	fmt.Println("   • Environment variables: PARROT_API_KEY, PARROT_DEBUG")
}