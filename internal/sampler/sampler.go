// Package sampler provides adaptive scan interval adjustment based on
// recent change activity. When ports change frequently the interval
// shrinks toward MinInterval; when the network is quiet it grows back
// toward MaxInterval.
package sampler

import (
	"sync"
	"time"
)

// Sampler tracks recent scan activity and recommends the next scan interval.
type Sampler struct {
	mu          sync.Mutex
	current     time.Duration
	min         time.Duration
	max         time.Duration
	stepUp      time.Duration // added when quiet
	stepDown     time.Duration // subtracted when active
	recentChange bool
}

// Config holds tuning parameters for the Sampler.
type Config struct {
	Initial  time.Duration
	Min      time.Duration
	Max      time.Duration
	StepUp   time.Duration
	StepDown time.Duration
}

// Default returns a Config with sensible defaults.
func Default() Config {
	return Config{
		Initial:  30 * time.Second,
		Min:      5 * time.Second,
		Max:      5 * time.Minute,
		StepUp:   15 * time.Second,
		StepDown: 10 * time.Second,
	}
}

// New creates a Sampler from the given Config.
func New(cfg Config) *Sampler {
	return &Sampler{
		current:  cfg.Initial,
		min:      cfg.Min,
		max:      cfg.Max,
		stepUp:   cfg.StepUp,
		stepDown: cfg.StepDown,
	}
}

// RecordChange signals that the last scan detected at least one port change.
func (s *Sampler) RecordChange() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.recentChange = true
}

// Next returns the recommended duration until the next scan and resets
// the internal state for the following cycle.
func (s *Sampler) Next() time.Duration {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.recentChange {
		s.current -= s.stepDown
	} else {
		s.current += s.stepUp
	}

	if s.current < s.min {
		s.current = s.min
	}
	if s.current > s.max {
		s.current = s.max
	}

	s.recentChange = false
	return s.current
}

// Current returns the current interval without advancing the state.
func (s *Sampler) Current() time.Duration {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.current
}
