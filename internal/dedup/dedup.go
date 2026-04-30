// Package dedup provides a deduplication layer that suppresses identical
// port-change diffs within a configurable time window, preventing redundant
// alerts when the same transition is observed across consecutive scans.
package dedup

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

// clock abstracts time for testing.
type clock interface {
	Now() time.Time
}

type realClock struct{}

func (realClock) Now() time.Time { return time.Now() }

// entry records when a diff signature was last seen.
type entry struct {
	sig  string
	seen time.Time
}

// Deduplicator filters out diffs whose signature was already emitted
// within the configured window.
type Deduplicator struct {
	mu     sync.Mutex
	clock  clock
	window time.Duration
	cache  map[string]entry
}

// New returns a Deduplicator with the given suppression window.
func New(window time.Duration) *Deduplicator {
	return newWithClock(window, realClock{})
}

func newWithClock(window time.Duration, c clock) *Deduplicator {
	return &Deduplicator{
		clock:  c,
		window: window,
		cache:  make(map[string]entry),
	}
}

// IsDuplicate reports whether d is a duplicate of a recently seen diff.
// If it is not a duplicate the diff is recorded and false is returned.
func (d *Deduplicator) IsDuplicate(diff snapshot.Diff) bool {
	sig := signature(diff)
	if sig == "" {
		return false
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	now := d.clock.Now()
	if e, ok := d.cache[sig]; ok && now.Sub(e.seen) < d.window {
		return true
	}
	d.cache[sig] = entry{sig: sig, seen: now}
	return false
}

// Purge removes all expired entries from the internal cache.
func (d *Deduplicator) Purge() {
	d.mu.Lock()
	defer d.mu.Unlock()
	now := d.clock.Now()
	for k, e := range d.cache {
		if now.Sub(e.seen) >= d.window {
			delete(d.cache, k)
		}
	}
}

// signature builds a compact, order-stable string key for a diff.
func signature(diff snapshot.Diff) string {
	if len(diff.Opened) == 0 && len(diff.Closed) == 0 {
		return ""
	}
	var b []byte
	for _, p := range diff.Opened {
		b = append(b, '+', byte(p>>8), byte(p))
	}
	for _, p := range diff.Closed {
		b = append(b, '-', byte(p>>8), byte(p))
	}
	return string(b)
}
