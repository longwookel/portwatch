// Package history maintains a rolling log of port change events
// so that users can review what changed over previous daemon cycles.
package history

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

// Event records a single diff event persisted to the history log.
type Event struct {
	Timestamp time.Time `json:"timestamp"`
	Opened    []uint16  `json:"opened,omitempty"`
	Closed    []uint16  `json:"closed,omitempty"`
}

// Log is a bounded, append-only history of port change events.
type Log struct {
	mu      sync.Mutex
	events  []Event
	maxSize int
	path    string
}

// New creates a Log that keeps at most maxSize events and persists them to path.
func New(path string, maxSize int) *Log {
	if maxSize <= 0 {
		maxSize = 100
	}
	return &Log{path: path, maxSize: maxSize}
}

// Record appends a new event derived from a snapshot.Diff result.
// Events with no opened or closed ports are silently dropped.
func (l *Log) Record(d snapshot.Diff) error {
	if len(d.Opened) == 0 && len(d.Closed) == 0 {
		return nil
	}

	e := Event{
		Timestamp: time.Now().UTC(),
		Opened:    portsToSlice(d.Opened),
		Closed:    portsToSlice(d.Closed),
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	l.events = append(l.events, e)
	if len(l.events) > l.maxSize {
		l.events = l.events[len(l.events)-l.maxSize:]
	}

	return l.flush()
}

// Events returns a copy of the current in-memory event slice.
func (l *Log) Events() []Event {
	l.mu.Lock()
	defer l.mu.Unlock()
	out := make([]Event, len(l.events))
	copy(out, l.events)
	return out
}

// flush writes the current event list to disk as newline-delimited JSON.
// Caller must hold l.mu.
func (l *Log) flush() error {
	if l.path == "" {
		return nil
	}
	f, err := os.Create(l.path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	for _, e := range l.events {
		if err := enc.Encode(e); err != nil {
			return err
		}
	}
	return nil
}

func portsToSlice(m map[uint16]struct{}) []uint16 {
	if len(m) == 0 {
		return nil
	}
	out := make([]uint16, 0, len(m))
	for p := range m {
		out = append(out, p)
	}
	return out
}
