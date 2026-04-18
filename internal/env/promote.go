package env

import "fmt"

// PromoteResult holds the result of a promote operation.
type PromoteResult struct {
	Added   []string
	Skipped []string
	Updated []string
}

func (r PromoteResult) Format() string {
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

// Promote copies entries from src into dst.
// If overwrite is true, existing keys in dst are updated.
// If keys is non-empty, only those keys are promoted.
func Promote(src, dst map[string]string, keys []string, overwrite bool) PromoteResult {
	result := PromoteResult{}

	targets := keys
	if len(targets) == 0 {
		for k := range src {
			targets = append(targets, k)
		}
	}

	for _, k := range targets {
		val, ok := src[k]
		if !ok {
			continue
		}
		if existing, exists := dst[k]; exists {
			if existing == val {
				result.Skipped = append(result.Skipped, k)
				continue
			}
			if !overwrite {
				result.Skipped = append(result.Skipped, k)
				continue
			}
			dst[k] = val
			result.Updated = append(result.Updated, k)
		} else {
			dst[k] = val
			result.Added = append(result.Added, k)
		}
	}
	return result
}

// PromoteFile promotes entries from srcPath to dstPath, writing the result.
func PromoteFile(srcPath, dstPath string, keys []string, overwrite bool) (PromoteResult, error) {
	src, err := ParseFile(srcPath)
	if err != nil {
		return PromoteResult{}, fmt.Errorf("reading source: %w", err)
	}
	dst, err := ParseFile(dstPath)
	if err != nil {
		dst = map[string]string{}
	}
	result := Promote(src, dst, keys, overwrite)
	if err := writeEnvFile(dstPath, dst); err != nil {
		return result, fmt.Errorf("writing destination: %w", err)
	}
	return result, nil
}
