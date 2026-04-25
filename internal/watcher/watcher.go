// Package watcher provides file-based state watching for portwatch,
// tracking when the snapshot file is modified on disk.
package watcher

import (
	"os"
	"time"
)

// FileWatcher watches a file for modification changes.
type FileWatcher struct {
	path     string
	lastMod  time.Time
	interval time.Duration
}

// New creates a FileWatcher for the given path, polling at the given interval.
func New(path string, interval time.Duration) (*FileWatcher, error) {
	fw := &FileWatcher{
		path:     path,
		interval: interval,
	}
	info, err := os.Stat(path)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	if err == nil {
		fw.lastMod = info.ModTime()
	}
	return fw, nil
}

// Changed reports whether the watched file has been modified since the last
// call to Changed (or since construction). It updates the internal timestamp
// on each invocation.
func (fw *FileWatcher) Changed() (bool, error) {
	info, err := os.Stat(fw.path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	mod := info.ModTime()
	if mod.After(fw.lastMod) {
		fw.lastMod = mod
		return true, nil
	}
	return false, nil
}

// Path returns the watched file path.
func (fw *FileWatcher) Path() string {
	return fw.path
}

// Interval returns the polling interval.
func (fw *FileWatcher) Interval() time.Duration {
	return fw.interval
}
