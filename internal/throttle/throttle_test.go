package throttle_test

import (
	"testing"
	"time"

	"github.com/yourorg/portwatch/internal/throttle"
)

// fakeClock returns a function that returns the pointed-to time value,
// allowing tests to advance time manually.
func fakeClock(t *time.Time) throttle.Clock {
	return func() time.Time { return *t }
}

func TestAllowFirstCallPermitted(t *testing.T) {
	now := time.Now()
	th := throttle.New(5*time.Second, fakeClock(&now))

	if !th.Allow("port:80") {
		t.Fatal("expected first Allow to return true")
	}
}

func TestAllowBlockedWithinCooldown(t *testing.T) {
	now := time.Now()
	th := throttle.New(10*time.Second, fakeClock(&now))

	th.Allow("port:443")

	// Advance time but stay within cooldown.
	now = now.Add(5 * time.Second)

	if th.Allow("port:443") {
		t.Fatal("expected Allow to be blocked within cooldown")
	}
}

func TestAllowPermittedAfterCooldown(t *testing.T) {
	now := time.Now()
	th := throttle.New(10*time.Second, fakeClock(&now))

	th.Allow("port:22")

	// Advance past the cooldown.
	now = now.Add(11 * time.Second)

	if !th.Allow("port:22") {
		t.Fatal("expected Allow to be permitted after cooldown expired")
	}
}

func TestAllowIndependentKeys(t *testing.T) {
	now := time.Now()
	th := throttle.New(30*time.Second, fakeClock(&now))

	th.Allow("port:80")

	// A different key should not be affected.
	if !th.Allow("port:8080") {
		t.Fatal("expected independent key to be allowed")
	}
}

func TestReset(t *testing.T) {
	now := time.Now()
	th := throttle.New(60*time.Second, fakeClock(&now))

	th.Allow("port:3306")
	th.Reset("port:3306")

	if !th.Allow("port:3306") {
		t.Fatal("expected Allow after Reset to return true")
	}
}

func TestRemaining(t *testing.T) {
	now := time.Now()
	th := throttle.New(20*time.Second, fakeClock(&now))

	if r := th.Remaining("port:5432"); r != 0 {
		t.Fatalf("expected 0 remaining for unseen key, got %v", r)
	}

	th.Allow("port:5432")
	now = now.Add(8 * time.Second)

	want := 12 * time.Second
	if got := th.Remaining("port:5432"); got != want {
		t.Fatalf("expected remaining=%v, got %v", want, got)
	}
}

func TestRemainingZeroAfterCooldown(t *testing.T) {
	now := time.Now()
	th := throttle.New(5*time.Second, fakeClock(&now))

	th.Allow("port:6379")
	now = now.Add(10 * time.Second)

	if r := th.Remaining("port:6379"); r != 0 {
		t.Fatalf("expected 0 remaining after cooldown, got %v", r)
	}
}
