package env

import (
	"fmt"
	"sort"
	"strings"
)

// PinResult holds the result of a pin operation.
type PinResult struct {
	Pinned  map[string]string
	Skipped []string
	DryRun  bool
}

// Pin locks specific keys to their current values by recording them.
// If keys is empty, all entries are pinned.
func Pin(entries []Entry, keys []string, dryRun bool) PinResult {
	keySet := make(map[string]bool, len(keys))
	for _, k := range keys {
		keySet[strings.ToUpper(k)] = true
	}

	pinned := make(map[string]string)
	var skipped []string

	for _, e := range entries {
		if e.Key == "" {
			continue
		}
		if len(keySet) == 0 || keySet[e.Key] {
			pinned[e.Key] = e.Value
		} else {
			skipped = append(skipped, e.Key)
		}
	}

	return PinResult{
		Pinned:  pinned,
		Skipped: skipped,
		DryRun:  dryRun,
	}
}

// Format returns a human-readable summary of the pin result.
func (r PinResult) Format() string {
	var sb strings.Builder
	if r.DryRun {
		sb.WriteString("[dry-run] ")
	}
	sb.WriteString(fmt.Sprintf("Pinned %d key(s)\n", len(r.Pinned)))

	keys := make([]string, 0, len(r.Pinned))
	for k := range r.Pinned {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := r.Pinned[k]
		if isSecret(k) {
			v = "***"
		}
		sb.WriteString(fmt.Sprintf("  pinned: %s=%s\n", k, v))
	}
	if len(r.Skipped) > 0 {
		sb.WriteString(fmt.Sprintf("  skipped: %s\n", strings.Join(r.Skipped, ", ")))
	}
	return sb.String()
}
