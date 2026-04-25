package debounce

import (
	"sync"
	"testing"
	"time"
)

// fakeClock is a controllable clock for deterministic tests.
type fakeClock struct {
	mu  sync.Mutex
	now time.Time
}

func (f *fakeClock) Now() time.Time {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.now
}

func (f *fakeClock) Advance(d time.Duration) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.now = f.now.Add(d)
}

func newFakeClock() *fakeClock {
	return &fakeClock{now: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)}
}

func TestAllowFirstCallPermitted(t *testing.T) {
	clk := newFakeClock()
	d := newWithClock(5*time.Second, clk)

	if !d.Allow("port:8080") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestAllowBlockedWithinQuietPeriod(t *testing.T) {
	clk := newFakeClock()
	d := newWithClock(5*time.Second, clk)

	d.Allow("port:8080")
	clk.Advance(3 * time.Second)

	if d.Allow("port:8080") {
		t.Fatal("expected call within quiet period to be blocked")
	}
}

func TestAllowPermittedAfterQuietPeriod(t *testing.T) {
	clk := newFakeClock()
	d := newWithClock(5*time.Second, clk)

	d.Allow("port:8080")
	clk.Advance(5 * time.Second)

	if !d.Allow("port:8080") {
		t.Fatal("expected call after quiet period to be allowed")
	}
}

func TestAllowIndependentKeys(t *testing.T) {
	clk := newFakeClock()
	d := newWithClock(5*time.Second, clk)

	d.Allow("port:8080")

	if !d.Allow("port:9090") {
		t.Fatal("expected independent key to be allowed")
	}
}

func TestResetClearsKey(t *testing.T) {
	clk := newFakeClock()
	d := newWithClock(5*time.Second, clk)

	d.Allow("port:8080")
	clk.Advance(1 * time.Second)
	d.Reset("port:8080")

	if !d.Allow("port:8080") {
		t.Fatal("expected Allow to return true after Reset")
	}
}

func TestLenTracksKeys(t *testing.T) {
	clk := newFakeClock()
	d := newWithClock(5*time.Second, clk)

	if d.Len() != 0 {
		t.Fatalf("expected 0 keys, got %d", d.Len())
	}
	d.Allow("port:8080")
	d.Allow("port:9090")
	if d.Len() != 2 {
		t.Fatalf("expected 2 keys, got %d", d.Len())
	}
	d.Reset("port:8080")
	if d.Len() != 1 {
		t.Fatalf("expected 1 key after reset, got %d", d.Len())
	}
}
