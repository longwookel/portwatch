package baseline_test

import (
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/user/portwatch/internal/baseline"
)

func tmpPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "baseline.json")
}

func TestNewBaselineIsEmpty(t *testing.T) {
	b := baseline.New(tmpPath(t))
	if got := b.Ports(); len(got) != 0 {
		t.Fatalf("expected empty baseline, got %v", got)
	}
}

func TestContainsAfterSet(t *testing.T) {
	b := baseline.New(tmpPath(t))
	b.Set([]uint16{80, 443, 8080})

	for _, p := range []uint16{80, 443, 8080} {
		if !b.Contains(p) {
			t.Errorf("expected port %d to be in baseline", p)
		}
	}
	if b.Contains(22) {
		t.Error("port 22 should not be in baseline")
	}
}

func TestSaveAndLoad(t *testing.T) {
	path := tmpPath(t)
	b := baseline.New(path)
	b.Set([]uint16{22, 80, 443})

	if err := b.Save(); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := baseline.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	got := loaded.Ports()
	sort.Slice(got, func(i, j int) bool { return got[i] < got[j] })
	want := []uint16{22, 80, 443}
	if len(got) != len(want) {
		t.Fatalf("want %v, got %v", want, got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("index %d: want %d, got %d", i, want[i], got[i])
		}
	}
}

func TestLoadMissingFileReturnsErrNotFound(t *testing.T) {
	_, err := baseline.Load("/nonexistent/path/baseline.json")
	if err != baseline.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestLoadInvalidJSONReturnsError(t *testing.T) {
	path := tmpPath(t)
	if err := os.WriteFile(path, []byte("not json"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := baseline.Load(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestSetReplacesExistingPorts(t *testing.T) {
	b := baseline.New(tmpPath(t))
	b.Set([]uint16{80, 443})
	b.Set([]uint16{22})

	if b.Contains(80) {
		t.Error("port 80 should have been replaced")
	}
	if !b.Contains(22) {
		t.Error("port 22 should be present after second Set")
	}
}
