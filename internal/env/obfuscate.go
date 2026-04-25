package env

import (
	"fmt"
	"strings"
)

// ObfuscateOptions controls how values are obfuscated.
type ObfuscateOptions struct {
	Keys        []string // specific keys to obfuscate; empty means all secrets
	Style       string   // "star", "hash", "partial" (default: "star")
	RevealChars int      // number of trailing chars to reveal (partial style only)
	DryRun      bool
}

// ObfuscateEntry holds the original and obfuscated representation of a single entry.
type ObfuscateEntry struct {
	Key      string
	Original string
	Result   string
	Changed  bool
}

// ObfuscateResult holds the full output of an Obfuscate call.
type ObfuscateResult struct {
	Entries []ObfuscateEntry
	Total   int
	Changed int
}

// Format returns a human-readable summary of the obfuscation result.
func (r ObfuscateResult) Format() string {
	var sb strings.Builder
	for _, e := range r.Entries {
		if e.Changed {
			sb.WriteString(fmt.Sprintf("  ~ %s=%s\n", e.Key, e.Result))
		} else {
			sb.WriteString(fmt.Sprintf("    %s=%s\n", e.Key, e.Result))
		}
	}
	sb.WriteString(fmt.Sprintf("obfuscated %d/%d entries\n", r.Changed, r.Total))
	return sb.String()
}

// Obfuscate applies value obfuscation to the provided entries according to opts.
func Obfuscate(entries []Entry, opts ObfuscateOptions) ObfuscateResult {
	targetSet := make(map[string]bool, len(opts.Keys))
	for _, k := range opts.Keys {
		targetSet[strings.ToUpper(k)] = true
	}

	result := ObfuscateResult{Total: len(entries)}

	for _, e := range entries {
		if e.Comment {
			result.Entries = append(result.Entries, ObfuscateEntry{
				Key:    e.Key,
				Result: e.Value,
			})
			continue
		}

		shouldObfuscate := false
		if len(targetSet) > 0 {
			shouldObfuscate = targetSet[strings.ToUpper(e.Key)]
		} else {
			shouldObfuscate = isSecret(e.Key)
		}

		var obfuscated string
		if shouldObfuscate && e.Value != "" {
			obfuscated = applyObfuscation(e.Value, opts)
		} else {
			obfuscated = e.Value
		}

		changed := obfuscated != e.Value
		if changed {
			result.Changed++
		}
		result.Entries = append(result.Entries, ObfuscateEntry{
			Key:      e.Key,
			Original: e.Value,
			Result:   obfuscated,
			Changed:  changed,
		})
	}
	return result
}

func applyObfuscation(value string, opts ObfuscateOptions) string {
	switch opts.Style {
	case "hash":
		return strings.Repeat("#", len(value))
	case "partial":
		reveal := opts.RevealChars
		if reveal <= 0 {
			reveal = 4
		}
		if len(value) <= reveal {
			return strings.Repeat("*", len(value))
		}
		return strings.Repeat("*", len(value)-reveal) + value[len(value)-reveal:]
	default: // "star"
		return strings.Repeat("*", len(value))
	}
}
