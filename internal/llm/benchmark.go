package llm

import (
	"fmt"
	"math"
	"time"
)

// BenchmarkSample represents a real command failure with expected outputs
type BenchmarkSample struct {
	ID          string
	Command     string
	ExitCode    int
	Stderr      string
	Context     SmartFallbackContext
	Category    string // "git", "npm", "docker", etc.
	Description string
	GoldInsults []string // Human-written example insults
	Tags        []string // Expected tags for this scenario
}

// BenchmarkResults contains evaluation metrics
type BenchmarkResults struct {
	SystemName      string
	TotalSamples    int
	AvgRelevance    float64
	AvgLatency      time.Duration
	AvgConfidence   float64
	DiversityScore  float64
	FallbackRate    float64
	MemoryUsageKB   int
	DetailedScores  []SampleScore
}

// SampleScore contains per-sample evaluation
type SampleScore struct {
	SampleID       string
	GeneratedInsult string
	Relevance      float64 // 0-1: How relevant to the error
	Latency        time.Duration
	Confidence     float64
	NoveltyScore   float64
	Method         string // "semantic", "tag", "markov", "ensemble"
}

// Benchmark framework for systematic evaluation
type Benchmark struct {
	Name    string
	Samples []BenchmarkSample
}

// NewBenchmark creates a comprehensive benchmark dataset
func NewBenchmark() *Benchmark {
	return &Benchmark{
		Name:    "Parrot Insult Quality Benchmark v1.0",
		Samples: createBenchmarkSamples(),
	}
}

