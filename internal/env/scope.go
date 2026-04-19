package env

import "fmt"

// ScopeResult holds the result of scoping entries to a specific environment.
type ScopeResult struct {
	Scope   string
	Entries []Entry
	Total   int
}

func (r ScopeResult) Format(mask bool) string {
	if len(r.Entries) == 0 {
		return fmt.Sprintf("[scope:%s] no entries found\n", r.Scope)
	}
	out := fmt.Sprintf("[scope:%s] %d of %d entries\n", r.Scope, len(r.Entries), r.Total)
	for _, e := range r.Entries {
		val := e.Value
		if mask && isSecret(e.Key) {
			val = "****"
		}
		out += fmt.Sprintf("  %s=%s\n", e.Key, val)
	}
	return out
}

// ScopeOptions controls how Scope filters entries.
type ScopeOptions struct {
	// Prefix is the environment prefix, e.g. "PROD", "DEV".
	Prefix    string
	// StripPrefix removes the prefix from keys in the result.
	StripPrefix bool
}

// Scope filters entries whose keys start with the given prefix.
func Scope(entries []Entry, opts ScopeOptions) ScopeResult {
	prefix := opts.Prefix
	if prefix == "" {
		return ScopeResult{Scope: "", Entries: entries, Total: len(entries)}
	}
	qualifier := prefix + "_"
	var matched []Entry
	for _, e := range entries {
		if len(e.Key) >= len(qualifier) && e.Key[:len(qualifier)] == qualifier {
			key := e.Key
			if opts.StripPrefix {
				key = e.Key[len(qualifier):]
			}
			matched = append(matched, Entry{Key: key, Value: e.Value, Raw: e.Raw})
		}
	}
	return ScopeResult{Scope: prefix, Entries: matched, Total: len(entries)}
}
