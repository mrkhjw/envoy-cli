package env

import (
	"fmt"
	"sort"
	"strings"
)

// SortOptions controls how entries are sorted.
type SortOptions struct {
	Reverse    bool
	ByValue    bool
	SecretsLast bool
}

// SortResult holds the sorted entries and a summary.
type SortResult struct {
	Entries []Entry
	Total   int
}

// Entry represents a single key-value pair.
type Entry struct {
	Key   string
	Value string
}

// Sort sorts the given env map according to options.
func Sort(env map[string]string, opts SortOptions) SortResult {
	entries := make([]Entry, 0, len(env))
	for k, v := range env {
		entries = append(entries, Entry{Key: k, Value: v})
	}

	sort.Slice(entries, func(i, j int) bool {
		if opts.SecretsLast {
			iSecret := isSecret(entries[i].Key)
			jSecret := isSecret(entries[j].Key)
			if iSecret != jSecret {
				return !iSecret
			}
		}
		var less bool
		if opts.ByValue {
			less = strings.ToLower(entries[i].Value) < strings.ToLower(entries[j].Value)
		} else {
			less = strings.ToLower(entries[i].Key) < strings.ToLower(entries[j].Key)
		}
		if opts.Reverse {
			return !less
		}
		return less
	})

	return SortResult{Entries: entries, Total: len(entries)}
}

// Format returns a human-readable summary.
func (r SortResult) Format(maskSecrets bool) string {
	var sb strings.Builder
	for _, e := range r.Entries {
		val := e.Value
		if maskSecrets && isSecret(e.Key) {
			val = "****"
		}
		sb.WriteString(fmt.Sprintf("%s=%s\n", e.Key, val))
	}
	return strings.TrimRight(sb.String(), "\n")
}
