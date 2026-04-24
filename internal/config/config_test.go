package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/user/portwatch/internal/config"
)

func TestDefaultConfig(t *testing.T) {
	cfg := config.Default()

	if cfg.Interval != 30*time.Second {
		t.Errorf("expected 30s interval, got %v", cfg.Interval)
	}
	if cfg.Ports.RangeMin != 1 || cfg.Ports.RangeMax != 65535 {
		t.Errorf("unexpected default port range: %d-%d", cfg.Ports.RangeMin, cfg.Ports.RangeMax)
	}
	if len(cfg.Ports.Protocols) != 1 || cfg.Ports.Protocols[0] != "tcp" {
		t.Errorf("unexpected default protocols: %v", cfg.Ports.Protocols)
	}
}

func TestLoadValidConfig(t *testing.T) {
	content := `
interval: 1m
snapshot_dir: /var/portwatch
ports:
  protocols: [tcp, udp]
  range_min: 1024
  range_max: 9000
alert:
  output: /var/log/portwatch.log
`
	f := writeTempFile(t, content)

	cfg, err := config.Load(f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Interval != time.Minute {
		t.Errorf("expected 1m, got %v", cfg.Interval)
	}
	if cfg.Ports.RangeMin != 1024 {
		t.Errorf("expected range_min 1024, got %d", cfg.Ports.RangeMin)
	}
}

func TestLoadIntervalTooShort(t *testing.T) {
	content := "interval: 500ms\n"
	f := writeTempFile(t, content)

	_, err := config.Load(f)
	if err != config.ErrIntervalTooShort {
		t.Errorf("expected ErrIntervalTooShort, got %v", err)
	}
}

func TestLoadInvalidPortRange(t *testing.T) {
	content := "ports:\n  range_min: 9000\n  range_max: 1000\n"
	f := writeTempFile(t, content)

	_, err := config.Load(f)
	if err != config.ErrInvalidPortRange {
		t.Errorf("expected ErrInvalidPortRange, got %v", err)
	}
}

func TestLoadMissingFile(t *testing.T) {
	_, err := config.Load("/nonexistent/portwatch.yaml")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "portwatch-cfg-*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	t.Cleanup(func() { os.Remove(f.Name()) })
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}
