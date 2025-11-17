package llm

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// ContextNode represents a failure context in the memory graph
type ContextNode struct {
	ID               string                 // Unique context identifier
	CommandPattern   string                 // Abstracted command pattern
	ErrorType        string                 // Error category
	ProjectType      string                 // Language/framework
	Timestamp        time.Time              // When this occurred
	Transitions      map[string]*Transition // Edges to other contexts
	InsultPool       []WeightedInsult       // Specialized insults for this path
	VisitCount       int                    // How often we've seen this
	SuccessRate      float64                // How often user moves on
}

// Transition represents a path between contexts
type Transition struct {
	ToContext      string    // Target context ID
	Count          int       // How many times we've seen this transition
	AvgTimeBetween float64   // Average time between failures (seconds)
	LastSeen       time.Time // Last time we saw this transition
	SpecialInsults []string  // Insults specific to this failure sequence
}

// WeightedInsult is an insult with dynamic weight
type WeightedInsult struct {
	Text          string
	BaseWeight    float64
	DynamicWeight float64 // Updated via RL
	UseCount      int
	LastUsed      time.Time
	Effectiveness float64 // RL score: how well it worked
}

// ContextualMemoryGraph tracks failure patterns and relationships
type ContextualMemoryGraph struct {
	mu               sync.RWMutex
	nodes            map[string]*ContextNode
	currentContext   string
	previousContext  string
	lastTransition   time.Time
	persistencePath  string

	// Configuration
	decayFactor      float64 // How quickly weights decay
	minEffectiveness float64 // Minimum effectiveness to keep insult
	maxPoolSize      int     // Maximum insults per pool
}

// NewContextualMemoryGraph creates a new memory graph
func NewContextualMemoryGraph() *ContextualMemoryGraph {
	homeDir, _ := os.UserHomeDir()
	persistPath := filepath.Join(homeDir, ".parrot", "context_graph.json")

	graph := &ContextualMemoryGraph{
		nodes:            make(map[string]*ContextNode),
		persistencePath:  persistPath,
		decayFactor:      0.98,  // Slow decay
		minEffectiveness: 0.3,   // Keep insults above 30% effectiveness
		maxPoolSize:      50,    // Max 50 specialized insults per context
	}

	// Try to load existing graph
	graph.Load()

	return graph
}

// RecordContext records a failure in the graph
func (cmg *ContextualMemoryGraph) RecordContext(ctx *SmartFallbackContext) string {
	cmg.mu.Lock()
	defer cmg.mu.Unlock()

	// Create context ID from key features
	contextID := cmg.generateContextID(ctx)

	// Get or create node
	node, exists := cmg.nodes[contextID]
	if !exists {
		node = &ContextNode{
			ID:             contextID,
			CommandPattern: abstractCommand(ctx.FullCommand),
			ErrorType:      ctx.ErrorPattern,
			ProjectType:    ctx.ProjectType,
			Timestamp:      time.Now(),
			Transitions:    make(map[string]*Transition),
			InsultPool:     make([]WeightedInsult, 0),
			VisitCount:     0,
			SuccessRate:    0.5, // Start neutral
		}
		cmg.nodes[contextID] = node
	}

	node.VisitCount++

	// Record transition from previous context
	if cmg.previousContext != "" && cmg.previousContext != contextID {
		timeSince := time.Since(cmg.lastTransition).Seconds()

		transition, exists := node.Transitions[cmg.previousContext]
		if !exists {
			transition = &Transition{
				ToContext:      contextID,
				Count:          0,
				AvgTimeBetween: 0,
				SpecialInsults: make([]string, 0),
			}
			node.Transitions[cmg.previousContext] = transition
		}

		// Update transition statistics (exponential moving average)
		alpha := 0.3
		transition.AvgTimeBetween = alpha*timeSince + (1-alpha)*transition.AvgTimeBetween
		transition.Count++
		transition.LastSeen = time.Now()

		// If this is a rapid repeat (< 30 seconds), it's likely the same issue
		if timeSince < 30 {
			// Generate special "repeated failure" insult for this path
			cmg.addTransitionInsult(transition, ctx)
		}
	}

	// Update context tracking
	cmg.previousContext = cmg.currentContext
	cmg.currentContext = contextID
	cmg.lastTransition = time.Now()

	// Periodically persist
	if node.VisitCount%10 == 0 {
		go cmg.Save()
	}

	return contextID
}

