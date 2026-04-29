// Package circuitbreaker implements a simple circuit breaker that prevents
// repeated alerting or scanning when a downstream component is consistently
// failing. It transitions between Closed (normal), Open (tripped), and
// Half-Open (probing) states.
package circuitbreaker

import (
	"errors"
	"sync"
	"time"
)

// ErrOpen is returned by Allow when the circuit is open.
var ErrOpen = errors.New("circuit breaker is open")

// State represents the current state of the circuit breaker.
type State int

const (
	StateClosed   State = iota // normal operation
	StateOpen                  // tripped; requests blocked
	StateHalfOpen              // one probe request allowed
)

// String returns a human-readable state name.
func (s State) String() string {
	switch s {
	case StateClosed:
		return "closed"
	case StateOpen:
		return "open"
	case StateHalfOpen:
		return "half-open"
	default:
		return "unknown"
	}
}

// Breaker is a thread-safe circuit breaker.
type Breaker struct {
	mu           sync.Mutex
	state        State
	failures     int
	threshold    int
	resetTimeout time.Duration
	openedAt     time.Time
	clock        func() time.Time
}

// New creates a Breaker that trips after threshold consecutive failures and
// attempts recovery after resetTimeout.
func New(threshold int, resetTimeout time.Duration) *Breaker {
	return newWithClock(threshold, resetTimeout, time.Now)
}

func newWithClock(threshold int, resetTimeout time.Duration, clock func() time.Time) *Breaker {
	return &Breaker{
		threshold:    threshold,
		resetTimeout: resetTimeout,
		clock:        clock,
	}
}

// Allow returns nil if the operation should proceed, or ErrOpen if blocked.
func (b *Breaker) Allow() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	switch b.state {
	case StateClosed:
		return nil
	case StateOpen:
		if b.clock().Sub(b.openedAt) >= b.resetTimeout {
			b.state = StateHalfOpen
			return nil
		}
		return ErrOpen
	case StateHalfOpen:
		return nil
	}
	return nil
}

// RecordSuccess resets the breaker to closed state.
func (b *Breaker) RecordSuccess() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures = 0
	b.state = StateClosed
}

// RecordFailure increments the failure count and trips the breaker if the
// threshold is reached.
func (b *Breaker) RecordFailure() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures++
	if b.state == StateHalfOpen || b.failures >= b.threshold {
		b.state = StateOpen
		b.openedAt = b.clock()
	}
}

// State returns the current state of the breaker.
func (b *Breaker) State() State {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.state
}
