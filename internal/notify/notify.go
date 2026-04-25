// Package notify provides pluggable notification backends for portwatch.
// Backends implement the Sender interface and can be composed via Multi.
package notify

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// Event represents a single notification payload.
type Event struct {
	Title   string
	Message string
	Level   Level
}

// Level indicates the severity of a notification.
type Level int

const (
	LevelInfo  Level = iota
	LevelWarn
	LevelError
)

func (l Level) String() string {
	switch l {
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// Sender is the interface that notification backends must implement.
type Sender interface {
	Send(e Event) error
}

// WriterSender writes events as formatted lines to an io.Writer.
type WriterSender struct {
	w io.Writer
}

// NewWriterSender returns a Sender that writes to w.
// If w is nil, os.Stdout is used.
func NewWriterSender(w io.Writer) *WriterSender {
	if w == nil {
		w = os.Stdout
	}
	return &WriterSender{w: w}
}

// Send formats the event and writes it to the underlying writer.
func (ws *WriterSender) Send(e Event) error {
	_, err := fmt.Fprintf(ws.w, "[%s] %s: %s\n", e.Level, e.Title, e.Message)
	return err
}

// Multi fans a single event out to multiple Senders.
// All senders are attempted; errors are joined and returned.
type Multi struct {
	senders []Sender
}

// NewMulti returns a Multi that dispatches to all provided senders.
func NewMulti(senders ...Sender) *Multi {
	return &Multi{senders: senders}
}

// Send delivers the event to every registered sender.
func (m *Multi) Send(e Event) error {
	var errs []string
	for _, s := range m.senders {
		if err := s.Send(e); err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("notify: %s", strings.Join(errs, "; "))
	}
	return nil
}
