package env

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// CheckpointEntry represents a saved checkpoint of env entries.
type CheckpointEntry struct {
	Timestamp time.Time  `json:"timestamp"`
	Label     string     `json:"label"`
	Entries   []Entry    `json:"entries"`
}

// CheckpointResult holds the outcome of a checkpoint operation.
type CheckpointResult struct {
	Label   string
	Count   int
	DryRun  bool
	OutPath string
}

func (r CheckpointResult) Format() string {
	if r.DryRun {
		return fmt.Sprintf("[dry-run] would save checkpoint %q with %d entries to %s", r.Label, r.Count, r.OutPath)
	}
	return fmt.Sprintf("checkpoint %q saved: %d entries → %s", r.Label, r.Count, r.OutPath)
}

// Checkpoint saves the given entries as a labeled checkpoint to outPath.
func Checkpoint(entries []Entry, label, outPath string, dryRun bool) (CheckpointResult, error) {
	result := CheckpointResult{
		Label:   label,
		Count:   len(entries),
		DryRun:  dryRun,
		OutPath: outPath,
	}
	if dryRun {
		return result, nil
	}
	cp := CheckpointEntry{
		Timestamp: time.Now().UTC(),
		Label:     label,
		Entries:   entries,
	}
	data, err := json.MarshalIndent(cp, "", "  ")
	if err != nil {
		return result, fmt.Errorf("checkpoint: marshal error: %w", err)
	}
	if err := os.WriteFile(outPath, data, 0600); err != nil {
		return result, fmt.Errorf("checkpoint: write error: %w", err)
	}
	return result, nil
}

// LoadCheckpoint reads a checkpoint file and returns its entries.
func LoadCheckpoint(path string) (CheckpointEntry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return CheckpointEntry{}, fmt.Errorf("checkpoint: read error: %w", err)
	}
	var cp CheckpointEntry
	if err := json.Unmarshal(data, &cp); err != nil {
		return CheckpointEntry{}, fmt.Errorf("checkpoint: parse error: %w", err)
	}
	return cp, nil
}
