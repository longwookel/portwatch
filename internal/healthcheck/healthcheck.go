// Package healthcheck provides a simple HTTP endpoint that exposes
// the current health and liveness status of the portwatch daemon.
package healthcheck

import (
	"encoding/json"
	"net/http"
	"sync/atomic"
	"time"
)

// Status holds the current health information reported by the daemon.
type Status struct {
	Alive     bool      `json:"alive"`
	Uptime    string    `json:"uptime"`
	LastScan  time.Time `json:"last_scan"`
	ScanCount int64     `json:"scan_count"`
}

// Server is a lightweight HTTP server that serves a /healthz endpoint.
type Server struct {
	addr      string
	start     time.Time
	lastScan  atomic.Value // stores time.Time
	scanCount atomic.Int64
}

// New creates a new healthcheck Server that will listen on addr.
func New(addr string) *Server {
	s := &Server{
		addr:  addr,
		start: time.Now(),
	}
	s.lastScan.Store(time.Time{})
	return s
}

// RecordScan updates the last-scan timestamp and increments the scan counter.
func (s *Server) RecordScan() {
	s.lastScan.Store(time.Now())
	s.scanCount.Add(1)
}

// ListenAndServe starts the HTTP server. It blocks until the server stops.
func (s *Server) ListenAndServe() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", s.handleHealth)
	return http.ListenAndServe(s.addr, mux)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	status := Status{
		Alive:     true,
		Uptime:    time.Since(s.start).Round(time.Second).String(),
		LastScan:  s.lastScan.Load().(time.Time),
		ScanCount: s.scanCount.Load(),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(status)
}
