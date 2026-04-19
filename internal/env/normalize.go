package env

import (
	"strings"
)

// NormalizeOptions controls how keys and values are normalized.
type NormalizeOptions struct {
	UppercaseKeys   bool
	TrimValues      bool
	QuoteValues     bool
	StripExported   bool
}

// NormalizeResult holds the result of a normalization pass.
type NormalizeResult struct {
	Entries  []Entry
	Changed  []string
}

// Normalize applies normalization rules to a list of entries.
func Normalize(entries []Entry, opts NormalizeOptions) NormalizeResult {
	result := NormalizeResult{}
	for _, e := range entries {
		origKey := e.Key
		origVal := e.Value

		if opts.StripExported {
			e.Key = strings.TrimPrefix(e.Key, "export ")
		}
		if opts.UppercaseKeys {
			e.Key = strings.ToUpper(e.Key)
		}
		if opts.TrimValues {
			e.Value = strings.TrimSpace(e.Value)
		}
		if opts.QuoteValues && !strings.HasPrefix(e.Value, `"`) {
			e.Value = `"` + e.Value + `"`
		}

		if e.Key != origKey || e.Value != origVal {
			result.Changed = append(result.Changed, origKey)
		}
		result.Entries = append(result.Entries, e)
	}
	return result
}

// Format returns a human-readable summary of normalization changes.
func (r NormalizeResult) Format() string {
	if len(r.Changed) == 0 {
		return "no changes"
	}
	var sb strings.Builder
	sb.WriteString("normalized keys:\n")
	for _, k := range r.Changed {
		sb.WriteString("  ~ " + k + "\n")
	}
	return strings.TrimRight(sb.String(), "\n")
}
