package alert

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/snapshot"
)

func TestNotifyOpened(t *testing.T) {
	var buf bytes.Buffer
	n := NewNotifier(&buf)

	diffs := []snapshot.Diff{
		{
			Kind: snapshot.DiffOpened,
			Port: scanner.PortState{Port: 8080, Proto: "tcp", PID: 1234},
		},
	}

	n.Notify(diffs)

	out := buf.String()
	if !strings.Contains(out, "ALERT") {
		t.Errorf("expected ALERT level in output, got: %s", out)
	}
	if !strings.Contains(out, "8080") {
		t.Errorf("expected port 8080 in output, got: %s", out)
	}
	if !strings.Contains(out, "OPENED") {
		t.Errorf("expected OPENED keyword in output, got: %s", out)
	}
}

func TestNotifyClosed(t *testing.T) {
	var buf bytes.Buffer
	n := NewNotifier(&buf)

	diffs := []snapshot.Diff{
		{
			Kind: snapshot.DiffClosed,
			Port: scanner.PortState{Port: 443, Proto: "tcp", PID: 999},
		},
	}

	n.Notify(diffs)

	out := buf.String()
	if !strings.Contains(out, "INFO") {
		t.Errorf("expected INFO level in output, got: %s", out)
	}
	if !strings.Contains(out, "443") {
		t.Errorf("expected port 443 in output, got: %s", out)
	}
	if !strings.Contains(out, "CLOSED") {
		t.Errorf("expected CLOSED keyword in output, got: %s", out)
	}
}

func TestNotifyEmpty(t *testing.T) {
	var buf bytes.Buffer
	n := NewNotifier(&buf)

	n.Notify([]snapshot.Diff{})

	if buf.Len() != 0 {
		t.Errorf("expected no output for empty diffs, got: %s", buf.String())
	}
}

func TestNewNotifierDefaultsToStdout(t *testing.T) {
	n := NewNotifier(nil)
	if n.out == nil {
		t.Error("expected non-nil writer when nil passed to NewNotifier")
	}
}
