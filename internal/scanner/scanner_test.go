package scanner

import (
	"net"
	"strconv"
	"testing"
)

func TestPortStateString(t *testing.T) {
	ps := PortState{Protocol: "tcp", Port: 8080, Address: "127.0.0.1"}
	expected := "127.0.0.1:8080 (tcp)"
	if ps.String() != expected {
		t.Errorf("expected %q, got %q", expected, ps.String())
	}
}

func TestScanDetectsOpenPort(t *testing.T) {
	// Start a temporary TCP listener on a random port.
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start listener: %v", err)
	}
	defer ln.Close()

	addr := ln.Addr().(*net.TCPAddr)
	port := addr.Port

	s := &Scanner{
		Protocols: []string{"tcp"},
		PortRange: [2]int{port, port},
	}

	results, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan returned error: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 open port, got %d", len(results))
	}
	if results[0].Port != port {
		t.Errorf("expected port %d, got %d", port, results[0].Port)
	}
	if results[0].Protocol != "tcp" {
		t.Errorf("expected protocol tcp, got %s", results[0].Protocol)
	}
}

func TestScanClosedPort(t *testing.T) {
	// Find a free port, then immediately close it so it's not listening.
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to bind: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	ln.Close()

	s := &Scanner{
		Protocols: []string{"tcp"},
		PortRange: [2]int{port, port},
	}

	results, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan returned error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 open ports, got %d (port %s)", len(results), strconv.Itoa(port))
	}
}
