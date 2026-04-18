package env

import (
	"fmt"
	"strings"
)

// ConvertFormat represents supported output formats for conversion
type ConvertFormat string

const (
	FormatEnv    ConvertFormat = "env"
	FormatExport ConvertFormat = "export"
	FormatYAML   ConvertFormat = "yaml"
	FormatTOML   ConvertFormat = "toml"
)

// ConvertResult holds the output of a conversion
type ConvertResult struct {
	Format ConvertFormat
	Output string
	Count  int
}

func (r ConvertResult) Format() string {
	return fmt.Sprintf("Converted %d key(s) to %s format.\n", r.Count, r.Format)
}

// Convert transforms a map of env vars into the specified format string
func Convert(entries map[string]string, format ConvertFormat, maskSecrets bool) (ConvertResult, error) {
	if maskSecrets {
		entries = MaskSecrets(entries)
	}

	var sb strings.Builder
	count := 0

	switch format {
	case FormatEnv:
		for k, v := range entries {
			fmt.Fprintf(&sb, "%s=%s\n", k, v)
			count++
		}
	case FormatExport:
		for k, v := range entries {
			fmt.Fprintf(&sb, "export %s=%s\n", k, v)
			count++
		}
	case FormatYAML:
		for k, v := range entries {
			fmt.Fprintf(&sb, "%s: \"%s\"\n", strings.ToLower(k), v)
			count++
		}
	case FormatTOML:
		for k, v := range entries {
			fmt.Fprintf(&sb, "%s = \"%s\"\n", strings.ToLower(k), v)
			count++
		}
	default:
		return ConvertResult{}, fmt.Errorf("unsupported format: %s", format)
	}

	return ConvertResult{
		Format: format,
		Output: sb.String(),
		Count:  count,
	}, nil
}
