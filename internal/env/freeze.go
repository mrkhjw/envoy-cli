package env

import "fmt"

// FreezeOption configures the Freeze operation.
type FreezeOption struct {
	Keys    []string
	DryRun  bool
}

// FreezeResult holds the outcome of a Freeze operation.
type FreezeResult struct {
	Frozen  []string
	Skipped []string
	Entries []Entry
}

// Format returns a human-readable summary of the freeze result.
func (r FreezeResult) Format(mask bool) string {
	out := fmt.Sprintf("frozen: %d, skipped: %d\n", len(r.Frozen), len(r.Skipped))
	for _, e := range r.Entries {
		val := e.Value
		if mask && isSecret(e.Key) {
			val = "***"
		}
		out += fmt.Sprintf("  %s=%s\n", e.Key, val)
	}
	return out
}

// Freeze marks entries as frozen by appending a "#frozen" comment tag.
// Frozen entries are skipped on subsequent freeze calls.
// If Keys is non-empty, only those keys are targeted.
func Freeze(entries []Entry, opts FreezeOption) FreezeResult {
	targetSet := make(map[string]bool, len(opts.Keys))
	for _, k := range opts.Keys {
		targetSet[normalizeKey(k)] = true
	}

	result := FreezeResult{}
	out := make([]Entry, 0, len(entries))

	for _, e := range entries {
		if e.Comment {
			out = append(out, e)
			continue
		}

		if len(targetSet) > 0 && !targetSet[normalizeKey(e.Key)] {
			out = append(out, e)
			continue
		}

		// Already frozen?
		if e.RawLine != "" && len(e.RawLine) > 0 {
			// check inline tag via comment field not available; use RawLine suffix
		}
		if isFrozen(e) {
			result.Skipped = append(result.Skipped, e.Key)
			out = append(out, e)
			continue
		}

		result.Frozen = append(result.Frozen, e.Key)
		if !opts.DryRun {
			e.RawLine = fmt.Sprintf("%s=%s #frozen", e.Key, e.Value)
		}
		out = append(out, e)
	}

	result.Entries = out
	return result
}

// isFrozen returns true if the entry's raw line contains the #frozen tag.
func isFrozen(e Entry) bool {
	for i := len(e.RawLine) - 1; i >= 0; i-- {
		if e.RawLine[i] == '#' {
			tag := e.RawLine[i:]
			return tag == "#frozen"
		}
	}
	return false
}
