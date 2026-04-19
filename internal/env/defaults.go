package env

import "fmt"

// DefaultsResult holds the result of applying defaults to env entries.
type DefaultsResult struct {
	Entries  []Entry
	Applied  []string
	Skipped  []string
}

// Defaults applies default values to entries where the value is empty.
// If overwrite is true, it replaces all values with the defaults.
func Defaults(entries []Entry, defaults map[string]string, overwrite bool) DefaultsResult {
	result := DefaultsResult{}
	index := map[string]int{}

	for i, e := range entries {
		index[e.Key] = i
	}

	out := make([]Entry, len(entries))
	copy(out, entries)

	for key, defVal := range defaults {
		if i, exists := index[key]; exists {
			if out[i].Value == "" || overwrite {
				out[i].Value = defVal
				result.Applied = append(result.Applied, key)
			} else {
				result.Skipped = append(result.Skipped, key)
			}
		} else {
			out = append(out, Entry{Key: key, Value: defVal})
			result.Applied = append(result.Applied, key)
		}
	}

	result.Entries = out
	return result
}

// Format returns a human-readable summary of the defaults result.
func (r DefaultsResult) Format() string {
	if len(r.Applied) == 0 && len(r.Skipped) == 0 {
		return "defaults: nothing to apply"
	}
	s := fmt.Sprintf("defaults: applied=%d skipped=%d", len(r.Applied), len(r.Skipped))
	for _, k := range r.Applied {
		s += fmt.Sprintf("\n  + %s", k)
	}
	for _, k := range r.Skipped {
		s += fmt.Sprintf("\n  ~ %s (skipped)", k)
	}
	return s
}
