package env

import "fmt"

// DedupeResult holds the result of a deduplication operation.
type DedupeResult struct {
	Entries  []Entry
	Removed  []string
	Total    int
	Dupes    int
}

// Entry represents a key-value pair with its original line order.
type Entry struct {
	Key   string
	Value string
}

// Dedupe removes duplicate keys from a slice of parsed entries,
// keeping the last occurrence of each key.
func Dedupe(entries []Entry) DedupeResult {
	seen := make(map[string]int)
	for i, e := range entries {
		seen[e.Key] = i
	}

	var deduped []Entry
	removed := []string{}
	included := make(map[string]bool)

	for i, e := range entries {
		if seen[e.Key] == i {
			deduped = append(deduped, e)
			included[e.Key] = true
		} else {
			removed = append(removed, e.Key)
		}
	}

	return DedupeResult{
		Entries: deduped,
		Removed: removed,
		Total:   len(entries),
		Dupes:   len(removed),
	}
}

// Format returns a human-readable summary of the deduplication result.
func (r DedupeResult) Format() string {
	if r.Dupes == 0 {
		return fmt.Sprintf("no duplicates found (%d entries)", r.Total)
	}
	out := fmt.Sprintf("removed %d duplicate(s) from %d entries:\n", r.Dupes, r.Total)
	for _, k := range r.Removed {
		out += fmt.Sprintf("  - %s\n", k)
	}
	return out
}
