package env

import (
	"fmt"
	"os"
	"strings"
)

// SyncResult holds the outcome of a sync operation.
type SyncResult struct {
	Applied []string
	Skipped []string
}

// SyncOptions controls sync behaviour.
type SyncOptions struct {
	Overwrite bool
	DryRun    bool
}

// Sync merges src into dst, writing the result to dstPath.
// Keys present in src but missing in dst are added.
// Keys present in both are updated only when opts.Overwrite is true.
func Sync(dst map[string]string, src map[string]string, dstPath string, opts SyncOptions) (SyncResult, error) {
	result := SyncResult{}

	merged := make(map[string]string, len(dst))
	for k, v := range dst {
		merged[k] = v
	}

	for k, v := range src {
		if existing, ok := merged[k]; ok {
			if opts.Overwrite && existing != v {
				merged[k] = v
				result.Applied = append(result.Applied, k)
			} else {
				result.Skipped = append(result.Skipped, k)
			}
		} else {
			merged[k] = v
			result.Applied = append(result.Applied, k)
		}
	}

	if opts.DryRun {
		return result, nil
	}

	if err := writeEnvFile(dstPath, merged); err != nil {
		return result, fmt.Errorf("sync: write failed: %w", err)
	}

	return result, nil
}

// writeEnvFile serialises a map to a .env file at path.
func writeEnvFile(path string, env map[string]string) error {
	var sb strings.Builder
	for k, v := range env {
		if strings.ContainsAny(v, " \t#") {
			fmt.Fprintf(&sb, "%s=\"%s\"\n", k, v)
		} else {
			fmt.Fprintf(&sb, "%s=%s\n", k, v)
		}
	}
	return os.WriteFile(path, []byte(sb.String()), 0644)
}
