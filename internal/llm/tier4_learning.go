package llm

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

// UserFailureHistory tracks persistent failure patterns across sessions
type UserFailureHistory struct {
	TotalFailures     int                `json:"total_failures"`
	CurrentStreak     int                `json:"current_streak"`
	LongestStreak     int                `json:"longest_streak"`
	LastFailureTime   time.Time          `json:"last_failure_time"`
	CommandFrequency  map[string]int     `json:"command_frequency"`
	HourlyDistribution map[int]int       `json:"hourly_distribution"`
	BranchDisasters   map[string]int     `json:"branch_disasters"`
	ProjectFailures   map[string]int     `json:"project_failures"`
	ExitCodeFrequency map[int]int        `json:"exit_code_frequency"`
	WorstDay          string             `json:"worst_day"`
	DailyFailures     map[string]int     `json:"daily_failures"`

	mu sync.Mutex `json:"-"` // Mutex for thread-safe operations
}

var (
	userHistory     *UserFailureHistory
	historyFile     string
	historyLoadOnce sync.Once
)

// getHistoryFilePath returns the path to the failure history file
func getHistoryFilePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	parrotDir := filepath.Join(home, ".parrot")
	os.MkdirAll(parrotDir, 0755)

	return filepath.Join(parrotDir, "failures.json")
}

// LoadUserHistory loads failure history from disk
func LoadUserHistory() *UserFailureHistory {
	historyLoadOnce.Do(func() {
		historyFile = getHistoryFilePath()
		if historyFile == "" {
			userHistory = newUserHistory()
			return
		}

		data, err := os.ReadFile(historyFile)
		if err != nil {
			// File doesn't exist yet, create new history
			userHistory = newUserHistory()
			return
		}

		var history UserFailureHistory
		if err := json.Unmarshal(data, &history); err != nil {
			// Corrupted file, start fresh
			userHistory = newUserHistory()
			return
		}

		userHistory = &history
		if userHistory.CommandFrequency == nil {
			userHistory.CommandFrequency = make(map[string]int)
		}
		if userHistory.HourlyDistribution == nil {
			userHistory.HourlyDistribution = make(map[int]int)
		}
		if userHistory.BranchDisasters == nil {
			userHistory.BranchDisasters = make(map[string]int)
		}
		if userHistory.ProjectFailures == nil {
			userHistory.ProjectFailures = make(map[string]int)
		}
		if userHistory.ExitCodeFrequency == nil {
			userHistory.ExitCodeFrequency = make(map[int]int)
		}
		if userHistory.DailyFailures == nil {
			userHistory.DailyFailures = make(map[string]int)
		}

		// Check if streak should be reset (more than 1 hour since last failure)
		if time.Since(userHistory.LastFailureTime) > time.Hour {
			userHistory.CurrentStreak = 0
		}
	})

	return userHistory
}

// newUserHistory creates a new empty history
func newUserHistory() *UserFailureHistory {
	return &UserFailureHistory{
		CommandFrequency:   make(map[string]int),
		HourlyDistribution: make(map[int]int),
		BranchDisasters:    make(map[string]int),
		ProjectFailures:    make(map[string]int),
		ExitCodeFrequency:  make(map[int]int),
		DailyFailures:      make(map[string]int),
	}
}

