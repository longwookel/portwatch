package watchlist

import "github.com/user/portwatch/internal/snapshot"

// Violation describes a port that violated watchlist expectations.
type Violation struct {
	Port   uint16
	Reason ViolationReason
}

// ViolationReason classifies why a port was flagged.
type ViolationReason int

const (
	// UnexpectedOpen means a port is open but not on the watchlist.
	UnexpectedOpen ViolationReason = iota
	// ExpectedClosed means a watchlisted port is no longer open.
	ExpectedClosed
)

// String returns a human-readable description of the reason.
func (r ViolationReason) String() string {
	switch r {
	case UnexpectedOpen:
		return "unexpected_open"
	case ExpectedClosed:
		return "expected_closed"
	default:
		return "unknown"
	}
}

// Check compares a snapshot diff against the watchlist and returns any
// violations. Ports that opened but are not watchlisted are flagged as
// UnexpectedOpen; watchlisted ports that closed are flagged as ExpectedClosed.
func (wl *Watchlist) Check(diff snapshot.Diff) []Violation {
	var violations []Violation

	for _, p := range diff.Opened {
		if !wl.Contains(p) {
			violations = append(violations, Violation{
				Port:   p,
				Reason: UnexpectedOpen,
			})
		}
	}

	for _, p := range diff.Closed {
		if wl.Contains(p) {
			violations = append(violations, Violation{
				Port:   p,
				Reason: ExpectedClosed,
			})
		}
	}

	return violations
}
