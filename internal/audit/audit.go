// Package audit provides a structured audit log for port change events,
// recording who (process/user context) triggered a scan and what changed.
package audit

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

// Entry represents a single audit log record.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Event     string    `json:"event"`
	Ports     []uint16  `json:"ports"`
	Hostname  string    `json:"hostname"`
}

// Logger writes audit entries to an io.Writer as newline-delimited JSON.
type Logger struct {
	w        io.Writer
	hostname string
	now      func() time.Time
}

// New returns a Logger that writes to w.
// If w is nil, os.Stdout is used.
func New(w io.Writer) *Logger {
	if w == nil {
		w = os.Stdout
	}
	host, _ := os.Hostname()
	return &Logger{w: w, hostname: host, now: time.Now}
}

// Record writes an audit entry for each direction of change in diff.
func (l *Logger) Record(diff snapshot.Diff) error {
	if len(diff.Opened) > 0 {
		if err := l.write("opened", diff.Opened); err != nil {
			return err
		}
	}
	if len(diff.Closed) > 0 {
		if err := l.write("closed", diff.Closed); err != nil {
			return err
		}
	}
	return nil
}

func (l *Logger) write(event string, ports []uint16) error {
	e := Entry{
		Timestamp: l.now().UTC(),
		Event:     event,
		Ports:     ports,
		Hostname:  l.hostname,
	}
	data, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("audit: marshal: %w", err)
	}
	_, err = fmt.Fprintf(l.w, "%s\n", data)
	return err
}
