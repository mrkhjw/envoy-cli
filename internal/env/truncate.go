package env

import "fmt"

// TruncateOptions controls how values are truncated.
type TruncateOptions struct {
	MaxLen  int
	Suffix  string
	Keys    []string
	DryRun  bool
}

// TruncateResult holds the outcome of a Truncate operation.
type TruncateResult struct {
	Entries   []Entry
	Truncated []string
	DryRun    bool
}

// Format returns a human-readable summary of the truncation result.
func (r TruncateResult) Format() string {
	if len(r.Truncated) == 0 {
		return "no values truncated"
	}
	tag := ""
	if r.DryRun {
		tag = " (dry run)"
	}
	return fmt.Sprintf("truncated %d key(s)%s: %v", len(r.Truncated), tag, r.Truncated)
}

// Truncate shortens values that exceed MaxLen, appending Suffix.
// If Keys is non-empty, only those keys are considered.
func Truncate(entries []Entry, opts TruncateOptions) TruncateResult {
	if opts.MaxLen <= 0 {
		opts.MaxLen = 64
	}
	if opts.Suffix == "" {
		opts.Suffix = "..."
	}

	keySet := make(map[string]bool, len(opts.Keys))
	for _, k := range opts.Keys {
		keySet[normalizeKey(k)] = true
	}

	result := make([]Entry, 0, len(entries))
	var truncated []string

	for _, e := range entries {
		if e.Comment || e.Key == "" {
			result = append(result, e)
			continue
		}
		if len(keySet) > 0 && !keySet[normalizeKey(e.Key)] {
			result = append(result, e)
			continue
		}
		if len(e.Value) > opts.MaxLen {
			if !opts.DryRun {
				e.Value = e.Value[:opts.MaxLen] + opts.Suffix
			}
			truncated = append(truncated, e.Key)
		}
		result = append(result, e)
	}

	return TruncateResult{
		Entries:   result,
		Truncated: truncated,
		DryRun:    opts.DryRun,
	}
}
