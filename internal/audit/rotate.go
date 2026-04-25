package audit

import (
	"fmt"
	"io"
	"os"
	"sync"
)

// RotatingFile is an io.WriteCloser that rolls over to a new file once
// the current file exceeds maxBytes. Old files are renamed with a .1 suffix.
type RotatingFile struct {
	mu       sync.Mutex
	path     string
	maxBytes int64
	file     *os.File
	written  int64
}

// OpenRotating opens (or creates) the log file at path.
// maxBytes is the size threshold that triggers rotation.
func OpenRotating(path string, maxBytes int64) (*RotatingFile, error) {
	if maxBytes <= 0 {
		maxBytes = 10 * 1024 * 1024 // 10 MiB default
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o640)
	if err != nil {
		return nil, fmt.Errorf("audit rotate: open %s: %w", path, err)
	}
	info, err := f.Stat()
	if err != nil {
		_ = f.Close()
		return nil, fmt.Errorf("audit rotate: stat: %w", err)
	}
	return &RotatingFile{
		path:     path,
		maxBytes: maxBytes,
		file:     f,
		written:  info.Size(),
	}, nil
}

// Write implements io.Writer. It rotates the file when the size threshold
// is exceeded before writing.
func (r *RotatingFile) Write(p []byte) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.written+int64(len(p)) > r.maxBytes {
		if err := r.rotate(); err != nil {
			return 0, err
		}
	}
	n, err := r.file.Write(p)
	r.written += int64(n)
	return n, err
}

// Close closes the underlying file.
func (r *RotatingFile) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.file.Close()
}

// rotate renames the current file to path.1 and opens a fresh file.
func (r *RotatingFile) rotate() error {
	if err := r.file.Close(); err != nil {
		return fmt.Errorf("audit rotate: close: %w", err)
	}
	if err := os.Rename(r.path, r.path+".1"); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("audit rotate: rename: %w", err)
	}
	f, err := os.OpenFile(r.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o640)
	if err != nil {
		return fmt.Errorf("audit rotate: reopen: %w", err)
	}
	r.file = f
	r.written = 0
	return nil
}

// Ensure RotatingFile satisfies io.WriteCloser.
var _ io.WriteCloser = (*RotatingFile)(nil)
