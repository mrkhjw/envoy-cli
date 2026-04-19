package env

import (
	"fmt"
	"strings"
)

type TransformOpts struct {
	UppercaseKeys bool
	LowercaseKeys bool
	UppercaseValues bool
	LowercaseValues bool
	TrimValues bool
	Keys []string // if empty, apply to all
}

type TransformResult struct {
	Entries  []Entry
	Modified int
}

func (r TransformResult) Format(mask bool) string {
	var sb strings.Builder
	for _, e := range r.Entries {
		val := e.Value
		if mask && isSecret(e.Key) {
			val = "****"
		}
		sb.WriteString(fmt.Sprintf("%s=%s\n", e.Key, val))
	}
	sb.WriteString(fmt.Sprintf("# %d key(s) transformed\n", r.Modified))
	return sb.String()
}

func Transform(entries []Entry, opts TransformOpts) TransformResult {
	targetSet := make(map[string]bool)
	for _, k := range opts.Keys {
		targetSet[strings.ToUpper(k)] = true
	}

	applyToAll := len(opts.Keys) == 0

	result := make([]Entry, 0, len(entries))
	modified := 0

	for _, e := range entries {
		origKey := e.Key
		origVal := e.Value

		inScope := applyToAll || targetSet[strings.ToUpper(e.Key)]

		if inScope {
			if opts.UppercaseKeys {
				e.Key = strings.ToUpper(e.Key)
			} else if opts.LowercaseKeys {
				e.Key = strings.ToLower(e.Key)
			}
			if opts.UppercaseValues {
				e.Value = strings.ToUpper(e.Value)
			} else if opts.LowercaseValues {
				e.Value = strings.ToLower(e.Value)
			}
			if opts.TrimValues {
				e.Value = strings.TrimSpace(e.Value)
			}
			if e.Key != origKey || e.Value != origVal {
				modified++
			}
		}

		result = append(result, e)
	}

	return TransformResult{Entries: result, Modified: modified}
}
