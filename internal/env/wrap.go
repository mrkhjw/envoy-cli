package env

import (
	"fmt"
	"strings"
)

// WrapOptions controls how values are wrapped.
type WrapOptions struct {
	MaxLength int
	Quote     bool
	Keys      []string
	DryRun    bool
}

// WrapResult holds the outcome of a Wrap operation.
type WrapResult struct {
	Wrapped  []Entry
	Modified int
	DryRun   bool
}

// Wrap ensures long values are quoted and optionally truncated with a suffix.
func Wrap(entries []Entry, opts WrapOptions) WrapResult {
	if opts.MaxLength <= 0 {
		opts.MaxLength = 80
	}

	keySet := make(map[string]bool, len(opts.Keys))
	for _, k := range opts.Keys {
		keySet[strings.ToUpper(k)] = true
	}

	result := WrapResult{DryRun: opts.DryRun}
	out := make([]Entry, 0, len(entries))

	for _, e := range entries {
		if e.Comment || e.Key == "" {
			out = append(out, e)
			continue
		}
		if len(keySet) > 0 && !keySet[strings.ToUpper(e.Key)] {
			out = append(out, e)
			continue
		}
		newVal := e.Value
		changed := false
		if len(e.Value) > opts.MaxLength {
			newVal = e.Value[:opts.MaxLength]
			changed = true
		}
		if opts.Quote && !strings.HasPrefix(newVal, "\"") {
			newVal = fmt.Sprintf("%q", newVal)
			changed = true
		}
		if changed {
			result.Modified++
		}
		newEntry := e
		if !opts.DryRun {
			newEntry.Value = newVal
		}
		out = append(out, newEntry)
	}

	result.Wrapped = out
	return result
}

// Format returns a human-readable summary of the WrapResult.
func (r WrapResult) Format() string {
	status := ""
	if r.DryRun {
		status = " (dry run)"
	}
	return fmt.Sprintf("wrap: %d value(s) modified%s", r.Modified, status)
}
