package digest_test

import (
	"testing"

	"github.com/yourorg/portwatch/internal/digest"
)

func TestComputeEmptySlice(t *testing.T) {
	d := digest.Compute(nil)
	if d != digest.Empty {
		t.Fatalf("expected Empty digest, got %s", d)
	}
}

func TestComputeEmptySliceExplicit(t *testing.T) {
	d := digest.Compute([]uint16{})
	if d != digest.Empty {
		t.Fatalf("expected Empty digest for empty slice, got %s", d)
	}
}

func TestComputeIsDeterministic(t *testing.T) {
	ports := []uint16{80, 443, 8080}
	a := digest.Compute(ports)
	b := digest.Compute(ports)
	if !digest.Equal(a, b) {
		t.Fatalf("same input produced different digests: %s vs %s", a, b)
	}
}

func TestComputeIsOrderIndependent(t *testing.T) {
	a := digest.Compute([]uint16{80, 443, 8080})
	b := digest.Compute([]uint16{8080, 80, 443})
	if !digest.Equal(a, b) {
		t.Fatalf("order should not affect digest: %s vs %s", a, b)
	}
}

func TestComputeDifferentPortsDifferentDigest(t *testing.T) {
	a := digest.Compute([]uint16{80})
	b := digest.Compute([]uint16{443})
	if digest.Equal(a, b) {
		t.Fatal("different port sets should produce different digests")
	}
}

func TestComputeDoesNotMutateInput(t *testing.T) {
	ports := []uint16{9000, 80, 443}
	original := make([]uint16, len(ports))
	copy(original, ports)
	digest.Compute(ports)
	for i, p := range ports {
		if p != original[i] {
			t.Fatalf("input slice was mutated at index %d", i)
		}
	}
}

func TestEqualSymmetric(t *testing.T) {
	a := digest.Compute([]uint16{22, 80})
	b := digest.Compute([]uint16{22, 80})
	if !digest.Equal(a, b) {
		t.Fatal("Equal should be symmetric")
	}
	if !digest.Equal(b, a) {
		t.Fatal("Equal should be symmetric (reversed)")
	}
}
