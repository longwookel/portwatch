// Package throttle provides a per-key cooldown mechanism to suppress
// repeated alerts for the same port within a configurable window.
package throttle

import (
	"sync"
	"time"
)

// Clock abstracts time so tests can inject a fake.
type Clock func() time.Time

// Throttle tracks the last alert time per key and suppresses duplicates
// within the cooldown window.
type Throttle struct {
	mu       sync.Mutex
	cooldown time.Duration
	clock    Clock
	last     map[string]time.Time
}

// New creates a Throttle with the given cooldown duration.
// Pass a nil clock to use the real wall clock.
func New(cooldown time.Duration, clock Clock) *Throttle {
	if clock == nil {
		clock = time.Now
	}
	return &Throttle{
		cooldown: cooldown,
		clock:    clock,
		last:     make(map[string]time.Time),
	}
}

// Allow returns true if the key has not been seen within the cooldown window,
// and records the current time for that key. Returns false if the key is still
// within its cooldown period.
func (t *Throttle) Allow(key string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := t.clock()
	if prev, ok := t.last[key]; ok {
		if now.Sub(prev) < t.cooldown {
			return false
		}
	}
	t.last[key] = now
	return true
}

// Reset clears the recorded time for a specific key, allowing it to fire
// immediately on the next call to Allow.
func (t *Throttle) Reset(key string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.last, key)
}

// Remaining returns how much cooldown time is left for the given key.
// Returns 0 if the key is not throttled.
func (t *Throttle) Remaining(key string) time.Duration {
	t.mu.Lock()
	defer t.mu.Unlock()

	prev, ok := t.last[key]
	if !ok {
		return 0
	}
	elapsed := t.clock().Sub(prev)
	if elapsed >= t.cooldown {
		return 0
	}
	return t.cooldown - elapsed
}
