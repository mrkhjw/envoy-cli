package env

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// LoadSchema reads a schema file where each line is:
// KEY [required|optional] [default=VALUE]
func LoadSchema(path string) ([]SchemaEntry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("schema: cannot open %s: %w", path, err)
	}
	defer f.Close()

	var entries []SchemaEntry
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}
		e := SchemaEntry{Key: parts[0]}
		for _, p := range parts[1:] {
			switch {
			case p == "required":
				e.Required = true
			case strings.HasPrefix(p, "default="):
				e.Default = strings.TrimPrefix(p, "default=")
			}
		}
		entries = append(entries, e)
	}
	return entries, scanner.Err()
}
