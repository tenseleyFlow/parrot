package llm

import (
	"crypto/sha256"
	"encoding/hex"
	"sync"
	"time"
)

// ResponseCache provides an LRU cache for LLM responses to avoid redundant calls
type ResponseCache struct {
	mu       sync.RWMutex
	entries  map[string]*cacheEntry
	maxSize  int
	ttl      time.Duration
}

type cacheEntry struct {
	response  string
	timestamp time.Time
}

// Global cache instance
var (
	responseCache     *ResponseCache
	responseCacheOnce sync.Once
)

// GetResponseCache returns the singleton response cache
func GetResponseCache() *ResponseCache {
	responseCacheOnce.Do(func() {
		responseCache = NewResponseCache(100, 5*time.Minute) // 100 entries, 5 min TTL
	})
	return responseCache
}

// NewResponseCache creates a new response cache
func NewResponseCache(maxSize int, ttl time.Duration) *ResponseCache {
	cache := &ResponseCache{
		entries: make(map[string]*cacheEntry),
		maxSize: maxSize,
		ttl:     ttl,
	}

	// Start background cleanup goroutine
	go cache.cleanupLoop()

	return cache
}

// generateKey creates a cache key from command signature
func (c *ResponseCache) generateKey(command, commandType, exitCode, mode string) string {
	// Create a hash of the command signature
	h := sha256.New()
	h.Write([]byte(command))
	h.Write([]byte("|"))
	h.Write([]byte(commandType))
	h.Write([]byte("|"))
	h.Write([]byte(exitCode))
	h.Write([]byte("|"))
	h.Write([]byte(mode))
	return hex.EncodeToString(h.Sum(nil))[:16] // Use first 16 chars of hash
}

// Get retrieves a cached response if available and not expired
func (c *ResponseCache) Get(command, commandType, exitCode, mode string) (string, bool) {
	key := c.generateKey(command, commandType, exitCode, mode)

	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.entries[key]
	if !exists {
		return "", false
	}

	// Check if expired
	if time.Since(entry.timestamp) > c.ttl {
		return "", false
	}

	return entry.response, true
}

// Set stores a response in the cache
func (c *ResponseCache) Set(command, commandType, exitCode, mode, response string) {
	key := c.generateKey(command, commandType, exitCode, mode)

	c.mu.Lock()
	defer c.mu.Unlock()

	// Evict oldest entries if at capacity
	if len(c.entries) >= c.maxSize {
		c.evictOldest()
	}

	c.entries[key] = &cacheEntry{
		response:  response,
		timestamp: time.Now(),
	}
}

// evictOldest removes the oldest entry (must be called with lock held)
func (c *ResponseCache) evictOldest() {
	var oldestKey string
	var oldestTime time.Time

	for key, entry := range c.entries {
		if oldestKey == "" || entry.timestamp.Before(oldestTime) {
			oldestKey = key
			oldestTime = entry.timestamp
		}
	}

	if oldestKey != "" {
		delete(c.entries, oldestKey)
	}
}

// cleanupLoop periodically removes expired entries
func (c *ResponseCache) cleanupLoop() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.cleanup()
	}
}

// cleanup removes all expired entries
func (c *ResponseCache) cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, entry := range c.entries {
		if now.Sub(entry.timestamp) > c.ttl {
			delete(c.entries, key)
		}
	}
}

// Stats returns cache statistics
func (c *ResponseCache) Stats() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return map[string]interface{}{
		"size":     len(c.entries),
		"max_size": c.maxSize,
		"ttl_secs": c.ttl.Seconds(),
	}
}
