package env

import (
	"fmt"
	"strings"
)

// UnsetResult holds the outcome of an Unset operation.
type UnsetResult struct {
	Removed []string
	Skipped []string
	Entries []Entry
	DryRun  bool
}

// Format returns a human-readable summary of the unset operation.
func (r UnsetResult) Format() string {
	var sb strings.Builder
	if r.DryRun {
		sb.WriteString("[dry-run] ")
	}
	sb.WriteString(fmt.Sprintf("removed: %d, skipped: %d\n", len(r.Removed), len(r.Skipped)))
	for _, k := range r.Removed {
		sb.WriteString(fmt.Sprintf("  - %s\n", k))
	}
	for _, k := range r.Skipped {
		sb.WriteString(fmt.Sprintf("  ~ %s (not found)\n", k))
	}
	return strings.TrimRight(sb.String(), "\n")
}

// Unset removes the specified keys from the given entries.
// If dryRun is true, no changes are applied but the result reflects what would happen.
func Unset(entries []Entry, keys []string, dryRun bool) UnsetResult {
	keySet := make(map[string]bool, len(keys))
	for _, k := range keys {
		keySet[strings.ToUpper(k)] = true
	}

	var result UnsetResult
	result.DryRun = dryRun

	removedSet := make(map[string]bool)
	var kept []Entry

	for _, e := range entries {
		norm := strings.ToUpper(e.Key)
		if e.Key != "" && keySet[norm] {
			removedSet[strings.ToUpper(e.Key)] = true
			result.Removed = append(result.Removed, e.Key)
			if dryRun {
				kept = append(kept, e)
			}
		} else {
			kept = append(kept, e)
		}
	}

	for _, k := range keys {
		if !removedSet[strings.ToUpper(k)] {
			result.Skipped = append(result.Skipped, k)
		}
	}

	if dryRun {
		result.Entries = entries
	} else {
		result.Entries = kept
	}

	return result
}
