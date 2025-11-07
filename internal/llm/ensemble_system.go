package llm

import (
	"math"
	"sort"
)

// EnsembleSystem combines multiple ML techniques for optimal insult selection
type EnsembleSystem struct {
	tfidfEngine      *TFIDFEngine
	bm25Engine       *BM25Engine  // NEW: Industry-standard BM25 ranking
	markovGen        *MarkovGenerator
	insultScorer     *InsultScorer
	database         *InsultDatabase
	history          *InsultHistory

	// Ensemble weights
	semanticWeight   float64
	tagWeight        float64
	markovWeight     float64
	historicalWeight float64

	// Quality thresholds
	minSemanticScore  float64
	minTagScore       float64
	minEnsembleScore  float64

	// Configuration
	useBM25    bool  // Use BM25 instead of TF-IDF (recommended)
	trained    bool  // Training state
}

// EnsembleScore represents a comprehensive scoring of an insult candidate
type EnsembleScore struct {
	Insult           string
	SemanticScore    float64 // TF-IDF cosine similarity
	TagScore         float64 // Tag-based matching
	HistoricalScore  float64 // Historical pattern matching
	NoveltyScore     float64 // Avoid repetition
	PersonalityScore float64 // Personality fit
	EnsembleScore    float64 // Weighted combination
	Confidence       float64 // Confidence calibration
	Source           string  // "semantic", "tag", "markov", "ensemble"
}

// NewEnsembleSystem creates a new ensemble learning system
func NewEnsembleSystem(db *InsultDatabase, scorer *InsultScorer, hist *InsultHistory) *EnsembleSystem {
	return &EnsembleSystem{
		tfidfEngine:      NewTFIDFEngine(),
		bm25Engine:       NewBM25Engine(),
		markovGen:        NewMarkovGenerator(2), // Bigram model
		insultScorer:     scorer,
		database:         db,
		history:          hist,

		// Default ensemble weights (can be tuned)
		semanticWeight:   0.35,
		tagWeight:        0.30,
		markovWeight:     0.20,
		historicalWeight: 0.15,

		// Quality thresholds
		minSemanticScore: 0.25,
		minTagScore:      0.30,
		minEnsembleScore: 0.40,

		// Use BM25 by default (proven better than TF-IDF)
		useBM25: true,
		trained: false,
	}
}

// Train trains all ML components on the insult database
func (es *EnsembleSystem) Train() {
	if es.trained {
		return // Already trained
	}

	// Collect all insult texts
	insults := make([]string, 0, len(es.database.Insults))
	for _, insult := range es.database.Insults {
		insults = append(insults, insult.Text)
	}

	// Train TF-IDF engine
	es.tfidfEngine.BuildCorpus(insults)

	// Train BM25 engine (improved ranking algorithm)
	es.bm25Engine.BuildCorpus(insults)

	// Train Markov generator
	es.markovGen.Train(insults)

	es.trained = true
}

// GenerateInsult generates the best possible insult using ensemble methods
func (es *EnsembleSystem) GenerateInsult(
	ctx *SmartFallbackContext,
	personality string,
) string {
	// Ensure training is done
	if !es.trained {
		es.Train()
	}

	// Get candidates from multiple sources
	candidates := es.getAllCandidates(ctx, personality)

	if len(candidates) == 0 {
		// Last resort: generate using Markov
		return es.markovGen.Blend(ctx)
	}

	// Sort by ensemble score
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].EnsembleScore > candidates[j].EnsembleScore
	})

	// Get best candidate
	best := candidates[0]

	// If best score is still low, try Markov generation
	if best.EnsembleScore < es.minEnsembleScore {
		markovInsult := es.markovGen.Blend(ctx)
		if markovInsult != "" && len(markovInsult) > 20 {
			// Record and return Markov-generated insult
			es.history.RecordInsult(markovInsult, ctx.FullCommand, 0.5)
			return markovInsult
		}
	}

	// Record selected insult
	es.history.RecordInsult(best.Insult, ctx.FullCommand, best.EnsembleScore)

	return best.Insult
}

