package env

import "fmt"

// DiffFile compares two .env files on disk and returns a DiffResult.
func DiffFile(pathA, pathB string) (DiffResult, error) {
	entriesA, err := ParseFile(pathA)
	if err != nil {
		return DiffResult{}, fmt.Errorf("reading %s: %w", pathA, err)
	}

	entriesB, err := ParseFile(pathB)
	if err != nil {
		return DiffResult{}, fmt.Errorf("reading %s: %w", pathB, err)
	}

	mapA := ToMap(entriesA)
	mapB := ToMap(entriesB)

	return Diff(mapA, mapB), nil
}
