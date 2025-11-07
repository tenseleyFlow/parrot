package llm

import (
	"math"
)

// BM25Engine implements BM25 ranking algorithm (superior to basic TF-IDF)
// BM25 is the industry standard for text search and ranking
type BM25Engine struct {
	vocabulary    map[string]int      // word -> index
	idf           map[string]float64  // word -> inverse document frequency
	docLengths    []int               // document lengths
	avgDocLength  float64             // average document length
	documentCount int

	// BM25 parameters (tunable)
	k1 float64 // term frequency saturation parameter (typical: 1.2-2.0)
	b  float64 // document length normalization (typical: 0.75)

	ngramRange [2]int // min and max n-gram size
}

// NewBM25Engine creates a new BM25 ranking engine
func NewBM25Engine() *BM25Engine {
	return &BM25Engine{
		vocabulary:    make(map[string]int),
		idf:           make(map[string]float64),
		docLengths:    make([]int, 0),
		documentCount: 0,

		// Standard BM25 parameters (Okapi BM25)
		k1: 1.5,  // Typical range: 1.2-2.0
		b:  0.75, // Typical range: 0.5-0.9

		ngramRange: [2]int{1, 3}, // unigrams, bigrams, trigrams
	}
}

// SetParameters allows tuning of BM25 parameters
func (engine *BM25Engine) SetParameters(k1, b float64) {
	engine.k1 = k1
	engine.b = b
}

// BuildCorpus builds the BM25 corpus from documents
func (engine *BM25Engine) BuildCorpus(documents []string) {
	// First pass: extract terms and calculate document frequencies
	documentFreq := make(map[string]int)
	engine.docLengths = make([]int, len(documents))
	totalLength := 0

	for docIdx, doc := range documents {
		tokens := engine.extractNGrams(doc)
		engine.docLengths[docIdx] = len(tokens)
		totalLength += len(tokens)

		// Track which terms appear in this document
		seen := make(map[string]bool)
		for _, token := range tokens {
			if !seen[token] {
				documentFreq[token]++
				seen[token] = true
			}

			if _, exists := engine.vocabulary[token]; !exists {
				engine.vocabulary[token] = len(engine.vocabulary)
			}
		}
	}

	engine.documentCount = len(documents)
	engine.avgDocLength = float64(totalLength) / float64(engine.documentCount)

	// Calculate IDF for each term using BM25 IDF formula
	// IDF = log((N - df + 0.5) / (df + 0.5) + 1)
	// This is the Robertson-Sparck Jones formula
	for term, df := range documentFreq {
		N := float64(engine.documentCount)
		numerator := N - float64(df) + 0.5
		denominator := float64(df) + 0.5
		engine.idf[term] = math.Log((numerator / denominator) + 1.0)
	}
}

// extractNGrams extracts n-grams from text (same as TF-IDF)
func (engine *BM25Engine) extractNGrams(text string) []string {
	text = toLowerSimple(text)
	words := engine.tokenize(text)

	var ngrams []string

	for n := engine.ngramRange[0]; n <= engine.ngramRange[1]; n++ {
		if n > len(words) {
			break
		}

		for i := 0; i <= len(words)-n; i++ {
			ngram := joinWords(words[i:i+n], " ")
			ngrams = append(ngrams, ngram)
		}
	}

	return ngrams
}

// tokenize splits text into words
func (engine *BM25Engine) tokenize(text string) []string {
	var words []string
	var currentWord string

	for _, r := range text {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			currentWord += string(r)
		} else {
			if len(currentWord) > 1 { // Skip single characters
				words = append(words, currentWord)
			}
			currentWord = ""
		}
	}

	if len(currentWord) > 1 {
		words = append(words, currentWord)
	}

	return words
}

// Score calculates BM25 score for a query against a document
func (engine *BM25Engine) Score(query string, document string) float64 {
	queryTerms := engine.extractNGrams(query)
	docTerms := engine.extractNGrams(document)

	// Calculate term frequencies in document
	termFreq := make(map[string]int)
	for _, term := range docTerms {
		termFreq[term]++
	}

	docLength := len(docTerms)

	// Calculate BM25 score
	score := 0.0

	// Track query terms we've processed (unique terms only)
	seenQuery := make(map[string]bool)

	for _, queryTerm := range queryTerms {
		if seenQuery[queryTerm] {
			continue
		}
		seenQuery[queryTerm] = true

		// Get IDF for this term
		idf, exists := engine.idf[queryTerm]
		if !exists {
			// Term not in vocabulary - use a small IDF
			idf = math.Log(float64(engine.documentCount) + 1.0)
		}

		// Get term frequency in document
		tf := float64(termFreq[queryTerm])

		// BM25 formula:
		// score = IDF(qi) × (f(qi, D) × (k1 + 1)) / (f(qi, D) + k1 × (1 - b + b × |D| / avgdl))

		numerator := tf * (engine.k1 + 1.0)
		denominator := tf + engine.k1*(1.0-engine.b+engine.b*float64(docLength)/engine.avgDocLength)

		termScore := idf * (numerator / denominator)
		score += termScore
	}

	return score
}

// ScoreMultiple scores a query against multiple documents
func (engine *BM25Engine) ScoreMultiple(query string, documents []string) []BM25Score {
	scores := make([]BM25Score, len(documents))

	for i, doc := range documents {
		scores[i] = BM25Score{
			Index:    i,
			Document: doc,
			Score:    engine.Score(query, doc),
		}
	}

	// Sort by score descending
	sortBM25Scores(scores)

	return scores
}

