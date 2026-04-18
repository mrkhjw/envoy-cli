package env

import (
	"fmt"
	"sort"
	"strings"
)

// FlattenOptions controls how flattening is performed.
type FlattenOptions struct {
	Prefix    string
	Separator string
	Uppercase bool
}

// FlattenResult holds the result of a flatten operation.
type FlattenResult struct {
	Entries []Entry
	Renamed int
}

// Entry represents a single key-value pair.
type Entry struct {
	Key   string
	Value string
}

// Flatten normalizes env map keys by applying prefix, separator, and casing rules.
func Flatten(env map[string]string, opts FlattenOptions) FlattenResult {
	if opts.Separator == "" {
		opts.Separator = "_"
	}

	result := FlattenResult{}
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := env[k]
		newKey := k

		if opts.Uppercase {
			newKey = strings.ToUpper(newKey)
		}

		if opts.Prefix != "" && !strings.HasPrefix(newKey, opts.Prefix) {
			newKey = opts.Prefix + opts.Separator + newKey
			result.Renamed++
		}

		result.Entries = append(result.Entries, Entry{Key: newKey, Value: v})
	}

	return result
}

// Format returns a human-readable summary of the flatten result.
func (r FlattenResult) Format() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Flattened %d keys (%d renamed)\n", len(r.Entries), r.Renamed))
	for _, e := range r.Entries {
		val := e.Value
		if isSecret(e.Key) {
			val = "****"
		}
		sb.WriteString(fmt.Sprintf("  %s=%s\n", e.Key, val))
	}
	return sb.String()
}
