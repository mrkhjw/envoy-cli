package env

import (
	"fmt"
	"os"
	"strings"
)

// AnnotateFile reads a .env file, annotates matching keys, and writes the result back.
func AnnotateFile(path, annotation string, keys []string, dryRun bool) (AnnotateResult, error) {
	entries, err := ParseFile(path)
	if err != nil {
		return AnnotateResult{}, fmt.Errorf("annotate: failed to parse %s: %w", path, err)
	}

	result := Annotate(entries, annotation, keys, dryRun)

	if dryRun {
		return result, nil
	}

	var sb strings.Builder
	for _, ae := range result.Entries {
		if ae.Key == "" {
			sb.WriteString(ae.Value + "\n")
			continue
		}
		line := fmt.Sprintf("%s=%s", ae.Key, ae.Value)
		if ae.Annotation != "" {
			line += fmt.Sprintf(" # %s", ae.Annotation)
		}
		sb.WriteString(line + "\n")
	}

	if err := os.WriteFile(path, []byte(sb.String()), 0644); err != nil {
		return AnnotateResult{}, fmt.Errorf("annotate: failed to write %s: %w", path, err)
	}

	return result, nil
}
