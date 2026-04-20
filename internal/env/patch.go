package env

import "fmt"

// PatchOp represents a single patch operation.
type PatchOp struct {
	Op    string // set, delete, rename
	Key   string
	Value string
	NewKey string
}

// PatchResult holds the outcome of a patch operation.
type PatchResult struct {
	Applied []string
	Skipped []string
	DryRun  bool
}

func (r PatchResult) Format() string {
	out := ""
	if r.DryRun {
		out += "[dry-run] "
	}
	out += fmt.Sprintf("applied: %d, skipped: %d\n", len(r.Applied), len(r.Skipped))
	for _, a := range r.Applied {
		out += fmt.Sprintf("  ~ %s\n", a)
	}
	for _, s := range r.Skipped {
		out += fmt.Sprintf("  - %s (skipped)\n", s)
	}
	return out
}

// Patch applies a list of patch operations to a slice of entries.
func Patch(entries []Entry, ops []PatchOp, dryRun bool) ([]Entry, PatchResult) {
	result := PatchResult{DryRun: dryRun}
	index := map[string]int{}
	for i, e := range entries {
		index[e.Key] = i
	}

	for _, op := range ops {
		switch op.Op {
		case "set":
			if i, ok := index[op.Key]; ok {
				if !dryRun {
					entries[i].Value = op.Value
				}
				result.Applied = append(result.Applied, fmt.Sprintf("set %s", op.Key))
			} else {
				if !dryRun {
					entries = append(entries, Entry{Key: op.Key, Value: op.Value})
				}
				result.Applied = append(result.Applied, fmt.Sprintf("add %s", op.Key))
			}
		case "delete":
			if i, ok := index[op.Key]; ok {
				if !dryRun {
					entries = append(entries[:i], entries[i+1:]...)
				}
				result.Applied = append(result.Applied, fmt.Sprintf("delete %s", op.Key))
			} else {
				result.Skipped = append(result.Skipped, op.Key)
			}
		case "rename":
			if i, ok := index[op.Key]; ok {
				if !dryRun {
					entries[i].Key = op.NewKey
				}
				result.Applied = append(result.Applied, fmt.Sprintf("rename %s -> %s", op.Key, op.NewKey))
			} else {
				result.Skipped = append(result.Skipped, op.Key)
			}
		}
	}
	return entries, result
}