// GetContextualInsult retrieves the best insult for current context
func (cmg *ContextualMemoryGraph) GetContextualInsult(contextID string) string {
	cmg.mu.RLock()
	defer cmg.mu.RUnlock()

	node, exists := cmg.nodes[contextID]
	if !exists || len(node.InsultPool) == 0 {
		return "" // No specialized insult
	}

	// Select insult using weighted random selection with novelty penalty
	bestInsult := ""
	bestScore := 0.0

	now := time.Now()
	for _, weighted := range node.InsultPool {
		// Calculate score based on effectiveness and recency
		score := weighted.DynamicWeight * weighted.Effectiveness

		// Penalty for recent use (novelty)
		hoursSinceUse := now.Sub(weighted.LastUsed).Hours()
		noveltyBonus := 1.0
		if hoursSinceUse < 1 {
			noveltyBonus = 0.1 // Heavy penalty
		} else if hoursSinceUse < 24 {
			noveltyBonus = 0.5 // Medium penalty
		}

		score *= noveltyBonus

		if score > bestScore {
			bestScore = score
			bestInsult = weighted.Text
		}
	}

	return bestInsult
}

// GetTransitionInsult gets insult specific to failure sequence
func (cmg *ContextualMemoryGraph) GetTransitionInsult(fromContext, toContext string) string {
	cmg.mu.RLock()
	defer cmg.mu.RUnlock()

	node, exists := cmg.nodes[toContext]
	if !exists {
		return ""
	}

	transition, exists := node.Transitions[fromContext]
	if !exists || len(transition.SpecialInsults) == 0 {
		return ""
	}

	// Return a transition-specific insult
	idx := transition.Count % len(transition.SpecialInsults)
	return transition.SpecialInsults[idx]
}

// RecordInsultUse records that an insult was used
func (cmg *ContextualMemoryGraph) RecordInsultUse(contextID string, insult string) {
	cmg.mu.Lock()
	defer cmg.mu.Unlock()

	node, exists := cmg.nodes[contextID]
	if !exists {
		return
	}

	// Find and update the insult
	for i := range node.InsultPool {
		if node.InsultPool[i].Text == insult {
			node.InsultPool[i].UseCount++
			node.InsultPool[i].LastUsed = time.Now()
			break
		}
	}
}

// RecordSuccess records that user moved past this failure
func (cmg *ContextualMemoryGraph) RecordSuccess(contextID string) {
	cmg.mu.Lock()
	defer cmg.mu.Unlock()

	node, exists := cmg.nodes[contextID]
	if !exists {
		return
	}

	// Update success rate (exponential moving average)
	alpha := 0.2
	node.SuccessRate = alpha*1.0 + (1-alpha)*node.SuccessRate

	// Boost effectiveness of recently used insults (they "worked")
	now := time.Now()
	for i := range node.InsultPool {
		if now.Sub(node.InsultPool[i].LastUsed).Minutes() < 5 {
			// This insult was used recently and user succeeded - boost it!
			node.InsultPool[i].Effectiveness =
				0.2*1.0 + 0.8*node.InsultPool[i].Effectiveness
		}
	}
}

// AddInsultToPool adds an insult to context-specific pool
func (cmg *ContextualMemoryGraph) AddInsultToPool(contextID string, insult string, weight float64) {
	cmg.mu.Lock()
	defer cmg.mu.Unlock()

	node, exists := cmg.nodes[contextID]
	if !exists {
		return
	}

	// Check if already exists
	for i := range node.InsultPool {
		if node.InsultPool[i].Text == insult {
			// Already exists, just update weight
			node.InsultPool[i].DynamicWeight =
				0.5*weight + 0.5*node.InsultPool[i].DynamicWeight
			return
		}
	}

	// Add new insult
	weighted := WeightedInsult{
		Text:          insult,
		BaseWeight:    weight,
		DynamicWeight: weight,
		UseCount:      0,
		LastUsed:      time.Now().Add(-24 * time.Hour), // Start as "old"
		Effectiveness: 0.5, // Start neutral
	}

	node.InsultPool = append(node.InsultPool, weighted)

	// Prune pool if too large (remove least effective)
	if len(node.InsultPool) > cmg.maxPoolSize {
		cmg.pruneInsultPool(node)
	}
}

// pruneInsultPool removes least effective insults
func (cmg *ContextualMemoryGraph) pruneInsultPool(node *ContextNode) {
	// Sort by effectiveness
	pool := node.InsultPool

	// Simple selection sort to find worst
	for i := 0; i < len(pool)-1; i++ {
		minIdx := i
		for j := i + 1; j < len(pool); j++ {
			if pool[j].Effectiveness < pool[minIdx].Effectiveness {
				minIdx = j
			}
		}
		if minIdx != i {
			pool[i], pool[minIdx] = pool[minIdx], pool[i]
		}
	}

	// Remove bottom 10%
	removeCount := cmg.maxPoolSize / 10
	if removeCount < 1 {
		removeCount = 1
	}

	node.InsultPool = pool[removeCount:]
}

