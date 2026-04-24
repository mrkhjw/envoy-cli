package env

import "fmt"

// ProtectOptions configures the Protect operation.
type ProtectOptions struct {
	Keys     []string
	ReadOnly bool
	DryRun   bool
}

// ProtectResult holds the outcome of a Protect operation.
type ProtectResult struct {
	Protected []string
	Skipped   []string
	DryRun    bool
}

// Format returns a human-readable summary of the protect result.
func (r ProtectResult) Format() string {
	if len(r.Protected) == 0 && len(r.Skipped) == 0 {
		return "no keys matched for protection"
	}
	out := ""
	if r.DryRun {
		out += "[dry-run] "
	}
	for _, k := range r.Protected {
		mode := "locked"
		out += fmt.Sprintf("  protected: %s (%s)\n", k, mode)
	}
	for _, k := range r.Skipped {
		out += fmt.Sprintf("  skipped:   %s (already protected)\n", k)
	}
	out += fmt.Sprintf("%d protected, %d skipped", len(r.Protected), len(r.Skipped))
	return out
}

// Protect marks specified keys as protected by appending a "#protected" comment
// marker after each matching key=value line. If no keys are specified, all
// secret keys are targeted.
func Protect(entries []Entry, opts ProtectOptions) ([]Entry, ProtectResult) {
	targetAll := len(opts.Keys) == 0
	keySet := make(map[string]bool, len(opts.Keys))
	for _, k := range opts.Keys {
		keySet[normalizeKey(k)] = true
	}

	result := ProtectResult{DryRun: opts.DryRun}
	out := make([]Entry, 0, len(entries))

	for _, e := range entries {
		if e.Comment {
			out = append(out, e)
			continue
		}
		nk := normalizeKey(e.Key)
		wantProtect := (targetAll && isSecret(e.Key)) || keySet[nk]
		alreadyProtected := e.Tags != nil && e.Tags["protected"]

		if wantProtect && alreadyProtected {
			result.Skipped = append(result.Skipped, e.Key)
			out = append(out, e)
			continue
		}

		if wantProtect && !opts.DryRun {
			if e.Tags == nil {
				e.Tags = map[string]bool{}
			}
			e.Tags["protected"] = true
			result.Protected = append(result.Protected, e.Key)
		} else if wantProtect {
			result.Protected = append(result.Protected, e.Key)
		}
		out = append(out, e)
	}
	return out, result
}
