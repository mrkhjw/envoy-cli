package env

import "fmt"

// CompactOptions controls how compaction is performed.
type CompactOptions struct {
	RemoveComments bool
	RemoveEmpty    bool
	DryRun         bool
}

// CompactResult holds the result of a Compact operation.
type CompactResult struct {
	Original  []Entry
	Compacted []Entry
	Removed   int
	DryRun    bool
}

// Format returns a human-readable summary of the compact result.
func (r CompactResult) Format() string {
	if r.Removed == 0 {
		return "nothing to compact"
	}
	status := ""
	if r.DryRun {
		status = " (dry run)"
	}
	return fmt.Sprintf("compacted %d entr%s%s",
		r.Removed,
		map[bool]string{true: "y", false: "ies"}[r.Removed == 1],
		status,
	)
}

// Compact removes comments and/or empty-value entries from the given entries.
func Compact(entries []Entry, opts CompactOptions) CompactResult {
	result := CompactResult{
		Original: entries,
		DryRun:   opts.DryRun,
	}

	var compacted []Entry
	for _, e := range entries {
		if opts.RemoveComments && e.Comment {
			result.Removed++
			continue
		}
		if opts.RemoveEmpty && !e.Comment && e.Value == "" {
			result.Removed++
			continue
		}
		compacted = append(compacted, e)
	}

	if opts.DryRun {
		result.Compacted = entries
	} else {
		result.Compacted = compacted
	}

	return result
}
