package cooldown

import (
	"testing"
	"time"
)

// fakeClock is an injectable time source for deterministic tests.
type fakeClock struct{ now time.Time }

func (f *fakeClock) Now() time.Time { return f.now }
func (f *fakeClock) Advance(d time.Duration) { f.now = f.now.Add(d) }

func newFake(base, max time.Duration) (*Cooldown, *fakeClock) {
	clk := &fakeClock{now: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)}
	cd := newWithClock(base, max, clk)
	return cd, clk
}

func TestAllowFirstCallPermitted(t *testing.T) {
	cd, _ := newFake(time.Second, time.Minute)
	if !cd.Allow("port:80") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestAllowBlockedWithinWindow(t *testing.T) {
	cd, clk := newFake(time.Second, time.Minute)
	cd.Allow("port:80")
	clk.Advance(500 * time.Millisecond)
	if cd.Allow("port:80") {
		t.Fatal("expected call within window to be blocked")
	}
}

func TestAllowPermittedAfterWindow(t *testing.T) {
	cd, clk := newFake(time.Second, time.Minute)
	cd.Allow("port:80")
	clk.Advance(2 * time.Second)
	if !cd.Allow("port:80") {
		t.Fatal("expected call after window to be allowed")
	}
}

func TestWindowDoublesOnSuccessiveAlerts(t *testing.T) {
	base := time.Second
	cd, clk := newFake(base, time.Hour)

	// First allow — window becomes base.
	cd.Allow("k")

	// Advance just past base so the second call is allowed.
	clk.Advance(base + time.Millisecond)
	cd.Allow("k") // window now 2*base

	// Advance by base — should still be blocked (window is now 2*base).
	clk.Advance(base)
	if cd.Allow("k") {
		t.Fatal("expected call to be blocked after window doubled")
	}

	// Advance the remainder of the doubled window.
	clk.Advance(base + time.Millisecond)
	if !cd.Allow("k") {
		t.Fatal("expected call to be allowed after doubled window elapsed")
	}
}

func TestWindowCappedAtMax(t *testing.T) {
	base := time.Second
	max := 3 * time.Second
	cd, clk := newFake(base, max)

	for i := 0; i < 5; i++ {
		clk.Advance(max + time.Millisecond)
		cd.Allow("k")
	}

	// Window must not exceed max; blocking within max should hold.
	clk.Advance(max / 2)
	if cd.Allow("k") {
		t.Fatal("expected call to be blocked; window should be capped at max")
	}
}

func TestIndependentKeys(t *testing.T) {
	cd, clk := newFake(time.Second, time.Minute)
	cd.Allow("a")
	cd.Allow("b")
	clk.Advance(500 * time.Millisecond)

	if cd.Allow("a") {
		t.Error("key 'a' should be blocked")
	}
	if cd.Allow("b") {
		t.Error("key 'b' should be blocked")
	}

	clk.Advance(600 * time.Millisecond)
	if !cd.Allow("a") {
		t.Error("key 'a' should be allowed after window")
	}
}

func TestResetClearsState(t *testing.T) {
	cd, _ := newFake(time.Second, time.Minute)
	cd.Allow("port:443")
	cd.Reset("port:443")
	if !cd.Allow("port:443") {
		t.Fatal("expected Allow to succeed after Reset")
	}
}
