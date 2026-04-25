package filter_test

import (
	"testing"

	"github.com/yourorg/portwatch/internal/filter"
)

func TestAllowNoRules(t *testing.T) {
	f, err := filter.New("", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, port := range []uint16{1, 80, 443, 8080, 65535} {
		if !f.Allow(port) {
			t.Errorf("expected port %d to be allowed with no rules", port)
		}
	}
}

func TestIncludeSinglePort(t *testing.T) {
	f, err := filter.New("443", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !f.Allow(443) {
		t.Error("expected 443 to be allowed")
	}
	if f.Allow(80) {
		t.Error("expected 80 to be excluded (not in include list)")
	}
}

func TestIncludeRange(t *testing.T) {
	f, err := filter.New("8000-8100", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !f.Allow(8000) {
		t.Error("expected 8000 to be allowed")
	}
	if !f.Allow(8050) {
		t.Error("expected 8050 to be allowed")
	}
	if !f.Allow(8100) {
		t.Error("expected 8100 to be allowed")
	}
	if f.Allow(7999) {
		t.Error("expected 7999 to be excluded")
	}
	if f.Allow(8101) {
		t.Error("expected 8101 to be excluded")
	}
}

func TestExcludePort(t *testing.T) {
	f, err := filter.New("", "22,23")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Allow(22) {
		t.Error("expected 22 to be excluded")
	}
	if f.Allow(23) {
		t.Error("expected 23 to be excluded")
	}
	if !f.Allow(80) {
		t.Error("expected 80 to be allowed")
	}
}

func TestIncludeAndExcludeCombined(t *testing.T) {
	// include 80-90 but exclude 85
	f, err := filter.New("80-90", "85")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !f.Allow(80) {
		t.Error("expected 80 to be allowed")
	}
	if f.Allow(85) {
		t.Error("expected 85 to be excluded")
	}
	if !f.Allow(90) {
		t.Error("expected 90 to be allowed")
	}
}

func TestInvalidRange(t *testing.T) {
	_, err := filter.New("9000-8000", "")
	if err == nil {
		t.Error("expected error for inverted range")
	}
}

func TestInvalidPort(t *testing.T) {
	_, err := filter.New("abc", "")
	if err == nil {
		t.Error("expected error for non-numeric port")
	}
}

func TestZeroPort(t *testing.T) {
	_, err := filter.New("0", "")
	if err == nil {
		t.Error("expected error for port 0")
	}
}
