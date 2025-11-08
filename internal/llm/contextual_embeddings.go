package llm

import (
	"math"
	"strings"
)

// ContextualEmbedding represents a context as a vector in semantic space
// This is a simplified embedding system that works without external ML libraries
type ContextualEmbedding struct {
	Vector     []float64
	ContextID  string
	Features   map[string]float64 // Named features for interpretability
	Magnitude  float64            // Vector magnitude (cached)
}

// EmbeddingEngine creates and compares contextual embeddings
type EmbeddingEngine struct {
	dimensions int
	vocabulary map[string]int // Word to index mapping
	idf        map[string]float64 // IDF scores for weighting
}

// NewEmbeddingEngine creates a new embedding engine
func NewEmbeddingEngine() *EmbeddingEngine {
	return &EmbeddingEngine{
		dimensions: 32, // 32-dimensional embedding space
		vocabulary: make(map[string]int),
		idf:        make(map[string]float64),
	}
}

// CreateEmbedding creates a vector representation of a context
func (ee *EmbeddingEngine) CreateEmbedding(ctx *SmartFallbackContext) *ContextualEmbedding {
	embedding := &ContextualEmbedding{
		Vector:    make([]float64, ee.dimensions),
		ContextID: ctx.CommandType + "_" + ctx.ErrorPattern,
		Features:  make(map[string]float64),
	}

	// Feature extraction and encoding
	// Each feature maps to specific dimensions in the vector

	// Dimension 0-7: Command type encoding
	ee.encodeCommandType(ctx, embedding, 0)

	// Dimension 8-15: Error pattern encoding
	ee.encodeErrorPattern(ctx, embedding, 8)

	// Dimension 16-19: Project context
	ee.encodeProjectContext(ctx, embedding, 16)

	// Dimension 20-23: Temporal context
	ee.encodeTemporalContext(ctx, embedding, 20)

	// Dimension 24-27: Environmental context
	ee.encodeEnvironmentalContext(ctx, embedding, 24)

	// Dimension 28-31: Behavioral patterns
	ee.encodeBehavioralContext(ctx, embedding, 28)

	// Normalize vector
	ee.normalizeVector(embedding)

	return embedding
}

// FindSimilarContexts finds contexts similar to the given one
func (ee *EmbeddingEngine) FindSimilarContexts(
	target *ContextualEmbedding,
	candidates []*ContextualEmbedding,
	topK int,
) []SimilarContext {
	similarities := make([]SimilarContext, 0, len(candidates))

	for _, candidate := range candidates {
		similarity := ee.CosineSimilarity(target, candidate)
		similarities = append(similarities, SimilarContext{
			Embedding:  candidate,
			Similarity: similarity,
		})
	}

	// Sort by similarity (descending)
	sortSimilarContexts(similarities)

	if len(similarities) > topK {
		similarities = similarities[:topK]
	}

	return similarities
}

// CosineSimilarity calculates cosine similarity between two embeddings
func (ee *EmbeddingEngine) CosineSimilarity(e1, e2 *ContextualEmbedding) float64 {
	if len(e1.Vector) != len(e2.Vector) {
		return 0.0
	}

	dotProduct := 0.0
	for i := range e1.Vector {
		dotProduct += e1.Vector[i] * e2.Vector[i]
	}

	// Use cached magnitudes
	magnitude1 := e1.Magnitude
	magnitude2 := e2.Magnitude

	if magnitude1 == 0 || magnitude2 == 0 {
		return 0.0
	}

	return dotProduct / (magnitude1 * magnitude2)
}

// EuclideanDistance calculates Euclidean distance between embeddings
func (ee *EmbeddingEngine) EuclideanDistance(e1, e2 *ContextualEmbedding) float64 {
	if len(e1.Vector) != len(e2.Vector) {
		return math.Inf(1)
	}

	sum := 0.0
	for i := range e1.Vector {
		diff := e1.Vector[i] - e2.Vector[i]
		sum += diff * diff
	}

	return math.Sqrt(sum)
}

// Encoding functions for different context aspects

func (ee *EmbeddingEngine) encodeCommandType(ctx *SmartFallbackContext, emb *ContextualEmbedding, offset int) {
	// One-hot-like encoding for common command types
	commandTypes := map[string]int{
		"git":        0,
		"docker":     1,
		"nodejs":     2,
		"python":     3,
		"rust":       4,
		"golang":     5,
		"kubernetes": 6,
		"database":   7,
	}

	if idx, exists := commandTypes[ctx.CommandType]; exists {
		emb.Vector[offset+idx] = 1.0
		emb.Features["command_type_"+ctx.CommandType] = 1.0
	}
}

