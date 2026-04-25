package watcher_test

import (
	"os"
	"testing"
	"time"

	"github.com/example/portwatch/internal/watcher"
)

func TestNewWatcher_NonExistentFile(t *testing.T) {
	fw, err := watcher.New("/tmp/portwatch_nonexistent_test_file", time.Second)
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if fw == nil {
		t.Fatal("expected non-nil FileWatcher")
	}
}

func TestChangedReturnsFalseInitially(t *testing.T) {
	f, err := os.CreateTemp("", "portwatch_watch_*.snap")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	f.Close()

	fw, err := watcher.New(f.Name(), time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	changed, err := fw.Changed()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if changed {
		t.Error("expected Changed to be false immediately after construction")
	}
}

func TestChangedDetectsModification(t *testing.T) {
	f, err := os.CreateTemp("", "portwatch_watch_*.snap")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	f.Close()

	fw, err := watcher.New(f.Name(), time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Ensure mtime changes by sleeping briefly then writing.
	time.Sleep(10 * time.Millisecond)
	if err := os.WriteFile(f.Name(), []byte("data"), 0o644); err != nil {
		t.Fatal(err)
	}

	changed, err := fw.Changed()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !changed {
		t.Error("expected Changed to be true after file modification")
	}

	// Second call should return false (mtime unchanged).
	changed, err = fw.Changed()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if changed {
		t.Error("expected Changed to be false on second call without modification")
	}
}

func TestPathAndInterval(t *testing.T) {
	path := "/tmp/portwatch_meta_test"
	interval := 5 * time.Second
	fw, _ := watcher.New(path, interval)
	if fw.Path() != path {
		t.Errorf("expected path %q, got %q", path, fw.Path())
	}
	if fw.Interval() != interval {
		t.Errorf("expected interval %v, got %v", interval, fw.Interval())
	}
}
