package env

import (
	"fmt"
	"os"
	"strings"
)

// SortFile reads a .env file, sorts its entries, and writes back to dest.
func SortFile(src, dest string, opts SortOptions) (SortResult, error) {
	entries, err := ParseFile(src)
	if err != nil {
		return SortResult{}, fmt.Errorf("failed to parse %s: %w", src, err)
	}

	result := Sort(entries, opts)

	var sb strings.Builder
	for _, e := range result.Entries {
		sb.WriteString(fmt.Sprintf("%s=%s\n", e.Key, e.Value))
	}

	if err := os.WriteFile(dest, []byte(sb.String()), 0644); err != nil {
		return SortResult{}, fmt.Errorf("failed to write %s: %w", dest, err)
	}

	return result, nil
}
