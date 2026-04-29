package digest_test

import (
	"sync"
	"testing"

	"github.com/yourorg/portwatch/internal/digest"
)

func TestChangedFirstCallAlwaysTrue(t *testing.T) {
	c := digest.NewCache()
	d := digest.Compute([]uint16{80, 443})
	if !c.Changed("host1", d) {
		t.Fatal("first call for new key should always report changed")
	}
}

func TestChangedSameDigestReturnsFalse(t *testing.T) {
	c := digest.NewCache()
	d := digest.Compute([]uint16{80})
	c.Changed("host1", d) // seed
	if c.Changed("host1", d) {
		t.Fatal("same digest should not be reported as changed")
	}
}

func TestChangedDifferentDigestReturnsTrue(t *testing.T) {
	c := digest.NewCache()
	d1 := digest.Compute([]uint16{80})
	d2 := digest.Compute([]uint16{443})
	c.Changed("host1", d1)
	if !c.Changed("host1", d2) {
		t.Fatal("different digest should be reported as changed")
	}
}

func TestResetCausesNextChangedToBeTrue(t *testing.T) {
	c := digest.NewCache()
	d := digest.Compute([]uint16{22})
	c.Changed("host1", d)
	c.Reset("host1")
	if !c.Changed("host1", d) {
		t.Fatal("after Reset, next Changed should return true")
	}
}

func TestPeekMissingKeyReturnsFalse(t *testing.T) {
	c := digest.NewCache()
	_, ok := c.Peek("missing")
	if ok {
		t.Fatal("Peek on unknown key should return false")
	}
}

func TestPeekReturnsStoredDigest(t *testing.T) {
	c := digest.NewCache()
	d := digest.Compute([]uint16{8080})
	c.Changed("svc", d)
	got, ok := c.Peek("svc")
	if !ok {
		t.Fatal("Peek should find stored digest")
	}
	if !digest.Equal(got, d) {
		t.Fatalf("Peek returned wrong digest: %s", got)
	}
}

func TestChangedConcurrentAccess(t *testing.T) {
	c := digest.NewCache()
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			d := digest.Compute([]uint16{uint16(n)})
			c.Changed("shared", d)
		}(i)
	}
	wg.Wait() // should not race
}
