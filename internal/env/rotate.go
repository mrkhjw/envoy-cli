package env

import (
	"fmt"
	"strings"
	"time"
)

// RotateResult holds the result of a key rotation operation.
type RotateResult struct {
	Rotated []string
	Skipped []string
	DryRun  bool
}

func (r RotateResult) Format() string {
	var sb strings.Builder
	if r.DryRun {
		sb.WriteString("[dry-run] ")
	}
	sb.WriteString(fmt.Sprintf("Rotated: %d, Skipped: %d\n", len(r.Rotated), len(r.Skipped)))
	for _, k := range r.Rotated {
		sb.WriteString(fmt.Sprintf("  ~ %s\n", k))
	}
	for _, k := range r.Skipped {
		sb.WriteString(fmt.Sprintf("  - %s (skipped)\n", k))
	}
	return sb.String()
}

// RotateOptions controls rotation behaviour.
type RotateOptions struct {
	Keys     []string // if empty, rotate all secret keys
	DryRun   bool
	Timestamp bool
}

// Rotate replaces values for secret keys (or specified keys) with a placeholder
// or timestamped sentinel, simulating a rotation workflow.
func Rotate(env map[string]string, opts RotateOptions) (map[string]string, RotateResult) {
	result := RotateResult{DryRun: opts.DryRun}
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = v
	}

	targets := opts.Keys
	if len(targets) == 0 {
		for k := range env {
			if isSecret(k) {
				targets = append(targets, k)
			}
		}
	}

	for _, k := range targets {
		if _, ok := env[k]; !ok {
			result.Skipped = append(result.Skipped, k)
			continue
		}
		newVal := "ROTATED"
		if opts.Timestamp {
			newVal = fmt.Sprintf("ROTATED_%s", time.Now().UTC().Format("20060102T150405"))
		}
		if !opts.DryRun {
			out[k] = newVal
		}
		result.Rotated = append(result.Rotated, k)
	}
	return out, result
}