func (ee *EmbeddingEngine) encodeErrorPattern(ctx *SmartFallbackContext, emb *ContextualEmbedding, offset int) {
	// Encode error pattern
	errorPatterns := map[string]int{
		"permission_denied": 0,
		"command_not_found": 1,
		"network_error":     2,
		"timeout":           3,
		"syntax_error":      4,
		"merge_conflict":    5,
		"build_failure":     6,
		"test_failure":      7,
	}

	if idx, exists := errorPatterns[ctx.ErrorPattern]; exists {
		emb.Vector[offset+idx] = 1.0
		emb.Features["error_"+ctx.ErrorPattern] = 1.0
	}

	// Encode exit code (normalized)
	if ctx.ExitCode > 0 {
		normalizedExitCode := math.Min(float64(ctx.ExitCode)/255.0, 1.0)
		emb.Vector[offset+7] = normalizedExitCode
		emb.Features["exit_code"] = float64(ctx.ExitCode)
	}
}

func (ee *EmbeddingEngine) encodeProjectContext(ctx *SmartFallbackContext, emb *ContextualEmbedding, offset int) {
	// Encode project type
	projectTypes := map[string]float64{
		"node":   0.0,
		"rust":   0.25,
		"go":     0.5,
		"python": 0.75,
		"java":   1.0,
	}

	if val, exists := projectTypes[ctx.ProjectType]; exists {
		emb.Vector[offset] = val
		emb.Features["project_type"] = val
	}

	// Encode git branch awareness
	if ctx.GitBranch != "" {
		branchScore := 0.0
		if ctx.GitBranch == "main" || ctx.GitBranch == "master" {
			branchScore = 1.0 // High risk
		} else if strings.Contains(ctx.GitBranch, "prod") {
			branchScore = 0.9
		} else if strings.Contains(ctx.GitBranch, "develop") {
			branchScore = 0.5
		}
		emb.Vector[offset+1] = branchScore
		emb.Features["branch_risk"] = branchScore
	}

	// Encode dependency complexity
	if ctx.DependencyCount > 0 {
		// Normalize dependency count (log scale)
		normalized := math.Log(float64(ctx.DependencyCount)+1) / math.Log(100)
		emb.Vector[offset+2] = math.Min(normalized, 1.0)
		emb.Features["dependency_complexity"] = normalized
	}

	// Has build files
	if ctx.HasDockerfile || ctx.HasMakefile {
		emb.Vector[offset+3] = 1.0
		emb.Features["has_build_system"] = 1.0
	}
}

func (ee *EmbeddingEngine) encodeTemporalContext(ctx *SmartFallbackContext, emb *ContextualEmbedding, offset int) {
	// Encode time of day (cyclical encoding using sin/cos)
	hour := float64(ctx.TimeOfDay)
	hourRadian := (hour / 24.0) * 2.0 * math.Pi

	emb.Vector[offset] = math.Sin(hourRadian)
	emb.Vector[offset+1] = math.Cos(hourRadian)
	emb.Features["time_sin"] = emb.Vector[offset]
	emb.Features["time_cos"] = emb.Vector[offset+1]

	// Late night coding indicator (high value between 22-4)
	lateNight := 0.0
	if hour >= 22 || hour <= 4 {
		lateNight = 1.0
	}
	emb.Vector[offset+2] = lateNight
	emb.Features["late_night"] = lateNight

	// Repeated failure indicator
	if ctx.IsRepeatedFailure {
		emb.Vector[offset+3] = 1.0
		emb.Features["repeated_failure"] = 1.0
	}
}

func (ee *EmbeddingEngine) encodeEnvironmentalContext(ctx *SmartFallbackContext, emb *ContextualEmbedding, offset int) {
	// CI/CD environment
	if ctx.IsCI {
		emb.Vector[offset] = 1.0
		emb.Features["is_ci"] = 1.0

		// Encode CI provider
		ciProviders := map[string]float64{
			"github": 0.2,
			"gitlab": 0.4,
			"jenkins": 0.6,
			"circle": 0.8,
		}
		if val, exists := ciProviders[ctx.CIProvider]; exists {
			emb.Vector[offset+1] = val
			emb.Features["ci_provider"] = val
		}
	}

	// Shell type encoding
	shells := map[string]float64{
		"bash": 0.2,
		"zsh":  0.4,
		"fish": 0.6,
		"sh":   0.8,
	}
	if val, exists := shells[ctx.Shell]; exists {
		emb.Vector[offset+2] = val
		emb.Features["shell"] = val
	}

	// Command complexity
	complexityScore := 0.0
	if ctx.HasPipes {
		complexityScore += 0.3
	}
	if ctx.HasChaining {
		complexityScore += 0.3
	}
	if ctx.CommandLength > 100 {
		complexityScore += 0.4
	}
	emb.Vector[offset+3] = math.Min(complexityScore, 1.0)
	emb.Features["complexity"] = complexityScore
}

