package env

import (
	"fmt"
	"strings"
)

// RedactResult holds the redacted output and metadata.
type RedactResult struct {
	Lines    []string
	Redacted int
}

// Redact replaces secret values in a map with a placeholder and returns
// a slice of KEY=VALUE lines safe for display or logging.
func Redact(vars map[string]string, placeholder string) RedactResult {
	if placeholder == "" {
		placeholder = "***"
	}
	result := RedactResult{}
	for k, v := range vars {
		if isSecret(k) {
			result.Lines = append(result.Lines, fmt.Sprintf("%s=%s", k, placeholder))
			result.Redacted++
		} else {
			result.Lines = append(result.Lines, fmt.Sprintf("%s=%s", k, v))
		}
	}
	return result
}

// RedactString scans a raw .env-formatted string and redacts secret lines.
func RedactString(input, placeholder string) string {
	if placeholder == "" {
		placeholder = "***"
	}
	var out []string
	for _, line := range strings.Split(input, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			out = append(out, line)
			continue
		}
		parts := strings.SplitN(trimmed, "=", 2)
		if len(parts) == 2 && isSecret(parts[0]) {
			out = append(out, fmt.Sprintf("%s=%s", parts[0], placeholder))
		} else {
			out = append(out, line)
		}
	}
	return strings.Join(out, "\n")
}
