package env

import (
	"fmt"
	"os"
)

// CloneResult holds the result of a clone operation.
type CloneResult struct {
	Source      string
	Destination string
	KeysCopied  int
	Skipped     int
}

func (r CloneResult) Summary() string {
	return fmt.Sprintf("Cloned %d keys from %s to %s (%d skipped)",
		r.KeysCopied, r.Source, r.Destination, r.Skipped)
}

// Clone copies key-value pairs from src to dst file.
// If overwrite is false, existing keys in dst are preserved.
// If mask is true, secret values are masked in the result summary.
func Clone(srcPath, dstPath string, overwrite bool) (CloneResult, error) {
	src, err := ParseFile(srcPath)
	if err != nil {
		return CloneResult{}, fmt.Errorf("reading source: %w", err)
	}

	var dst map[string]string
	if _, err := os.Stat(dstPath); err == nil {
		dst, err = ParseFile(dstPath)
		if err != nil {
			return CloneResult{}, fmt.Errorf("reading destination: %w", err)
		}
	} else {
		dst = make(map[string]string)
	}

	result := CloneResult{Source: srcPath, Destination: dstPath}

	for k, v := range src {
		if _, exists := dst[k]; exists && !overwrite {
			result.Skipped++
			continue
		}
		dst[k] = v
		result.KeysCopied++
	}

	if err := writeEnvFile(dstPath, dst); err != nil {
		return CloneResult{}, fmt.Errorf("writing destination: %w", err)
	}

	return result, nil
}
