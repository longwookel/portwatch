package circuitbreaker

import (
	"sync"
	"testing"
	"time"
)

func newFakeClock(t time.Time) (func() time.Time, func(d time.Duration)) {
	var mu sync.Mutex
	current := t
	return func() time.Time {
			mu.Lock()
			defer mu.Unlock()
			return current
		}, func(d time.Duration) {
			mu.Lock()
			defer mu.Unlock()
			current = current.Add(d)
		}
}

func TestAllowWhenClosed(t *testing.T) {
	b := New(3, 5*time.Second)
	if err := b.Allow(); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestTripsAfterThreshold(t *testing.T) {
	clock, _ := newFakeClock(time.Now())
	b := newWithClock(3, 5*time.Second, clock)

	b.RecordFailure()
	b.RecordFailure()
	if b.State() != StateClosed {
		t.Fatal("should still be closed after 2 failures")
	}
	b.RecordFailure()
	if b.State() != StateOpen {
		t.Fatalf("expected open, got %s", b.State())
	}
}

func TestOpenBlocksRequests(t *testing.T) {
	clock, _ := newFakeClock(time.Now())
	b := newWithClock(1, 5*time.Second, clock)
	b.RecordFailure()

	if err := b.Allow(); err != ErrOpen {
		t.Fatalf("expected ErrOpen, got %v", err)
	}
}

func TestTransitionsToHalfOpenAfterTimeout(t *testing.T) {
	clock, advance := newFakeClock(time.Now())
	b := newWithClock(1, 5*time.Second, clock)
	b.RecordFailure()

	advance(6 * time.Second)
	if err := b.Allow(); err != nil {
		t.Fatalf("expected nil in half-open, got %v", err)
	}
	if b.State() != StateHalfOpen {
		t.Fatalf("expected half-open, got %s", b.State())
	}
}

func TestSuccessClosesBreakerFromHalfOpen(t *testing.T) {
	clock, advance := newFakeClock(time.Now())
	b := newWithClock(1, 5*time.Second, clock)
	b.RecordFailure()
	advance(6 * time.Second)
	b.Allow() // transition to half-open
	b.RecordSuccess()

	if b.State() != StateClosed {
		t.Fatalf("expected closed after success, got %s", b.State())
	}
}

func TestFailureInHalfOpenReopens(t *testing.T) {
	clock, advance := newFakeClock(time.Now())
	b := newWithClock(1, 5*time.Second, clock)
	b.RecordFailure()
	advance(6 * time.Second)
	b.Allow() // transition to half-open
	b.RecordFailure()

	if b.State() != StateOpen {
		t.Fatalf("expected open after failure in half-open, got %s", b.State())
	}
}

func TestStateString(t *testing.T) {
	cases := []struct {
		s    State
		want string
	}{
		{StateClosed, "closed"},
		{StateOpen, "open"},
		{StateHalfOpen, "half-open"},
		{State(99), "unknown"},
	}
	for _, c := range cases {
		if got := c.s.String(); got != c.want {
			t.Errorf("State(%d).String() = %q, want %q", c.s, got, c.want)
		}
	}
}
