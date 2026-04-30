package sampler

import (
	"testing"
	"time"
)

func defaultSampler() *Sampler {
	return New(Config{
		Initial:  30 * time.Second,
		Min:      5 * time.Second,
		Max:      2 * time.Minute,
		StepUp:   10 * time.Second,
		StepDown: 10 * time.Second,
	})
}

func TestInitialCurrentMatchesConfig(t *testing.T) {
	s := defaultSampler()
	if got := s.Current(); got != 30*time.Second {
		t.Fatalf("expected 30s, got %v", got)
	}
}

func TestNextQuietIncreasesInterval(t *testing.T) {
	s := defaultSampler()
	next := s.Next() // no RecordChange → quiet
	if next != 40*time.Second {
		t.Fatalf("expected 40s, got %v", next)
	}
}

func TestNextActiveDecreasesInterval(t *testing.T) {
	s := defaultSampler()
	s.RecordChange()
	next := s.Next()
	if next != 20*time.Second {
		t.Fatalf("expected 20s, got %v", next)
	}
}

func TestNextDoesNotDropBelowMin(t *testing.T) {
	s := New(Config{
		Initial:  6 * time.Second,
		Min:      5 * time.Second,
		Max:      2 * time.Minute,
		StepUp:   10 * time.Second,
		StepDown: 10 * time.Second,
	})
	s.RecordChange()
	next := s.Next()
	if next != 5*time.Second {
		t.Fatalf("expected min 5s, got %v", next)
	}
}

func TestNextDoesNotExceedMax(t *testing.T) {
	s := New(Config{
		Initial:  115 * time.Second,
		Min:      5 * time.Second,
		Max:      2 * time.Minute,
		StepUp:   10 * time.Second,
		StepDown: 10 * time.Second,
	})
	next := s.Next()
	if next != 2*time.Minute {
		t.Fatalf("expected max 2m, got %v", next)
	}
}

func TestRecordChangeIsConsumedAfterNext(t *testing.T) {
	s := defaultSampler()
	s.RecordChange()
	s.Next() // consumes the change flag
	next := s.Next() // should now be quiet → step up
	if next != 30*time.Second {
		// after first Next: 30-10=20; second Next (quiet): 20+10=30
		t.Fatalf("expected 30s after flag reset, got %v", next)
	}
}

func TestDefaultConfigSaneValues(t *testing.T) {
	cfg := Default()
	if cfg.Min >= cfg.Max {
		t.Fatal("Min must be less than Max")
	}
	if cfg.Initial < cfg.Min || cfg.Initial > cfg.Max {
		t.Fatal("Initial must be between Min and Max")
	}
}
