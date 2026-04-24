package env

import (
	"fmt"
	"strings"
)

// AnnotateEntry holds a key-value pair with an associated comment annotation.
type AnnotateEntry struct {
	Key        string
	Value      string
	Annotation string
	IsSecret   bool
}

// AnnotateResult holds the outcome of an annotate operation.
type AnnotateResult struct {
	Entries  []AnnotateEntry
	Modified int
}

// Format returns a human-readable summary of the annotation result.
func (r AnnotateResult) Format() string {
	if r.Modified == 0 {
		return "No keys annotated."
	}
	return fmt.Sprintf("Annotated %d key(s).", r.Modified)
}

// Annotate adds inline comments to matching keys in the provided entries.
// If keys is empty, all entries are annotated with the given text.
// Dry-run mode returns the result without modifying entries in place.
func Annotate(entries []Entry, annotation string, keys []string, dryRun bool) AnnotateResult {
	keySet := make(map[string]bool, len(keys))
	for _, k := range keys {
		keySet[strings.ToUpper(k)] = true
	}

	result := AnnotateResult{}

	for _, e := range entries {
		if e.Comment {
			result.Entries = append(result.Entries, AnnotateEntry{
				Key:        "",
				Value:      e.Value,
				Annotation: "",
			})
			continue
		}

		target := len(keys) == 0 || keySet[strings.ToUpper(e.Key)]
		ae := AnnotateEntry{
			Key:      e.Key,
			Value:    e.Value,
			IsSecret: isSecret(e.Key),
		}

		if target && annotation != "" {
			ae.Annotation = annotation
			if !dryRun {
				result.Modified++
			}
		}

		result.Entries = append(result.Entries, ae)
	}

	return result
}
