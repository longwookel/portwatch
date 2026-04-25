package ratelimit

import (
	"testing"
	"time"
)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestAllowConsumesTokens(t *testing.T) {
	l := New(1, 3)
	base := time.Now()
	l.clock = fixedClock(base)
	l.lastTick = base

	for i := 0; i < 3; i++ {
		if !l.Allow() {
			t.Fatalf("expected Allow()=true on call %d", i+1)
		}
	}
	if l.Allow() {
		t.Fatal("expected Allow()=false after burst exhausted")
	}
}

func TestAllowRefillsOverTime(t *testing.T) {
	l := New(2, 2) // 2 tokens/sec, burst 2
	base := time.Now()
	l.clock = fixedClock(base)
	l.lastTick = base

	// Drain the bucket.
	l.Allow()
	l.Allow()

	// Advance clock by 1 second → should refill 2 tokens.
	l.clock = fixedClock(base.Add(time.Second))

	if !l.Allow() {
		t.Fatal("expected Allow()=true after refill")
	}
}

func TestAllowDoesNotExceedBurst(t *testing.T) {
	l := New(10, 3) // high refill rate, burst capped at 3
	base := time.Now()
	l.clock = fixedClock(base)
	l.lastTick = base

	// Advance a long time — tokens should not exceed burst.
	l.clock = fixedClock(base.Add(10 * time.Second))

	count := 0
	for l.Allow() {
		count++
	}
	if count != 3 {
		t.Fatalf("expected 3 tokens after refill cap, got %d", count)
	}
}

func TestRemaining(t *testing.T) {
	l := New(1, 5)
	base := time.Now()
	l.clock = fixedClock(base)
	l.lastTick = base

	if r := l.Remaining(); r != 5 {
		t.Fatalf("expected 5 remaining, got %d", r)
	}
	l.Allow()
	l.Allow()
	if r := l.Remaining(); r != 3 {
		t.Fatalf("expected 3 remaining, got %d", r)
	}
}
