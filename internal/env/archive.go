package env

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// ArchiveEntry represents a single archived snapshot of an env file.
type ArchiveEntry struct {
	Timestamp time.Time
	Label     string
	Entries   []Entry
}

// ArchiveResult holds the result of an archive operation.
type ArchiveResult struct {
	Archived int
	Label    string
	Path     string
}

func (r ArchiveResult) Format() string {
	return fmt.Sprintf("archived %d entries to %s (label: %s)", r.Archived, r.Path, r.Label)
}

// Archive saves the given entries to a JSON archive file.
func Archive(entries []Entry, dest, label string) (ArchiveResult, error) {
	ae := ArchiveEntry{
		Timestamp: time.Now().UTC(),
		Label:     label,
		Entries:   entries,
	}

	data, err := json.MarshalIndent(ae, "", "  ")
	if err != nil {
		return ArchiveResult{}, fmt.Errorf("marshal error: %w", err)
	}

	if err := os.WriteFile(dest, data, 0600); err != nil {
		return ArchiveResult{}, fmt.Errorf("write error: %w", err)
	}

	return ArchiveResult{Archived: len(entries), Label: label, Path: dest}, nil
}

// LoadArchive reads an archive file and returns its entry.
func LoadArchive(path string) (ArchiveEntry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return ArchiveEntry{}, fmt.Errorf("read error: %w", err)
	}

	var ae ArchiveEntry
	if err := json.Unmarshal(data, &ae); err != nil {
		return ArchiveEntry{}, fmt.Errorf("unmarshal error: %w", err)
	}

	return ae, nil
}
