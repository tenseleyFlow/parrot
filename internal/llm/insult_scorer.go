package llm

import (
	"math"
	"sort"
	"strings"
)

// InsultScore represents a scored insult with relevance metrics
type InsultScore struct {
	Insult            TaggedInsult
	TotalScore        float64
	TagMatchScore     float64
	ErrorMatchScore   float64
	ContextScore      float64
	NoveltyScore      float64
	PersonalityScore  float64
}

// InsultScorer ranks insults based on multiple factors
type InsultScorer struct {
	database         *InsultDatabase
	errorClassifier  *ErrorClassifier
	intentParser     *IntentParser
	recentInsults    []string // Track recent insults to avoid repetition
	maxRecentHistory int
}

// ScoringWeights defines the weight of each factor in the final score
type ScoringWeights struct {
	TagMatch     float64 // How well tags match the context
	ErrorMatch   float64 // How well error type matches
	Context      float64 // Environmental context relevance
	Novelty      float64 // Avoid recent repetition
	Personality  float64 // Match personality preference
}

// DefaultWeights returns the default scoring weights
func DefaultWeights() ScoringWeights {
	return ScoringWeights{
		TagMatch:     0.35,
		ErrorMatch:   0.30,
		Context:      0.20,
		Novelty:      0.10,
		Personality:  0.05,
	}
}

// NewInsultScorer creates a new insult scorer
func NewInsultScorer(database *InsultDatabase) *InsultScorer {
	return &InsultScorer{
		database:         database,
		errorClassifier:  NewErrorClassifier(),
		intentParser:     NewIntentParser(),
		recentInsults:    make([]string, 0, 20),
		maxRecentHistory: 20,
	}
}

// ScoreAndRank analyzes context and returns top-ranked insults
func (is *InsultScorer) ScoreAndRank(
	ctx *SmartFallbackContext,
	personality string,
	topN int,
) []InsultScore {
	// Parse command intent
	intent := is.intentParser.ParseIntent(ctx.FullCommand)

	// Classify error type
	errorCategories := is.errorClassifier.ClassifyError(
		ctx.FullCommand,
		ctx.ExitCode,
		ctx.ErrorPattern,
	)

	// Generate contextual tags
	contextTags := ContextualTags(ctx, intent)

	// Convert error categories to tags
	errorTags := errorCategoriesToTags(errorCategories)

	// Get all relevant insults
	allInsults := is.database.Insults

	// Score each insult
	scores := make([]InsultScore, 0, len(allInsults))
	weights := DefaultWeights()

	for _, insult := range allInsults {
		score := InsultScore{
			Insult: insult,
		}

		// 1. Tag matching score (35%)
		score.TagMatchScore = is.calculateTagMatchScore(
			insult.Tags,
			contextTags,
			errorTags,
		)

		// 2. Error type matching (30%)
		score.ErrorMatchScore = is.calculateErrorMatchScore(
			insult.Tags,
			errorTags,
		)

		// 3. Context relevance (20%)
		score.ContextScore = is.calculateContextScore(
			insult,
			ctx,
			intent,
		)

		// 4. Novelty score (10%) - penalize recent insults
		score.NoveltyScore = is.calculateNoveltyScore(insult.Text)

		// 5. Personality match (5%)
		score.PersonalityScore = is.calculatePersonalityScore(
			insult,
			personality,
		)

		// Calculate weighted total
		score.TotalScore = (score.TagMatchScore * weights.TagMatch) +
			(score.ErrorMatchScore * weights.ErrorMatch) +
			(score.ContextScore * weights.Context) +
			(score.NoveltyScore * weights.Novelty) +
			(score.PersonalityScore * weights.Personality)

		// Apply base weight from insult
		score.TotalScore *= insult.Weight

		scores = append(scores, score)
	}

	// Sort by score descending
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].TotalScore > scores[j].TotalScore
	})

	// Return top N
	if len(scores) > topN {
		scores = scores[:topN]
	}

	return scores
}

