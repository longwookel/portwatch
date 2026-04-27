package rollup_test

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/example/portwatch/internal/rollup"
	"github.com/example/portwatch/internal/scanner"
	"github.com/example/portwatch/internal/snapshot"
)

// TestRollupTimerResetOnLateAdd verifies that a late-arriving Add resets the
// quiet-period timer so the flush is delayed until after the last event.
func TestRollupTimerResetOnLateAdd(t *testing.T) {
	var flushCount int32

	r := rollup.New(50*time.Millisecond, func(d snapshot.Diff) {
		atomic.AddInt32(&flushCount, 1)
	})

	r.Add(snapshot.Diff{Opened: []scanner.PortState{ps(80)}})
	time.Sleep(30 * time.Millisecond) // before window expires
	r.Add(snapshot.Diff{Opened: []scanner.PortState{ps(443)}}) // resets timer
	time.Sleep(30 * time.Millisecond) // still within new window

	if n := atomic.LoadInt32(&flushCount); n != 0 {
		t.Errorf("expected 0 flushes mid-window, got %d", n)
	}

	time.Sleep(40 * time.Millisecond) // now past the reset window

	if n := atomic.LoadInt32(&flushCount); n != 1 {
		t.Errorf("expected 1 flush after window, got %d", n)
	}
}

// TestRollupConcurrentAdds ensures that concurrent Add calls do not race and
// that all ports are represented in the final merged flush.
func TestRollupConcurrentAdds(t *testing.T) {
	resultCh := make(chan snapshot.Diff, 1)

	r := rollup.New(60*time.Millisecond, func(d snapshot.Diff) {
		resultCh <- d
	})

	ports := []uint16{21, 22, 25, 53, 80, 110, 143, 443, 8080, 8443}

	for _, port := range ports {
		port := port
		go r.Add(snapshot.Diff{Opened: []scanner.PortState{ps(port)}})
	}

	select {
	case d := <-resultCh:
		if len(d.Opened) != len(ports) {
			t.Errorf("expected %d opened ports, got %d", len(ports), len(d.Opened))
		}
	case <-time.After(300 * time.Millisecond):
		t.Fatal("timed out waiting for rollup flush")
	}
}
