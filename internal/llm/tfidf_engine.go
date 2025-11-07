package llm

import (
	"math"
	"strings"
	"unicode"
)

// TFIDFEngine implements semantic similarity using TF-IDF vectors
type TFIDFEngine struct {
	vocabulary     map[string]int    // word -> index
	idf            map[string]float64 // word -> inverse document frequency
	documentCount  int
	ngramRange     [2]int // min and max n-gram size
}

// Document represents a text document with its TF-IDF vector
type Document struct {
	Text   string
	Vector map[string]float64 // sparse vector representation
}

// NewTFIDFEngine creates a new TF-IDF engine
func NewTFIDFEngine() *TFIDFEngine {
	return &TFIDFEngine{
		vocabulary:    make(map[string]int),
		idf:           make(map[string]float64),
		documentCount: 0,
		ngramRange:    [2]int{1, 3}, // unigrams, bigrams, trigrams
	}
}

// BuildCorpus builds the TF-IDF corpus from a collection of documents
func (engine *TFIDFEngine) BuildCorpus(documents []string) {
	// First pass: build vocabulary and count document frequencies
	documentFreq := make(map[string]int)

	for _, doc := range documents {
		tokens := engine.extractNGrams(doc)
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

	// Calculate IDF for each term
	for term, docFreq := range documentFreq {
		// IDF = log(N / df) where N is total docs, df is docs containing term
		engine.idf[term] = math.Log(float64(engine.documentCount) / float64(docFreq))
	}
}

// extractNGrams extracts n-grams from text
func (engine *TFIDFEngine) extractNGrams(text string) []string {
	text = strings.ToLower(text)
	words := engine.tokenize(text)

	var ngrams []string

	// Generate n-grams for all sizes in range
	for n := engine.ngramRange[0]; n <= engine.ngramRange[1]; n++ {
		if n > len(words) {
			break
		}

		for i := 0; i <= len(words)-n; i++ {
			ngram := strings.Join(words[i:i+n], " ")
			ngrams = append(ngrams, ngram)
		}
	}

	return ngrams
}

// tokenize splits text into words
func (engine *TFIDFEngine) tokenize(text string) []string {
	var words []string
	var currentWord strings.Builder

	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsNumber(r) || r == '-' || r == '_' {
			currentWord.WriteRune(r)
		} else {
			if currentWord.Len() > 0 {
				word := currentWord.String()
				if len(word) > 1 { // Skip single characters
					words = append(words, word)
				}
				currentWord.Reset()
			}
		}
	}

	if currentWord.Len() > 0 {
		word := currentWord.String()
		if len(word) > 1 {
			words = append(words, word)
		}
	}

	return words
}

// Vectorize converts text to TF-IDF vector
func (engine *TFIDFEngine) Vectorize(text string) map[string]float64 {
	vector := make(map[string]float64)
	tokens := engine.extractNGrams(text)

	// Calculate term frequencies
	termFreq := make(map[string]int)
	for _, token := range tokens {
		termFreq[token]++
	}

	// Calculate TF-IDF for each term
	totalTerms := len(tokens)
	for term, freq := range termFreq {
		// TF = freq / total_terms
		tf := float64(freq) / float64(totalTerms)

		// Get IDF (use 1.0 if term not in vocabulary - rare term)
		idf := 1.0
		if val, exists := engine.idf[term]; exists {
			idf = val
		}

		// TF-IDF = TF * IDF
		vector[term] = tf * idf
	}

	// Normalize vector
	return engine.normalizeVector(vector)
}

// normalizeVector normalizes a vector to unit length
func (engine *TFIDFEngine) normalizeVector(vector map[string]float64) map[string]float64 {
	// Calculate magnitude
	var sumSquares float64
	for _, value := range vector {
		sumSquares += value * value
	}
	magnitude := math.Sqrt(sumSquares)

	if magnitude == 0 {
		return vector
	}

	// Normalize
	normalized := make(map[string]float64)
	for term, value := range vector {
		normalized[term] = value / magnitude
	}

	return normalized
}

// CosineSimilarity calculates cosine similarity between two vectors
func (engine *TFIDFEngine) CosineSimilarity(vec1, vec2 map[string]float64) float64 {
	// Calculate dot product
	var dotProduct float64
	for term, val1 := range vec1 {
		if val2, exists := vec2[term]; exists {
			dotProduct += val1 * val2
		}
	}

	// Vectors are already normalized, so similarity = dot product
	return dotProduct
}

// FindMostSimilar finds the most similar documents to a query
func (engine *TFIDFEngine) FindMostSimilar(
	query string,
	documents []Document,
	topK int,
) []SimilarityScore {
	queryVec := engine.Vectorize(query)

	scores := make([]SimilarityScore, 0, len(documents))
	for i, doc := range documents {
		similarity := engine.CosineSimilarity(queryVec, doc.Vector)
		scores = append(scores, SimilarityScore{
			Index:      i,
			Similarity: similarity,
			Text:       doc.Text,
		})
	}

	// Sort by similarity descending
	sortSimilarityScores(scores)

	// Return top K
	if len(scores) > topK {
		scores = scores[:topK]
	}

	return scores
}

// SimilarityScore represents a similarity score for a document
type SimilarityScore struct {
	Index      int
	Similarity float64
	Text       string
}

// sortSimilarityScores sorts scores in descending order
func sortSimilarityScores(scores []SimilarityScore) {
	// Simple bubble sort (good enough for small datasets)
	n := len(scores)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if scores[j].Similarity < scores[j+1].Similarity {
				scores[j], scores[j+1] = scores[j+1], scores[j]
			}
		}
	}
}

// ExtractKeyPhrases extracts important phrases from text using TF-IDF
func (engine *TFIDFEngine) ExtractKeyPhrases(text string, topN int) []string {
	vector := engine.Vectorize(text)

	// Convert to sorted list
	type termScore struct {
		term  string
		score float64
	}

	scores := make([]termScore, 0, len(vector))
	for term, score := range vector {
		scores = append(scores, termScore{term, score})
	}

	// Sort by score descending
	for i := 0; i < len(scores)-1; i++ {
		for j := 0; j < len(scores)-i-1; j++ {
			if scores[j].score < scores[j+1].score {
				scores[j], scores[j+1] = scores[j+1], scores[j]
			}
		}
	}

	// Extract top N terms
	result := make([]string, 0, topN)
	for i := 0; i < topN && i < len(scores); i++ {
		result = append(result, scores[i].term)
	}

	return result
}

// CalculateSemanticScore calculates semantic similarity between command and insult
func (engine *TFIDFEngine) CalculateSemanticScore(
	command string,
	insult string,
) float64 {
	cmdVec := engine.Vectorize(command)
	insultVec := engine.Vectorize(insult)

	return engine.CosineSimilarity(cmdVec, insultVec)
}
