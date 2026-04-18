package env

import "strings"

// FilterResult holds the result of a filter operation.
type FilterResult struct {
	Matched map[string]string
	Skipped int
}

// FilterOptions controls how filtering is applied.
type FilterOptions struct {
	Prefix    string
	Suffix    string
	Keys      []string
	SecretsOnly bool
}

// Filter returns entries from the map that match the given options.
func Filter(entries map[string]string, opts FilterOptions) FilterResult {
	matched := make(map[string]string)
	skipped := 0

	keySet := make(map[string]bool, len(opts.Keys))
	for _, k := range opts.Keys {
		keySet[strings.ToUpper(k)] = true
	}

	for k, v := range entries {
		if opts.Prefix != "" && !strings.HasPrefix(k, strings.ToUpper(opts.Prefix)) {
			skipped++
			continue
		}
		if opts.Suffix != "" && !strings.HasSuffix(k, strings.ToUpper(opts.Suffix)) {
			skipped++
			continue
		}
		if len(keySet) > 0 && !keySet[k] {
			skipped++
			continue
		}
		if opts.SecretsOnly && !isSecret(k) {
			skipped++
			continue
		}
		matched[k] = v
	}

	return FilterResult{Matched: matched, Skipped: skipped}
}

// FilterFile parses a .env file and applies the given filter options.
func FilterFile(path string, opts FilterOptions) (FilterResult, error) {
	entries, err := ParseFile(path)
	if err != nil {
		return FilterResult{}, err
	}
	m := make(map[string]string, len(entries))
	for _, e := range entries {
		m[e.Key] = e.Value
	}
	return Filter(m, opts), nil
}
