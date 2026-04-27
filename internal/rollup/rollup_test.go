package rollup_test

import (
	"sync"
	"testing"
	"time"

	"github.com/example/portwatch/internal/rollup"
	"github.com/example/portwatch/internal/scanner"
	"github.com/example/portwatch/internal/snapshot"
)

func collectFlusher(t *testing.T) (func(snapshot.Diff), *[]snapshot.Diff, *sync.Mutex) {
	t.Helper()
	var mu sync.Mutex
	var got []snapshot.Diff
	return func(d snapshot.Diff) {
		mu.Lock()
		got = append(got, d)
		mu.Unlock()
	}, &got, &mu
}

func ps(port uint16) scanner.PortState {
	return scanner.PortState{Port: port, Open: true}
}

func TestSingleDiffFlushedAfterWindow(t *testing.T) {
	f, got, mu := collectFlusher(t)
	r := rollup.New(30*time.Millisecond, f)

	r.Add(snapshot.Diff{Opened: []scanner.PortState{ps(80)}})

	time.Sleep(60 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if len(*got) != 1 {
		t.Fatalf("expected 1 flush, got %d", len(*got))
	}
	if len((*got)[0].Opened) != 1 || (*got)[0].Opened[0].Port != 80 {
		t.Errorf("unexpected diff: %+v", (*got)[0])
	}
}

func TestMultipleDiffsMergedIntoOne(t *testing.T) {
	f, got, mu := collectFlusher(t)
	r := rollup.New(40*time.Millisecond, f)

	r.Add(snapshot.Diff{Opened: []scanner.PortState{ps(80)}})
	r.Add(snapshot.Diff{Opened: []scanner.PortState{ps(443)}})
	r.Add(snapshot.Diff{Closed: []scanner.PortState{ps(8080)}})

	time.Sleep(80 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if len(*got) != 1 {
		t.Fatalf("expected 1 merged flush, got %d", len(*got))
	}
	if len((*got)[0].Opened) != 2 {
		t.Errorf("expected 2 opened ports, got %d", len((*got)[0].Opened))
	}
	if len((*got)[0].Closed) != 1 {
		t.Errorf("expected 1 closed port, got %d", len((*got)[0].Closed))
	}
}

func TestFlushForcesImmediateDelivery(t *testing.T) {
	f, got, mu := collectFlusher(t)
	r := rollup.New(5*time.Second, f) // long window — won't fire naturally

	r.Add(snapshot.Diff{Opened: []scanner.PortState{ps(22)}})
	r.Flush()

	mu.Lock()
	defer mu.Unlock()
	if len(*got) != 1 {
		t.Fatalf("expected 1 flush after Flush(), got %d", len(*got))
	}
}

func TestEmptyDiffIsIgnored(t *testing.T) {
	f, got, mu := collectFlusher(t)
	r := rollup.New(20*time.Millisecond, f)

	r.Add(snapshot.Diff{})
	time.Sleep(40 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if len(*got) != 0 {
		t.Errorf("expected no flush for empty diff, got %d", len(*got))
	}
}

func TestFlushWithNothingPendingIsNoop(t *testing.T) {
	f, got, mu := collectFlusher(t)
	r := rollup.New(20*time.Millisecond, f)

	r.Flush() // nothing pending

	mu.Lock()
	defer mu.Unlock()
	if len(*got) != 0 {
		t.Errorf("expected no flush, got %d", len(*got))
	}
}
