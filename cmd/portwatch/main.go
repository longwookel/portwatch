// Command portwatch is a lightweight CLI daemon that monitors open ports
// and alerts on unexpected changes.
package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/daemon"
)

func main() {
	cfgPath := flag.String("config", "config.yaml", "path to config file")
	snapshotPath := flag.String("snapshot", "/tmp/portwatch_snapshot.json", "path to snapshot file")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		log.Printf("config load warning: %v — using defaults", err)
		cfg = config.Default()
	}

	d, err := daemon.New(cfg, *snapshotPath)
	if err != nil {
		log.Fatalf("failed to initialise daemon: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := d.Run(ctx); err != nil && err != context.Canceled {
		log.Fatalf("daemon exited with error: %v", err)
	}
}
