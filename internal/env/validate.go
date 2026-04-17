package env

import (
	"fmt"
	"strings"
)

// ValidationError represents a single validation issue.
type ValidationError struct {
	Line    int
	Key     string
	Message string
}

func (e ValidationError) Error() string {
	if e.Line > 0 {
		return fmt.Sprintf("line %d: %s: %s", e.Line, e.Key, e.Message)
	}
	return fmt.Sprintf("%s: %s", e.Key, e.Message)
}

// ValidationResult holds all errors found during validation.
type ValidationResult struct {
	Errors []ValidationError
}

func (r *ValidationResult) Valid() bool {
	return len(r.Errors) == 0
}

func (r *ValidationResult) Summary() string {
	if r.Valid() {
		return "✔ No validation errors found."
	}
	lines := []string{fmt.Sprintf("✖ %d validation error(s) found:", len(r.Errors))}
	for _, e := range r.Errors {
		lines = append(lines, "  - "+e.Error())
	}
	return strings.Join(lines, "\n")
}

// Validate checks an env map and raw lines for common issues.
func Validate(vars map[string]string, lines []string) *ValidationResult {
	result := &ValidationResult{}
	seen := map[string]int{}

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		parts := strings.SplitN(trimmed, "=", 2)
		if len(parts) != 2 {
			result.Errors = append(result.Errors, ValidationError{
				Line:    i + 1,
				Message: "invalid format, expected KEY=VALUE",
			})
			continue
		}
		key := strings.TrimSpace(parts[0])
		if key == "" {
			result.Errors = append(result.Errors, ValidationError{
				Line:    i + 1,
				Message: "empty key",
			})
			continue
		}
		if prev, dup := seen[key]; dup {
			result.Errors = append(result.Errors, ValidationError{
				Line:    i + 1,
				Key:     key,
				Message: fmt.Sprintf("duplicate key (first seen on line %d)", prev),
			})
		} else {
			seen[key] = i + 1
		}
		value := strings.TrimSpace(parts[1])
		if value == "" {
			result.Errors = append(result.Errors, ValidationError{
				Line:    i + 1,
				Key:     key,
				Message: "empty value",
			})
		}
	}
	_ = vars
	return result
}