// getAllCandidates gets scored candidates from all sources
func (es *EnsembleSystem) getAllCandidates(
	ctx *SmartFallbackContext,
	personality string,
) []EnsembleScore {
	candidates := make([]EnsembleScore, 0, len(es.database.Insults))

	// Score all insults in database using ensemble
	for _, insult := range es.database.Insults {
		score := es.scoreInsult(insult, ctx, personality)

		// Only include if above minimum thresholds
		if score.EnsembleScore >= es.minEnsembleScore {
			candidates = append(candidates, score)
		}
	}

	return candidates
}

// scoreInsult scores a single insult using ensemble methods
func (es *EnsembleSystem) scoreInsult(
	insult TaggedInsult,
	ctx *SmartFallbackContext,
	personality string,
) EnsembleScore {
	score := EnsembleScore{
		Insult: insult.Text,
		Source: "ensemble",
	}

	// 1. Semantic similarity score (TF-IDF)
	score.SemanticScore = es.calculateSemanticScore(ctx, insult)

	// 2. Tag-based score (existing system)
	score.TagScore = es.calculateTagScore(ctx, insult)

	// 3. Historical pattern score
	score.HistoricalScore = es.calculateHistoricalScore(ctx, insult)

	// 4. Novelty score (avoid repetition)
	score.NoveltyScore = es.history.GetNoveltyScore(insult.Text)

	// 5. Personality fit score
	score.PersonalityScore = es.calculatePersonalityScore(insult, personality)

	// Calculate weighted ensemble score
	score.EnsembleScore = (score.SemanticScore * es.semanticWeight) +
		(score.TagScore * es.tagWeight) +
		(score.HistoricalScore * es.historicalWeight) +
		(score.NoveltyScore * 0.10) +
		(score.PersonalityScore * 0.05)

	// Apply insult base weight
	score.EnsembleScore *= insult.Weight

	// Calculate confidence (how much methods agree)
	score.Confidence = es.calculateConfidence(score)

	// Boost score if high confidence
	if score.Confidence > 0.8 {
		score.EnsembleScore *= 1.1
	}

	return score
}

// calculateSemanticScore uses BM25 or TF-IDF for semantic similarity
func (es *EnsembleSystem) calculateSemanticScore(
	ctx *SmartFallbackContext,
	insult TaggedInsult,
) float64 {
	// Create a rich context description
	contextText := es.buildContextText(ctx)

	var score float64

	if es.useBM25 {
		// Use BM25 (industry standard, proven better)
		// BM25 scores are typically in range 0-10, normalize to 0-1
		rawScore := es.bm25Engine.Score(contextText, insult.Text)
		score = math.Min(rawScore/10.0, 1.0)
	} else {
		// Use TF-IDF (for comparison)
		similarity := es.tfidfEngine.CalculateSemanticScore(contextText, insult.Text)
		score = sigmoid(similarity * 2.0)
	}

	return score
}

// buildContextText creates rich text representation of context
func (es *EnsembleSystem) buildContextText(ctx *SmartFallbackContext) string {
	var parts []string

	// Add command and type
	parts = append(parts, ctx.FullCommand)
	parts = append(parts, ctx.CommandType)
	parts = append(parts, ctx.Command)

	// Add error pattern
	if ctx.ErrorPattern != "" {
		parts = append(parts, ctx.ErrorPattern)
	}

	// Add project type
	if ctx.ProjectType != "" {
		parts = append(parts, ctx.ProjectType)
	}

	// Add git branch
	if ctx.GitBranch != "" {
		parts = append(parts, ctx.GitBranch)
	}

	// Add time context
	if ctx.TimeOfDay >= 22 || ctx.TimeOfDay <= 4 {
		parts = append(parts, "late night coding")
	}

	// Add CI context
	if ctx.IsCI {
		parts = append(parts, "continuous integration", "ci pipeline")
	}

	// Add repeated failure context
	if ctx.IsRepeatedFailure {
		parts = append(parts, "repeated failure", "again", "still failing")
	}

	return join(parts, " ")
}

// calculateTagScore uses the existing tag-based system
func (es *EnsembleSystem) calculateTagScore(
	ctx *SmartFallbackContext,
	insult TaggedInsult,
) float64 {
	// Parse intent
	parser := NewIntentParser()
	intent := parser.ParseIntent(ctx.FullCommand)

	// Generate contextual tags
	contextTags := ContextualTags(ctx, intent)

	// Classify error
	classifier := NewErrorClassifier()
	errorCategories := classifier.ClassifyError(ctx.FullCommand, ctx.ExitCode, ctx.ErrorPattern)
	errorTags := errorCategoriesToTags(errorCategories)

	// Combine tags
	allTags := append(contextTags, errorTags...)

	// Count matches
	matches := 0
	for _, contextTag := range allTags {
		for _, insultTag := range insult.Tags {
			if contextTag == insultTag {
				matches++
			}
		}
	}

	if len(allTags) == 0 {
		return 0.5
	}

	// Calculate match ratio
	score := float64(matches) / float64(len(allTags))

	// Bonus for multiple matches
	if matches > 2 {
		score = math.Min(1.0, score*1.2)
	}

	return score
}