// calculateTagMatchScore measures how well insult tags match context
func (is *InsultScorer) calculateTagMatchScore(
	insultTags []InsultTag,
	contextTags []InsultTag,
	errorTags []InsultTag,
) float64 {
	if len(contextTags) == 0 && len(errorTags) == 0 {
		return 0.5 // Neutral score if no context
	}

	allSearchTags := append(contextTags, errorTags...)
	matches := 0
	totalSearchTags := len(allSearchTags)

	for _, searchTag := range allSearchTags {
		for _, insultTag := range insultTags {
			if searchTag == insultTag {
				matches++
				break
			}
		}
	}

	// Calculate percentage of search tags that matched
	score := float64(matches) / float64(totalSearchTags)

	// Bonus for multiple matches
	if matches > 2 {
		score = math.Min(1.0, score*1.2)
	}

	return score
}

// calculateErrorMatchScore focuses specifically on error type matching
func (is *InsultScorer) calculateErrorMatchScore(
	insultTags []InsultTag,
	errorTags []InsultTag,
) float64 {
	if len(errorTags) == 0 {
		return 0.5 // Neutral if no specific error
	}

	matches := 0
	for _, errorTag := range errorTags {
		for _, insultTag := range insultTags {
			if errorTag == insultTag {
				matches++
			}
		}
	}

	// Strong match for error-specific insults
	if matches > 0 {
		return math.Min(1.0, float64(matches)/float64(len(errorTags))*1.5)
	}

	return 0.0
}

// calculateContextScore evaluates environmental and situational relevance
func (is *InsultScorer) calculateContextScore(
	insult TaggedInsult,
	ctx *SmartFallbackContext,
	intent CommandIntent,
) float64 {
	score := 0.5 // Base score

	// Time relevance
	hour := ctx.TimeOfDay
	if (hour >= 22 || hour <= 4) && hasTag(insult.Tags, TagLateNight) {
		score += 0.2
	}

	// CI context
	if ctx.IsCI && hasTag(insult.Tags, TagCI) {
		score += 0.2
	}

	// Main branch failures (more serious)
	if (ctx.GitBranch == "main" || ctx.GitBranch == "master") &&
	   hasTag(insult.Tags, TagMainBranch) {
		score += 0.15
	}

	// Repeated failures
	if ctx.IsRepeatedFailure && hasTag(insult.Tags, TagRepeated) {
		score += 0.15
	}

	// Complexity match
	if intent.Complexity == "simple" && hasTag(insult.Tags, TagSimple) {
		score += 0.1
	} else if intent.Complexity == "complex" && hasTag(insult.Tags, TagComplex) {
		score += 0.1
	}

	// High risk operations
	if intent.RiskLevel == "high" &&
	   (hasTag(insult.Tags, TagProduction) || hasTag(insult.Tags, TagMainBranch)) {
		score += 0.15
	}

	return math.Min(1.0, score)
}

// calculateNoveltyScore penalizes recently shown insults
func (is *InsultScorer) calculateNoveltyScore(insultText string) float64 {
	// Check if this insult was shown recently
	for i, recent := range is.recentInsults {
		if recent == insultText {
			// More recent = lower score
			recency := float64(len(is.recentInsults)-i) / float64(len(is.recentInsults))
			return 1.0 - (recency * 0.8) // Up to 80% penalty
		}
	}

	return 1.0 // Full score for novel insults
}

// calculatePersonalityScore ensures insult matches personality setting
func (is *InsultScorer) calculatePersonalityScore(
	insult TaggedInsult,
	personality string,
) float64 {
	switch personality {
	case "mild":
		if hasTag(insult.Tags, TagMild) {
			return 1.0
		}
		if insult.Severity <= 4 {
			return 0.8
		}
		return 0.3

	case "sarcastic":
		if hasTag(insult.Tags, TagSarcastic) {
			return 1.0
		}
		if insult.Severity >= 4 && insult.Severity <= 7 {
			return 0.8
		}
		return 0.5

	case "savage":
		if hasTag(insult.Tags, TagSavage) {
			return 1.0
		}
		if insult.Severity >= 6 {
			return 0.8
		}
		return 0.4

	default:
		return 0.7 // Neutral score
	}
}

