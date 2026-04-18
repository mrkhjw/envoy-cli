package env

import (
	"fmt"
	"regexp"
	"strings"
)

var interpolatePattern = regexp.MustCompile(`\$\{([A-Za-z_][A-Za-z0-9_]*)\}`)

// InterpolateResult holds the result of an interpolation operation.
type InterpolateResult struct {
	Resolved map[string]string
	Unresolved []string
}

func (r InterpolateResult) Format() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Resolved: %d key(s)\n", len(r.Resolved)))
	if len(r.Unresolved) > 0 {
		sb.WriteString(fmt.Sprintf("Unresolved references: %s\n", strings.Join(r.Unresolved, ", ")))
	}
	return strings.TrimRight(sb.String(), "\n")
}

// Interpolate expands ${VAR} references in values using other keys in the map.
// It performs a single-pass expansion and collects any unresolved references.
func Interpolate(env map[string]string) InterpolateResult {
	resolved := make(map[string]string, len(env))
	unresolvedSet := map[string]struct{}{}

	for k, v := range env {
		expanded := interpolatePattern.ReplaceAllStringFunc(v, func(match string) string {
			inner := match[2 : len(match)-1]
			if val, ok := env[inner]; ok {
				return val
			}
			unresolvedSet[inner] = struct{}{}
			return match
		})
		resolved[k] = expanded
	}

	unresolved := make([]string, 0, len(unresolvedSet))
	for k := range unresolvedSet {
		unresolved = append(unresolved, k)
	}

	return InterpolateResult{
		Resolved:   resolved,
		Unresolved: unresolved,
	}
}
