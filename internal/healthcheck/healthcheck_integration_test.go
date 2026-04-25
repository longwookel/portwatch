package healthcheck_test

import (
	"encoding/json"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/user/portwatch/internal/healthcheck"
)

func TestListenAndServeIntegration(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("could not bind: %v", err)
	}
	addr := ln.Addr().String()
	ln.Close() // release; server will re-bind

	s := healthcheck.New(addr)
	s.RecordScan()
	s.RecordScan()

	errCh := make(chan error, 1)
	go func() {
		errCh <- s.ListenAndServe()
	}()

	// Give the server a moment to start.
	var resp *http.Response
	for i := 0; i < 20; i++ {
		time.Sleep(20 * time.Millisecond)
		resp, err = http.Get("http://" + addr + "/healthz")
		if err == nil {
			break
		}
	}
	if err != nil {
		t.Fatalf("server did not start: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var status struct {
		Alive     bool   `json:"alive"`
		Uptime    string `json:"uptime"`
		ScanCount int64  `json:"scan_count"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if !status.Alive {
		t.Error("expected alive=true")
	}
	if status.ScanCount != 2 {
		t.Errorf("expected scan_count=2, got %d", status.ScanCount)
	}
	if status.Uptime == "" {
		t.Error("expected non-empty uptime")
	}
}
