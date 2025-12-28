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
	shellName := filepath.Base(shell)

	if shellName == "fish" {
		// Fish uses a different config directory structure
		rcFile = filepath.Join(homeDir, ".config/fish/conf.d/parrot.fish")
		// Create directory if needed
		confDir := filepath.Join(homeDir, ".config/fish/conf.d")
		if err := os.MkdirAll(confDir, 0755); err != nil {
			fmt.Printf("❌ Error creating fish config directory: %v\n", err)
			return
		}
	} else if shellName == "zsh" {
		rcFile = filepath.Join(homeDir, ".zshrc")
	} else {
		rcFile = filepath.Join(homeDir, ".bashrc")
	}

	// Try standard installation paths for the hook script
	var hookPath string
	var hookFile string
	if shellName == "fish" {
		hookFile = "parrot-hook.fish"
	} else {
		hookFile = "parrot-hook.sh"
	}

	// Build list of possible paths, starting with most specific
	possiblePaths := []string{}

	// Check relative to executable (works for NixOS, Homebrew, and other non-FHS systems)
	if exePath, err := os.Executable(); err == nil {
		// Resolve symlinks to get actual binary location
		if realPath, err := filepath.EvalSymlinks(exePath); err == nil {
			exePath = realPath
		}
		exeDir := filepath.Dir(exePath)
		// Try ../share/parrot/ relative to bin directory
		possiblePaths = append(possiblePaths, filepath.Join(exeDir, "..", "share", "parrot", hookFile))
	}

	// Standard FHS and user paths
	possiblePaths = append(possiblePaths,
		filepath.Join("/usr/share/parrot", hookFile),            // RPM/system installation
		filepath.Join("/usr/local/share/parrot", hookFile),      // Manual system installation
		filepath.Join(homeDir, ".local/share/parrot", hookFile), // make install (user local)
		filepath.Join(".", hookFile),                            // Development
	)
	
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

	fmt.Printf("🦜 Installing parrot hooks to: %s\n", rcFile)
	fmt.Printf("🔧 Configuring Ollama for better performance\n")

	if shellName == "fish" {
		// For fish, copy the hook file directly to conf.d
		// Fish automatically sources files in conf.d
		hookContent, err := os.ReadFile(hookPath)
		if err != nil {
			fmt.Printf("❌ Error reading hook file: %v\n", err)
			return
		}

		// Check if already installed
		if _, err := os.Stat(rcFile); err == nil {
			fmt.Println("✅ Parrot hooks already installed!")
			return
		}

		// Write hook file with Ollama configuration
		file, err := os.OpenFile(rcFile, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("❌ Error creating %s: %v\n", rcFile, err)
			return
		}
		defer file.Close()

		fishConfig := fmt.Sprintf(`# Parrot CLI hooks and configuration
# Keep AI models loaded for better performance
set -gx OLLAMA_KEEP_ALIVE "1h"

%s`, string(hookContent))

		_, err = file.WriteString(fishConfig)
		if err != nil {
			fmt.Printf("❌ Error writing to %s: %v\n", rcFile, err)
			return
		}

		fmt.Println("✅ Parrot hooks installed successfully!")
		fmt.Println("🔄 Restart fish or run 'source ~/.config/fish/config.fish' to activate.")
	} else {
		// For bash/zsh, add source line to RC file
		sourceLine := fmt.Sprintf("source \"%s\"", hookPath)
		fmt.Printf("📝 Adding hook: %s\n", sourceLine)

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
		if shellName == "zsh" {
			fmt.Println("🔄 Run 'source ~/.zshrc' to activate, or start a new shell session.")
		} else {
			fmt.Println("🔄 Run 'source ~/.bashrc' to activate, or start a new shell session.")
		}
	}
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