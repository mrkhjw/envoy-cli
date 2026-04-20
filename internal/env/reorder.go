package env

import "sort"

// ReorderOptions controls how entries are reordered.
type ReorderOptions struct {
	Keys    []string // desired key order; unspecified keys follow at the end
	DryRun  bool
}

// ReorderResult holds the result of a Reorder operation.
type ReorderResult struct {
	Entries  []Entry
	Moved    int
	DryRun   bool
}

// Reorder rearranges entries so that the specified keys appear first, in order.
func Reorder(entries []Entry, opts ReorderOptions) ReorderResult {
	indexOf := make(map[string]int, len(opts.Keys))
	for i, k := range opts.Keys {
		indexOf[k] = i
	}

	pinned := make([]Entry, len(opts.Keys))
	pinnedSet := make([]bool, len(opts.Keys))
	var rest []Entry

	for _, e := range entries {
		if idx, ok := indexOf[e.Key]; ok {
			pinned[idx] = e
			pinnedSet[idx] = true
		} else {
			rest = append(rest, e)
		}
	}

	sort.SliceStable(rest, func(i, j int) bool {
		return rest[i].Key < rest[j].Key
	})

	result := make([]Entry, 0, len(entries))
	moved := 0
	for i, e := range pinned {
		if pinnedSet[i] {
			result = append(result, e)
			moved++
		}
	}
	result = append(result, rest...)

	return ReorderResult{
		Entries: result,
		Moved:   moved,
		DryRun:  opts.DryRun,
	}
}

// Format returns a human-readable summary of the reorder result.
func (r ReorderResult) Format() string {
	if r.DryRun {
		return fmt.Sprintf("[dry-run] would reorder %d key(s)", r.Moved)
	}
	return fmt.Sprintf("reordered %d key(s)", r.Moved)
}
