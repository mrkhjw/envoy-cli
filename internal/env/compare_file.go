package env

import "fmt"

// CompareFile parses two .env files and returns a CompareResult.
func CompareFile(path1, path2 string) (CompareResult, error) {
	entries1, err := ParseFile(path1)
	if err != nil {
		return CompareResult{}, fmt.Errorf("reading %s: %w", path1, err)
	}

	entries2, err := ParseFile(path2)
	if err != nil {
		return CompareResult{}, fmt.Errorf("reading %s: %w", path2, err)
	}

	map1 := make(map[string]string, len(entries1))
	for _, e := range entries1 {
		map1[e.Key] = e.Value
	}

	map2 := make(map[string]string, len(entries2))
	for _, e := range entries2 {
		map2[e.Key] = e.Value
	}

	return Compare(map1, map2), nil
}
