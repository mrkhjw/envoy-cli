package env

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Entry represents a single key-value pair from a .env file.
type Entry struct {
	Key    string
	Value  string
	Secret bool
}

// ParseFile reads a .env file and returns a slice of entries.
func ParseFile(path string) ([]Entry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()

	var entries []Entry
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.Trim(strings.TrimSpace(parts[1]), `"`)
		entries = append(entries, Entry{
			Key:    key,
			Value:  value,
			Secret: isSecret(key),
		})
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan %s: %w", path, err)
	}
	return entries, nil
}

// isSecret returns true if the key looks like it holds sensitive data.
func isSecret(key string) bool {
	secretKeywords := []string{"SECRET", "PASSWORD", "PASS", "TOKEN", "KEY", "PRIVATE", "CREDENTIAL"}
	upper := strings.ToUpper(key)
	for _, kw := range secretKeywords {
		if strings.Contains(upper, kw) {
			return true
		}
	}
	return false
}
