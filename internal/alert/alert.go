package alert

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelAlert Level = "ALERT"
)

// Alert represents a single port change notification.
type Alert struct {
	Timestamp time.Time
	Level     Level
	Message   string
}

// Notifier sends alerts to a configured output.
type Notifier struct {
	out io.Writer
}

// NewNotifier creates a Notifier that writes to the given writer.
// If w is nil, os.Stdout is used.
func NewNotifier(w io.Writer) *Notifier {
	if w == nil {
		w = os.Stdout
	}
	return &Notifier{out: w}
}

// Notify processes a slice of Diff entries and emits alerts for each change.
func (n *Notifier) Notify(diffs []snapshot.Diff) {
	for _, d := range diffs {
		a := buildAlert(d)
		fmt.Fprintf(n.out, "[%s] %s %s\n",
			a.Timestamp.Format(time.RFC3339),
			a.Level,
			a.Message,
		)
	}
}

// buildAlert converts a Diff into an Alert.
func buildAlert(d snapshot.Diff) Alert {
	var level Level
	var msg string

	switch d.Kind {
	case snapshot.DiffOpened:
		level = LevelAlert
		msg = fmt.Sprintf("Port %d/%s newly OPENED (pid %d)", d.Port.Port, d.Port.Proto, d.Port.PID)
	case snapshot.DiffClosed:
		level = LevelInfo
		msg = fmt.Sprintf("Port %d/%s CLOSED (was pid %d)", d.Port.Port, d.Port.Proto, d.Port.PID)
	default:
		level = LevelInfo
		msg = fmt.Sprintf("Port %d/%s changed", d.Port.Port, d.Port.Proto)
	}

	return Alert{
		Timestamp: time.Now(),
		Level:     level,
		Message:   msg,
	}
}
