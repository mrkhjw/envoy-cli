package env

import (
	"fmt"
	"strings"
)

// SearchOptions configures how the search is performed.
type SearchOptions struct {
	Key        string
	Value      string
	CaseSensitive bool
}

// SearchResult holds matched entries.
type SearchResult struct {
	Matches []Entry
}

// Entry represents a single key-value pair.
type Entry struct {
	Key   string
	Value string
}

// Search scans an env map for entries matching the given options.
func Search(env map[string]string, opts SearchOptions) SearchResult {
	var matches []Entry

	for k, v := range env {
		keyMatch := true
		valMatch := true

		if opts.Key != "" {
			if opts.CaseSensitive {
				keyMatch = strings.Contains(k, opts.Key)
			} else {
				keyMatch = strings.Contains(strings.ToLower(k), strings.ToLower(opts.Key))
			}
		}

		if opts.Value != "" {
			if opts.CaseSensitive {
				valMatch = strings.Contains(v, opts.Value)
			} else {
				valMatch = strings.Contains(strings.ToLower(v), strings.ToLower(opts.Value))
			}
		}

		if keyMatch && valMatch {
			matches = append(matches, Entry{Key: k, Value: v})
		}
	}

	return SearchResult{Matches: matches}
}

// Format returns a human-readable summary of search results.
func (r SearchResult) Format(mask bool) string {
	if len(r.Matches) == 0 {
		return "no matches found"
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d match(es):\n", len(r.Matches))
	for _, e := range r.Matches {
		val := e.Value
		if mask && isSecret(e.Key) {
			val = "***"
		}
		fmt.Fprintf(&sb, "  %s=%s\n", e.Key, val)
	}
	return strings.TrimRight(sb.String(), "\n")
}
