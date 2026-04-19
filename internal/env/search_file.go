package env

import "fmt"

// SearchFile parses the given file and searches its entries.
func SearchFile(path string, opts SearchOptions, mask bool) (SearchResult, error) {
	entries, err := ParseFile(path)
	if err != nil {
		return SearchResult{}, fmt.Errorf("search: %w", err)
	}

	envMap := make(map[string]string, len(entries))
	for _, e := range entries {
		envMap[e.Key] = e.Value
	}

	return Search(envMap, opts), nil
}
