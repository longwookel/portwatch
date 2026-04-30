package dedup

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

// fakeClock is a controllable clock for tests.
type fakeClock struct{ now time.Time }

func (f *fakeClock) Now() time.Time { return f.now }
func (f *fakeClock) Advance(d time.Duration) { f.now = f.now.Add(d) }

func newFake() (*Deduplicator, *fakeClock) {
	c := &fakeClock{now: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)}
	return newWithClock(30*time.Second, c), c
}

func makeDiff(opened, closed []uint16) snapshot.Diff {
	return snapshot.Diff{Opened: opened, Closed: closed}
}

func TestFirstCallIsNotDuplicate(t *testing.T) {
	d, _ := newFake()
	diff := makeDiff([]uint16{8080}, nil)
	if d.IsDuplicate(diff) {
		t.Fatal("first call should not be a duplicate")
	}
}

func TestSecondCallWithinWindowIsDuplicate(t *testing.T) {
	d, c := newFake()
	diff := makeDiff([]uint16{8080}, nil)
	d.IsDuplicate(diff)
	c.Advance(10 * time.Second)
	if !d.IsDuplicate(diff) {
		t.Fatal("same diff within window should be a duplicate")
	}
}

func TestCallAfterWindowIsNotDuplicate(t *testing.T) {
	d, c := newFake()
	diff := makeDiff([]uint16{8080}, nil)
	d.IsDuplicate(diff)
	c.Advance(31 * time.Second)
	if d.IsDuplicate(diff) {
		t.Fatal("same diff after window expiry should not be a duplicate")
	}
}

func TestDistinctDiffsAreIndependent(t *testing.T) {
	d, _ := newFake()
	diff1 := makeDiff([]uint16{80}, nil)
	diff2 := makeDiff([]uint16{443}, nil)
	d.IsDuplicate(diff1)
	if d.IsDuplicate(diff2) {
		t.Fatal("different diffs should not be considered duplicates of each other")
	}
}

func TestEmptyDiffIsNeverDuplicate(t *testing.T) {
	d, _ := newFake()
	empty := makeDiff(nil, nil)
	d.IsDuplicate(empty)
	if d.IsDuplicate(empty) {
		t.Fatal("empty diff should never be treated as a duplicate")
	}
}

func TestPurgeRemovesExpiredEntries(t *testing.T) {
	d, c := newFake()
	diff := makeDiff([]uint16{22}, nil)
	d.IsDuplicate(diff)
	c.Advance(60 * time.Second)
	d.Purge()
	d.mu.Lock()
	l := len(d.cache)
	d.mu.Unlock()
	if l != 0 {
		t.Fatalf("expected empty cache after purge, got %d entries", l)
	}
}

func TestPurgeKeepsActiveEntries(t *testing.T) {
	d, c := newFake()
	diff := makeDiff([]uint16{22}, nil)
	d.IsDuplicate(diff)
	c.Advance(10 * time.Second)
	d.Purge()
	d.mu.Lock()
	l := len(d.cache)
	d.mu.Unlock()
	if l != 1 {
		t.Fatalf("expected 1 active entry after purge, got %d", l)
	}
}
