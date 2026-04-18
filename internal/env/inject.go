package env

import (
	"fmt"
	"os"
	"strings"
)

// InjectResult holds the result of an inject operation.
type InjectResult struct {
	Injected []string
	Skipped  []string
}

func (r InjectResult) Format() string {
	var sb strings.Builder
	if len(r.Injected) > 0 {
		sb.WriteString(fmt.Sprintf("Injected %d variable(s):\n", len(r.Injected)))
		for _, k := range r.Injected {
			sb.WriteString(fmt.Sprintf("  + %s\n", k))
		}
	}
	if len(r.Skipped) > 0 {
		sb.WriteString(fmt.Sprintf("Skipped %d variable(s) (already set):\n", len(r.Skipped)))
		for _, k := range r.Skipped {
			sb.WriteString(fmt.Sprintf("  ~ %s\n", k))
		}
	}
	if len(r.Injected) == 0 && len(r.Skipped) == 0 {
		sb.WriteString("Nothing to inject.\n")
	}
	return sb.String()
}

// Inject sets environment variables from the given map into the current process.
// If overwrite is false, existing env vars are skipped.
func Inject(vars map[string]string, overwrite bool) InjectResult {
	result := InjectResult{}
	for k, v := range vars {
		_, exists := os.LookupEnv(k)
		if exists && !overwrite {
			result.Skipped = append(result.Skipped, k)
			continue
		}
		os.Setenv(k, v)
		result.Injected = append(result.Injected, k)
	}
	return result
}

// InjectFile parses a .env file and injects its variables into the current process.
func InjectFile(path string, overwrite bool) (InjectResult, error) {
	vars, err := ParseFile(path)
	if err != nil {
		return InjectResult{}, fmt.Errorf("inject: %w", err)
	}
	return Inject(vars, overwrite), nil
}
