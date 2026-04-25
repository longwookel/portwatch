package suppress

import (
	"fmt"
	"testing"
	"time"
)

type fakeClock struct {
	current time.Time
}

func (f *fakeClock) now() time.Time { return f.current }
func (f *fakeClock) advance(d time.Duration) { f.current = f.current.Add(d) }

func newFake() (*Suppressor, *fakeClock) {
	fc := &fakeClock{current: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)}
	s := newWithClock(5*time.Second, fc.now)
	return s, fc
}

func TestAllowFirstCallPermitted(t *testing.T) {
	s, _ := newFake()
	if !s.Allow("port:80:opened") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestAllowBlockedWithinWindow(t *testing.T) {
	s, fc := newFake()
	s.Allow("port:80:opened")
	fc.advance(3 * time.Second)
	if s.Allow("port:80:opened") {
		t.Fatal("expected call within window to be suppressed")
	}
}

func TestAllowPermittedAfterWindow(t *testing.T) {
	s, fc := newFake()
	s.Allow("port:80:opened")
	fc.advance(6 * time.Second)
	if !s.Allow("port:80:opened") {
		t.Fatal("expected call after window to be allowed")
	}
}

func TestAllowIndependentKeys(t *testing.T) {
	s, _ := newFake()
	if !s.Allow("port:80:opened") {
		t.Fatal("key 80 should be allowed")
	}
	if !s.Allow("port:443:opened") {
		t.Fatal("key 443 should be allowed independently")
	}
}

func TestPruneRemovesExpiredEntries(t *testing.T) {
	s, fc := newFake()
	for i := 0; i < 10; i++ {
		s.Allow(fmt.Sprintf("port:%d:opened", i))
	}
	fc.advance(10 * time.Second)
	// trigger prune via a new Allow call
	s.Allow("port:9999:opened")
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.seen) != 1 {
		t.Fatalf("expected 1 entry after prune, got %d", len(s.seen))
	}
}

func TestResetClearsState(t *testing.T) {
	s, _ := newFake()
	s.Allow("port:80:opened")
	s.Reset()
	if !s.Allow("port:80:opened") {
		t.Fatal("expected allow after reset")
	}
}
