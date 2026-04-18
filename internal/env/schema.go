package env

import (
	"fmt"
	"strings"
)

// SchemaEntry defines expected metadata for a single env key.
type SchemaEntry struct {
	Key      string
	Required bool
	Default  string
}

// SchemaResult holds the outcome of a schema validation.
type SchemaResult struct {
	Missing  []string
	Extra    []string
	Defaults map[string]string
}

func (r SchemaResult) Format() string {
	var sb strings.Builder
	if len(r.Missing) == 0 && len(r.Extra) == 0 {
		sb.WriteString("schema: all required keys present\n")
	}
	for _, k := range r.Missing {
		sb.WriteString(fmt.Sprintf("missing: %s\n", k))
	}
	for _, k := range r.Extra {
		sb.WriteString(fmt.Sprintf("extra:   %s\n", k))
	}
	return sb.String()
}

// ValidateSchema checks env entries against a schema definition.
func ValidateSchema(entries []Entry, schema []SchemaEntry) SchemaResult {
	result := SchemaResult{Defaults: make(map[string]string)}

	envKeys := make(map[string]bool)
	for _, e := range entries {
		envKeys[e.Key] = true
	}

	schemaKeys := make(map[string]bool)
	for _, s := range schema {
		schemaKeys[s.Key] = true
		if !envKeys[s.Key] {
			if s.Required {
				result.Missing = append(result.Missing, s.Key)
			} else if s.Default != "" {
				result.Defaults[s.Key] = s.Default
			}
		}
	}

	for _, e := range entries {
		if !schemaKeys[e.Key] {
			result.Extra = append(result.Extra, e.Key)
		}
	}

	return result
}
