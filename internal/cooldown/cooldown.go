// Package cooldown provides per-key exponential back-off suppression.
// After each allowed event the quiet window doubles (up to a configurable
// maximum), so a port that flaps repeatedly generates progressively fewer
// alerts.
package cooldown

import (
	"sync"
	"time"
)

// Clock is a thin interface so tests can inject a fake time source.
type Clock interface {
	Now() time.Time
}

type realClock struct{}

func (realClock) Now() time.Time { return time.Now() }

// state tracks the back-off window for a single key.
type state struct {
	nextAllowed time.Time
	window      time.Duration
}

// Cooldown enforces exponential back-off on a per-key basis.
type Cooldown struct {
	mu      sync.Mutex
	clock   Clock
	base    time.Duration
	max     time.Duration
	entries map[string]*state
}

// New returns a Cooldown whose initial quiet window is base and whose
// maximum quiet window is max.
func New(base, max time.Duration) *Cooldown {
	return newWithClock(base, max, realClock{})
}

func newWithClock(base, max time.Duration, c Clock) *Cooldown {
	return &Cooldown{
		clock:   c,
		base:    base,
		max:     max,
		entries: make(map[string]*state),
	}
}

// Allow returns true when the event identified by key is permitted.
// On each allowed call the quiet window for that key is doubled (capped at
// max). The window resets automatically once the key has been quiet for
// longer than the current window.
func (cd *Cooldown) Allow(key string) bool {
	cd.mu.Lock()
	defer cd.mu.Unlock()

	now := cd.clock.Now()
	s, ok := cd.entries[key]
	if !ok {
		// First time we see this key — allow and start tracking.
		cd.entries[key] = &state{
			nextAllowed: now.Add(cd.base),
			window:      cd.base,
		}
		return true
	}

	if now.Before(s.nextAllowed) {
		return false
	}

	// The key has been quiet long enough — reset window if it has been
	// idle for more than twice the current window (natural decay).
	idle := now.Sub(s.nextAllowed)
	if idle >= s.window*2 {
		s.window = cd.base
	} else {
		s.window *= 2
		if s.window > cd.max {
			s.window = cd.max
		}
	}
	s.nextAllowed = now.Add(s.window)
	return true
}

// Reset clears the back-off state for key, as if it had never been seen.
func (cd *Cooldown) Reset(key string) {
	cd.mu.Lock()
	defer cd.mu.Unlock()
	delete(cd.entries, key)
}
