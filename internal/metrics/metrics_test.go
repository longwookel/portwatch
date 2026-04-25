package metrics_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/metrics"
)

func TestNewMetricsInitialState(t *testing.T) {
	m := metrics.New()
	snap := m.Snapshot()

	if snap.ScanCount != 0 {
		t.Errorf("expected ScanCount 0, got %d", snap.ScanCount)
	}
	if snap.AlertCount != 0 {
		t.Errorf("expected AlertCount 0, got %d", snap.AlertCount)
	}
	if snap.StartedAt.IsZero() {
		t.Error("expected StartedAt to be set")
	}
}

func TestRecordScanIncrementsCount(t *testing.T) {
	m := metrics.New()
	m.RecordScan(42 * time.Millisecond)
	m.RecordScan(10 * time.Millisecond)

	snap := m.Snapshot()
	if snap.ScanCount != 2 {
		t.Errorf("expected ScanCount 2, got %d", snap.ScanCount)
	}
	if snap.LastScanElapsed != 10*time.Millisecond {
		t.Errorf("expected LastScanElapsed 10ms, got %v", snap.LastScanElapsed)
	}
	if snap.LastScanAt.IsZero() {
		t.Error("expected LastScanAt to be set after RecordScan")
	}
}

func TestRecordAlertAccumulates(t *testing.T) {
	m := metrics.New()
	m.RecordAlert(3)
	m.RecordAlert(2)

	snap := m.Snapshot()
	if snap.AlertCount != 5 {
		t.Errorf("expected AlertCount 5, got %d", snap.AlertCount)
	}
}

func TestUptimeIsPositive(t *testing.T) {
	m := metrics.New()
	time.Sleep(1 * time.Millisecond)
	up := m.Uptime()
	if up <= 0 {
		t.Errorf("expected positive uptime, got %v", up)
	}
}

func TestSnapshotIsIndependent(t *testing.T) {
	m := metrics.New()
	snap1 := m.Snapshot()
	m.RecordScan(5 * time.Millisecond)
	snap2 := m.Snapshot()

	if snap1.ScanCount != 0 {
		t.Errorf("snap1 should not be affected by later mutations")
	}
	if snap2.ScanCount != 1 {
		t.Errorf("snap2 should reflect the new scan")
	}
}
