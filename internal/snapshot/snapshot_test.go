package snapshot

import (
	"os"
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

var (
	port80  = scanner.PortState{Protocol: "tcp", Port: 80, Address: "127.0.0.1"}
	port443 = scanner.PortState{Protocol: "tcp", Port: 443, Address: "127.0.0.1"}
	port22  = scanner.PortState{Protocol: "tcp", Port: 22, Address: "127.0.0.1"}
)

func TestSaveAndLoad(t *testing.T) {
	snap := New([]scanner.PortState{port80, port443})

	tmp, err := os.CreateTemp("", "portwatch-*.json")
	if err != nil {
		t.Fatal(err)
	}
	tmp.Close()
	defer os.Remove(tmp.Name())

	if err := snap.Save(tmp.Name()); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := Load(tmp.Name())
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if len(loaded.Ports) != 2 {
		t.Errorf("expected 2 ports, got %d", len(loaded.Ports))
	}
}

func TestDiffOpened(t *testing.T) {
	prev := New([]scanner.PortState{port80})
	curr := New([]scanner.PortState{port80, port443})

	opened, closed := Diff(prev, curr)
	if len(opened) != 1 {
		t.Errorf("expected 1 opened port, got %d", len(opened))
	}
	if len(closed) != 0 {
		t.Errorf("expected 0 closed ports, got %d", len(closed))
	}
}

func TestDiffClosed(t *testing.T) {
	prev := New([]scanner.PortState{port80, port22})
	curr := New([]scanner.PortState{port80})

	opened, closed := Diff(prev, curr)
	if len(opened) != 0 {
		t.Errorf("expected 0 opened ports, got %d", len(opened))
	}
	if len(closed) != 1 {
		t.Errorf("expected 1 closed port, got %d", len(closed))
	}
}

func TestDiffNoChange(t *testing.T) {
	prev := New([]scanner.PortState{port80, port443})
	curr := New([]scanner.PortState{port80, port443})

	opened, closed := Diff(prev, curr)
	if len(opened) != 0 || len(closed) != 0 {
		t.Errorf("expected no diff, got opened=%d closed=%d", len(opened), len(closed))
	}
}