// createBenchmarkSamples creates a comprehensive test dataset
func createBenchmarkSamples() []BenchmarkSample {
	samples := []BenchmarkSample{}

	// Git failures
	samples = append(samples, BenchmarkSample{
		ID:       "git-001",
		Command:  "git push origin main",
		ExitCode: 1,
		Stderr:   "error: failed to push some refs\nTo github.com:user/repo.git\n ! [rejected] main -> main (fetch first)",
		Context: SmartFallbackContext{
			CommandType:       "git",
			Command:           "git",
			Subcommand:        "push",
			GitBranch:         "main",
			ErrorPattern:      "permission_denied",
			IsRepeatedFailure: false,
		},
		Category:    "git",
		Description: "Git push rejected on main branch",
		GoldInsults: []string{
			"Push rejected. Did you forget to pull first?",
			"The remote has standards. Your code doesn't meet them.",
		},
		Tags: []string{"git", "push", "main_branch"},
	})

	samples = append(samples, BenchmarkSample{
		ID:       "git-002",
		Command:  "git merge feature/new-ui",
		ExitCode: 1,
		Stderr:   "CONFLICT (content): Merge conflict in src/app.js\nAutomatic merge failed; fix conflicts and then commit the result.",
		Context: SmartFallbackContext{
			CommandType:       "git",
			Command:           "git",
			Subcommand:        "merge",
			GitBranch:         "main",
			ErrorPattern:      "merge_conflict",
			IsRepeatedFailure: false,
		},
		Category:    "git",
		Description: "Merge conflict",
		GoldInsults: []string{
			"Merge conflict. Maybe communicate with your team?",
			"<<<<<<< HEAD is not a valid merge resolution strategy",
		},
		Tags: []string{"git", "merge", "merge_conflict"},
	})

	samples = append(samples, BenchmarkSample{
		ID:       "git-003",
		Command:  "git push --force origin main",
		ExitCode: 1,
		Stderr:   "error: refusing to update checked out branch: refs/heads/main",
		Context: SmartFallbackContext{
			CommandType:       "git",
			Command:           "git",
			Subcommand:        "push",
			GitBranch:         "main",
			ErrorPattern:      "permission_denied",
			IsRepeatedFailure: true,
			TimeOfDay:         2,
		},
		Category:    "git",
		Description: "Force push to main at 2 AM (repeated failure)",
		GoldInsults: []string{
			"Force pushing to main at 2 AM? Bold strategy.",
			"--force won't force competence into you",
		},
		Tags: []string{"git", "push", "main_branch", "late_night", "repeated"},
	})

	// NPM failures
	samples = append(samples, BenchmarkSample{
		ID:       "npm-001",
		Command:  "npm install",
		ExitCode: 1,
		Stderr:   "npm ERR! code ENOENT\nnpm ERR! syscall open\nnpm ERR! path /home/user/project/package.json\nnpm ERR! errno -2",
		Context: SmartFallbackContext{
			CommandType:  "nodejs",
			Command:      "npm",
			Subcommand:   "install",
			ProjectType:  "node",
			ErrorPattern: "not_found",
		},
		Category:    "npm",
		Description: "Missing package.json",
		GoldInsults: []string{
			"package.json not found. Neither is your organizational skill.",
			"Are you in the right directory? Rhetorical question.",
		},
		Tags: []string{"npm", "install", "not_found"},
	})

	samples = append(samples, BenchmarkSample{
		ID:       "npm-002",
		Command:  "npm install typescript --save-dev",
		ExitCode: 1,
		Stderr:   "npm ERR! code ERESOLVE\nnpm ERR! ERESOLVE unable to resolve dependency tree\nnpm ERR! peer dep missing: react@^18.0.0",
		Context: SmartFallbackContext{
			CommandType:  "nodejs",
			Command:      "npm",
			Subcommand:   "install",
			ProjectType:  "node",
			ErrorPattern: "dependency",
		},
		Category:    "npm",
		Description: "Dependency resolution failure",
		GoldInsults: []string{
			"Dependency hell. You're everyone's least favorite dependency.",
			"ERESOLVE: Can't resolve your incompetence either",
		},
		Tags: []string{"npm", "install", "dependency"},
	})

	samples = append(samples, BenchmarkSample{
		ID:       "npm-003",
		Command:  "npm test",
		ExitCode: 1,
		Stderr:   "FAIL src/components/App.test.js\n  ● App › renders correctly\n    expect(received).toEqual(expected)\n    Expected: true\n    Received: false",
		Context: SmartFallbackContext{
			CommandType:  "nodejs",
			Command:      "npm",
			Subcommand:   "test",
			ProjectType:  "node",
			ErrorPattern: "test_failure",
			IsCI:         true,
			CIProvider:   "github",
		},
		Category:    "npm",
		Description: "Test failure in CI",
		GoldInsults: []string{
			"Tests failed. Shocking absolutely no one who read your code",
			"Did you test this before committing? Oh wait, that's what CI is for",
		},
		Tags: []string{"npm", "test", "test_failure", "ci"},
	})

	// Docker failures
	samples = append(samples, BenchmarkSample{
		ID:       "docker-001",
		Command:  "docker build -t myapp .",
		ExitCode: 1,
		Stderr:   "Step 5/10 : RUN npm install\nERROR [5/10] RUN npm install\nfailed to solve with frontend dockerfile.v0",
		Context: SmartFallbackContext{
			CommandType:    "docker",
			Command:        "docker",
			Subcommand:     "build",
			HasDockerfile:  true,
			ErrorPattern:   "build_failure",
		},
		Category:    "docker",
		Description: "Docker build failure",
		GoldInsults: []string{
			"Docker build failed. Can't containerize disaster.",
			"FROM scratch. You are scratch.",
		},
		Tags: []string{"docker", "build", "build_failure"},
	})

	samples = append(samples, BenchmarkSample{
		ID:       "docker-002",
		Command:  "docker run -p 3000:3000 myapp",
		ExitCode: 125,
		Stderr:   "docker: Error response from daemon: driver failed programming external connectivity on endpoint\nError starting userland proxy: listen tcp4 0.0.0.0:3000: bind: address already in use.",
		Context: SmartFallbackContext{
			CommandType:  "docker",
			Command:      "docker",
			Subcommand:   "run",
			ErrorPattern: "port_in_use",
			NumericArgs:  []int{3000},
		},
		Category:    "docker",
		Description: "Port already in use",
		GoldInsults: []string{
			"Port 3000 already in use. By someone competent, probably.",
			"Port conflict. Your existence is a conflict.",
		},
		Tags: []string{"docker", "run", "network"},
	})

	// Python failures
	samples = append(samples, BenchmarkSample{
		ID:       "python-001",
		Command:  "python app.py",
		ExitCode: 1,
		Stderr:   "Traceback (most recent call last):\n  File \"app.py\", line 5, in <module>\n    import requests\nModuleNotFoundError: No module named 'requests'",
		Context: SmartFallbackContext{
			CommandType:  "python",
			Command:      "python",
			ProjectType:  "python",
			ErrorPattern: "dependency",
			FileExtensions: []string{".py"},
		},
		Category:    "python",
		Description: "Missing Python module",
		GoldInsults: []string{
			"ModuleNotFoundError: Module 'brain' not found",
			"Did you activate your venv? Don't answer, I know you didn't",
		},
		Tags: []string{"python", "dependency"},
	})

	samples = append(samples, BenchmarkSample{
		ID:       "python-002",
		Command:  "python script.py",
		ExitCode: 1,
		Stderr:   "  File \"script.py\", line 15\n    if x == 5\nSyntaxError: invalid syntax",
		Context: SmartFallbackContext{
			CommandType:    "python",
			Command:        "python",
			ProjectType:    "python",
			ErrorPattern:   "syntax_error",
			FileExtensions: []string{".py"},
		},
		Category:    "python",
		Description: "Python syntax error",
		GoldInsults: []string{
			"SyntaxError: Invalid syntax, invalid developer",
			"Python is trying to tell you something. Maybe listen for once?",
		},
		Tags: []string{"python", "syntax"},
	})

	// Rust failures
	samples = append(samples, BenchmarkSample{
		ID:       "rust-001",
		Command:  "cargo build",
		ExitCode: 101,
		Stderr:   "error[E0502]: cannot borrow `x` as mutable because it is also borrowed as immutable\n  --> src/main.rs:10:5",
		Context: SmartFallbackContext{
			CommandType:  "rust",
			Command:      "cargo",
			Subcommand:   "build",
			ProjectType:  "rust",
			ErrorPattern: "borrow_checker",
		},
		Category:    "rust",
		Description: "Borrow checker error",
		GoldInsults: []string{
			"Borrow checker says no. And honestly, it has a point.",
			"Fighting the borrow checker? The borrow checker always wins.",
		},
		Tags: []string{"rust", "build", "borrow_checker"},
	})

	// Permission errors
	samples = append(samples, BenchmarkSample{
		ID:       "perm-001",
		Command:  "chmod 777 /etc/passwd",
		ExitCode: 1,
		Stderr:   "chmod: changing permissions of '/etc/passwd': Operation not permitted",
		Context: SmartFallbackContext{
			Command:      "chmod",
			ErrorPattern: "permission_denied",
			NumericArgs:  []int{777},
		},
		Category:    "permission",
		Description: "Permission denied with chmod 777",
		GoldInsults: []string{
			"chmod 777 isn't the answer this time, though I admire your optimism",
			"777: Jackpot of incompetence",
		},
		Tags: []string{"permission", "chmod"},
	})

	// Late night scenarios
	samples = append(samples, BenchmarkSample{
		ID:       "time-001",
		Command:  "make build",
		ExitCode: 2,
		Stderr:   "make: *** [Makefile:15: build] Error 2",
		Context: SmartFallbackContext{
			Command:      "make",
			ErrorPattern: "build_failure",
			TimeOfDay:    3,
			HasMakefile:  true,
		},
		Category:    "build",
		Description: "Build failure at 3 AM",
		GoldInsults: []string{
			"It's 3 AM. The bugs aren't the only thing that needs fixing",
			"Late night debugging? Tomorrow-you is going to hate today-you",
		},
		Tags: []string{"build", "late_night"},
	})

	return samples
}

