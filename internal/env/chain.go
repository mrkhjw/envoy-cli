package env

import "fmt"

// ChainEntry represents a single step in a chain of env transformations.
type ChainEntry struct {
	Step   int
	Action string
	Key    string
	Before string
	After  string
}

// ChainResult holds the output of a Chain operation.
type ChainResult struct {
	Entries []Entry
	Log     []ChainEntry
}

// ChainOptions controls which transformations are applied in sequence.
type ChainOptions struct {
	TrimValues  bool
	UpperKeys   bool
	MaskSecrets bool
	DryRun      bool
}

// Chain applies a sequence of transformations to a slice of entries and
// records each change in a log for auditability.
func Chain(entries []Entry, opts ChainOptions) ChainResult {
	result := make([]Entry, len(entries))
	copy(result, entries)

	var log []ChainEntry
	step := 0

	if opts.TrimValues {
		step++
		for i, e := range result {
			if e.Comment {
				continue
			}
			trimmed := trimWhitespace(e.Value)
			if trimmed != e.Value {
				log = append(log, ChainEntry{Step: step, Action: "trim", Key: e.Key, Before: e.Value, After: trimmed})
				result[i].Value = trimmed
			}
		}
	}

	if opts.UpperKeys {
		step++
		for i, e := range result {
			if e.Comment {
				continue
			}
			upper := toUpper(e.Key)
			if upper != e.Key {
				log = append(log, ChainEntry{Step: step, Action: "uppercase_key", Key: e.Key, Before: e.Key, After: upper})
				result[i].Key = upper
			}
		}
	}

	if opts.MaskSecrets {
		step++
		for i, e := range result {
			if e.Comment || !isSecret(e.Key) {
				continue
			}
			masked := "****"
			log = append(log, ChainEntry{Step: step, Action: "mask", Key: e.Key, Before: e.Value, After: masked})
			result[i].Value = masked
		}
	}

	return ChainResult{Entries: result, Log: log}
}

// Format returns a human-readable summary of the chain log.
func (r ChainResult) Format() string {
	if len(r.Log) == 0 {
		return "chain: no transformations applied"
	}
	out := fmt.Sprintf("chain: %d transformation(s) applied\n", len(r.Log))
	for _, l := range r.Log {
		out += fmt.Sprintf("  step %d [%s] %s: %q -> %q\n", l.Step, l.Action, l.Key, l.Before, l.After)
	}
	return out
}

func trimWhitespace(s string) string {
	start, end := 0, len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t') {
		end--
	}
	return s[start:end]
}

func toUpper(s string) string {
	b := []byte(s)
	for i, c := range b {
		if c >= 'a' && c <= 'z' {
			b[i] = c - 32
		}
	}
	return string(b)
}
