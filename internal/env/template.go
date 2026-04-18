package env

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// TemplateResult holds the result of a template render operation.
type TemplateResult struct {
	Rendered  string
	Missing   []string
	Replaced  int
}

func (r TemplateResult) Format() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Replaced: %d variable(s)\n", r.Replaced)
	if len(r.Missing) > 0 {
		fmt.Fprintf(&sb, "Missing:  %s\n", strings.Join(r.Missing, ", "))
	}
	return sb.String()
}

var templateVarRe = regexp.MustCompile(`\$\{([A-Z_][A-Z0-9_]*)\}`)

// Template renders a template string by substituting ${KEY} placeholders
// with values from the provided env map.
func Template(tmpl string, env map[string]string) TemplateResult {
	missing := []string{}
	replaced := 0
	seen := map[string]bool{}

	result := templateVarRe.ReplaceAllStringFunc(tmpl, func(match string) string {
		key := templateVarRe.FindStringSubmatch(match)[1]
		if val, ok := env[key]; ok {
			replaced++
			return val
		}
		if !seen[key] {
			missing = append(missing, key)
			seen[key] = true
		}
		return match
	})

	return TemplateResult{Rendered: result, Missing: missing, Replaced: replaced}
}

// TemplateFile reads a template file, renders it with the given env map,
// and writes the result to outPath.
func TemplateFile(tmplPath string, env map[string]string, outPath string) (TemplateResult, error) {
	data, err := os.ReadFile(tmplPath)
	if err != nil {
		return TemplateResult{}, fmt.Errorf("read template: %w", err)
	}

	res := Template(string(data), env)

	if err := os.WriteFile(outPath, []byte(res.Rendered), 0644); err != nil {
		return TemplateResult{}, fmt.Errorf("write output: %w", err)
	}

	return res, nil
}
