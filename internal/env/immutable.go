package env

import "fmt"

// ImmutableResult holds the result of marking entries as immutable (read-only).
type ImmutableResult struct {
	Entries  []Entry
	Marked   []string
	Skipped  []string
	DryRun   bool
}

// ImmutableOptions controls the behaviour of the Immutable operation.
type ImmutableOptions struct {
	// Keys restricts immutability tagging to specific keys.
	// If empty, all secret keys are targeted.
	Keys []string

	// DryRun previews changes without modifying entries.
	DryRun bool

	// Overwrite re-marks entries that are already tagged as immutable.
	Overwrite bool
}

const immutableTag = "#immutable"

// isImmutable reports whether an entry already carries the immutable tag.
func isImmutable(e Entry) bool {
	return e.Comment == immutableTag
}

// Immutable marks the selected entries as immutable by appending a special
// inline comment tag. Subsequent tooling (e.g. patch, sync) should honour
// this tag and refuse to modify tagged entries.
//
// When Keys is empty, all entries whose key is considered a secret are
// targeted. DryRun returns the result without mutating the slice.
func Immutable(entries []Entry, opts ImmutableOptions) ImmutableResult {
	targetSet := make(map[string]bool, len(opts.Keys))
	for _, k := range opts.Keys {
		targetSet[normalizeKey(k)] = true
	}

	useAll := len(opts.Keys) == 0

	result := ImmutableResult{
		Entries: make([]Entry, len(entries)),
		DryRun:  opts.DryRun,
	}
	copy(result.Entries, entries)

	for i, e := range result.Entries {
		if e.Key == "" {
			continue
		}

		wanted := useAll && isSecret(e.Key) || targetSet[normalizeKey(e.Key)]
		if !wanted {
			continue
		}

		if isImmutable(e) && !opts.Overwrite {
			result.Skipped = append(result.Skipped, e.Key)
			continue
		}

		result.Marked = append(result.Marked, e.Key)
		if !opts.DryRun {
			result.Entries[i].Comment = immutableTag
		}
	}

	return result
}

// Format returns a human-readable summary of the ImmutableResult.
func (r ImmutableResult) Format() string {
	prefix := ""
	if r.DryRun {
		prefix = "[dry-run] "
	}

	if len(r.Marked) == 0 && len(r.Skipped) == 0 {
		return prefix + "no entries targeted"
	}

	out := ""
	for _, k := range r.Marked {
		out += fmt.Sprintf("%smarked immutable: %s\n", prefix, k)
	}
	for _, k := range r.Skipped {
		out += fmt.Sprintf("%sskipped (already immutable): %s\n", prefix, k)
	}
	return out
}
