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

	// Get current working directory to find the hook script
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("❌ Error getting working directory: %v\n", err)
		return
	}
	
	hookPath := filepath.Join(wd, "parrot-hook.sh")
	
	// Check if hook script exists
	if _, err := os.Stat(hookPath); os.IsNotExist(err) {
		fmt.Printf("❌ Hook script not found at: %s\n", hookPath)
		fmt.Println("Make sure you're running this from the parrot directory.")
		return
	}

	// Add source line to RC file
	sourceLine := fmt.Sprintf("source \"%s\"", hookPath)
	
	fmt.Printf("🦜 Installing parrot hooks to: %s\n", rcFile)
	fmt.Printf("📝 Adding line: %s\n", sourceLine)
	
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
	
	_, err = file.WriteString(fmt.Sprintf("\n# Parrot CLI hooks\n%s\n", sourceLine))
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
	
	return strings.Contains(string(content), sourceLine)
}