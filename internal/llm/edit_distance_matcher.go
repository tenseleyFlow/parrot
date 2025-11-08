package llm

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// EditDistanceMatcher finds similar past failures and adapts their insults
// Uses Levenshtein distance for command similarity matching
type EditDistanceMatcher struct {
	mu               sync.RWMutex
	commandHistory   []CommandRecord
	maxHistory       int
	similarityThresh float64
	persistencePath  string
}

// CommandRecord stores a past command failure with its successful insult
type CommandRecord struct {
	Command     string
	CommandType string
	ErrorPattern string
	ProjectType string
	Timestamp   time.Time
	Insult      string
	Effectiveness float64 // How well this insult worked (from RL)
}

// SimilarCommand represents a similar past command
type SimilarCommand struct {
	Record     CommandRecord
	Similarity float64
	Distance   int
}

// NewEditDistanceMatcher creates a new matcher
func NewEditDistanceMatcher() *EditDistanceMatcher {
	homeDir, _ := os.UserHomeDir()
	persistPath := filepath.Join(homeDir, ".parrot", "command_history.json")

	matcher := &EditDistanceMatcher{
		commandHistory:   make([]CommandRecord, 0),
		maxHistory:       1000, // Keep last 1000 commands
		similarityThresh: 0.7,  // 70% similarity required
		persistencePath:  persistPath,
	}

	// Try to load existing history
	matcher.Load()

	return matcher
}

// RecordCommand stores a command failure with its insult
func (edm *EditDistanceMatcher) RecordCommand(
	ctx *SmartFallbackContext,
	insult string,
	effectiveness float64,
) {
	edm.mu.Lock()
	defer edm.mu.Unlock()

	record := CommandRecord{
		Command:       ctx.FullCommand,
		CommandType:   ctx.CommandType,
		ErrorPattern:  ctx.ErrorPattern,
		ProjectType:   ctx.ProjectType,
		Timestamp:     time.Now(),
		Insult:        insult,
		Effectiveness: effectiveness,
	}

	edm.commandHistory = append(edm.commandHistory, record)

	// Prune if too large
	if len(edm.commandHistory) > edm.maxHistory {
		// Remove oldest entries
		edm.commandHistory = edm.commandHistory[len(edm.commandHistory)-edm.maxHistory:]
	}

	// Periodically persist
	if len(edm.commandHistory)%50 == 0 {
		go edm.Save()
	}
}

// FindSimilarCommands finds commands similar to the current one
func (edm *EditDistanceMatcher) FindSimilarCommands(
	command string,
	topK int,
) []SimilarCommand {
	edm.mu.RLock()
	defer edm.mu.RUnlock()

	if len(edm.commandHistory) == 0 {
		return nil
	}

	similarities := make([]SimilarCommand, 0)

	// Calculate similarity to each historical command
	for _, record := range edm.commandHistory {
		distance := levenshteinDistanceTier6(command, record.Command)
		maxLen := maxTier6(len(command), len(record.Command))

		// Calculate similarity (1.0 = identical, 0.0 = completely different)
		similarity := 1.0 - (float64(distance) / float64(maxLen))

		if similarity >= edm.similarityThresh {
			similarities = append(similarities, SimilarCommand{
				Record:     record,
				Similarity: similarity,
				Distance:   distance,
			})
		}
	}

	// Sort by similarity (descending)
	sortSimilarCommands(similarities)

	// Return top K
	if len(similarities) > topK {
		similarities = similarities[:topK]
	}

	return similarities
}

// GetAdaptedInsult gets an insult adapted from similar command
func (edm *EditDistanceMatcher) GetAdaptedInsult(ctx *SmartFallbackContext) string {
	similar := edm.FindSimilarCommands(ctx.FullCommand, 5)

	if len(similar) == 0 {
		return ""
	}

	// Pick the most effective similar command
	bestIdx := 0
	bestScore := similar[0].Similarity * similar[0].Record.Effectiveness

	for i, sim := range similar {
		score := sim.Similarity * sim.Record.Effectiveness
		if score > bestScore {
			bestScore = score
			bestIdx = i
		}
	}

	best := similar[bestIdx]

	// Adapt the insult to current context
	adapted := edm.adaptInsult(best.Record.Insult, best.Record.Command, ctx)

	return adapted
}

