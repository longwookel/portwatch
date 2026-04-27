package watchlist_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/watchlist"
)

func tmpPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "watchlist.json")
}

func TestNewWatchlistIsEmpty(t *testing.T) {
	wl := watchlist.New()
	if wl.Len() != 0 {
		t.Fatalf("expected empty watchlist, got len %d", wl.Len())
	}
}

func TestAddAndContains(t *testing.T) {
	wl := watchlist.New()
	wl.Add(80)
	if !wl.Contains(80) {
		t.Fatal("expected port 80 to be in watchlist")
	}
	if wl.Contains(443) {
		t.Fatal("expected port 443 to not be in watchlist")
	}
}

func TestRemove(t *testing.T) {
	wl := watchlist.New()
	wl.Add(8080)
	wl.Remove(8080)
	if wl.Contains(8080) {
		t.Fatal("expected port 8080 to be removed")
	}
}

func TestPortsSorted(t *testing.T) {
	wl := watchlist.New()
	wl.Add(443)
	wl.Add(80)
	wl.Add(22)
	ports := wl.Ports()
	expected := []uint16{22, 80, 443}
	for i, p := range expected {
		if ports[i] != p {
			t.Fatalf("index %d: want %d got %d", i, p, ports[i])
		}
	}
}

func TestSaveAndLoad(t *testing.T) {
	wl := watchlist.New()
	wl.Add(22)
	wl.Add(443)
	path := tmpPath(t)
	if err := wl.Save(path); err != nil {
		t.Fatalf("Save: %v", err)
	}
	loaded, err := watchlist.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if !loaded.Contains(22) || !loaded.Contains(443) {
		t.Fatal("loaded watchlist missing expected ports")
	}
	if loaded.Len() != 2 {
		t.Fatalf("expected 2 ports, got %d", loaded.Len())
	}
}

func TestLoadMissingFileReturnsErrNotFound(t *testing.T) {
	_, err := watchlist.Load("/nonexistent/path.json")
	if err != watchlist.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestLoadInvalidJSONReturnsError(t *testing.T) {
	path := tmpPath(t)
	if err := os.WriteFile(path, []byte("not-json"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := watchlist.Load(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestSaveProducesValidJSON(t *testing.T) {
	wl := watchlist.New()
	wl.Add(8080)
	path := tmpPath(t)
	if err := wl.Save(path); err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(path)
	var ports []uint16
	if err := json.Unmarshal(data, &ports); err != nil {
		t.Fatalf("saved file is not valid JSON: %v", err)
	}
}
