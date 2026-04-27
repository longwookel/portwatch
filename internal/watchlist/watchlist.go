// Package watchlist manages a set of ports that are explicitly expected
// to be open. Any port not on the watchlist that appears open will be
// flagged as unexpected, while any watchlisted port that closes will
// also trigger an alert.
package watchlist

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
)

// ErrNotFound is returned when the watchlist file does not exist.
var ErrNotFound = errors.New("watchlist: file not found")

// Watchlist holds a set of expected open ports.
type Watchlist struct {
	ports map[uint16]struct{}
}

// New returns an empty Watchlist.
func New() *Watchlist {
	return &Watchlist{ports: make(map[uint16]struct{})}
}

// Load reads a JSON array of port numbers from path and returns a Watchlist.
func Load(path string) (*Watchlist, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("watchlist: read %s: %w", path, err)
	}
	var ports []uint16
	if err := json.Unmarshal(data, &ports); err != nil {
		return nil, fmt.Errorf("watchlist: parse %s: %w", path, err)
	}
	wl := New()
	for _, p := range ports {
		wl.ports[p] = struct{}{}
	}
	return wl, nil
}

// Save writes the watchlist as a JSON array to path.
func (wl *Watchlist) Save(path string) error {
	data, err := json.MarshalIndent(wl.Ports(), "", "  ")
	if err != nil {
		return fmt.Errorf("watchlist: marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("watchlist: write %s: %w", path, err)
	}
	return nil
}

// Add adds a port to the watchlist.
func (wl *Watchlist) Add(port uint16) { wl.ports[port] = struct{}{} }

// Remove removes a port from the watchlist.
func (wl *Watchlist) Remove(port uint16) { delete(wl.ports, port) }

// Contains reports whether port is in the watchlist.
func (wl *Watchlist) Contains(port uint16) bool {
	_, ok := wl.ports[port]
	return ok
}

// Ports returns a sorted slice of all ports in the watchlist.
func (wl *Watchlist) Ports() []uint16 {
	out := make([]uint16, 0, len(wl.ports))
	for p := range wl.ports {
		out = append(out, p)
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}

// Len returns the number of ports in the watchlist.
func (wl *Watchlist) Len() int { return len(wl.ports) }
