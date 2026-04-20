package env

import "sort"

// EnvMap represents a map of environment variable key-value pairs.
type EnvMapResult struct {
	Entries []Entry
	Total   int
	Keys    []string
}

// BuildEnvMap constructs an ordered EnvMapResult from a slice of entries.
func BuildEnvMap(entries []Entry) EnvMapResult {
	keys := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.Key != "" && !e.IsComment {
			keys = append(keys, e.Key)
		}
	}
	sort.Strings(keys)
	return EnvMapResult{
		Entries: entries,
		Total:   len(keys),
		Keys:    keys,
	}
}

// ToMap converts a slice of entries into a plain string map.
func ToMap(entries []Entry) map[string]string {
	m := make(map[string]string, len(entries))
	for _, e := range entries {
		if e.Key != "" && !e.IsComment {
			m[e.Key] = e.Value
		}
	}
	return m
}

// FromMap converts a plain string map into a slice of entries (sorted by key).
func FromMap(m map[string]string) []Entry {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	entries := make([]Entry, 0, len(keys))
	for _, k := range keys {
		entries = append(entries, Entry{Key: k, Value: m[k]})
	}
	return entries
}
