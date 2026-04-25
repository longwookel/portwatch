package history_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/history"
	"github.com/user/portwatch/internal/snapshot"
)

func makeDiff(opened, closed []uint16) snapshot.Diff {
	d := snapshot.Diff{
		Opened: make(map[uint16]struct{}),
		Closed: make(map[uint16]struct{}),
	}
	for _, p := range opened {
		d.Opened[p] = struct{}{}
	}
	for _, p := range closed {
		d.Closed[p] = struct{}{}
	}
	return d
}

func TestRecordEmptyDiffIsIgnored(t *testing.T) {
	l := history.New("", 10)
	if err := l.Record(makeDiff(nil, nil)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := len(l.Events()); got != 0 {
		t.Fatalf("expected 0 events, got %d", got)
	}
}

func TestRecordStoresEvent(t *testing.T) {
	l := history.New("", 10)
	if err := l.Record(makeDiff([]uint16{80, 443}, nil)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	events := l.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Timestamp.IsZero() {
		t.Error("timestamp should not be zero")
	}
}

func TestRecordRespectsMaxSize(t *testing.T) {
	l := history.New("", 3)
	for i := 0; i < 5; i++ {
		if err := l.Record(makeDiff([]uint16{uint16(8000 + i)}, nil)); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}
	if got := len(l.Events()); got != 3 {
		t.Fatalf("expected 3 events after overflow, got %d", got)
	}
}

func TestRecordPersistsToDisk(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "history.ndjson")

	l := history.New(path, 10)
	if err := l.Record(makeDiff([]uint16{22}, []uint16{8080})); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("could not read history file: %v", err)
	}

	var e history.Event
	if err := json.Unmarshal(data, &e); err != nil {
		t.Fatalf("could not parse persisted event: %v", err)
	}
	if e.Timestamp.IsZero() {
		t.Error("persisted event has zero timestamp")
	}
}
