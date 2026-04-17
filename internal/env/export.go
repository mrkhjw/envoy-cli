package env

import (
	"fmt"
	"os"
	"strings"
)

// ExportFormat defines the output format for exported env vars
type ExportFormat string

const (
	FormatShell  ExportFormat = "shell"
	FormatDocker ExportFormat = "docker"
	FormatJSON   ExportFormat = "json"
)

// Export renders env vars in the specified format
func Export(vars map[string]string, format ExportFormat, maskSecrets bool) (string, error) {
	if maskSecrets {
		vars = MaskSecrets(vars)
	}

	switch format {
	case FormatShell:
		return exportShell(vars), nil
	case FormatDocker:
		return exportDocker(vars), nil
	case FormatJSON:
		return exportJSON(vars), nil
	default:
		return "", fmt.Errorf("unknown format: %s", format)
	}
}

// ExportFile reads a .env file and exports it in the given format
func ExportFile(path string, format ExportFormat, maskSecrets bool) (string, error) {
	vars, err := ParseFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to parse file: %w", err)
	}
	return Export(vars, format, maskSecrets)
}

func exportShell(vars map[string]string) string {
	var sb strings.Builder
	for k, v := range vars {
		fmt.Fprintf(&sb, "export %s=%q\n", k, v)
	}
	return sb.String()
}

func exportDocker(vars map[string]string) string {
	var sb strings.Builder
	for k, v := range vars {
		fmt.Fprintf(&sb, "-e %s=%s\n", k, v)
	}
	return sb.String()
}

func exportJSON(vars map[string]string) string {
	var sb strings.Builder
	sb.WriteString("{\n")
	i := 0
	for k, v := range vars {
		comma := ","
		if i == len(vars)-1 {
			comma = ""
		}
		fmt.Fprintf(&sb, "  %q: %q%s\n", k, v, comma)
		i++
	}
	sb.WriteString("}\n")
	return sb.String()
}

// WriteExport writes exported content to a file
func WriteExport(content, dest string) error {
	return os.WriteFile(dest, []byte(content), 0644)
}
