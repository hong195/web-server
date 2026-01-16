package cache

import "time"

// Cache defines the interface for caching operations.
type Cache interface {
	Get(key string) ([]byte, bool)
	Set(key string, value []byte, ttl time.Duration)
}
