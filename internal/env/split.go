package env

import "fmt"

// SplitOptions configures how entries are split into groups.
type SplitOptions struct {
	Keys    []string
	Invert  bool
	DryRun  bool
}

// SplitResult holds the two resulting groups after a split.
type SplitResult struct {
	Matched   []Entry
	Remainder []Entry
	DryRun    bool
}

// Split divides entries into matched and remainder groups based on keys.
func Split(entries []Entry, opts SplitOptions) SplitResult {
	keySet := make(map[string]bool, len(opts.Keys))
	for _, k := range opts.Keys {
		keySet[normalizeKey(k)] = true
	}

	var matched, remainder []Entry
	for _, e := range entries {
		inSet := keySet[normalizeKey(e.Key)]
		if inSet != opts.Invert {
			matched = append(matched, e)
		} else {
			remainder = append(remainder, e)
		}
	}

	return SplitResult{
		Matched:   matched,
		Remainder: remainder,
		DryRun:    opts.DryRun,
	}
}

func (r SplitResult) Format() string {
	out := fmt.Sprintf("matched: %d, remainder: %d", len(r.Matched), len(r.Remainder))
	if r.DryRun {
		out += " (dry run)"
	}
	return out
}

func normalizeKey(k string) string {
	return k
}