// FindTopK returns the top K documents for a query
func (engine *BM25Engine) FindTopK(query string, documents []string, k int) []BM25Score {
	scores := engine.ScoreMultiple(query, documents)

	if len(scores) > k {
		scores = scores[:k]
	}

	return scores
}

// BM25Score represents a scored document
type BM25Score struct {
	Index    int
	Document string
	Score    float64
}

// sortBM25Scores sorts scores in descending order (bubble sort for simplicity)
func sortBM25Scores(scores []BM25Score) {
	n := len(scores)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if scores[j].Score < scores[j+1].Score {
				scores[j], scores[j+1] = scores[j+1], scores[j]
			}
		}
	}
}

// Explanation generates a human-readable explanation of the score
func (engine *BM25Engine) Explanation(query string, document string) string {
	queryTerms := engine.extractNGrams(query)
	docTerms := engine.extractNGrams(document)

	termFreq := make(map[string]int)
	for _, term := range docTerms {
		termFreq[term]++
	}

	docLength := len(docTerms)

	explanation := "BM25 Score Breakdown:\n"
	explanation += "=====================\n\n"

	totalScore := 0.0
	seenQuery := make(map[string]bool)

	for _, queryTerm := range queryTerms {
		if seenQuery[queryTerm] {
			continue
		}
		seenQuery[queryTerm] = true

		if termFreq[queryTerm] == 0 {
			continue // Term not in document
		}

		idf := engine.idf[queryTerm]
		tf := float64(termFreq[queryTerm])

		numerator := tf * (engine.k1 + 1.0)
		denominator := tf + engine.k1*(1.0-engine.b+engine.b*float64(docLength)/engine.avgDocLength)
		termScore := idf * (numerator / denominator)

		totalScore += termScore

		explanation += formatString("Term: '%s'\n", queryTerm)
		explanation += formatString("  TF: %d (occurs %d times)\n", termFreq[queryTerm], termFreq[queryTerm])
		explanation += formatString("  IDF: %.4f\n", idf)
		explanation += formatString("  BM25 component: %.4f\n", termScore)
		explanation += "\n"
	}

	explanation += formatString("Total BM25 Score: %.4f\n", totalScore)
	explanation += formatString("Document length: %d (avg: %.1f)\n", docLength, engine.avgDocLength)

	return explanation
}

// CompareWithTFIDF compares BM25 with basic TF-IDF for analysis
func (engine *BM25Engine) CompareWithTFIDF(query string, document string, tfidfScore float64) string {
	bm25Score := engine.Score(query, document)

	comparison := "BM25 vs TF-IDF Comparison:\n"
	comparison += "===========================\n\n"
	comparison += formatString("BM25 Score:  %.4f\n", bm25Score)
	comparison += formatString("TF-IDF Score: %.4f\n", tfidfScore)

	diff := bm25Score - tfidfScore
	percentDiff := (diff / tfidfScore) * 100

	if diff > 0 {
		comparison += formatString("Difference: +%.4f (+%.1f%%)\n", diff, percentDiff)
		comparison += "✅ BM25 scores higher (better)\n"
	} else {
		comparison += formatString("Difference: %.4f (%.1f%%)\n", diff, percentDiff)
		comparison += "⚠️  TF-IDF scores higher\n"
	}

	comparison += "\nWhy BM25 is generally better:\n"
	comparison += "- Term frequency saturation (diminishing returns)\n"
	comparison += "- Document length normalization (fairer comparison)\n"
	comparison += "- More sophisticated IDF formula\n"
	comparison += "- Industry standard for search engines\n"

	return comparison
}

// Helper functions

func toLowerSimple(s string) string {
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

func joinWords(words []string, sep string) string {
	if len(words) == 0 {
		return ""
	}

	result := words[0]
	for i := 1; i < len(words); i++ {
		result += sep + words[i]
	}
	return result
}

func formatString(format string, args ...interface{}) string {
	// Simple sprintf equivalent for basic formatting
	// This is a simplified version - in production use fmt.Sprintf
	result := format
	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			result = replaceFirst(result, "%s", v)
		case int:
			result = replaceFirst(result, "%d", intToString(v))
		case float64:
			// Simple float formatting
			result = replaceFirst(result, "%.4f", floatToString(v, 4))
			result = replaceFirst(result, "%.1f", floatToString(v, 1))
		}
	}
	return result
}

func replaceFirst(s, old, new string) string {
	idx := findSubstring(s, old)
	if idx < 0 {
		return s
	}
	return s[:idx] + new + s[idx+len(old):]
}

func intToString(n int) string {
	if n == 0 {
		return "0"
	}

	negative := n < 0
	if negative {
		n = -n
	}

	digits := ""
	for n > 0 {
		digits = string('0'+rune(n%10)) + digits
		n /= 10
	}

	if negative {
		digits = "-" + digits
	}

	return digits
}

func floatToString(f float64, precision int) string {
	// Simple float to string conversion
	intPart := int(f)
	fracPart := f - float64(intPart)

	result := intToString(intPart) + "."

	for i := 0; i < precision; i++ {
		fracPart *= 10
		digit := int(fracPart) % 10
		result += string('0' + rune(digit))
	}

	return result
}