func (ee *EmbeddingEngine) encodeBehavioralContext(ctx *SmartFallbackContext, emb *ContextualEmbedding, offset int) {
	// Working directory patterns
	wdPatterns := map[string]float64{
		"tmp":       0.8,
		"downloads": 0.6,
		"desktop":   0.4,
	}

	wdLower := strings.ToLower(ctx.WorkingDir)
	for pattern, score := range wdPatterns {
		if strings.Contains(wdLower, pattern) {
			emb.Vector[offset] = score
			emb.Features["wd_pattern"] = score
			break
		}
	}

	// File extensions present
	if len(ctx.FileExtensions) > 0 {
		emb.Vector[offset+1] = 1.0
		emb.Features["has_files"] = 1.0
	}

	// Numeric arguments (ports, chmod values, etc.)
	if len(ctx.NumericArgs) > 0 {
		// Check for risky values
		riskyNums := map[int]float64{
			777: 1.0, // chmod 777
			666: 0.8, // chmod 666
			9:   0.6, // kill -9
		}
		for _, num := range ctx.NumericArgs {
			if score, exists := riskyNums[num]; exists {
				emb.Vector[offset+2] = score
				emb.Features["risky_numeric"] = score
				break
			}
		}
	}

	// Git-specific context
	if ctx.CommandType == "git" {
		if strings.Contains(strings.ToLower(ctx.FullCommand), "force") {
			emb.Vector[offset+3] = 1.0
			emb.Features["force_flag"] = 1.0
		}
	}
}

func (ee *EmbeddingEngine) normalizeVector(emb *ContextualEmbedding) {
	// Calculate magnitude
	sumSquares := 0.0
	for _, val := range emb.Vector {
		sumSquares += val * val
	}
	emb.Magnitude = math.Sqrt(sumSquares)

	// Normalize (make unit vector)
	if emb.Magnitude > 0 {
		for i := range emb.Vector {
			emb.Vector[i] /= emb.Magnitude
		}
		// Update magnitude to 1.0 after normalization
		emb.Magnitude = 1.0
	}
}

// SimilarContext represents a similar context with similarity score
type SimilarContext struct {
	Embedding  *ContextualEmbedding
	Similarity float64
}

func sortSimilarContexts(contexts []SimilarContext) {
	// Bubble sort by similarity (descending)
	n := len(contexts)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if contexts[j].Similarity < contexts[j+1].Similarity {
				contexts[j], contexts[j+1] = contexts[j+1], contexts[j]
			}
		}
	}
}

// GetFeatureImportance returns which features contributed most to the embedding
func (emb *ContextualEmbedding) GetFeatureImportance() []FeatureImportance {
	features := make([]FeatureImportance, 0, len(emb.Features))

	for name, value := range emb.Features {
		features = append(features, FeatureImportance{
			Name:  name,
			Value: value,
		})
	}

	// Sort by absolute value (descending)
	sortFeatureImportance(features)

	return features
}

// FeatureImportance represents a feature and its value
type FeatureImportance struct {
	Name  string
	Value float64
}

func sortFeatureImportance(features []FeatureImportance) {
	n := len(features)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			absI := features[i].Value
			if absI < 0 {
				absI = -absI
			}
			absJ := features[j].Value
			if absJ < 0 {
				absJ = -absJ
			}

			if absJ < absI {
				features[i], features[j] = features[j], features[i]
			}
		}
	}
}

// ExplainSimilarity explains why two contexts are similar
func (ee *EmbeddingEngine) ExplainSimilarity(e1, e2 *ContextualEmbedding) string {
	explanation := "Similarity Breakdown:\n"

	// Compare feature by feature
	sharedFeatures := make([]string, 0)
	for feature := range e1.Features {
		if _, exists := e2.Features[feature]; exists {
			sharedFeatures = append(sharedFeatures, feature)
		}
	}

	if len(sharedFeatures) > 0 {
		explanation += "Shared features: "
		for i, feature := range sharedFeatures {
			if i > 0 {
				explanation += ", "
			}
			explanation += feature
		}
		explanation += "\n"
	}

	// Calculate similarity
	similarity := ee.CosineSimilarity(e1, e2)
	explanation += "Cosine similarity: "
	explanation += formatFloat(similarity)
	explanation += "\n"

	return explanation
}
