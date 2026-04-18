package env

import (
	"fmt"
	"strings"
)

// TrimResult holds the result of a trim operation.
type TrimResult struct {
	Trimmed []string
	Skipped []string
}

// TrimOptions controls how trimming is applied.
type TrimOptions struct {
	Keys      []string // if empty, trim all
	TrimLeft  bool
	TrimRight bool
	DryRun    bool
}

// Trim removes leading/trailing whitespace from env values.
func Trim(entries []Entry, opts TrimOptions) ([]Entry, TrimResult) {
	keySet := make(map[string]bool)
	for _, k := range opts.Keys {
		keySet[strings.ToUpper(k)] = true
	}

	var result TrimResult
	out := make([]Entry, 0, len(entries))

	for _, e := range entries {
		target := len(opts.Keys) == 0 || keySet[strings.ToUpper(e.Key)]
		if !target {
			out = append(out, e)
			result.Skipped = append(result.Skipped, e.Key)
			continue
		}

		newVal := e.Value
		if opts.TrimLeft {
			newVal = strings.TrimLeft(newVal, " \t")
		}
		if opts.TrimRight {
			newVal = strings.TrimRight(newVal, " \t")
		}
		if !opts.TrimLeft && !opts.TrimRight {
			newVal = strings.TrimSpace(newVal)
		}

		if newVal != e.Value {
			result.Trimmed = append(result.Trimmed, e.Key)
		}

		if !opts.DryRun {
			e.Value = newVal
		}
		out = append(out, e)
	}

	return out, result
}

// Format returns a human-readable summary of the trim result.
func (r TrimResult) Format() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Trimmed: %d key(s)\n", len(r.Trimmed)))
	for _, k := range r.Trimmed {
		sb.WriteString(fmt.Sprintf("  ~ %s\n", k))
	}
	if len(r.Skipped) > 0 {
		sb.WriteString(fmt.Sprintf("Skipped: %d key(s)\n", len(r.Skipped)))
	}
	return sb.String()
}
