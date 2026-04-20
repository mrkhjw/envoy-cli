package env

import (
	"strings"
)

// SanitizeOptions controls sanitize behavior.
type SanitizeOptions struct {
	StripControlChars bool
	NormalizeLineEndings bool
	TrimKeys bool
	TrimValues bool
	RemoveNullBytes bool
}

// SanitizeResult holds the result of a sanitize operation.
type SanitizeResult struct {
	Entries  []Entry
	Cleaned  int
	Unchanged int
}

// Format returns a human-readable summary of the sanitize result.
func (r SanitizeResult) Format() string {
	var sb strings.Builder
	for _, e := range r.Entries {
		if isSecret(e.Key) {
			sb.WriteString(e.Key + "=***\n")
		} else {
			sb.WriteString(e.Key + "=" + e.Value + "\n")
		}
	}
	sb.WriteString("\n")
	sb.WriteString("sanitized: " + itoa(r.Cleaned) + ", unchanged: " + itoa(r.Unchanged))
	return sb.String()
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	result := ""
	for n > 0 {
		result = string(rune('0'+n%10)) + result
		n /= 10
	}
	return result
}

// Sanitize cleans entries by applying the given options.
func Sanitize(entries []Entry, opts SanitizeOptions) SanitizeResult {
	result := SanitizeResult{}
	for _, e := range entries {
		origKey := e.Key
		origVal := e.Value

		if opts.TrimKeys {
			e.Key = strings.TrimSpace(e.Key)
		}
		if opts.TrimValues {
			e.Value = strings.TrimSpace(e.Value)
		}
		if opts.NormalizeLineEndings {
			e.Value = strings.ReplaceAll(e.Value, "\r\n", "\n")
			e.Value = strings.ReplaceAll(e.Value, "\r", "\n")
		}
		if opts.RemoveNullBytes {
			e.Key = strings.ReplaceAll(e.Key, "\x00", "")
			e.Value = strings.ReplaceAll(e.Value, "\x00", "")
		}
		if opts.StripControlChars {
			e.Value = stripControlChars(e.Value)
		}

		if e.Key != origKey || e.Value != origVal {
			result.Cleaned++
		} else {
			result.Unchanged++
		}
		result.Entries = append(result.Entries, e)
	}
	return result
}

func stripControlChars(s string) string {
	var sb strings.Builder
	for _, r := range s {
		if r >= 32 || r == '\t' || r == '\n' {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}
