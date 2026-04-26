// Package baseline manages a trusted set of ports that are expected to be
// open. Ports in the baseline are excluded from alerting; only deviations
// from it are reported.
package baseline

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

// ErrNotFound is returned when the baseline file does not exist.
var ErrNotFound = errors.New("baseline: file not found")

// Baseline holds a set of trusted open ports.
type Baseline struct {
	mu    sync.RWMutex
	ports map[uint16]struct{}
	path  string
}

// New returns an empty Baseline backed by the given file path.
func New(path string) *Baseline {
	return &Baseline{
		ports: make(map[uint16]struct{}),
		path:  path,
	}
}

// Load reads the baseline from disk. Returns ErrNotFound if the file does
// not yet exist (treated as an empty baseline).
func Load(path string) (*Baseline, error) {
	b := New(path)
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return b, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	var ports []uint16
	if err := json.Unmarshal(data, &ports); err != nil {
		return nil, err
	}
	for _, p := range ports {
		b.ports[p] = struct{}{}
	}
	return b, nil
}

// Save persists the current baseline to disk.
func (b *Baseline) Save() error {
	b.mu.RLock()
	defer b.mu.RUnlock()
	ports := make([]uint16, 0, len(b.ports))
	for p := range b.ports {
		ports = append(ports, p)
	}
	data, err := json.MarshalIndent(ports, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(b.path, data, 0o644)
}

// Set replaces the entire baseline with the provided ports.
func (b *Baseline) Set(ports []uint16) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.ports = make(map[uint16]struct{}, len(ports))
	for _, p := range ports {
		b.ports[p] = struct{}{}
	}
}

// Contains reports whether port is part of the trusted baseline.
func (b *Baseline) Contains(port uint16) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	_, ok := b.ports[port]
	return ok
}

// Ports returns a copy of all ports in the baseline.
func (b *Baseline) Ports() []uint16 {
	b.mu.RLock()
	defer b.mu.RUnlock()
	out := make([]uint16, 0, len(b.ports))
	for p := range b.ports {
		out = append(out, p)
	}
	return out
}
