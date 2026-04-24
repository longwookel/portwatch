package snapshot

import (
	"encoding/json"
	"os"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Snapshot holds a timestamped record of open ports.
type Snapshot struct {
	Timestamp time.Time            `json:"timestamp"`
	Ports     []scanner.PortState  `json:"ports"`
}

// New creates a new Snapshot from the given port states.
func New(ports []scanner.PortState) *Snapshot {
	return &Snapshot{
		Timestamp: time.Now().UTC(),
		Ports:     ports,
	}
}

// Save writes the snapshot as JSON to the specified file path.
func (s *Snapshot) Save(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(s)
}

// Load reads a snapshot from a JSON file.
func Load(path string) (*Snapshot, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var snap Snapshot
	if err := json.NewDecoder(f).Decode(&snap); err != nil {
		return nil, err
	}
	return &snap, nil
}

// Diff compares two snapshots and returns newly opened and closed ports.
func Diff(previous, current *Snapshot) (opened, closed []scanner.PortState) {
	prev := index(previous.Ports)
	curr := index(current.Ports)

	for key, ps := range curr {
		if _, exists := prev[key]; !exists {
			opened = append(opened, ps)
		}
	}
	for key, ps := range prev {
		if _, exists := curr[key]; !exists {
			closed = append(closed, ps)
		}
	}
	return
}

func index(ports []scanner.PortState) map[string]scanner.PortState {
	m := make(map[string]scanner.PortState, len(ports))
	for _, p := range ports {
		key := p.Protocol + ":" + p.Address + ":" + string(rune(p.Port))
		m[key] = p
	}
	return m
}
