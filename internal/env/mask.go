package env

import (
	"fmt"
	"strings"
)

// MaskOptions controls how masking is applied
type MaskOptions struct {
	Placeholder  string
	RevealPrefix int // number of leading chars to reveal
}

// MaskResult holds the result of a masking operation
type MaskResult struct {
	Entries []Entry
	Masked  int
	Total   int
}

// Mask applies secret masking to entries based on options
func Mask(entries []Entry, opts MaskOptions) MaskResult {
	if opts.Placeholder == "" {
		opts.Placeholder = "***"
	}

	result := MaskResult{Total: len(entries)}
	masked := make([]Entry, 0, len(entries))

	for _, e := range entries {
		if isSecret(e.Key) {
			val := applyMask(e.Value, opts)
			masked = append(masked, Entry{Key: e.Key, Value: val, Raw: e.Raw})
			result.Masked++
		} else {
			masked = append(masked, e)
		}
	}

	result.Entries = masked
	return result
}

func applyMask(val string, opts MaskOptions) string {
	if opts.RevealPrefix > 0 && len(val) > opts.RevealPrefix {
		return val[:opts.RevealPrefix] + opts.Placeholder
	}
	return opts.Placeholder
}

// Format returns a human-readable summary of the mask operation
func (r MaskResult) Format() string {
	var sb strings.Builder
	sb.WriteString("Mask summary:\n")
	for _, e := range r.Entries {
		sb.WriteString(fmt.Sprintf("  %s=%s\n", e.Key, e.Value))
	}
	sb.WriteString(fmt.Sprintf("Masked %d/%d entries.", r.Masked, r.Total))
	return sb.String()
}
