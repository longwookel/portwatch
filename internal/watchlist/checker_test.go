package watchlist_test

import (
	"testing"

	"github.com/user/portwatch/internal/snapshot"
	"github.com/user/portwatch/internal/watchlist"
)

func makeDiff(opened, closed []uint16) snapshot.Diff {
	return snapshot.Diff{Opened: opened, Closed: closed}
}

func TestCheckNoViolationsWhenEmpty(t *testing.T) {
	wl := watchlist.New()
	vs := wl.Check(makeDiff(nil, nil))
	if len(vs) != 0 {
		t.Fatalf("expected no violations, got %d", len(vs))
	}
}

func TestCheckUnexpectedOpen(t *testing.T) {
	wl := watchlist.New()
	wl.Add(80)
	// Port 9000 opened but is not watchlisted.
	vs := wl.Check(makeDiff([]uint16{9000}, nil))
	if len(vs) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(vs))
	}
	if vs[0].Port != 9000 {
		t.Errorf("expected port 9000, got %d", vs[0].Port)
	}
	if vs[0].Reason != watchlist.UnexpectedOpen {
		t.Errorf("expected UnexpectedOpen, got %v", vs[0].Reason)
	}
}

func TestCheckWatchlistedOpenIsAllowed(t *testing.T) {
	wl := watchlist.New()
	wl.Add(443)
	vs := wl.Check(makeDiff([]uint16{443}, nil))
	if len(vs) != 0 {
		t.Fatalf("expected no violations for watchlisted port, got %d", len(vs))
	}
}

func TestCheckExpectedClosed(t *testing.T) {
	wl := watchlist.New()
	wl.Add(22)
	vs := wl.Check(makeDiff(nil, []uint16{22}))
	if len(vs) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(vs))
	}
	if vs[0].Reason != watchlist.ExpectedClosed {
		t.Errorf("expected ExpectedClosed, got %v", vs[0].Reason)
	}
}

func TestCheckNonWatchlistedClosedIsIgnored(t *testing.T) {
	wl := watchlist.New()
	wl.Add(22)
	// Port 8080 closed but was never watchlisted — not a violation.
	vs := wl.Check(makeDiff(nil, []uint16{8080}))
	if len(vs) != 0 {
		t.Fatalf("expected no violations, got %d", len(vs))
	}
}

func TestCheckMixedViolations(t *testing.T) {
	wl := watchlist.New()
	wl.Add(22)
	wl.Add(80)
	diff := makeDiff([]uint16{9999}, []uint16{22})
	vs := wl.Check(diff)
	if len(vs) != 2 {
		t.Fatalf("expected 2 violations, got %d", len(vs))
	}
}

func TestViolationReasonString(t *testing.T) {
	if watchlist.UnexpectedOpen.String() != "unexpected_open" {
		t.Errorf("unexpected string for UnexpectedOpen")
	}
	if watchlist.ExpectedClosed.String() != "expected_closed" {
		t.Errorf("unexpected string for ExpectedClosed")
	}
}
