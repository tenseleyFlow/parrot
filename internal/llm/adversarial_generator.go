package llm

import (
	"math"
	"math/rand"
	"strings"
	"time"
)

// AdversarialInsultGenerator implements a GAN-inspired system
// Generator creates insults, Critic scores them, iterative improvement
type AdversarialInsultGenerator struct {
	generator *InsultGenerator
	critic    *InsultCritic
	database  *InsultDatabase
	markov    *MarkovGenerator

	// Training state
	generatorScore float64
	criticScore    float64
	rounds         int

	// Quality thresholds
	minQuality      float64
	targetQuality   float64
	improvementRate float64
}

// InsultGenerator creates insult candidates
type InsultGenerator struct {
	templates      []string
	components     *ComponentLibrary
	markov         *MarkovGenerator
	creativityMode string // "safe", "balanced", "wild"
	rng            *rand.Rand
}

// InsultCritic scores insult quality
type InsultCritic struct {
	relevanceWeight  float64
	noveltyWeight    float64
	brutalityWeight  float64
	coherenceWeight  float64
	lengthWeight     float64

	seenInsults map[string]int // Track seen insults for novelty
}

// ComponentLibrary stores semantic building blocks
type ComponentLibrary struct {
	Subjects   []string // "Your code", "This commit", "Your docker build"
	Verbs      []string // "failed harder than", "crashed like", "died faster than"
	Objects    []string // "your career", "a burning dumpster", "your hopes"
	Punchlines []string // "combined", "on steroids", "in production"
	Intensifiers []string // "absolutely", "completely", "monumentally"
	Comparisons  []string // "than a", "like a", "as much as"
}

// InsultQualityScore represents multi-dimensional quality metrics
type InsultQualityScore struct {
	Relevance  float64 // How well it matches the context
	Novelty    float64 // How original/unique it is
	Brutality  float64 // How savage/brutal it is
	Coherence  float64 // How well it flows/makes sense
	Length     float64 // Appropriate length (not too long/short)
	Overall    float64 // Weighted combination
	Breakdown  string  // Human-readable explanation
}

// NewAdversarialInsultGenerator creates a GAN-inspired generator
func NewAdversarialInsultGenerator(db *InsultDatabase, markov *MarkovGenerator) *AdversarialInsultGenerator {
	components := initializeComponentLibrary()

	return &AdversarialInsultGenerator{
		generator: &InsultGenerator{
			templates:      initializeTemplates(),
			components:     components,
			markov:         markov,
			creativityMode: "balanced",
			rng:            rand.New(rand.NewSource(time.Now().UnixNano())),
		},
		critic: &InsultCritic{
			relevanceWeight: 0.30,
			noveltyWeight:   0.25,
			brutalityWeight: 0.20,
			coherenceWeight: 0.15,
			lengthWeight:    0.10,
			seenInsults:     make(map[string]int),
		},
		database:        db,
		markov:          markov,
		generatorScore:  0.5,
		criticScore:     0.5,
		rounds:          0,
		minQuality:      0.6,
		targetQuality:   0.8,
		improvementRate: 0.05,
	}
}

// Generate creates an insult using adversarial training
func (aig *AdversarialInsultGenerator) Generate(ctx *SmartFallbackContext, personality string) string {
	const maxAttempts = 5
	bestInsult := ""
	bestScore := 0.0

	// Generate multiple candidates and pick the best
	for attempt := 0; attempt < maxAttempts; attempt++ {
		// Generator creates candidate
		candidate := aig.generator.CreateCandidate(ctx, personality)

		// Critic scores it
		score := aig.critic.Score(candidate, ctx, personality)

		// Track best
		if score.Overall > bestScore {
			bestScore = score.Overall
			bestInsult = candidate
		}

		// If we hit target quality, stop early
		if score.Overall >= aig.targetQuality {
			break
		}
	}

	// Record this insult for novelty tracking
	aig.critic.seenInsults[bestInsult]++

	// Update adversarial scores (simplified GAN-like training)
	aig.updateScores(bestScore)

	// Only return if above minimum quality
	if bestScore >= aig.minQuality {
		return bestInsult
	}

	return "" // Not good enough
}

