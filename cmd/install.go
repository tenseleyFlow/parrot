package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install parrot shell hooks",
	Long:  "Adds parrot hooks to your shell configuration",
	Run:   installHooks,
}

func init() {
	rootCmd.AddCommand(installCmd)
}

func installHooks(cmd *cobra.Command, args []string) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("❌ Error getting home directory: %v\n", err)
		return
	}

	// Detect shell and appropriate RC file
	shell := os.Getenv("SHELL")
	var rcFile string
	
	if filepath.Base(shell) == "zsh" {
		rcFile = filepath.Join(homeDir, ".zshrc")
	} else {
		rcFile = filepath.Join(homeDir, ".bashrc")
	}

	// Try standard installation paths for the hook script
	var hookPath string
	possiblePaths := []string{
		"/usr/share/parrot/parrot-hook.sh",     // RPM installation
		"/usr/local/share/parrot/parrot-hook.sh", // Manual installation
		"./parrot-hook.sh",                     // Development
	}
	
	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			hookPath = path
			break
		}
	}
	
	if hookPath == "" {
		fmt.Println("❌ Hook script not found. Searched in:")
		for _, path := range possiblePaths {
			fmt.Printf("   - %s\n", path)
		}
		fmt.Println("Make sure parrot is properly installed.")
		return
	}

	// Add source line and environment setup to RC file
	sourceLine := fmt.Sprintf("source \"%s\"", hookPath)
	
	fmt.Printf("🦜 Installing parrot hooks to: %s\n", rcFile)
	fmt.Printf("📝 Adding hook: %s\n", sourceLine)
	fmt.Printf("🔧 Configuring Ollama for better performance\n")
	
	// Check if already installed
	if isAlreadyInstalled(rcFile, sourceLine) {
		fmt.Println("✅ Parrot hooks already installed!")
		return
	}
	
	// Append to RC file
	file, err := os.OpenFile(rcFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("❌ Error opening %s: %v\n", rcFile, err)
		return
	}
	defer file.Close()
	
	installContent := fmt.Sprintf(`
# Parrot CLI hooks and configuration
export OLLAMA_KEEP_ALIVE="1h"  # Keep AI models loaded for better performance
%s
`, sourceLine)
	
	_, err = file.WriteString(installContent)
	if err != nil {
		fmt.Printf("❌ Error writing to %s: %v\n", rcFile, err)
		return
	}
	
	fmt.Println("✅ Parrot hooks installed successfully!")
	fmt.Println("🔄 Run 'source ~/.bashrc' (or ~/.zshrc) to activate, or start a new shell session.")
}

func isAlreadyInstalled(rcFile, sourceLine string) bool {
	content, err := os.ReadFile(rcFile)
	if err != nil {
		return false
	}
	
	// Check for both the source line and OLLAMA_KEEP_ALIVE setting
	contentStr := string(content)
	return strings.Contains(contentStr, sourceLine) && strings.Contains(contentStr, "OLLAMA_KEEP_ALIVE")
}