// ApplyDecay applies decay to all insult weights (prevents stagnation)
func (cmg *ContextualMemoryGraph) ApplyDecay() {
	cmg.mu.Lock()
	defer cmg.mu.Unlock()

	for _, node := range cmg.nodes {
		for i := range node.InsultPool {
			// Decay dynamic weight back toward base weight
			node.InsultPool[i].DynamicWeight =
				cmg.decayFactor*node.InsultPool[i].DynamicWeight +
				(1-cmg.decayFactor)*node.InsultPool[i].BaseWeight
		}
	}
}

// GetStats returns graph statistics
func (cmg *ContextualMemoryGraph) GetStats() map[string]interface{} {
	cmg.mu.RLock()
	defer cmg.mu.RUnlock()

	totalInsults := 0
	totalTransitions := 0

	for _, node := range cmg.nodes {
		totalInsults += len(node.InsultPool)
		totalTransitions += len(node.Transitions)
	}

	return map[string]interface{}{
		"total_contexts":    len(cmg.nodes),
		"total_insults":     totalInsults,
		"total_transitions": totalTransitions,
		"current_context":   cmg.currentContext,
	}
}

// Save persists the graph to disk
func (cmg *ContextualMemoryGraph) Save() error {
	cmg.mu.RLock()
	defer cmg.mu.RUnlock()

	// Ensure directory exists
	dir := filepath.Dir(cmg.persistencePath)
	os.MkdirAll(dir, 0755)

	// Marshal to JSON
	data, err := json.MarshalIndent(cmg.nodes, "", "  ")
	if err != nil {
		return err
	}

	// Write atomically
	tmpPath := cmg.persistencePath + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		return err
	}

	return os.Rename(tmpPath, cmg.persistencePath)
}

// Load restores the graph from disk
func (cmg *ContextualMemoryGraph) Load() error {
	cmg.mu.Lock()
	defer cmg.mu.Unlock()

	data, err := os.ReadFile(cmg.persistencePath)
	if err != nil {
		return err // File might not exist yet
	}

	nodes := make(map[string]*ContextNode)
	if err := json.Unmarshal(data, &nodes); err != nil {
		return err
	}

	cmg.nodes = nodes
	return nil
}

// Helper functions

func (cmg *ContextualMemoryGraph) generateContextID(ctx *SmartFallbackContext) string {
	// Create a unique but abstract ID
	// Format: commandType_errorPattern_projectType
	parts := []string{
		ctx.CommandType,
		ctx.ErrorPattern,
		ctx.ProjectType,
	}

	id := ""
	for i, part := range parts {
		if part == "" {
			part = "unknown"
		}
		if i > 0 {
			id += "_"
		}
		id += part
	}

	return id
}

func abstractCommand(command string) string {
	// Replace specific values with placeholders
	// "git push origin feature-123" -> "git push origin <branch>"

	abstract := command

	// Replace file paths
	abstract = replacePattern(abstract, `[\w/]+\.(js|ts|py|go|rs|java|cpp)`, "<file>")

	// Replace branch names
	abstract = replacePattern(abstract, `(feature|bugfix|hotfix)/[\w-]+`, "<branch>")

	// Replace version numbers
	abstract = replacePattern(abstract, `\d+\.\d+\.\d+`, "<version>")

	// Replace ports
	abstract = replacePattern(abstract, `:\d{2,5}`, ":<port>")

	return abstract
}

func replacePattern(text string, pattern string, replacement string) string {
	// Simple pattern replacement (in production, use regexp)
	// This is a placeholder - implement properly
	return text
}

func (cmg *ContextualMemoryGraph) addTransitionInsult(transition *Transition, ctx *SmartFallbackContext) {
	// Generate insult specific to this failure sequence
	insults := []string{
		"Same mistake, different hour. Classic.",
		"Trying again won't make it work. Won't make you smarter either.",
		"Definition of insanity: Detected.",
		"Still failing? Maybe it's not the computer.",
		"Repeated failure #" + intToStr(transition.Count) + ": The saga continues.",
		"Your consistency is admirable. Consistently wrong.",
	}

	// Add if not already present
	for _, insult := range insults {
		found := false
		for _, existing := range transition.SpecialInsults {
			if existing == insult {
				found = true
				break
			}
		}
		if !found && len(transition.SpecialInsults) < 10 {
			transition.SpecialInsults = append(transition.SpecialInsults, insult)
		}
	}
}

func intToStr(n int) string {
	if n == 0 {
		return "0"
	}

	digits := ""
	for n > 0 {
		digits = string('0'+rune(n%10)) + digits
		n /= 10
	}
	return digits
}

