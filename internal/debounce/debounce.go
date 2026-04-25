// Package debounce provides a mechanism to suppress repeated events within
// a configurable quiet period, preventing alert storms when ports flap.
package debounce

import (
	"sync"
	"time"
)

// Clock is an interface for time operations, allowing tests to inject a fake.
type Clock interface {
	Now() time.Time
}

type realClock struct{}

func (realClock) Now() time.Time { return time.Now() }

// Debouncer tracks the last time a key was seen and reports whether enough
// time has elapsed since the previous occurrence.
type Debouncer struct {
	mu       sync.Mutex
	quiet    time.Duration
	clock    Clock
	lastSeen map[string]time.Time
}

// New returns a Debouncer that suppresses repeated keys within quiet duration.
func New(quiet time.Duration) *Debouncer {
	return newWithClock(quiet, realClock{})
}

func newWithClock(quiet time.Duration, clk Clock) *Debouncer {
	return &Debouncer{
		quiet:    quiet,
		clock:    clk,
		lastSeen: make(map[string]time.Time),
	}
}

// Allow returns true if the key has not been seen within the quiet period.
// Calling Allow always records the current time for the key.
func (d *Debouncer) Allow(key string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := d.clock.Now()
	last, seen := d.lastSeen[key]
	d.lastSeen[key] = now

	if !seen {
		return true
	}
	return now.Sub(last) >= d.quiet
}

// Reset clears the recorded time for key, so the next call to Allow returns true.
func (d *Debouncer) Reset(key string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.lastSeen, key)
}

// Len returns the number of keys currently tracked.
func (d *Debouncer) Len() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.lastSeen)
}
