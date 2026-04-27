package env

import "fmt"

// VersionDiffFile compares two .env files identified by version labels.
// It parses both files and returns a VersionDiffResult.
func VersionDiffFile(v1Label, file1, v2Label, file2 string) (VersionDiffResult, error) {
	entries1, err := ParseFile(file1)
	if err != nil {
		return VersionDiffResult{}, fmt.Errorf("reading %s: %w", file1, err)
	}
	entries2, err := ParseFile(file2)
	if err != nil {
		return VersionDiffResult{}, fmt.Errorf("reading %s: %w", file2, err)
	}
	return VersionDiff(v1Label, entries1, v2Label, entries2), nil
}
