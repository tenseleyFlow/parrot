package llm

import (
	"math/rand"
	"strings"
	"time"
)

// MarkovGenerator generates novel insults using Markov chains
type MarkovGenerator struct {
	chains      map[string]map[string]int // state -> next_word -> count
	starters    []string                   // possible starting words
	order       int                        // n-gram order (2 = bigram)
	minLength   int                        // minimum generated text length
	maxLength   int                        // maximum generated text length
	rng         *rand.Rand
}

// NewMarkovGenerator creates a new Markov chain generator
func NewMarkovGenerator(order int) *MarkovGenerator {
	return &MarkovGenerator{
		chains:    make(map[string]map[string]int),
		starters:  make([]string, 0),
		order:     order,
		minLength: 30,  // Minimum 30 characters
		maxLength: 150, // Maximum 150 characters
		rng:       rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Train trains the Markov chain on a corpus of insults
func (mg *MarkovGenerator) Train(insults []string) {
	for _, insult := range insults {
		mg.trainOnText(insult)
	}
}

// trainOnText trains on a single text
func (mg *MarkovGenerator) trainOnText(text string) {
	words := mg.tokenize(text)
	if len(words) < mg.order+1 {
		return
	}

	// Add first state as starter
	state := strings.Join(words[:mg.order], " ")
	mg.starters = append(mg.starters, state)

	// Build chain
	for i := 0; i < len(words)-mg.order; i++ {
		state := strings.Join(words[i:i+mg.order], " ")
		nextWord := words[i+mg.order]

		if _, exists := mg.chains[state]; !exists {
			mg.chains[state] = make(map[string]int)
		}

		mg.chains[state][nextWord]++
	}
}

// tokenize splits text into words
func (mg *MarkovGenerator) tokenize(text string) []string {
	// Split on spaces and punctuation, but keep punctuation
	var words []string
	var currentWord strings.Builder

	for _, r := range text {
		if r == ' ' || r == '\n' || r == '\t' {
			if currentWord.Len() > 0 {
				words = append(words, currentWord.String())
				currentWord.Reset()
			}
		} else if r == '.' || r == '!' || r == '?' || r == ',' || r == ':' || r == ';' {
			if currentWord.Len() > 0 {
				words = append(words, currentWord.String())
				currentWord.Reset()
			}
			words = append(words, string(r))
		} else {
			currentWord.WriteRune(r)
		}
	}

	if currentWord.Len() > 0 {
		words = append(words, currentWord.String())
	}

	return words
}

// Generate generates a novel insult
func (mg *MarkovGenerator) Generate() string {
	if len(mg.starters) == 0 || len(mg.chains) == 0 {
		return "" // Not trained yet
	}

	// Pick a random starting state
	state := mg.starters[mg.rng.Intn(len(mg.starters))]
	words := strings.Split(state, " ")

	// Generate until we hit max length or a terminal state
	attempts := 0
	maxAttempts := 100

	for len(strings.Join(words, " ")) < mg.maxLength && attempts < maxAttempts {
		attempts++

		// Get next word choices
		nextWords := mg.chains[state]
		if len(nextWords) == 0 {
			break // Terminal state
		}

		// Choose next word based on frequency
		nextWord := mg.weightedChoice(nextWords)
		words = append(words, nextWord)

		// Update state
		if len(words) >= mg.order {
			state = strings.Join(words[len(words)-mg.order:], " ")
		}

		// Stop at sentence endings if we've generated enough
		if (nextWord == "." || nextWord == "!" || nextWord == "?") &&
			len(strings.Join(words, " ")) >= mg.minLength {
			break
		}
	}

	// Reconstruct text with proper spacing
	return mg.reconstructText(words)
}

// weightedChoice selects a word based on frequency weights
func (mg *MarkovGenerator) weightedChoice(choices map[string]int) string {
	// Calculate total weight
	totalWeight := 0
	for _, count := range choices {
		totalWeight += count
	}

	// Random selection
	r := mg.rng.Intn(totalWeight)
	cumulative := 0

	for word, count := range choices {
		cumulative += count
		if r < cumulative {
			return word
		}
	}

	// Fallback (shouldn't reach here)
	for word := range choices {
		return word
	}

	return ""
}

// reconstructText reconstructs text with proper spacing around punctuation
func (mg *MarkovGenerator) reconstructText(words []string) string {
	var result strings.Builder

	for i, word := range words {
		// Don't add space before punctuation
		if i > 0 && !mg.isPunctuation(word) {
			result.WriteString(" ")
		}

		result.WriteString(word)
	}

	return result.String()
}

// isPunctuation checks if a word is punctuation
func (mg *MarkovGenerator) isPunctuation(word string) bool {
	return word == "." || word == "!" || word == "?" ||
		word == "," || word == ":" || word == ";" ||
		word == "(" || word == ")"
}

// GenerateContextual generates an insult with context hints
func (mg *MarkovGenerator) GenerateContextual(seedWords []string) string {
	if len(mg.chains) == 0 {
		return ""
	}

	// Find states that contain any of the seed words
	var matchingStarters []string
	for _, starter := range mg.starters {
		for _, seed := range seedWords {
			if strings.Contains(strings.ToLower(starter), strings.ToLower(seed)) {
				matchingStarters = append(matchingStarters, starter)
				break
			}
		}
	}

	// If we found matching starters, use them; otherwise use any starter
	if len(matchingStarters) == 0 {
		matchingStarters = mg.starters
	}

	// Pick a random matching starter
	state := matchingStarters[mg.rng.Intn(len(matchingStarters))]
	words := strings.Split(state, " ")

	// Generate as normal
	attempts := 0
	maxAttempts := 100

	for len(strings.Join(words, " ")) < mg.maxLength && attempts < maxAttempts {
		attempts++

		nextWords := mg.chains[state]
		if len(nextWords) == 0 {
			break
		}

		nextWord := mg.weightedChoice(nextWords)
		words = append(words, nextWord)

		if len(words) >= mg.order {
			state = strings.Join(words[len(words)-mg.order:], " ")
		}

		if (nextWord == "." || nextWord == "!" || nextWord == "?") &&
			len(strings.Join(words, " ")) >= mg.minLength {
			break
		}
	}

	return mg.reconstructText(words)
}

// GenerateWithTemplate generates using a template with variable slots
func (mg *MarkovGenerator) GenerateWithTemplate(template string, variables map[string]string) string {
	result := template

	for key, value := range variables {
		placeholder := "{" + key + "}"
		result = strings.ReplaceAll(result, placeholder, value)
	}

	// Fill remaining slots with Markov-generated content
	if strings.Contains(result, "{random}") {
		generated := mg.Generate()
		result = strings.ReplaceAll(result, "{random}", generated)
	}

	return result
}

// Blend creates a hybrid insult by blending Markov generation with templates
func (mg *MarkovGenerator) Blend(ctx *SmartFallbackContext) string {
	// Extract key terms from the context
	seedWords := []string{}

	// Add command type
	if ctx.CommandType != "" {
		seedWords = append(seedWords, ctx.CommandType)
	}

	// Add command
	if ctx.Command != "" {
		seedWords = append(seedWords, ctx.Command)
	}

	// Add error pattern
	if ctx.ErrorPattern != "" {
		seedWords = append(seedWords, strings.ReplaceAll(ctx.ErrorPattern, "_", " "))
	}

	// Generate contextual insult
	generated := mg.GenerateContextual(seedWords)

	// Post-process: ensure it's not too similar to training data
	if mg.tooSimilarToTraining(generated) {
		// Try again with different seed
		return mg.Generate()
	}

	return generated
}

// tooSimilarToTraining checks if generated text is too close to training data
func (mg *MarkovGenerator) tooSimilarToTraining(text string) bool {
	// Simple heuristic: if the text is very short or contains many consecutive
	// words from a single training example, it's too similar
	return len(text) < mg.minLength
}

// HybridGenerate combines Markov with template system for best results
func (mg *MarkovGenerator) HybridGenerate(
	ctx *SmartFallbackContext,
	templates []string,
) string {
	// 50% chance to use pure Markov, 50% template + Markov
	if mg.rng.Float64() < 0.5 {
		return mg.Blend(ctx)
	}

	// Pick a random template
	if len(templates) == 0 {
		return mg.Blend(ctx)
	}

	template := templates[mg.rng.Intn(len(templates))]

	// Fill template variables
	variables := map[string]string{
		"command":     ctx.Command,
		"commandType": ctx.CommandType,
		"exitCode":    string(rune(ctx.ExitCode)),
		"error":       ctx.ErrorPattern,
	}

	return mg.GenerateWithTemplate(template, variables)
}

// GetStats returns statistics about the trained model
func (mg *MarkovGenerator) GetStats() map[string]interface{} {
	return map[string]interface{}{
		"states":        len(mg.chains),
		"starters":      len(mg.starters),
		"order":         mg.order,
		"vocabulary":    mg.countVocabulary(),
		"avg_choices":   mg.averageChoices(),
	}
}

func (mg *MarkovGenerator) countVocabulary() int {
	vocab := make(map[string]bool)
	for state := range mg.chains {
		words := strings.Split(state, " ")
		for _, word := range words {
			vocab[word] = true
		}
	}
	return len(vocab)
}

func (mg *MarkovGenerator) averageChoices() float64 {
	if len(mg.chains) == 0 {
		return 0
	}

	total := 0
	for _, choices := range mg.chains {
		total += len(choices)
	}

	return float64(total) / float64(len(mg.chains))
}
