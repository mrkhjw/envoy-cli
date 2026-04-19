package env

import "fmt"

// CopyResult holds the result of a copy operation.
type CopyResult struct {
	Copied  []string
	Skipped []string
}

// Format returns a human-readable summary of the copy result.
func (r CopyResult) Format(mask bool) string {
	var out string
	if len(r.Copied) == 0 && len(r.Skipped) == 0 {
		return "nothing to copy"
	}
	for _, k := range r.Copied {
		out += fmt.Sprintf("  copied: %s\n", k)
	}
	for _, k := range r.Skipped {
		out += fmt.Sprintf("  skipped: %s\n", k)
	}
	out += fmt.Sprintf("total: %d copied, %d skipped", len(r.Copied), len(r.Skipped))
	return out
}

// CopyOptions controls the behaviour of Copy.
type CopyOptions struct {
	Keys      []string // if empty, copy all keys
	Overwrite bool
	DryRun    bool
}

// Copy copies keys from src into dst according to opts.
func Copy(src, dst []Entry, opts CopyOptions) ([]Entry, CopyResult) {
	result := CopyResult{}

	srcMap := make(map[string]string)
	for _, e := range src {
		srcMap[e.Key] = e.Value
	}

	dstMap := make(map[string]string)
	for _, e := range dst {
		dstMap[e.Key] = e.Value
	}

	targetKeys := opts.Keys
	if len(targetKeys) == 0 {
		for k := range srcMap {
			targetKeys = append(targetKeys, k)
		}
	}

	for _, k := range targetKeys {
		v, ok := srcMap[k]
		if !ok {
			result.Skipped = append(result.Skipped, k)
			continue
		}
		if _, exists := dstMap[k]; exists && !opts.Overwrite {
			result.Skipped = append(result.Skipped, k)
			continue
		}
		if !opts.DryRun {
			dstMap[k] = v
		}
		result.Copied = append(result.Copied, k)
	}

	if opts.DryRun {
		return dst, result
	}

	var out []Entry
	seen := make(map[string]bool)
	for _, e := range dst {
		e.Value = dstMap[e.Key]
		out = append(out, e)
		seen[e.Key] = true
	}
	for _, k := range result.Copied {
		if !seen[k] {
			out = append(out, Entry{Key: k, Value: dstMap[k]})
		}
	}
	return out, result
}
