// Package metrics tracks runtime statistics for the portwatch daemon,
// such as scan counts, alert counts, and last scan duration.
package metrics

import (
	"sync"
	"time"
)

// Metrics holds runtime counters and timing information.
type Metrics struct {
	mu              sync.RWMutex
	ScanCount       int64
	AlertCount      int64
	LastScanAt      time.Time
	LastScanElapsed time.Duration
	StartedAt       time.Time
}

// New returns a new Metrics instance with StartedAt set to now.
func New() *Metrics {
	return &Metrics{
		StartedAt: time.Now(),
	}
}

// RecordScan records a completed scan, incrementing ScanCount and storing
// the timestamp and elapsed duration.
func (m *Metrics) RecordScan(elapsed time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ScanCount++
	m.LastScanAt = time.Now()
	m.LastScanElapsed = elapsed
}

// RecordAlert increments the alert counter by n.
func (m *Metrics) RecordAlert(n int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.AlertCount += int64(n)
}

// Snapshot returns a copy of the current metrics, safe for reading outside
// the lock.
func (m *Metrics) Snapshot() Metrics {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return Metrics{
		ScanCount:       m.ScanCount,
		AlertCount:      m.AlertCount,
		LastScanAt:      m.LastScanAt,
		LastScanElapsed: m.LastScanElapsed,
		StartedAt:       m.StartedAt,
	}
}

// Uptime returns the duration since the metrics were created.
func (m *Metrics) Uptime() time.Duration {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return time.Since(m.StartedAt)
}
