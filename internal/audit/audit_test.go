package audit_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/audit"
	"github.com/user/portwatch/internal/snapshot"
)

func TestRecordOpenedWritesEntry(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)

	diff := snapshot.Diff{Opened: []uint16{8080, 9090}}
	if err := l.Record(diff); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}

	var entry audit.Entry
	if err := json.Unmarshal([]byte(lines[0]), &entry); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if entry.Event != "opened" {
		t.Errorf("expected event=opened, got %s", entry.Event)
	}
	if len(entry.Ports) != 2 {
		t.Errorf("expected 2 ports, got %d", len(entry.Ports))
	}
}

func TestRecordClosedWritesEntry(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)

	diff := snapshot.Diff{Closed: []uint16{443}}
	if err := l.Record(diff); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var entry audit.Entry
	if err := json.Unmarshal([]byte(strings.TrimSpace(buf.String())), &entry); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if entry.Event != "closed" {
		t.Errorf("expected event=closed, got %s", entry.Event)
	}
}

func TestRecordBothDirectionsWritesTwoLines(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)

	diff := snapshot.Diff{Opened: []uint16{80}, Closed: []uint16{8080}}
	if err := l.Record(diff); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
}

func TestRecordEmptyDiffWritesNothing(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)

	if err := l.Record(snapshot.Diff{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected no output for empty diff")
	}
}

func TestEntryTimestampIsUTC(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)

	before := time.Now().UTC()
	_ = l.Record(snapshot.Diff{Opened: []uint16{22}})
	after := time.Now().UTC()

	var entry audit.Entry
	_ = json.Unmarshal([]byte(strings.TrimSpace(buf.String())), &entry)

	if entry.Timestamp.Before(before) || entry.Timestamp.After(after) {
		t.Errorf("timestamp %v not in expected range", entry.Timestamp)
	}
}

func TestNewDefaultsToStdout(t *testing.T) {
	// Ensure New(nil) does not panic.
	l := audit.New(nil)
	if l == nil {
		t.Fatal("expected non-nil logger")
	}
}
