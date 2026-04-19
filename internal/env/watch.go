package env

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"time"
)

// WatchResult holds the result of a file change detection.
type WatchResult struct {
	File    string
	Changed bool
	OldHash string
	NewHash string
}

func (r WatchResult) Format() string {
	if !r.Changed {
		return fmt.Sprintf("[watch] %s: no changes detected", r.File)
	}
	return fmt.Sprintf("[watch] %s: changed (old=%s, new=%s)", r.File, r.OldHash[:8], r.NewHash[:8])
}

// HashFile returns an MD5 hash of the file contents.
func HashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("cannot open file: %w", err)
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// Watch polls a file for changes at the given interval, calling onChange when a change is detected.
// It stops when the done channel is closed.
func Watch(path string, interval time.Duration, done <-chan struct{}, onChange func(WatchResult)) error {
	current, err := HashFile(path)
	if err != nil {
		return err
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return nil
		case <-ticker.C:
			next, err := HashFile(path)
			if err != nil {
				continue
			}
			if next != current {
				onChange(WatchResult{File: path, Changed: true, OldHash: current, NewHash: next})
				current = next
			}
		}
	}
}