// EvaluateSystem runs the benchmark against a system
func (b *Benchmark) EvaluateSystem(system *EnsembleSystem) BenchmarkResults {
	results := BenchmarkResults{
		SystemName:     "Ensemble ML System",
		TotalSamples:   len(b.Samples),
		DetailedScores: make([]SampleScore, 0, len(b.Samples)),
	}

	var totalRelevance float64
	var totalLatency time.Duration
	var totalConfidence float64
	var fallbackCount int

	for _, sample := range b.Samples {
		start := time.Now()
		insult := system.GenerateInsult(&sample.Context, "sarcastic")
		latency := time.Since(start)

		// Calculate relevance score
		relevance := calculateRelevanceScore(sample, insult)

		// Determine if it was a Markov fallback
		isFallback := len(insult) > 0 && !containsInsult(system.database.Insults, insult)

		if isFallback {
			fallbackCount++
		}

		score := SampleScore{
			SampleID:        sample.ID,
			GeneratedInsult: insult,
			Relevance:       relevance,
			Latency:         latency,
			Confidence:      0.75, // Placeholder
			NoveltyScore:    1.0,
			Method:          determineMethod(isFallback),
		}

		results.DetailedScores = append(results.DetailedScores, score)

		totalRelevance += relevance
		totalLatency += latency
		totalConfidence += score.Confidence
	}

	results.AvgRelevance = totalRelevance / float64(len(b.Samples))
	results.AvgLatency = totalLatency / time.Duration(len(b.Samples))
	results.AvgConfidence = totalConfidence / float64(len(b.Samples))
	results.FallbackRate = float64(fallbackCount) / float64(len(b.Samples))
	results.DiversityScore = calculateDiversityScore(results.DetailedScores)

	return results
}

