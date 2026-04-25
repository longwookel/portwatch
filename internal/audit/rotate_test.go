package audit_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/audit"
)

func TestRotatingFileWritesData(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.log")

	rf, err := audit.OpenRotating(path, 1024)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer rf.Close()

	msg := "hello audit\n"
	if _, err := rf.Write([]byte(msg)); err != nil {
		t.Fatalf("write: %v", err)
	}

	data, _ := os.ReadFile(path)
	if !strings.Contains(string(data), "hello audit") {
		t.Errorf("expected written data in file")
	}
}

func TestRotatingFileRotatesOnSizeExceeded(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.log")

	// maxBytes = 10 so any realistic write triggers rotation.
	rf, err := audit.OpenRotating(path, 10)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer rf.Close()

	line := strings.Repeat("x", 20) + "\n"
	// First write — file is empty so this goes through then sets written > max.
	if _, err := rf.Write([]byte(line)); err != nil {
		t.Fatalf("first write: %v", err)
	}
	// Second write — should trigger rotation first.
	if _, err := rf.Write([]byte(line)); err != nil {
		t.Fatalf("second write: %v", err)
	}

	if _, err := os.Stat(path + ".1"); os.IsNotExist(err) {
		t.Errorf("expected rotated file %s.1 to exist", path)
	}
}

func TestRotatingFileDefaultMaxBytes(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.log")

	// maxBytes <= 0 should use the default (10 MiB) without error.
	rf, err := audit.OpenRotating(path, 0)
	if err != nil {
		t.Fatalf("open with zero maxBytes: %v", err)
	}
	_ = rf.Close()
}

func TestOpenRotatingMissingDir(t *testing.T) {
	_, err := audit.OpenRotating("/nonexistent/dir/audit.log", 1024)
	if err == nil {
		t.Error("expected error for missing directory")
	}
}
