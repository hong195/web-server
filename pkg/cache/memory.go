package cache

import (
	"sync"
	"time"
)

// cacheEntry represents a single cache entry with expiration.
type cacheEntry struct {
	value     []byte
	expiresAt time.Time
}

// MemoryCache is an in-memory implementation of Cache.
type MemoryCache struct {
	data map[string]cacheEntry
	mu   sync.RWMutex
}

// NewMemoryCache creates a new MemoryCache instance.
func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		data: make(map[string]cacheEntry),
	}
}

// Get retrieves a value from cache if it exists and hasn't expired.
func (c *MemoryCache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	entry, exists := c.data[key]
	c.mu.RUnlock()

	if !exists {
		return nil, false
	}

	if time.Now().After(entry.expiresAt) {
		c.mu.Lock()
		delete(c.data, key)
		c.mu.Unlock()
		return nil, false
	}

	return entry.value, true
}

// Set stores a value in cache with the given TTL.
func (c *MemoryCache) Set(key string, value []byte, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = cacheEntry{
		value:     value,
		expiresAt: time.Now().Add(ttl),
	}
}
