package reporter_test

import (
	"strings"
	"testing"

	"github.com/user/portwatch/internal/reporter"
	"github.com/user/portwatch/internal/snapshot"
)

func openedDiff() snapshot.Diff {
	return snapshot.Diff{
		Opened: []snapshot.PortState{{Protocol: "tcp", Port: 8080}},
	}
}

func closedDiff() snapshot.Diff {
	return snapshot.Diff{
		Closed: []snapshot.PortState{{Protocol: "udp", Port: 53}},
	}
}

func TestReportTextOpened(t *testing.T) {
	var buf strings.Builder
	r := reporter.New(&buf, reporter.FormatText)
	if err := r.Report(openedDiff()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "OPENED") {
		t.Errorf("expected OPENED in output, got: %s", out)
	}
	if !strings.Contains(out, "tcp/8080") {
		t.Errorf("expected tcp/8080 in output, got: %s", out)
	}
}

func TestReportTextClosed(t *testing.T) {
	var buf strings.Builder
	r := reporter.New(&buf, reporter.FormatText)
	if err := r.Report(closedDiff()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "CLOSED") {
		t.Errorf("expected CLOSED in output, got: %s", out)
	}
	if !strings.Contains(out, "udp/53") {
		t.Errorf("expected udp/53 in output, got: %s", out)
	}
}

func TestReportJSONOpened(t *testing.T) {
	var buf strings.Builder
	r := reporter.New(&buf, reporter.FormatJSON)
	if err := r.Report(openedDiff()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"opened"`) {
		t.Errorf("expected opened key in JSON, got: %s", out)
	}
	if !strings.Contains(out, `"port":8080`) {
		t.Errorf("expected port 8080 in JSON, got: %s", out)
	}
}

func TestReportEmptyDiffIsNoop(t *testing.T) {
	var buf strings.Builder
	r := reporter.New(&buf, reporter.FormatText)
	if err := r.Report(snapshot.Diff{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected no output for empty diff, got: %s", buf.String())
	}
}

func TestNewReporterDefaultsToStdout(t *testing.T) {
	// Should not panic when w is nil.
	r := reporter.New(nil, reporter.FormatText)
	if r == nil {
		t.Fatal("expected non-nil reporter")
	}
}
