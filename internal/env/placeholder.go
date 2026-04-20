package env

import (
	"fmt"
	"strings"
)

// PlaceholderResult holds the result of a placeholder check
type PlaceholderResult struct {
	Entries  []Entry
	Missing  []string
	Filled   int
}

// PlaceholderOptions configures placeholder filling behavior
type PlaceholderOptions struct {
	Token     string // default: "CHANGEME"
	Overwrite bool
	DryRun    bool
}

// FillPlaceholders replaces placeholder values with provided replacements
func FillPlaceholders(entries []Entry, replacements map[string]string, opts PlaceholderOptions) PlaceholderResult {
	token := opts.Token
	if token == "" {
		token = "CHANGEME"
	}

	result := PlaceholderResult{}
	for _, e := range entries {
		if strings.EqualFold(e.Value, token) {
			if val, ok := replacements[e.Key]; ok {
				if !opts.DryRun {
					e.Value = val
				}
				result.Filled++
			} else {
				result.Missing = append(result.Missing, e.Key)
			}
		}
		result.Entries = append(result.Entries, e)
	}
	return result
}

// Format returns a human-readable summary of the placeholder result
func (r PlaceholderResult) Format() string {
	var sb strings.Builder
	if r.Filled > 0 {
		fmt.Fprintf(&sb, "filled %d placeholder(s)\n", r.Filled)
	}
	for _, k := range r.Missing {
		fmt.Fprintf(&sb, "missing replacement for: %s\n", k)
	}
	if r.Filled == 0 && len(r.Missing) == 0 {
		sb.WriteString("no placeholders found\n")
	}
	return strings.TrimRight(sb.String(), "\n")
}
