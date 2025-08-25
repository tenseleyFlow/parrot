package colors

import (
	"fmt"
	"os"
	"strings"
)

// ANSI color codes
const (
	Reset  = "\033[0m"
	Bold   = "\033[1m"
	Dim    = "\033[2m"
	
	// Regular colors
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"
	
	// Bright colors
	BrightRed     = "\033[91m"
	BrightGreen   = "\033[92m"
	BrightYellow  = "\033[93m"
	BrightBlue    = "\033[94m"
	BrightMagenta = "\033[95m"
	BrightCyan    = "\033[96m"
	BrightWhite   = "\033[97m"
)

// ParrotStyle represents different visual styles for parrot output
type ParrotStyle struct {
	Parrot   string // Color for the parrot emoji/prefix
	Response string // Color for the response text
	Accent   string // Color for emphasis
}

// Predefined styles for different personalities
var Styles = map[string]ParrotStyle{
	"mild": {
		Parrot:   BrightBlue,
		Response: Blue,
		Accent:   BrightCyan,
	},
	"sarcastic": {
		Parrot:   BrightYellow,
		Response: Yellow,
		Accent:   BrightMagenta,
	},
	"savage": {
		Parrot:   BrightRed,
		Response: Red,
		Accent:   BrightYellow,
	},
	"default": {
		Parrot:   BrightGreen,
		Response: Green,
		Accent:   BrightCyan,
	},
}

// ColorEnabled checks if color output should be enabled
func ColorEnabled() bool {
	// Disable colors if NO_COLOR is set
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	
	// Disable colors if not a terminal
	if os.Getenv("TERM") == "dumb" {
		return false
	}
	
	// Check if stdout is a terminal (simplified check)
	stat, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	
	// Check if it's a character device (terminal)
	return (stat.Mode() & os.ModeCharDevice) != 0
}

// Colorize adds color to text if colors are enabled
func Colorize(color, text string) string {
	if !ColorEnabled() {
		return text
	}
	return color + text + Reset
}

// FormatParrotOutput formats the parrot response with personality-based colors
func FormatParrotOutput(personality, response string, enhanced bool) string {
	if !ColorEnabled() {
		return fmt.Sprintf("🦜 %s", response)
	}
	
	style, exists := Styles[personality]
	if !exists {
		style = Styles["default"]
	}
	
	if enhanced {
		return formatEnhancedOutput(style, response)
	}
	
	// Simple colored output
	parrotEmoji := Colorize(style.Parrot, "🦜")
	coloredResponse := Colorize(style.Response, response)
	
	return fmt.Sprintf("%s %s", parrotEmoji, coloredResponse)
}

// formatEnhancedOutput creates fancy formatted output with personality-specific styling
func formatEnhancedOutput(style ParrotStyle, response string) string {
	var output strings.Builder
	
	// Fancy parrot prefix with personality
	parrotPrefix := Colorize(style.Parrot, "🦜 ▶") 
	
	// Add some visual flair based on personality
	border := Colorize(style.Accent, "━")
	
	// Format the response with potential emphasis
	coloredResponse := enhanceResponseText(style, response)
	
	output.WriteString(fmt.Sprintf("%s %s %s", border, parrotPrefix, coloredResponse))
	
	return output.String()
}

// enhanceResponseText adds emphasis to certain words in responses
func enhanceResponseText(style ParrotStyle, response string) string {
	if !ColorEnabled() {
		return response
	}
	
	// Words to emphasize for extra sass
	emphasisWords := []string{
		"failed", "error", "disaster", "incompetent", "broken", 
		"genius", "classic", "impressive", "amazing", "brilliant",
		"404", "rejected", "crashed", "destroyed",
	}
	
	result := response
	for _, word := range emphasisWords {
		// Case-insensitive replacement with emphasis
		lowerWord := strings.ToLower(word)
		if strings.Contains(strings.ToLower(result), lowerWord) {
			// Find and replace with emphasized version
			result = replaceWordWithEmphasis(result, word, style.Accent)
		}
	}
	
	// Color the main response
	return Colorize(style.Response, result)
}

// replaceWordWithEmphasis replaces words with emphasized versions (case-insensitive)
func replaceWordWithEmphasis(text, word, accentColor string) string {
	words := strings.Fields(text)
	for i, w := range words {
		// Remove punctuation for comparison
		cleanWord := strings.Trim(strings.ToLower(w), ".,!?;:")
		if cleanWord == strings.ToLower(word) {
			// Keep original punctuation, but emphasize the word
			punctuation := ""
			if len(w) > len(cleanWord) {
				punctuation = w[len(w)-1:]
			}
			words[i] = Colorize(accentColor+Bold, strings.ToUpper(word)) + punctuation + Reset
		}
	}
	return strings.Join(words, " ")
}

// GetAvailableStyles returns list of available color styles
func GetAvailableStyles() []string {
	styles := make([]string, 0, len(Styles))
	for style := range Styles {
		if style != "default" {
			styles = append(styles, style)
		}
	}
	return styles
}

// Demo shows color samples for all personalities
func Demo() {
	fmt.Println("🎨 Parrot Color Demo")
	fmt.Println("━━━━━━━━━━━━━━━━━━━")
	
	responses := map[string]string{
		"mild":      "Git command failed. Maybe check your remote branch?",
		"sarcastic": "Git good? More like git wrecked!",
		"savage":    "Your git skills are as non-existent as your social life.",
	}
	
	for personality, response := range responses {
		fmt.Printf("\n%s personality:\n", strings.Title(personality))
		fmt.Printf("  Simple: %s\n", FormatParrotOutput(personality, response, false))
		fmt.Printf("  Enhanced: %s\n", FormatParrotOutput(personality, response, true))
	}
	
	fmt.Printf("\nColors enabled: %t\n", ColorEnabled())
}