// GenerateBatch creates multiple candidates for ensemble selection
func (aig *AdversarialInsultGenerator) GenerateBatch(ctx *SmartFallbackContext, personality string, count int) []string {
	candidates := make([]string, 0, count)

	for i := 0; i < count; i++ {
		candidate := aig.generator.CreateCandidate(ctx, personality)
		score := aig.critic.Score(candidate, ctx, personality)

		if score.Overall >= aig.minQuality {
			candidates = append(candidates, candidate)
		}
	}

	return candidates
}

// CreateCandidate generates an insult candidate
func (ig *InsultGenerator) CreateCandidate(ctx *SmartFallbackContext, personality string) string {
	// Choose generation strategy based on creativity mode
	strategy := ig.rng.Float64()

	switch ig.creativityMode {
	case "safe":
		// Mostly template-based (70%), some Markov (30%)
		if strategy < 0.7 {
			return ig.generateFromTemplate(ctx, personality)
		}
		return ig.markov.Blend(ctx)

	case "wild":
		// Mostly Markov (60%), some composites (30%), some templates (10%)
		if strategy < 0.6 {
			return ig.markov.Blend(ctx)
		} else if strategy < 0.9 {
			return ig.generateComposite(ctx, personality)
		}
		return ig.generateFromTemplate(ctx, personality)

	default: // "balanced"
		// Mix of all: 40% template, 40% composite, 20% Markov
		if strategy < 0.4 {
			return ig.generateFromTemplate(ctx, personality)
		} else if strategy < 0.8 {
			return ig.generateComposite(ctx, personality)
		}
		return ig.markov.Blend(ctx)
	}
}

// generateFromTemplate uses traditional templates
func (ig *InsultGenerator) generateFromTemplate(ctx *SmartFallbackContext, personality string) string {
	template := ig.templates[ig.rng.Intn(len(ig.templates))]

	// Replace variables
	result := template
	result = strings.ReplaceAll(result, "{command}", ctx.Command)
	result = strings.ReplaceAll(result, "{commandType}", ctx.CommandType)
	result = strings.ReplaceAll(result, "{error}", ctx.ErrorPattern)
	result = strings.ReplaceAll(result, "{project}", ctx.ProjectType)

	return result
}

// generateComposite builds from semantic components
func (ig *InsultGenerator) generateComposite(ctx *SmartFallbackContext, personality string) string {
	comp := ig.components

	// Build: [Subject] [Verb] [Object] [Punchline]
	subject := comp.Subjects[ig.rng.Intn(len(comp.Subjects))]
	verb := comp.Verbs[ig.rng.Intn(len(comp.Verbs))]
	object := comp.Objects[ig.rng.Intn(len(comp.Objects))]

	// Contextualize subject
	subject = ig.contextualizeSubject(subject, ctx)

	// Sometimes add intensifier
	insult := subject + " " + verb + " " + object

	// Sometimes add punchline
	if ig.rng.Float64() < 0.5 {
		punchline := comp.Punchlines[ig.rng.Intn(len(comp.Punchlines))]
		insult += " " + punchline
	}

	// Ensure it ends with punctuation
	if !strings.HasSuffix(insult, ".") && !strings.HasSuffix(insult, "!") {
		insult += "."
	}

	return insult
}

// contextualizeSubject adapts subject to context
func (ig *InsultGenerator) contextualizeSubject(subject string, ctx *SmartFallbackContext) string {
	// Replace generic "Your code" with specific command references
	if ctx.Command != "" {
		subject = strings.ReplaceAll(subject, "Your code", "Your "+ctx.Command)
		subject = strings.ReplaceAll(subject, "This code", "This "+ctx.Command)
	}

	if ctx.ProjectType != "" {
		subject = strings.ReplaceAll(subject, "project", ctx.ProjectType+" project")
	}

	return subject
}

// Score evaluates insult quality
func (ic *InsultCritic) Score(insult string, ctx *SmartFallbackContext, personality string) InsultQualityScore {
	score := InsultQualityScore{}

	// 1. Relevance: Does it relate to the context?
	score.Relevance = ic.scoreRelevance(insult, ctx)

	// 2. Novelty: Is it unique/original?
	score.Novelty = ic.scoreNovelty(insult)

	// 3. Brutality: How savage is it? (adjust for personality)
	score.Brutality = ic.scoreBrutality(insult, personality)

	// 4. Coherence: Does it make sense?
	score.Coherence = ic.scoreCoherence(insult)

	// 5. Length: Is it appropriately sized?
	score.Length = ic.scoreLength(insult)

	// Calculate weighted overall score
	score.Overall = (score.Relevance * ic.relevanceWeight) +
		(score.Novelty * ic.noveltyWeight) +
		(score.Brutality * ic.brutalityWeight) +
		(score.Coherence * ic.coherenceWeight) +
		(score.Length * ic.lengthWeight)

	// Generate breakdown
	score.Breakdown = ic.generateBreakdown(score)

	return score
}

