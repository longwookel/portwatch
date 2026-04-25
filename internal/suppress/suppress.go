// Package suppress provides a mechanism to suppress repeated alerts
// for the same port change within a configurable time window.
package suppress

import (
	"sync"
	"time"
)

// clock allows time to be injected for testing.
type clock func() time.Time

// Suppressor tracks recently seen port-change keys and suppresses
// duplicate notifications until the quiet window expires.
type Suppressor struct {
	mu      sync.Mutex
	window  time.Duration
	seen    map[string]time.Time
	now     clock
}

// New returns a Suppressor that silences repeated events for the same
// key within window.
func New(window time.Duration) *Suppressor {
	return newWithClock(window, time.Now)
}

func newWithClock(window time.Duration, now clock) *Suppressor {
	return &Suppressor{
		window: window,
		seen:   make(map[string]time.Time),
		now:    now,
	}
}

// Allow returns true the first time a key is seen, and again only
// after the quiet window has elapsed since it was last allowed.
// Expired entries are pruned on each call to keep memory bounded.
func (s *Suppressor) Allow(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := s.now()
	s.prune(now)

	if last, ok := s.seen[key]; ok && now.Sub(last) < s.window {
		return false
	}

	s.seen[key] = now
	return true
}

// prune removes entries whose window has expired. Must be called with s.mu held.
func (s *Suppressor) prune(now time.Time) {
	for k, t := range s.seen {
		if now.Sub(t) >= s.window {
			delete(s.seen, k)
		}
	}
}

// Reset clears all suppression state, allowing all keys through immediately.
func (s *Suppressor) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.seen = make(map[string]time.Time)
}
