// Package sighandler provides OS signal handling for graceful shutdown.
package sighandler

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// WaitFunc is a function that blocks until a shutdown signal is received,
// then cancels the provided context.
type WaitFunc func(ctx context.Context) context.Context

// Handler listens for OS signals and triggers context cancellation.
type Handler struct {
	signals []os.Signal
	notify  func(chan<- os.Signal, ...os.Signal)
	stop    func(chan<- os.Signal)
}

// New returns a Handler that reacts to SIGINT and SIGTERM by default.
func New() *Handler {
	return &Handler{
		signals: []os.Signal{syscall.SIGINT, syscall.SIGTERM},
		notify:  signal.Notify,
		stop:    signal.Stop,
	}
}

// WithSignals overrides the set of signals the Handler listens for.
func (h *Handler) WithSignals(sigs ...os.Signal) *Handler {
	h.signals = sigs
	return h
}

// Attach starts a goroutine that cancels ctx when a watched signal arrives.
// It returns a derived context and its cancel function. The caller should
// defer the returned cancel to free resources.
func (h *Handler) Attach(parent context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(parent)

	ch := make(chan os.Signal, 1)
	h.notify(ch, h.signals...)

	go func() {
		defer h.stop(ch)
		select {
		case <-ch:
			cancel()
		case <-ctx.Done():
		}
	}()

	return ctx, cancel
}
