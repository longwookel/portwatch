// Package ratelimit provides a simple token-bucket rate limiter used to
// suppress alert floods when many ports change state in a short window.
package ratelimit

import (
	"sync"
	"time"
)

// Limiter is a token-bucket rate limiter.
type Limiter struct {
	mu       sync.Mutex
	tokens   float64
	max      float64
	refill   float64 // tokens per second
	lastTick time.Time
	clock    func() time.Time
}

// New creates a Limiter that allows up to burst events and refills at rate
// tokens per second.
func New(rate float64, burst int) *Limiter {
	return &Limiter{
		tokens:   float64(burst),
		max:      float64(burst),
		refill:   rate,
		lastTick: time.Now(),
		clock:    time.Now,
	}
}

// Allow reports whether an event may proceed. It consumes one token if
// available, refilling the bucket based on elapsed time first.
func (l *Limiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.clock()
	elapsed := now.Sub(l.lastTick).Seconds()
	l.lastTick = now

	l.tokens += elapsed * l.refill
	if l.tokens > l.max {
		l.tokens = l.max
	}

	if l.tokens < 1 {
		return false
	}
	l.tokens--
	return true
}

// Remaining returns the current number of available tokens (floored to zero).
func (l *Limiter) Remaining() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.tokens < 0 {
		return 0
	}
	return int(l.tokens)
}