// calculateHistoricalScore uses historical patterns
func (es *EnsembleSystem) calculateHistoricalScore(
	ctx *SmartFallbackContext,
	insult TaggedInsult,
) float64 {
	// Check if similar commands have been failed before
	// For now, use a simple heuristic based on command type

	baseScore := 0.5

	// Boost for matching command type
	for _, tag := range insult.Tags {
		if string(tag) == ctx.CommandType {
			baseScore += 0.2
		}
	}

	// Boost for matching error pattern
	if ctx.ErrorPattern != "" {
		for _, tag := range insult.Tags {
			if string(tag) == ctx.ErrorPattern {
				baseScore += 0.3
			}
		}
	}

	return math.Min(1.0, baseScore)
}

// calculatePersonalityScore ensures insult matches personality
func (es *EnsembleSystem) calculatePersonalityScore(
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
		return 0.7
	}
}

// calculateConfidence measures how much different methods agree
func (es *EnsembleSystem) calculateConfidence(score EnsembleScore) float64 {
	scores := []float64{
		score.SemanticScore,
		score.TagScore,
		score.HistoricalScore,
		score.NoveltyScore,
		score.PersonalityScore,
	}

	// Calculate variance
	mean := 0.0
	for _, s := range scores {
		mean += s
	}
	mean /= float64(len(scores))

	variance := 0.0
	for _, s := range scores {
		variance += (s - mean) * (s - mean)
	}
	variance /= float64(len(scores))

	// Low variance = high confidence (methods agree)
	// Convert variance to confidence (0-1)
	confidence := 1.0 - math.Min(variance*4.0, 1.0)

	return confidence
}

// GenerateMarkovInsult generates a novel insult using Markov chains
func (es *EnsembleSystem) GenerateMarkovInsult(ctx *SmartFallbackContext) string {
	if !es.trained {
		es.Train()
	}

	return es.markovGen.Blend(ctx)
}

// AnalyzeScoring provides detailed scoring breakdown for debugging
func (es *EnsembleSystem) AnalyzeScoring(
	ctx *SmartFallbackContext,
	personality string,
	topN int,
) []EnsembleScore {
	if !es.trained {
		es.Train()
	}

	candidates := es.getAllCandidates(ctx, personality)

	// Sort by ensemble score
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].EnsembleScore > candidates[j].EnsembleScore
	})

	if len(candidates) > topN {
		candidates = candidates[:topN]
	}

	return candidates
}

// UpdateWeights allows dynamic weight tuning based on feedback
func (es *EnsembleSystem) UpdateWeights(
	semanticW, tagW, markovW, historicalW float64,
) {
	total := semanticW + tagW + markovW + historicalW

	es.semanticWeight = semanticW / total
	es.tagWeight = tagW / total
	es.markovWeight = markovW / total
	es.historicalWeight = historicalW / total
}

// GetStats returns ensemble system statistics
func (es *EnsembleSystem) GetStats() map[string]interface{} {
	stats := make(map[string]interface{})

	stats["trained"] = es.trained
	stats["database_size"] = len(es.database.Insults)

	if es.trained {
		stats["tfidf_vocabulary"] = len(es.tfidfEngine.vocabulary)
		stats["markov_stats"] = es.markovGen.GetStats()
	}

	stats["weights"] = map[string]float64{
		"semantic":   es.semanticWeight,
		"tag":        es.tagWeight,
		"markov":     es.markovWeight,
		"historical": es.historicalWeight,
	}

	return stats
}

// Helper functions

func sigmoid(x float64) float64 {
	return 1.0 / (1.0 + math.Exp(-x))
}

func join(parts []string, sep string) string {
	result := ""
	for i, part := range parts {
		if i > 0 {
			result += sep
		}
		result += part
	}
	return result
}
