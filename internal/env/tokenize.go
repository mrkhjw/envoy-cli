package env

import (
	"fmt"
	"strings"
)

// TokenizeOptions controls tokenization behavior.
type TokenizeOptions struct {
	Delimiter string
	Keys      []string
	MaskSecrets bool
}

// TokenizeResult holds the output of a tokenization operation.
type TokenizeResult struct {
	Tokens  map[string][]string
	Total   int
	Skipped int
}

// Format returns a human-readable summary of the tokenization result.
func (r TokenizeResult) Format(mask bool) string {
	var sb strings.Builder
	for key, tokens := range r.Tokens {
		if mask && isSecret(key) {
			sb.WriteString(fmt.Sprintf("%s = [REDACTED]\n", key))
			continue
		}
		sb.WriteString(fmt.Sprintf("%s = [%s]\n", key, strings.Join(tokens, ", ")))
	}
	sb.WriteString(fmt.Sprintf("total=%d skipped=%d", r.Total, r.Skipped))
	return sb.String()
}

// Tokenize splits env entry values by a delimiter and returns token lists per key.
func Tokenize(entries []Entry, opts TokenizeOptions) TokenizeResult {
	delim := opts.Delimiter
	if delim == "" {
		delim = ","
	}

	keySet := make(map[string]bool)
	for _, k := range opts.Keys {
		keySet[strings.ToUpper(k)] = true
	}

	tokens := make(map[string][]string)
	total := 0
	skipped := 0

	for _, e := range entries {
		if e.Comment || e.Key == "" {
			continue
		}
		if len(keySet) > 0 && !keySet[strings.ToUpper(e.Key)] {
			skipped++
			continue
		}
		parts := strings.Split(e.Value, delim)
		cleaned := make([]string, 0, len(parts))
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p != "" {
				cleaned = append(cleaned, p)
			}
		}
		tokens[e.Key] = cleaned
		total++
	}

	return TokenizeResult{
		Tokens:  tokens,
		Total:   total,
		Skipped: skipped,
	}
}
