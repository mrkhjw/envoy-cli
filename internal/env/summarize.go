package env

import (
	"fmt"
	"sort"
	"strings"
)

// SummarizeResult holds the output of a Summarize operation.
type SummarizeResult struct {
	Total      int
	Secrets    int
	Empty      int
	Comments   int
	Groups     map[string]int // key prefix -> count
	Entries    []Entry
}

// Format returns a human-readable summary string.
func (r SummarizeResult) Format(maskSecrets bool) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Total keys   : %d\n", r.Total))
	sb.WriteString(fmt.Sprintf("Secrets      : %d\n", r.Secrets))
	sb.WriteString(fmt.Sprintf("Empty values : %d\n", r.Empty))
	sb.WriteString(fmt.Sprintf("Comments     : %d\n", r.Comments))
	if len(r.Groups) > 0 {
		sb.WriteString("Groups:\n")
		// Sort group prefixes for deterministic output
		prefixes := make([]string, 0, len(r.Groups))
		for prefix := range r.Groups {
			prefixes = append(prefixes, prefix)
		}
		sort.Strings(prefixes)
		for _, prefix := range prefixes {
			sb.WriteString(fmt.Sprintf("  %-20s %d\n", prefix+"_*", r.Groups[prefix]))
		}
	}
	if len(r.Entries) > 0 {
		sb.WriteString("\nEntries:\n")
		for _, e := range r.Entries {
			val := e.Value
			if maskSecrets && isSecret(e.Key) {
				val = "****"
			}
			sb.WriteString(fmt.Sprintf("  %s=%s\n", e.Key, val))
		}
	}
	return strings.TrimRight(sb.String(), "\n")
}

// Summarize analyses a slice of Entry values and returns a SummarizeResult.
func Summarize(entries []Entry, separator string) SummarizeResult {
	if separator == "" {
		separator = "_"
	}

	result := SummarizeResult{
		Groups:  make(map[string]int),
		Entries: entries,
	}

	for _, e := range entries {
		if strings.HasPrefix(e.Key, "#") {
			result.Comments++
			continue
		}
		result.Total++
		if isSecret(e.Key) {
			result.Secrets++
		}
		if e.Value == "" {
			result.Empty++
		}
		if idx := strings.Index(e.Key, separator); idx > 0 {
			prefix := e.Key[:idx]
			result.Groups[prefix]++
		}
	}
	return result
}
