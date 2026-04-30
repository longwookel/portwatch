package backoff_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/user/portwatch/internal/backoff"
)

var errTransient = errors.New("transient")

func TestDoSucceedsImmediately(t *testing.T) {
	p := backoff.Policy{
		InitialInterval: time.Millisecond,
		MaxInterval:     10 * time.Millisecond,
		Multiplier:      2.0,
		MaxAttempts:     3,
	}
	calls := 0
	err := p.Do(context.Background(), func() error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestDoRetriesUpToMaxAttempts(t *testing.T) {
	p := backoff.Policy{
		InitialInterval: time.Millisecond,
		MaxInterval:     5 * time.Millisecond,
		Multiplier:      1.5,
		MaxAttempts:     4,
	}
	calls := 0
	err := p.Do(context.Background(), func() error {
		calls++
		return errTransient
	})
	if !errors.Is(err, errTransient) {
		t.Fatalf("expected errTransient, got %v", err)
	}
	if calls != 4 {
		t.Fatalf("expected 4 calls, got %d", calls)
	}
}

func TestDoSucceedsOnSecondAttempt(t *testing.T) {
	p := backoff.Policy{
		InitialInterval: time.Millisecond,
		MaxInterval:     10 * time.Millisecond,
		Multiplier:      2.0,
		MaxAttempts:     5,
	}
	calls := 0
	err := p.Do(context.Background(), func() error {
		calls++
		if calls < 2 {
			return errTransient
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if calls != 2 {
		t.Fatalf("expected 2 calls, got %d", calls)
	}
}

func TestDoCancelledContext(t *testing.T) {
	p := backoff.Policy{
		InitialInterval: 50 * time.Millisecond,
		MaxInterval:     500 * time.Millisecond,
		Multiplier:      2.0,
		MaxAttempts:     0, // unlimited
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	err := p.Do(ctx, func() error { return errTransient })
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected DeadlineExceeded, got %v", err)
	}
}

func TestDefaultPolicyHasSaneValues(t *testing.T) {
	p := backoff.Default()
	if p.InitialInterval <= 0 {
		t.Error("InitialInterval must be positive")
	}
	if p.MaxInterval < p.InitialInterval {
		t.Error("MaxInterval must be >= InitialInterval")
	}
	if p.Multiplier <= 1.0 {
		t.Error("Multiplier must be > 1")
	}
	if p.MaxAttempts <= 0 {
		t.Error("MaxAttempts must be positive for Default policy")
	}
}
