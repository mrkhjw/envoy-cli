package env

import "fmt"

// MergeResult holds the outcome of a merge operation.
type MergeResult struct {
	Added    []string
	Skipped  []string
	Updated  []string
}

func (r MergeResult) Format() string {
	out := ""
	for _, k := range r.Added {
		out += fmt.Sprintf("+ %s\n", k)
	}
	for _, k := range r.Updated {
		out += fmt.Sprintf("~ %s\n", k)
	}
	for _, k := range r.Skipped {
		out += fmt.Sprintf("= %s (skipped)\n", k)
	}
	return out
}

// Merge combines src into dst. If overwrite is true, existing keys are updated.
func Merge(dst, src map[string]string, overwrite bool) (map[string]string, MergeResult) {
	result := MergeResult{}
	out := make(map[string]string)
	for k, v := range dst {
		out[k] = v
	}
	for k, v := range src {
		if existing, ok := out[k]; ok {
			if overwrite && existing != v {
				out[k] = v
				result.Updated = append(result.Updated, k)
			} else {
				result.Skipped = append(result.Skipped, k)
			}
		} else {
			out[k] = v
			result.Added = append(result.Added, k)
		}
	}
	return out, result
}
