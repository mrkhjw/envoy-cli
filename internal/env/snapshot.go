package env

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot represents a saved state of an env file.
type Snapshot struct {
	Timestamp time.Time         `json:"timestamp"`
	Source    string            `json:"source"`
	Entries   map[string]string `json:"entries"`
}

// SnapshotResult holds the outcome of a snapshot operation.
type SnapshotResult struct {
	Path      string
	KeyCount  int
	Timestamp time.Time
}

func (r SnapshotResult) Format() string {
	return fmt.Sprintf("Snapshot saved to %s (%d keys) at %s", r.Path, r.KeyCount, r.Timestamp.Format(time.RFC3339))
}

// TakeSnapshot parses the given env file and writes a JSON snapshot to dest.
func TakeSnapshot(source, dest string) (SnapshotResult, error) {
	entries, err := ParseFile(source)
	if err != nil {
		return SnapshotResult{}, fmt.Errorf("parse error: %w", err)
	}

	snap := Snapshot{
		Timestamp: time.Now().UTC(),
		Source:    source,
		Entries:   entries,
	}

	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return SnapshotResult{}, fmt.Errorf("marshal error: %w", err)
	}

	if err := os.WriteFile(dest, data, 0600); err != nil {
		return SnapshotResult{}, fmt.Errorf("write error: %w", err)
	}

	return SnapshotResult{Path: dest, KeyCount: len(entries), Timestamp: snap.Timestamp}, nil
}

// LoadSnapshot reads a snapshot JSON file and returns the Snapshot.
func LoadSnapshot(path string) (Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Snapshot{}, fmt.Errorf("read error: %w", err)
	}
	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return Snapshot{}, fmt.Errorf("unmarshal error: %w", err)
	}
	return snap, nil
}