// RecordShownInsult adds an insult to recent history
func (is *InsultScorer) RecordShownInsult(insultText string) {
	is.recentInsults = append(is.recentInsults, insultText)

	// Keep only recent N insults
	if len(is.recentInsults) > is.maxRecentHistory {
		is.recentInsults = is.recentInsults[1:]
	}
}

// GetBestInsult returns the top-ranked insult
func (is *InsultScorer) GetBestInsult(
	ctx *SmartFallbackContext,
	personality string,
) string {
	scores := is.ScoreAndRank(ctx, personality, 5)

	if len(scores) == 0 {
		return "Something went wrong. How ironic."
	}

	bestInsult := scores[0].Insult.Text
	is.RecordShownInsult(bestInsult)

	return bestInsult
}

// GetTopInsults returns multiple top candidates (useful for variety)
func (is *InsultScorer) GetTopInsults(
	ctx *SmartFallbackContext,
	personality string,
	count int,
) []string {
	scores := is.ScoreAndRank(ctx, personality, count*2)

	results := make([]string, 0, count)
	for i := 0; i < count && i < len(scores); i++ {
		results = append(results, scores[i].Insult.Text)
	}

	return results
}

// Helper functions

func hasTag(tags []InsultTag, tag InsultTag) bool {
	for _, t := range tags {
		if t == tag {
			return true
		}
	}
	return false
}

func errorCategoriesToTags(categories []ErrorCategory) []InsultTag {
	tags := make([]InsultTag, 0, len(categories))
	for _, cat := range categories {
		tag := errorCategoryToTag(cat)
		if tag != "" {
			tags = append(tags, InsultTag(tag))
		}
	}
	return tags
}

func errorCategoryToTag(category ErrorCategory) string {
	mapping := map[ErrorCategory]string{
		ErrorPermission:      "permission",
		ErrorSyntax:          "syntax",
		ErrorNetwork:         "network",
		ErrorDependency:      "dependency",
		ErrorMergeConflict:   "merge_conflict",
		ErrorTestFailure:     "test_failure",
		ErrorBuildFailure:    "build_failure",
		ErrorTimeout:         "timeout",
		ErrorAuthentication:  "authentication",
		ErrorDiskSpace:       "disk_space",
		ErrorMemory:          "memory",
		ErrorSegfault:        "segfault",
		ErrorLinting:         "linting",
		ErrorTypeMismatch:    "typing",
		ErrorDeprecated:      "deprecated",
	}

	if tag, exists := mapping[category]; exists {
		return tag
	}
	return ""
}

// AnalyzeScoring provides detailed scoring breakdown (useful for debugging)
func (is *InsultScorer) AnalyzeScoring(
	ctx *SmartFallbackContext,
	personality string,
	topN int,
) []InsultScore {
	return is.ScoreAndRank(ctx, personality, topN)
}

// ClearHistory resets the recent insults history
func (is *InsultScorer) ClearHistory() {
	is.recentInsults = make([]string, 0, is.maxRecentHistory)
}

// LevenshteinDistance calculates similarity between strings (for fuzzy matching)
func levenshteinDistance(s1, s2 string) int {
	s1Lower := strings.ToLower(s1)
	s2Lower := strings.ToLower(s2)

	m := len(s1Lower)
	n := len(s2Lower)

	if m == 0 {
		return n
	}
	if n == 0 {
		return m
	}

	d := make([][]int, m+1)
	for i := range d {
		d[i] = make([]int, n+1)
		d[i][0] = i
	}
	for j := range d[0] {
		d[0][j] = j
	}

	for j := 1; j <= n; j++ {
		for i := 1; i <= m; i++ {
			cost := 1
			if s1Lower[i-1] == s2Lower[j-1] {
				cost = 0
			}
			d[i][j] = min(
				d[i-1][j]+1,      // deletion
				d[i][j-1]+1,      // insertion
				d[i-1][j-1]+cost, // substitution
			)
		}
	}

	return d[m][n]
}

func min(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}
