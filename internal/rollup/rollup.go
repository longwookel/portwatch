// Package rollup coalesces rapid-fire port-change events into a single
// summarised notification, reducing alert noise during port storms.
package rollup

import (
	"sync"
	"time"

	"github.com/example/portwatch/internal/snapshot"
)

// Flusher is called with the merged Diff when the window closes.
type Flusher func(d snapshot.Diff)

// Rollup batches Diff events that arrive within a quiet window and
// delivers a single merged Diff to the Flusher once the window expires.
type Rollup struct {
	mu       sync.Mutex
	window   time.Duration
	flush    Flusher
	pending  *snapshot.Diff
	timer    *time.Timer
	nowFunc  func() time.Time
}

// New creates a Rollup that waits window duration after the last Add
// before calling flush with the merged result.
func New(window time.Duration, flush Flusher) *Rollup {
	return &Rollup{
		window:  window,
		flush:   flush,
		nowFunc: time.Now,
	}
}

// Add merges d into the pending Diff and (re)starts the quiet-period timer.
// It is safe to call from multiple goroutines.
func (r *Rollup) Add(d snapshot.Diff) {
	if len(d.Opened) == 0 && len(d.Closed) == 0 {
		return
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if r.pending == nil {
		r.pending = &snapshot.Diff{}
	}

	for _, p := range d.Opened {
		r.pending.Opened = append(r.pending.Opened, p)
	}
	for _, p := range d.Closed {
		r.pending.Closed = append(r.pending.Closed, p)
	}

	if r.timer != nil {
		r.timer.Stop()
	}
	r.timer = time.AfterFunc(r.window, r.fire)
}

// Flush forces an immediate delivery of any pending Diff, bypassing the
// quiet-period timer. It is a no-op when there is nothing pending.
func (r *Rollup) Flush() {
	r.mu.Lock()
	if r.timer != nil {
		r.timer.Stop()
		r.timer = nil
	}
	r.mu.Unlock()
	r.fire()
}

func (r *Rollup) fire() {
	r.mu.Lock()
	d := r.pending
	r.pending = nil
	r.timer = nil
	r.mu.Unlock()

	if d != nil && (len(d.Opened) > 0 || len(d.Closed) > 0) {
		r.flush(*d)
	}
}
