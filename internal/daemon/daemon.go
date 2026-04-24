// Package daemon provides the main watch loop for portwatch.
package daemon

import (
	"context"
	"log"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/snapshot"
)

// Daemon runs the port-monitoring loop.
type Daemon struct {
	cfg      *config.Config
	scanner  *scanner.Scanner
	notifier *alert.Notifier
	snapshotPath string
}

// New creates a Daemon with the provided configuration.
func New(cfg *config.Config, snapshotPath string) (*Daemon, error) {
	s, err := scanner.NewScanner(cfg.Ports)
	if err != nil {
		return nil, err
	}
	n, err := alert.NewNotifier(cfg)
	if err != nil {
		return nil, err
	}
	return &Daemon{
		cfg:          cfg,
		scanner:      s,
		notifier:     n,
		snapshotPath: snapshotPath,
	}, nil
}

// Run starts the watch loop and blocks until ctx is cancelled.
func (d *Daemon) Run(ctx context.Context) error {
	log.Printf("portwatch started (interval: %s)", d.cfg.Interval)

	if err := d.tick(); err != nil {
		log.Printf("initial scan error: %v", err)
	}

	ticker := time.NewTicker(d.cfg.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("portwatch stopping")
			return ctx.Err()
		case <-ticker.C:
			if err := d.tick(); err != nil {
				log.Printf("scan error: %v", err)
			}
		}
	}
}

func (d *Daemon) tick() error {
	current, err := d.scanner.Scan()
	if err != nil {
		return err
	}

	prev, err := snapshot.Load(d.snapshotPath)
	if err != nil {
		// No previous snapshot — save baseline silently.
		return snapshot.New(current).Save(d.snapshotPath)
	}

	opened, closed := snapshot.Diff(prev, snapshot.New(current))
	if err := d.notifier.Notify(opened, closed); err != nil {
		log.Printf("notify error: %v", err)
	}

	return snapshot.New(current).Save(d.snapshotPath)
}