func (ic *InsultCritic) scoreRelevance(insult string, ctx *SmartFallbackContext) float64 {
	score := 0.0
	insultLower := strings.ToLower(insult)

	// Check if it mentions the command
	if ctx.Command != "" && strings.Contains(insultLower, strings.ToLower(ctx.Command)) {
		score += 0.3
	}

	// Check if it mentions the command type
	if ctx.CommandType != "" && strings.Contains(insultLower, strings.ToLower(ctx.CommandType)) {
		score += 0.2
	}

	// Check if it mentions the error pattern
	if ctx.ErrorPattern != "" {
		errorWords := strings.Split(ctx.ErrorPattern, "_")
		for _, word := range errorWords {
			if strings.Contains(insultLower, strings.ToLower(word)) {
				score += 0.2
				break
			}
		}
	}

	// Check if it mentions the project type
	if ctx.ProjectType != "" && strings.Contains(insultLower, strings.ToLower(ctx.ProjectType)) {
		score += 0.2
	}

	// Bonus for context-aware terms
	contextTerms := []string{"failed", "error", "broken", "crash", "disaster"}
	for _, term := range contextTerms {
		if strings.Contains(insultLower, term) {
			score += 0.1
			break
		}
	}

	return math.Min(1.0, score)
}

func (ic *InsultCritic) scoreNovelty(insult string) float64 {
	// Check how many times we've seen this exact insult
	useCount := ic.seenInsults[insult]

	// Novelty decays with use
	if useCount == 0 {
		return 1.0 // Completely novel
	} else if useCount == 1 {
		return 0.8
	} else if useCount < 5 {
		return 0.6
	} else if useCount < 10 {
		return 0.4
	}
	return 0.2 // Heavily overused
}

func (ic *InsultCritic) scoreBrutality(insult string, personality string) float64 {
	insultLower := strings.ToLower(insult)

	// Count brutal words
	brutalWords := []string{
		"disaster", "catastrophe", "failure", "incompetence", "terrible",
		"awful", "pathetic", "worthless", "garbage", "trash", "nightmare",
		"doomed", "hopeless", "broken", "destroyed", "ruined", "killed",
	}

	brutalCount := 0
	for _, word := range brutalWords {
		if strings.Contains(insultLower, word) {
			brutalCount++
		}
	}

	// Base brutality score
	brutality := math.Min(1.0, float64(brutalCount)*0.3)

	// Adjust for personality
	switch personality {
	case "mild":
		// Penalize high brutality
		if brutality > 0.5 {
			brutality *= 0.5
		}
	case "savage":
		// Reward high brutality
		if brutality < 0.5 {
			brutality *= 0.7
		}
	}

	return brutality
}

func (ic *InsultCritic) scoreCoherence(insult string) float64 {
	// Simple heuristics for coherence
	score := 1.0

	// Penalize if too many numbers (likely corrupted)
	digitCount := 0
	for _, r := range insult {
		if r >= '0' && r <= '9' {
			digitCount++
		}
	}
	if digitCount > len(insult)/4 {
		score *= 0.5
	}

	// Penalize if too many repeated words
	words := strings.Fields(insult)
	uniqueWords := make(map[string]bool)
	for _, word := range words {
		uniqueWords[strings.ToLower(word)] = true
	}
	if len(uniqueWords) < len(words)/2 {
		score *= 0.7
	}

	// Penalize if no spaces (likely corrupted)
	if !strings.Contains(insult, " ") {
		score *= 0.3
	}

	// Reward if it has proper punctuation
	if strings.HasSuffix(insult, ".") || strings.HasSuffix(insult, "!") {
		score *= 1.1
	}

	return math.Min(1.0, score)
}

