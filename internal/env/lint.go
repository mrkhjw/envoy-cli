package env

import (
	"fmt"
	"strings"
)

type LintIssue struct {
	Line    int
	Key     string
	Message string
	Severity string
}

type LintResult struct {
	Issues []LintIssue
}

func (r *LintResult) HasErrors() bool {
	for _, issue := range r.Issues {
		if issue.Severity == "error" {
			return true
		}
	}
	return false
}

func (r *LintResult) Format() string {
	if len(r.Issues) == 0 {
		return "No lint issues found."
	}
	var sb strings.Builder
	for _, issue := range r.Issues {
		sb.WriteString(fmt.Sprintf("[%s] line %d: %s\n", strings.ToUpper(issue.Severity), issue.Line, issue.Message))
	}
	return strings.TrimRight(sb.String(), "\n")
}

func Lint(lines []string) *LintResult {
	result := &LintResult{}
	seen := map[string]int{}

	for i, raw := range lines {
		lineNum := i + 1
		line := strings.TrimSpace(raw)

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if !strings.Contains(line, "=") {
			result.Issues = append(result.Issues, LintIssue{
				Line:     lineNum,
				Message:  fmt.Sprintf("invalid format: %q", line),
				Severity: "error",
			})
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])

		if key != strings.ToUpper(key) {
			result.Issues = append(result.Issues, LintIssue{
				Line:     lineNum,
				Key:      key,
				Message:  fmt.Sprintf("key %q should be uppercase", key),
				Severity: "warning",
			})
		}

		if val == "" {
			result.Issues = append(result.Issues, LintIssue{
				Line:     lineNum,
				Key:      key,
				Message:  fmt.Sprintf("key %q has empty value", key),
				Severity: "warning",
			})
		}

		if prev, ok := seen[key]; ok {
			result.Issues = append(result.Issues, LintIssue{
				Line:     lineNum,
				Key:      key,
				Message:  fmt.Sprintf("duplicate key %q (first seen on line %d)", key, prev),
				Severity: "error",
			})
		} else {
			seen[key] = lineNum
		}
	}

	return result
}
