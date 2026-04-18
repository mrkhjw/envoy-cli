package env

import (
	"fmt"
	"sort"
	"strings"
)

// GroupResult holds keys grouped by a common prefix.
type GroupResult struct {
	Groups map[string][]string
}

// Group organizes env keys by their prefix (split by separator).
func Group(entries []Entry, separator string) GroupResult {
	if separator == "" {
		separator = "_"
	}
	groups := make(map[string][]string)
	for _, e := range entries {
		parts := strings.SplitN(e.Key, separator, 2)
		prefix := parts[0]
		if len(parts) == 1 {
			prefix = "(ungrouped)"
		}
		groups[prefix] = append(groups[prefix], e.Key)
	}
	return GroupResult{Groups: groups}
}

// Format returns a human-readable summary of grouped keys.
func (g GroupResult) Format(entries []Entry, maskSecrets bool) string {
	valueMap := make(map[string]string, len(entries))
	for _, e := range entries {
		valueMap[e.Key] = e.Value
	}

	prefixes := make([]string, 0, len(g.Groups))
	for p := range g.Groups {
		prefixes = append(prefixes, p)
	}
	sort.Strings(prefixes)

	var sb strings.Builder
	for _, prefix := range prefixes {
		fmt.Fprintf(&sb, "[%s]\n", prefix)
		keys := g.Groups[prefix]
		sort.Strings(keys)
		for _, k := range keys {
			v := valueMap[k]
			if maskSecrets && isSecret(k) {
				v = "***"
			}
			fmt.Fprintf(&sb, "  %s=%s\n", k, v)
		}
	}
	return sb.String()
}
