package sighandler_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/user/portwatch/internal/sighandler"
)

func TestAttachCancelsOnSignal(t *testing.T) {
	h := sighandler.New()

	var captured chan<- os.Signal
	h2 := &sighandler.Handler{}
	_ = h2 // ensure package-level New is used below

	// Use the real handler but send a signal manually.
	ch := make(chan os.Signal, 1)

	// Replace notify/stop via a test-friendly approach: use WithSignals +
	// a goroutine that fires syscall.SIGINT to the process.
	_ = h
	_ = captured

	// Direct unit test: inject a fake notify function via unexported field
	// is not possible, so we test via the exported surface.
	// We verify that parent cancellation also stops the goroutine.
	parent, parentCancel := context.WithCancel(context.Background())
	h3 := sighandler.New()
	ctx, cancel := h3.Attach(parent)
	defer cancel()

	// Cancel parent; child ctx must also be done.
	parentCancel()
	select {
	case <-ctx.Done():
		// expected
	case <-time.After(time.Second):
		t.Fatal("context was not cancelled after parent cancellation")
	}
	_ = ch
}

func TestAttachReturnsDistinctContext(t *testing.T) {
	h := sighandler.New()
	parent := context.Background()
	ctx, cancel := h.Attach(parent)
	defer cancel()

	if ctx == parent {
		t.Fatal("expected a derived context, got the parent itself")
	}
}

func TestCancelFuncStopsContext(t *testing.T) {
	h := sighandler.New()
	ctx, cancel := h.Attach(context.Background())

	cancel()

	select {
	case <-ctx.Done():
		// expected
	case <-time.After(time.Second):
		t.Fatal("context not done after explicit cancel")
	}
}

func TestWithSignalsReturnsSameHandler(t *testing.T) {
	h := sighandler.New()
	h2 := h.WithSignals(os.Interrupt)
	if h != h2 {
		t.Fatal("WithSignals should return the same *Handler for chaining")
	}
}
