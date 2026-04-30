// Package backoff implements exponential back-off with optional jitter for
// retrying transient failures (e.g. webhook delivery, file I/O).
package backoff

import (
	"context"
	"math"
	"math/rand"
	"time"
)

// Policy holds the configuration for an exponential back-off sequence.
type Policy struct {
	// InitialInterval is the wait time before the first retry.
	InitialInterval time.Duration
	// MaxInterval caps the computed interval.
	MaxInterval time.Duration
	// Multiplier is the growth factor applied after each attempt.
	Multiplier float64
	// Jitter, when true, adds a random fraction of the interval to spread load.
	Jitter bool
	// MaxAttempts is the maximum number of retries (0 = unlimited).
	MaxAttempts int
}

// Default returns a Policy with sensible defaults suitable for network calls.
func Default() Policy {
	return Policy{
		InitialInterval: 500 * time.Millisecond,
		MaxInterval:     30 * time.Second,
		Multiplier:      2.0,
		Jitter:          true,
		MaxAttempts:     5,
	}
}

// Do executes fn, retrying on non-nil error according to p.
// It respects ctx cancellation between attempts.
// Returns the last error if all attempts are exhausted.
func (p Policy) Do(ctx context.Context, fn func() error) error {
	interval := p.InitialInterval
	for attempt := 0; ; attempt++ {
		err := fn()
		if err == nil {
			return nil
		}
		if p.MaxAttempts > 0 && attempt+1 >= p.MaxAttempts {
			return err
		}
		wait := interval
		if p.Jitter {
			//nolint:gosec // non-cryptographic jitter is intentional
			wait += time.Duration(rand.Int63n(int64(interval) / 2))
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(wait):
		}
		interval = time.Duration(math.Min(
			float64(interval)*p.Multiplier,
			float64(p.MaxInterval),
		))
	}
}
