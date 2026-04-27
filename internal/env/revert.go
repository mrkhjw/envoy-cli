package env

import (
	"fmt"
	"strings"
)

// RevertResult holds the outcome of a revert operation.
type RevertResult struct {
	Reverted []string
	Skipped  []string
	DryRun   bool
}

// RevertOptions configures the Revert operation.
type RevertOptions struct {
	Keys      []string
	DryRun    bool
	Overwrite bool
}

// Revert restores entries in current to their values from baseline.
// If Keys is specified, only those keys are reverted.
func Revert(current, baseline []Entry, opts RevertOptions) ([]Entry, RevertResult) {
	baseMap := make(map[string]string)
	for _, e := range baseline {
		if !e.IsComment && e.Key != "" {
			baseMap[strings.ToUpper(e.Key)] = e.Value
		}
	}

	keySet := make(map[string]bool)
	for _, k := range opts.Keys {
		keySet[strings.ToUpper(k)] = true
	}

	result := RevertResult{DryRun: opts.DryRun}
	out := make([]Entry, 0, len(current))

	for _, e := range current {
		if e.IsComment || e.Key == "" {
			out = append(out, e)
			continue
		}

		norm := strings.ToUpper(e.Key)
		if len(keySet) > 0 && !keySet[norm] {
			out = append(out, e)
			continue
		}

		baseVal, exists := baseMap[norm]
		if !exists {
			result.Skipped = append(result.Skipped, e.Key)
			out = append(out, e)
			continue
		}

		if e.Value == baseVal && !opts.Overwrite {
			result.Skipped = append(result.Skipped, e.Key)
			out = append(out, e)
			continue
		}

		result.Reverted = append(result.Reverted, e.Key)
		if !opts.DryRun {
			e.Value = baseVal
		}
		out = append(out, e)
	}

	return out, result
}

// Format returns a human-readable summary of the revert result.
func (r RevertResult) Format() string {
	var sb strings.Builder
	if r.DryRun {
		sb.WriteString("[dry-run] ")
	}
	sb.WriteString(fmt.Sprintf("Reverted: %d, Skipped: %d\n", len(r.Reverted), len(r.Skipped)))
	for _, k := range r.Reverted {
		sb.WriteString(fmt.Sprintf("  reverted: %s\n", k))
	}
	for _, k := range r.Skipped {
		sb.WriteString(fmt.Sprintf("  skipped:  %s\n", k))
	}
	return sb.String()
}
