package env

import "fmt"

// HealthIssue represents a single health check finding.
type HealthIssue struct {
	Key      string
	Severity string // "error", "warn", "info"
	Message  string
}

// HealthResult holds the outcome of a health check.
type HealthResult struct {
	Issues []HealthIssue
	Total  int
	Errors int
	Warns  int
}

func (r HealthResult) Format() string {
	if len(r.Issues) == 0 {
		return "✔ No health issues found."
	}
	out := fmt.Sprintf("Health check: %d issue(s) found (%d errors, %d warnings)\n", len(r.Issues), r.Errors, r.Warns)
	for _, issue := range r.Issues {
		var icon string
		switch issue.Severity {
		case "error":
			icon = "✖"
		case "warn":
			icon = "⚠"
		default:
			icon = "ℹ"
		}
		out += fmt.Sprintf("  %s [%s] %s: %s\n", icon, issue.Severity, issue.Key, issue.Message)
	}
	return out
}

// Health checks a map of env entries for common issues.
func Health(entries []Entry) HealthResult {
	result := HealthResult{Total: len(entries)}
	seen := map[string]bool{}

	for _, e := range entries {
		if e.Key == "" {
			continue
		}
		if seen[e.Key] {
			result.Issues = append(result.Issues, HealthIssue{Key: e.Key, Severity: "error", Message: "duplicate key"})
			result.Errors++
		}
		seen[e.Key] = true

		if e.Value == "" {
			sev := "warn"
			if isSecret(e.Key) {
				sev = "error"
				result.Errors++
			} else {
				result.Warns++
			}
			result.Issues = append(result.Issues, HealthIssue{Key: e.Key, Severity: sev, Message: "empty value"})
		}

		if isSecret(e.Key) && len(e.Value) < 8 && e.Value != "" {
			result.Issues = append(result.Issues, HealthIssue{Key: e.Key, Severity: "warn", Message: "secret value may be too short"})
			result.Warns++
		}
	}
	return result
}
