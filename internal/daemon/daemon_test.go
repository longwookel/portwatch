package daemon_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/daemon"
)

func defaultCfg() *config.Config {
	cfg := config.Default()
	cfg.Interval = 50 * time.Millisecond
	return cfg
}

func TestNewDaemon(t *testing.T) {
	cfg := defaultCfg()
	d, err := daemon.New(cfg, t.TempDir()+"/snap.json")
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}
	if d == nil {
		t.Fatal("expected non-nil Daemon")
	}
}

func TestRunCancelStops(t *testing.T) {
	cfg := defaultCfg()
	snapshotFile := t.TempDir() + "/snap.json"

	d, err := daemon.New(cfg, snapshotFile)
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Millisecond)
	defer cancel()

	err = d.Run(ctx)
	if err != context.DeadlineExceeded {
		t.Errorf("expected DeadlineExceeded, got %v", err)
	}
}

func TestRunCreatesSnapshot(t *testing.T) {
	cfg := defaultCfg()
	snapshotFile := t.TempDir() + "/snap.json"

	d, err := daemon.New(cfg, snapshotFile)
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 80*time.Millisecond)
	defer cancel()
	d.Run(ctx) //nolint:errcheck

	if _, err := os.Stat(snapshotFile); os.IsNotExist(err) {
		t.Error("expected snapshot file to exist after run")
	}
}
