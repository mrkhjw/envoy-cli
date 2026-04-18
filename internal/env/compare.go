package env

import "fmt"

// CompareResult holds the result of comparing two env files
type CompareResult struct {
	File1    string
	File2    string
	OnlyIn1  map[string]string
	OnlyIn2  map[string]string
	InBoth   map[string]string
	Conflict map[string][2]string
}

// Compare compares two parsed env maps and returns a CompareResult
func Compare(file1, file2 string, env1, env2 map[string]string) CompareResult {
	result := CompareResult{
		File1:    file1,
		File2:    file2,
		OnlyIn1:  make(map[string]string),
		OnlyIn2:  make(map[string]string),
		InBoth:   make(map[string]string),
		Conflict: make(map[string][2]string),
	}

	for k, v1 := range env1 {
		if v2, ok := env2[k]; ok {
			if v1 == v2 {
				result.InBoth[k] = v1
			} else {
				result.Conflict[k] = [2]string{v1, v2}
			}
		} else {
			result.OnlyIn1[k] = v1
		}
	}

	for k, v2 := range env2 {
		if _, ok := env1[k]; !ok {
			result.OnlyIn2[k] = v2
		}
	}

	return result
}

// Summary returns a human-readable summary of the comparison
func (c CompareResult) Summary(maskSecrets bool) string {
	out := fmt.Sprintf("Comparing %s <-> %s\n", c.File1, c.File2)
	out += fmt.Sprintf("  Identical keys : %d\n", len(c.InBoth))
	out += fmt.Sprintf("  Only in %s: %d\n", c.File1, len(c.OnlyIn1))
	out += fmt.Sprintf("  Only in %s: %d\n", c.File2, len(c.OnlyIn2))
	out += fmt.Sprintf("  Conflicts      : %d\n", len(c.Conflict))

	if len(c.Conflict) > 0 {
		out += "\nConflicting keys:\n"
		for k, vals := range c.Conflict {
			v1, v2 := vals[0], vals[1]
			if maskSecrets && isSecret(k) {
				v1, v2 = "***", "***"
			}
			out += fmt.Sprintf("  %s: %q vs %q\n", k, v1, v2)
		}
	}

	return out
}
