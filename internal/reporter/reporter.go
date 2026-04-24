// Package reporter provides formatted output for port change reports.
package reporter

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

// Format represents the output format for reports.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Reporter writes port change summaries to an output destination.
type Reporter struct {
	out    io.Writer
	format Format
}

// New creates a Reporter writing to w in the given format.
// If w is nil, os.Stdout is used.
func New(w io.Writer, format Format) *Reporter {
	if w == nil {
		w = os.Stdout
	}
	return &Reporter{out: w, format: format}
}

// Report writes a summary of the given diff to the reporter's output.
func (r *Reporter) Report(diff snapshot.Diff) error {
	if len(diff.Opened) == 0 && len(diff.Closed) == 0 {
		return nil
	}
	switch r.format {
	case FormatJSON:
		return r.writeJSON(diff)
	default:
		return r.writeText(diff)
	}
}

func (r *Reporter) writeText(diff snapshot.Diff) error {
	timestamp := time.Now().Format(time.RFC3339)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("[%s] Port changes detected:\n", timestamp))
	for _, p := range diff.Opened {
		sb.WriteString(fmt.Sprintf("  + OPENED  %s/%d\n", p.Protocol, p.Port))
	}
	for _, p := range diff.Closed {
		sb.WriteString(fmt.Sprintf("  - CLOSED  %s/%d\n", p.Protocol, p.Port))
	}
	_, err := fmt.Fprint(r.out, sb.String())
	return err
}

func (r *Reporter) writeJSON(diff snapshot.Diff) error {
	timestamp := time.Now().Format(time.RFC3339)
	opened := formatPorts(diff.Opened)
	closed := formatPorts(diff.Closed)
	line := fmt.Sprintf(
		`{"timestamp":%q,"opened":[%s],"closed":[%s]}\n`,
		timestamp, opened, closed,
	)
	_, err := fmt.Fprint(r.out, line)
	return err
}

func formatPorts(ports []snapshot.PortState) string {
	parts := make([]string, len(ports))
	for i, p := range ports {
		parts[i] = fmt.Sprintf(`{"protocol":%q,"port":%d}`, p.Protocol, p.Port)
	}
	return strings.Join(parts, ",")
}
