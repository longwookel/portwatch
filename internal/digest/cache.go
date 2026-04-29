package digest

import "sync"

// Cache stores the last-seen Digest per scan target (e.g. host or interface)
// so callers can detect whether a new scan result differs from the previous one
// without storing the full port list.
type Cache struct {
	mu    sync.Mutex
	store map[string]Digest
}

// NewCache returns an initialised, empty Cache.
func NewCache() *Cache {
	return &Cache{store: make(map[string]Digest)}
}

// Changed reports whether the supplied digest differs from the previously
// stored value for key. It atomically updates the stored digest to next.
// On the first call for a given key it always returns true so that an initial
// snapshot is treated as a change.
func (c *Cache) Changed(key string, next Digest) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	prev, exists := c.store[key]
	c.store[key] = next
	if !exists {
		return true
	}
	return !Equal(prev, next)
}

// Reset removes the stored digest for key, causing the next call to Changed
// to return true regardless of the digest value.
func (c *Cache) Reset(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.store, key)
}

// Peek returns the last stored digest for key and whether it exists.
func (c *Cache) Peek(key string) (Digest, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	d, ok := c.store[key]
	return d, ok
}