// calculateRelevanceScore measures how relevant the insult is to the error
func calculateRelevanceScore(sample BenchmarkSample, insult string) float64 {
	score := 0.0

	// Check for keyword matches
	keywords := extractKeywords(sample)
	for _, keyword := range keywords {
		if containsWord(insult, keyword) {
			score += 0.2
		}
	}

	// Check for tag matches
	for _, tag := range sample.Tags {
		if containsWord(insult, tag) {
			score += 0.15
		}
	}

	// Check similarity to gold insults
	if len(sample.GoldInsults) > 0 {
		maxSimilarity := 0.0
		for _, gold := range sample.GoldInsults {
			sim := simpleStringSimilarity(insult, gold)
			if sim > maxSimilarity {
				maxSimilarity = sim
			}
		}
		score += maxSimilarity * 0.3
	}

	return math.Min(1.0, score)
}

// extractKeywords extracts key terms from sample
func extractKeywords(sample BenchmarkSample) []string {
	keywords := []string{
		sample.Context.Command,
		sample.Context.Subcommand,
		sample.Context.CommandType,
		sample.Context.ErrorPattern,
	}

	if sample.Context.GitBranch != "" {
		keywords = append(keywords, sample.Context.GitBranch)
	}

	if sample.Context.ProjectType != "" {
		keywords = append(keywords, sample.Context.ProjectType)
	}

	return keywords
}

// containsWord checks if text contains word (case-insensitive)
func containsWord(text, word string) bool {
	textLower := toLower(text)
	wordLower := toLower(word)
	return contains(textLower, wordLower)
}

// simpleStringSimilarity calculates basic string similarity
func simpleStringSimilarity(s1, s2 string) float64 {
	// Simple word overlap metric
	words1 := splitWords(toLower(s1))
	words2 := splitWords(toLower(s2))

	if len(words1) == 0 || len(words2) == 0 {
		return 0.0
	}

	matches := 0
	for _, w1 := range words1 {
		for _, w2 := range words2 {
			if w1 == w2 && len(w1) > 2 { // Skip short words
				matches++
				break
			}
		}
	}

	return float64(matches) / float64(max(len(words1), len(words2)))
}

// calculateDiversityScore measures insult variety
func calculateDiversityScore(scores []SampleScore) float64 {
	if len(scores) < 2 {
		return 1.0
	}

	// Count unique insults
	unique := make(map[string]bool)
	for _, score := range scores {
		unique[score.GeneratedInsult] = true
	}

	return float64(len(unique)) / float64(len(scores))
}

// containsInsult checks if insult exists in database
func containsInsult(insults []TaggedInsult, target string) bool {
	for _, insult := range insults {
		if insult.Text == target {
			return true
		}
	}
	return false
}

// determineMethod identifies which method generated the insult
func determineMethod(isFallback bool) string {
	if isFallback {
		return "markov"
	}
	return "ensemble"
}

// PrintResults outputs benchmark results
func (r *BenchmarkResults) Print() {
	fmt.Println("╔═══════════════════════════════════════════════════════════╗")
	fmt.Printf("║ Benchmark Results: %-38s ║\n", r.SystemName)
	fmt.Println("╠═══════════════════════════════════════════════════════════╣")
	fmt.Printf("║ Total Samples:     %-41d ║\n", r.TotalSamples)
	fmt.Printf("║ Avg Relevance:     %-41.3f ║\n", r.AvgRelevance)
	fmt.Printf("║ Avg Latency:       %-41s ║\n", r.AvgLatency)
	fmt.Printf("║ Avg Confidence:    %-41.3f ║\n", r.AvgConfidence)
	fmt.Printf("║ Diversity Score:   %-41.3f ║\n", r.DiversityScore)
	fmt.Printf("║ Fallback Rate:     %-40.1f%% ║\n", r.FallbackRate*100)
	fmt.Println("╚═══════════════════════════════════════════════════════════╝")
}

// Helper functions
func toLower(s string) string {
	result := ""
	for _, r := range s {
		if r >= 'A' && r <= 'Z' {
			result += string(r + 32)
		} else {
			result += string(r)
		}
	}
	return result
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && findSubstring(s, substr) >= 0
}

func findSubstring(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func splitWords(s string) []string {
	var words []string
	var current string

	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			current += string(r)
		} else {
			if len(current) > 0 {
				words = append(words, current)
				current = ""
			}
		}
	}

	if len(current) > 0 {
		words = append(words, current)
	}

	return words
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