// adaptInsult adapts an insult from past command to current context
func (edm *EditDistanceMatcher) adaptInsult(
	insult string,
	oldCommand string,
	ctx *SmartFallbackContext,
) string {
	adapted := insult

	// Extract command differences
	oldParts := tokenizeCommand(oldCommand)
	newParts := tokenizeCommand(ctx.FullCommand)

	// Find what changed
	oldUnique := findUnique(oldParts, newParts)
	newUnique := findUnique(newParts, oldParts)

	// Replace old unique parts with new ones
	for i := 0; i < minTier6(len(oldUnique), len(newUnique)); i++ {
		adapted = replaceWordTier6(adapted, oldUnique[i], newUnique[i])
	}

	// Add context-specific references if missing
	if ctx.ProjectType != "" && !containsWordTier6(adapted, ctx.ProjectType) {
		// Append project context
		adapted = adapted + " In " + ctx.ProjectType + "."
	}

	return adapted
}


// Save persists command history to disk
func (edm *EditDistanceMatcher) Save() error {
	edm.mu.RLock()
	defer edm.mu.RUnlock()

	// Ensure directory exists
	dir := filepath.Dir(edm.persistencePath)
	os.MkdirAll(dir, 0755)

	// Marshal to JSON
	data, err := json.MarshalIndent(edm.commandHistory, "", "  ")
	if err != nil {
		return err
	}

	// Write atomically
	tmpPath := edm.persistencePath + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		return err
	}

	return os.Rename(tmpPath, edm.persistencePath)
}

// Load restores command history from disk
func (edm *EditDistanceMatcher) Load() error {
	edm.mu.Lock()
	defer edm.mu.Unlock()

	data, err := os.ReadFile(edm.persistencePath)
	if err != nil {
		return err // File might not exist yet
	}

	history := make([]CommandRecord, 0)
	if err := json.Unmarshal(data, &history); err != nil {
		return err
	}

	edm.commandHistory = history
	return nil
}

// GetStats returns matcher statistics
func (edm *EditDistanceMatcher) GetStats() map[string]interface{} {
	edm.mu.RLock()
	defer edm.mu.RUnlock()

	avgEffectiveness := 0.0
	if len(edm.commandHistory) > 0 {
		for _, record := range edm.commandHistory {
			avgEffectiveness += record.Effectiveness
		}
		avgEffectiveness /= float64(len(edm.commandHistory))
	}

	return map[string]interface{}{
		"total_commands":      len(edm.commandHistory),
		"avg_effectiveness":   avgEffectiveness,
		"similarity_threshold": edm.similarityThresh,
	}
}

// Helper functions

func tokenizeCommand(command string) []string {
	// Simple tokenization (split on spaces and special chars)
	tokens := make([]string, 0)
	current := ""

	for _, r := range command {
		if r == ' ' || r == '/' || r == '-' || r == '.' {
			if len(current) > 0 {
				tokens = append(tokens, current)
				current = ""
			}
		} else {
			current += string(r)
		}
	}

	if len(current) > 0 {
		tokens = append(tokens, current)
	}

	return tokens
}

func findUnique(slice1, slice2 []string) []string {
	unique := make([]string, 0)
	set2 := make(map[string]bool)

	for _, item := range slice2 {
		set2[item] = true
	}

	for _, item := range slice1 {
		if !set2[item] {
			unique = append(unique, item)
		}
	}

	return unique
}

func sortSimilarCommands(commands []SimilarCommand) {
	// Bubble sort by similarity (descending)
	n := len(commands)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if commands[j].Similarity < commands[j+1].Similarity {
				commands[j], commands[j+1] = commands[j+1], commands[j]
			}
		}
	}
}
