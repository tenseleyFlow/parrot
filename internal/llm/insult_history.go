package llm

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// InsultHistory tracks shown insults to avoid repetition
type InsultHistory struct {
	RecentInsults   []HistoryEntry `json:"recent_insults"`
	InsultFrequency map[string]int `json:"insult_frequency"` // How many times each shown
	LastCleanup     time.Time      `json:"last_cleanup"`
	mu              sync.RWMutex
	maxHistory      int
	persistPath     string
}

// HistoryEntry represents a single shown insult
type HistoryEntry struct {
	Text      string    `json:"text"`
	Timestamp time.Time `json:"timestamp"`
	Context   string    `json:"context"` // Command that triggered it
	Score     float64   `json:"score"`   // Relevance score
}

// NewInsultHistory creates a new history tracker
func NewInsultHistory(maxHistory int) *InsultHistory {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}

	persistPath := filepath.Join(homeDir, ".parrot", "insult_history.json")

	history := &InsultHistory{
		RecentInsults:   make([]HistoryEntry, 0, maxHistory),
		InsultFrequency: make(map[string]int),
		maxHistory:      maxHistory,
		persistPath:     persistPath,
		LastCleanup:     time.Now(),
	}

	// Try to load existing history
	history.Load()

	return history
}

// RecordInsult adds an insult to the history
func (h *InsultHistory) RecordInsult(text string, context string, score float64) {
	h.mu.Lock()
	defer h.mu.Unlock()

	entry := HistoryEntry{
		Text:      text,
		Timestamp: time.Now(),
		Context:   context,
		Score:     score,
	}

	h.RecentInsults = append(h.RecentInsults, entry)

	// Trim if exceeds max
	if len(h.RecentInsults) > h.maxHistory {
		h.RecentInsults = h.RecentInsults[1:]
	}

	// Update frequency
	h.InsultFrequency[text]++

	// Periodic cleanup (once per day)
	if time.Since(h.LastCleanup) > 24*time.Hour {
		h.cleanup()
	}

	// Persist to disk
	h.persist()
}

// WasRecentlyShown checks if an insult was shown recently
func (h *InsultHistory) WasRecentlyShown(text string, withinLast int) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if withinLast > len(h.RecentInsults) {
		withinLast = len(h.RecentInsults)
	}

	startIdx := len(h.RecentInsults) - withinLast
	if startIdx < 0 {
		startIdx = 0
	}

	for i := startIdx; i < len(h.RecentInsults); i++ {
		if h.RecentInsults[i].Text == text {
			return true
		}
	}

	return false
}

// GetRecencyScore returns a score (0-1) based on how recently insult was shown
// 0 = shown very recently, 1 = never shown or shown long ago
func (h *InsultHistory) GetRecencyScore(text string) float64 {
	h.mu.RLock()
	defer h.mu.RUnlock()

	// Find the most recent occurrence
	for i := len(h.RecentInsults) - 1; i >= 0; i-- {
		if h.RecentInsults[i].Text == text {
			// Calculate position from end (more recent = lower score)
			position := len(h.RecentInsults) - i
			// Normalize to 0-1 range
			recencyPenalty := float64(position) / float64(h.maxHistory)
			// Invert so recent = low score
			return 1.0 - recencyPenalty
		}
	}

	return 1.0 // Not found = full score
}

// GetFrequencyScore returns a score based on how often insult was used
// Less frequent = higher score
func (h *InsultHistory) GetFrequencyScore(text string) float64 {
	h.mu.RLock()
	defer h.mu.RUnlock()

	frequency := h.InsultFrequency[text]
	if frequency == 0 {
		return 1.0 // Never shown
	}

	// Find max frequency for normalization
	maxFreq := 1
	for _, freq := range h.InsultFrequency {
		if freq > maxFreq {
			maxFreq = freq
		}
	}

	// Normalize and invert (higher frequency = lower score)
	return 1.0 - (float64(frequency) / float64(maxFreq) * 0.7) // Max 70% penalty
}

// GetNoveltyScore combines recency and frequency for overall novelty
func (h *InsultHistory) GetNoveltyScore(text string) float64 {
	recency := h.GetRecencyScore(text)
	frequency := h.GetFrequencyScore(text)

	// Weight recency more heavily (70% recency, 30% frequency)
	return (recency * 0.7) + (frequency * 0.3)
}

// cleanup removes old entries and resets counters
func (h *InsultHistory) cleanup() {
	// Remove entries older than 7 days
	cutoff := time.Now().AddDate(0, 0, -7)
	validEntries := make([]HistoryEntry, 0)

	for _, entry := range h.RecentInsults {
		if entry.Timestamp.After(cutoff) {
			validEntries = append(validEntries, entry)
		}
	}

	h.RecentInsults = validEntries

	// Reset frequency counters for removed insults
	newFrequency := make(map[string]int)
	for _, entry := range h.RecentInsults {
		newFrequency[entry.Text]++
	}
	h.InsultFrequency = newFrequency

	h.LastCleanup = time.Now()
}

// persist saves history to disk
func (h *InsultHistory) persist() {
	// Create directory if it doesn't exist
	dir := filepath.Dir(h.persistPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return // Silently fail if can't create directory
	}

	data, err := json.MarshalIndent(h, "", "  ")
	if err != nil {
		return
	}

	// Write to temp file first, then rename (atomic)
	tempPath := h.persistPath + ".tmp"
	if err := os.WriteFile(tempPath, data, 0644); err != nil {
		return
	}

	os.Rename(tempPath, h.persistPath)
}

// Load reads history from disk
func (h *InsultHistory) Load() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	data, err := os.ReadFile(h.persistPath)
	if err != nil {
		return err // File doesn't exist or can't be read
	}

	// Create a temporary struct to unmarshal into
	var loaded struct {
		RecentInsults   []HistoryEntry `json:"recent_insults"`
		InsultFrequency map[string]int `json:"insult_frequency"`
		LastCleanup     time.Time      `json:"last_cleanup"`
	}

	if err := json.Unmarshal(data, &loaded); err != nil {
		return err
	}

	h.RecentInsults = loaded.RecentInsults
	h.InsultFrequency = loaded.InsultFrequency
	h.LastCleanup = loaded.LastCleanup

	return nil
}

// Clear removes all history
func (h *InsultHistory) Clear() {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.RecentInsults = make([]HistoryEntry, 0, h.maxHistory)
	h.InsultFrequency = make(map[string]int)
	h.LastCleanup = time.Now()

	h.persist()
}

// GetStats returns statistics about insult history
func (h *InsultHistory) GetStats() map[string]interface{} {
	h.mu.RLock()
	defer h.mu.RUnlock()

	// Find most frequent insult
	mostFrequent := ""
	maxFreq := 0
	for text, freq := range h.InsultFrequency {
		if freq > maxFreq {
			maxFreq = freq
			mostFrequent = text
		}
	}

	return map[string]interface{}{
		"total_insults_shown": len(h.RecentInsults),
		"unique_insults":      len(h.InsultFrequency),
		"most_frequent":       mostFrequent,
		"most_frequent_count": maxFreq,
		"last_cleanup":        h.LastCleanup,
	}
}
