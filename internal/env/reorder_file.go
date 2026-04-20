package env

import "fmt"

// ReorderFile reads src, reorders entries, and writes the result to dst.
func ReorderFile(src, dst string, opts ReorderOptions) (ReorderResult, error) {
	entries, err := ParseFile(src)
	if err != nil {
		return ReorderResult{}, fmt.Errorf("reorder: read %s: %w", src, err)
	}

	result := Reorder(entries, opts)

	if opts.DryRun {
		return result, nil
	}

	if err := writeEnvFile(dst, result.Entries); err != nil {
		return ReorderResult{}, fmt.Errorf("reorder: write %s: %w", dst, err)
	}

	return result, nil
}
