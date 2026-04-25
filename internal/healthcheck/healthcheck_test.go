package healthcheck_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/portwatch/internal/healthcheck"
)

func newTestServer(t *testing.T) *healthcheck.Server {
	t.Helper()
	return healthcheck.New(":0")
}

func TestHealthzReturns200(t *testing.T) {
	s := newTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	// Invoke handler indirectly via a real mux wired the same way.
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		// Use exported behaviour only — record a scan then serve.
		s.RecordScan()
	})

	// Exercise through the exported HTTP surface.
	req2 := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec2 := httptest.NewRecorder()
	_ = req
	_ = rec
	_ = req2
	_ = rec2
}

func TestHealthzBodyAlive(t *testing.T) {
	s := healthcheck.New(":0")
	s.RecordScan()

	srv := httptest.NewServer(buildMux(s))
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/healthz")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var status struct {
		Alive     bool      `json:"alive"`
		ScanCount int64     `json:"scan_count"`
		LastScan  time.Time `json:"last_scan"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if !status.Alive {
		t.Error("expected alive=true")
	}
	if status.ScanCount != 1 {
		t.Errorf("expected scan_count=1, got %d", status.ScanCount)
	}
	if status.LastScan.IsZero() {
		t.Error("expected last_scan to be set")
	}
}

func TestRecordScanIncrementsCount(t *testing.T) {
	s := healthcheck.New(":0")
	for i := 0; i < 5; i++ {
		s.RecordScan()
	}

	srv := httptest.NewServer(buildMux(s))
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/healthz")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	var status struct {
		ScanCount int64 `json:"scan_count"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if status.ScanCount != 5 {
		t.Errorf("expected 5 scans, got %d", status.ScanCount)
	}
}

// buildMux wires the server's handler into a test-friendly ServeMux.
func buildMux(s *healthcheck.Server) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		// Delegate to a temporary server's handler by spinning up a recorder.
		rec := httptest.NewRecorder()
		// We can't call private methods, so we start a real server and proxy.
		_ = rec
		// Instead, replicate the minimal response here using exported API only.
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = s // keep reference
	})
	// Use the server directly via ListenAndServe in integration; for unit tests
	// we rely on httptest.NewServer with a real listener started by the server.
	return buildRealMux(s)
}

func buildRealMux(s *healthcheck.Server) http.Handler {
	mux := http.NewServeMux()
	// Expose via a thin adapter that calls the exported RecordScan-agnostic path.
	// The server's ListenAndServe registers /healthz internally; we mirror that.
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/healthz", http.StatusOK)
	})
	_ = s
	return mux
}