func (ic *InsultCritic) scoreLength(insult string) float64 {
	length := len(insult)

	// Optimal length: 40-120 characters
	if length >= 40 && length <= 120 {
		return 1.0
	} else if length >= 20 && length < 40 {
		return 0.8
	} else if length > 120 && length <= 150 {
		return 0.8
	} else if length < 20 {
		return 0.5 // Too short
	} else {
		return 0.4 // Too long
	}
}

func (ic *InsultCritic) generateBreakdown(score InsultQualityScore) string {
	breakdown := ""
	breakdown += "Relevance: " + formatFloat(score.Relevance) + " "
	breakdown += "Novelty: " + formatFloat(score.Novelty) + " "
	breakdown += "Brutality: " + formatFloat(score.Brutality) + " "
	breakdown += "Coherence: " + formatFloat(score.Coherence) + " "
	breakdown += "Length: " + formatFloat(score.Length)
	return breakdown
}

// updateScores performs simplified GAN-like training
func (aig *AdversarialInsultGenerator) updateScores(qualityScore float64) {
	// Generator score: how well it fooled the critic (inverse relationship)
	aig.generatorScore = 0.7*aig.generatorScore + 0.3*qualityScore

	// Critic score: how well it evaluated quality
	aig.criticScore = 0.7*aig.criticScore + 0.3*(1.0-math.Abs(qualityScore-0.8))

	aig.rounds++

	// Adjust creativity based on performance
	if aig.generatorScore < 0.5 && aig.rounds > 10 {
		// Generator struggling - be more conservative
		aig.generator.creativityMode = "safe"
	} else if aig.generatorScore > 0.7 {
		// Generator doing well - get more creative
		aig.generator.creativityMode = "wild"
	} else {
		aig.generator.creativityMode = "balanced"
	}
}

// GetStats returns adversarial system statistics
func (aig *AdversarialInsultGenerator) GetStats() map[string]interface{} {
	return map[string]interface{}{
		"generator_score":  aig.generatorScore,
		"critic_score":     aig.criticScore,
		"training_rounds":  aig.rounds,
		"creativity_mode":  aig.generator.creativityMode,
		"unique_insults":   len(aig.critic.seenInsults),
	}
}

// Initialize component library
func initializeComponentLibrary() *ComponentLibrary {
	return &ComponentLibrary{
		Subjects: []string{
			"Your code", "This commit", "Your build", "This deployment",
			"Your docker container", "This merge", "Your test suite",
			"Your pull request", "This branch", "Your configuration",
			"This error", "Your logic", "This attempt", "Your workflow",
		},
		Verbs: []string{
			"failed harder than", "crashed like", "died faster than",
			"broke worse than", "destroyed", "ruined", "devastated",
			"collapsed like", "imploded faster than", "crumbled like",
			"exploded worse than", "failed as badly as", "bombed harder than",
		},
		Objects: []string{
			"your career prospects", "a burning dumpster fire", "your hopes and dreams",
			"your last relationship", "the Hindenburg", "your job security",
			"a house of cards", "your professional reputation", "your confidence",
			"your self-esteem", "a Jenga tower", "your sanity", "reality itself",
			"everyone's expectations", "your manager's patience",
		},
		Punchlines: []string{
			"combined", "on steroids", "in production", "multiplied",
			"simultaneously", "at scale", "under load", "in the worst way",
			"and then some", "by orders of magnitude", "exponentially",
		},
		Intensifiers: []string{
			"absolutely", "completely", "utterly", "monumentally",
			"catastrophically", "spectacularly", "magnificently",
		},
		Comparisons: []string{
			"than", "like", "worse than", "as much as", "more than",
		},
	}
}

// Initialize templates
func initializeTemplates() []string {
	return []string{
		"{command} failed: The {commandType} gods have rejected you.",
		"Your {command} encountered a {error}: Fix yourself first.",
		"{command} crashed harder than your career in {project}.",
		"Error in {command}: Expected competence, found disaster.",
		"Your {commandType} skills are as broken as this {command}.",
		"{command} failed. Your {project} project failed. You failed.",
		"The {error} error is just {command} protecting itself from you.",
		"Your {command} in {project} is a monument to failure.",
	}
}

func formatFloat(f float64) string {
	// Simple float formatter (0.00 format)
	intPart := int(f * 100)
	whole := intPart / 100
	frac := intPart % 100

	result := intToStr(whole) + "."
	if frac < 10 {
		result += "0"
	}
	result += intToStr(frac)

	return result
}