// RecordFailure records a new failure in the history
func (h *UserFailureHistory) RecordFailure(ctx SmartFallbackContext) {
	if h == nil {
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	// Update counters
	h.TotalFailures++

	// Update streak
	if time.Since(h.LastFailureTime) <= time.Hour {
		h.CurrentStreak++
	} else {
		h.CurrentStreak = 1
	}

	if h.CurrentStreak > h.LongestStreak {
		h.LongestStreak = h.CurrentStreak
	}

	h.LastFailureTime = time.Now()

	// Track command frequency
	if ctx.Command != "" {
		h.CommandFrequency[ctx.Command]++
	}

	// Track hourly distribution
	hour := time.Now().Hour()
	h.HourlyDistribution[hour]++

	// Track branch disasters
	if ctx.GitBranch != "" {
		h.BranchDisasters[ctx.GitBranch]++
	}

	// Track project failures
	if ctx.ProjectType != "" {
		h.ProjectFailures[ctx.ProjectType]++
	}

	// Track exit codes
	if ctx.ExitCode > 0 {
		h.ExitCodeFrequency[ctx.ExitCode]++
	}

	// Track daily failures
	today := time.Now().Format("2006-01-02")
	h.DailyFailures[today]++

	// Determine worst day
	maxFailures := 0
	for day, count := range h.DailyFailures {
		if count > maxFailures {
			maxFailures = count
			h.WorstDay = day
		}
	}

	// Save to disk
	h.Save()
}

// Save writes the history to disk
func (h *UserFailureHistory) Save() error {
	if h == nil || historyFile == "" {
		return nil
	}

	data, err := json.MarshalIndent(h, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(historyFile, data, 0644)
}

// GetMostFailedCommand returns the command with the most failures
func (h *UserFailureHistory) GetMostFailedCommand() (string, int) {
	if h == nil {
		return "", 0
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	maxCmd := ""
	maxCount := 0

	for cmd, count := range h.CommandFrequency {
		if count > maxCount {
			maxCount = count
			maxCmd = cmd
		}
	}

	return maxCmd, maxCount
}

// GetWorstHour returns the hour with the most failures
func (h *UserFailureHistory) GetWorstHour() int {
	if h == nil {
		return -1
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	maxHour := -1
	maxCount := 0

	for hour, count := range h.HourlyDistribution {
		if count > maxCount {
			maxCount = count
			maxHour = hour
		}
	}

	return maxHour
}

// GetMostDisastrousBranch returns the branch with the most failures
func (h *UserFailureHistory) GetMostDisastrousBranch() (string, int) {
	if h == nil {
		return "", 0
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	maxBranch := ""
	maxCount := 0

	for branch, count := range h.BranchDisasters {
		if count > maxCount {
			maxCount = count
			maxBranch = branch
		}
	}

	return maxBranch, maxCount
}

// GetFailureStats returns a summary of failure statistics
type FailureStats struct {
	TotalFailures     int
	CurrentStreak     int
	LongestStreak     int
	WorstCommand      string
	WorstCommandCount int
	WorstHour         int
	WorstBranch       string
	WorstBranchCount  int
	TodayFailures     int
	WorstDay          string
	WorstDayCount     int
}

func (h *UserFailureHistory) GetStats() FailureStats {
	if h == nil {
		return FailureStats{}
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	stats := FailureStats{
		TotalFailures: h.TotalFailures,
		CurrentStreak: h.CurrentStreak,
		LongestStreak: h.LongestStreak,
	}

	// Find worst command
	for cmd, count := range h.CommandFrequency {
		if count > stats.WorstCommandCount {
			stats.WorstCommandCount = count
			stats.WorstCommand = cmd
		}
	}

	// Find worst hour
	maxHourCount := 0
	for hour, count := range h.HourlyDistribution {
		if count > maxHourCount {
			maxHourCount = count
			stats.WorstHour = hour
		}
	}

	// Find worst branch
	for branch, count := range h.BranchDisasters {
		if count > stats.WorstBranchCount {
			stats.WorstBranchCount = count
			stats.WorstBranch = branch
		}
	}

	// Today's failures
	today := time.Now().Format("2006-01-02")
	stats.TodayFailures = h.DailyFailures[today]

	// Worst day
	stats.WorstDay = h.WorstDay
	if h.WorstDay != "" {
		stats.WorstDayCount = h.DailyFailures[h.WorstDay]
	}

	return stats
}

// GetTopFailedCommands returns the top N most failed commands
func (h *UserFailureHistory) GetTopFailedCommands(n int) []string {
	if h == nil {
		return nil
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	type cmdCount struct {
		cmd   string
		count int
	}

	var commands []cmdCount
	for cmd, count := range h.CommandFrequency {
		commands = append(commands, cmdCount{cmd, count})
	}

	sort.Slice(commands, func(i, j int) bool {
		return commands[i].count > commands[j].count
	})

	var result []string
	for i := 0; i < n && i < len(commands); i++ {
		result = append(result, commands[i].cmd)
	}

	return result
}
