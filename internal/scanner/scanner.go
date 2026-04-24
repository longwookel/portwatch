package scanner

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

// PortState represents the state of a single open port.
type PortState struct {
	Protocol string
	Port     int
	Address  string
}

// String returns a human-readable representation of the port state.
func (p PortState) String() string {
	return fmt.Sprintf("%s:%d (%s)", p.Address, p.Port, p.Protocol)
}

// Scanner scans for open ports on the local machine.
type Scanner struct {
	Protocols []string
	PortRange [2]int
}

// NewScanner creates a Scanner with sensible defaults.
func NewScanner() *Scanner {
	return &Scanner{
		Protocols: []string{"tcp", "udp"},
		PortRange: [2]int{1, 65535},
	}
}

// Scan probes the given protocol and port range, returning open ports.
func (s *Scanner) Scan() ([]PortState, error) {
	var open []PortState
	for _, proto := range s.Protocols {
		for port := s.PortRange[0]; port <= s.PortRange[1]; port++ {
			addr := net.JoinHostPort("127.0.0.1", strconv.Itoa(port))
			conn, err := net.Dial(proto, addr)
			if err != nil {
				continue
			}
			conn.Close()
			parts := strings.Split(conn.LocalAddr().String(), ":")
			ip := "127.0.0.1"
			if len(parts) >= 2 {
				ip = strings.Join(parts[:len(parts)-1], ":")
			}
			open = append(open, PortState{
				Protocol: proto,
				Port:     port,
				Address:  ip,
			})
		}
	}
	return open, nil
}
