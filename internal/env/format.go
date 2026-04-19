package env

import (
	"fmt"
	"strings"
)

// FormatOptions controls how entries are formatted for output.
type FormatOptions struct {
	QuoteValues  bool
	ExportPrefix bool
	UppercaseKeys bool
	MaskSecrets  bool
}

// FormatResult holds the formatted output lines.
type FormatResult struct {
	Lines []string
	Total int
}

// Format applies formatting rules to a list of env entries.
func Format(entries []Entry, opts FormatOptions) FormatResult {
	var lines []string
	for _, e := range entries {
		key := e.Key
		val := e.Value

		if opts.UppercaseKeys {
			key = strings.ToUpper(key)
		}

		if opts.MaskSecrets && isSecret(key) {
			val = "****"
		}

		if opts.QuoteValues && !strings.HasPrefix(val, "\"") {
			val = fmt.Sprintf("%q", val)
		}

		line := fmt.Sprintf("%s=%s", key, val)
		if opts.ExportPrefix {
			line = "export " + line
		}

		lines = append(lines, line)
	}
	return FormatResult{Lines: lines, Total: len(lines)}
}

// String returns all formatted lines joined by newlines.
func (r FormatResult) String() string {
	return strings.Join(r.Lines, "\n")
}
