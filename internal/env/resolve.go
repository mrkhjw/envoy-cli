package env

import "fmt"

// ResolveOption configures how resolution behaves.
type ResolveOption struct {
	Strict bool // fail if any reference is unresolved
}

// ResolveResult holds the outcome of a Resolve operation.
type ResolveResult struct {
	Resolved   []EnvEntry
	Unresolved []string
}

// Format returns a human-readable summary of the resolve result.
func (r ResolveResult) Format() string {
	out := fmt.Sprintf("Resolved: %d entries\n", len(r.Resolved))
	if len(r.Unresolved) == 0 {
		out += "All references resolved successfully.\n"
		return out
	}
	out += fmt.Sprintf("Unresolved references (%d):\n", len(r.Unresolved))
	for _, k := range r.Unresolved {
		out += fmt.Sprintf("  ! %s\n", k)
	}
	return out
}

// Resolve expands ${VAR} references in entry values using the provided env map.
// It performs a single-pass resolution and tracks any keys that remain unresolved.
func Resolve(entries []EnvEntry, opts ResolveOption) (ResolveResult, error) {
	lookup := make(map[string]string, len(entries))
	for _, e := range entries {
		lookup[e.Key] = e.Value
	}

	var result ResolveResult
	unresolvedSet := make(map[string]bool)

	for _, e := range entries {
		resolved, missing := expandValue(e.Value, lookup)
		for _, m := range missing {
			unresolvedSet[m] = true
		}
		e.Value = resolved
		result.Resolved = append(result.Resolved, e)
	}

	for k := range unresolvedSet {
		result.Unresolved = append(result.Unresolved, k)
	}

	if opts.Strict && len(result.Unresolved) > 0 {
		return result, fmt.Errorf("unresolved references: %v", result.Unresolved)
	}

	return result, nil
}

// expandValue replaces ${KEY} tokens in s using the lookup map.
// Returns the expanded string and a list of keys that were not found.
func expandValue(s string, lookup map[string]string) (string, []string) {
	var missing []string
	out := ""
	i := 0
	for i < len(s) {
		if i+1 < len(s) && s[i] == '$' && s[i+1] == '{' {
			end := -1
			for j := i + 2; j < len(s); j++ {
				if s[j] == '}' {
					end = j
					break
				}
			}
			if end == -1 {
				out += string(s[i])
				i++
				continue
			}
			key := s[i+2 : end]
			if val, ok := lookup[key]; ok {
				out += val
			} else {
				missing = append(missing, key)
				out += s[i : end+1]
			}
			i = end + 1
		} else {
			out += string(s[i])
			i++
		}
	}
	return out, missing
}
