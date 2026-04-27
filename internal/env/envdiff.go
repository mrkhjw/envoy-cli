package env

import "fmt"

// VersionedEntry holds an env entry with a version label.
type VersionedEntry struct {
	Entry
	Version string
}

// VersionDiffResult holds the result of comparing two versioned env snapshots.
type VersionDiffResult struct {
	Added    []VersionedEntry
	Removed  []VersionedEntry
	Changed  []VersionedEntry
	Version1 string
	Version2 string
}

// VersionDiff compares two labeled slices of entries and returns a VersionDiffResult.
func VersionDiff(v1Label string, entries1 []Entry, v2Label string, entries2 []Entry) VersionDiffResult {
	map1 := make(map[string]Entry)
	for _, e := range entries1 {
		if !e.IsComment {
			map1[e.Key] = e
		}
	}
	map2 := make(map[string]Entry)
	for _, e := range entries2 {
		if !e.IsComment {
			map2[e.Key] = e
		}
	}

	result := VersionDiffResult{Version1: v1Label, Version2: v2Label}

	for k, e2 := range map2 {
		if e1, ok := map1[k]; !ok {
			result.Added = append(result.Added, VersionedEntry{Entry: e2, Version: v2Label})
		} else if e1.Value != e2.Value {
			result.Changed = append(result.Changed, VersionedEntry{Entry: e2, Version: v2Label})
		}
	}

	for k, e1 := range map1 {
		if _, ok := map2[k]; !ok {
			result.Removed = append(result.Removed, VersionedEntry{Entry: e1, Version: v1Label})
		}
	}

	return result
}

// Format returns a human-readable summary of the VersionDiffResult.
func (r VersionDiffResult) Format(maskSecrets bool) string {
	out := fmt.Sprintf("diff %s..%s\n", r.Version1, r.Version2)
	for _, e := range r.Added {
		v := e.Value
		if maskSecrets && isSecret(e.Key) {
			v = "***"
		}
		out += fmt.Sprintf("+ %s=%s\n", e.Key, v)
	}
	for _, e := range r.Removed {
		v := e.Value
		if maskSecrets && isSecret(e.Key) {
			v = "***"
		}
		out += fmt.Sprintf("- %s=%s\n", e.Key, v)
	}
	for _, e := range r.Changed {
		v := e.Value
		if maskSecrets && isSecret(e.Key) {
			v = "***"
		}
		out += fmt.Sprintf("~ %s=%s\n", e.Key, v)
	}
	return out
}
