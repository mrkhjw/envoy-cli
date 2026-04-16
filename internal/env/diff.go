package env

import "fmt"

// DiffResult holds the differences between two env maps.
type DiffResult struct {
	Added   map[string]string
	Removed map[string]string
	Changed map[string][2]string // [0] = old, [1] = new
}

// Diff compares two env maps and returns the differences.
func Diff(base, target map[string]string) DiffResult {
	result := DiffResult{
		Added:   make(map[string]string),
		Removed: make(map[string]string),
		Changed: make(map[string][2]string),
	}

	for k, v := range target {
		baseVal, exists := base[k]
		if !exists {
			result.Added[k] = v
		} else if baseVal != v {
			result.Changed[k] = [2]string{baseVal, v}
		}
	}

	for k, v := range base {
		if _, exists := target[k]; !exists {
			result.Removed[k] = v
		}
	}

	return result
}

// Format returns a human-readable string of the diff.
func (d DiffResult) Format(maskSecrets bool) string {
	out := ""
	for k, v := range d.Added {
		val := v
		if maskSecrets && isSecret(k) {
			val = "***"
		}
		out += fmt.Sprintf("+ %s=%s\n", k, val)
	}
	for k, v := range d.Removed {
		val := v
		if maskSecrets && isSecret(k) {
			val = "***"
		}
		out += fmt.Sprintf("- %s=%s\n", k, val)
	}
	for k, v := range d.Changed {
		oldVal, newVal := v[0], v[1]
		if maskSecrets && isSecret(k) {
			oldVal = "***"
			newVal = "***"
		}
		out += fmt.Sprintf("~ %s: %s -> %s\n", k, oldVal, newVal)
	}
	return out
}
