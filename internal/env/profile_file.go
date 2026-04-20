package env

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// LoadProfiles reads a multi-profile env file.
// Profiles are delimited by "# @profile <name>" comment headers.
// Lines before the first header belong to the "default" profile.
func LoadProfiles(path string) (map[string][]Entry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open profile file: %w", err)
	}
	defer f.Close()

	profiles := make(map[string][]Entry)
	current := "default"

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "# @profile ") {
			current = strings.TrimPrefix(line, "# @profile ")
			current = strings.TrimSpace(current)
			continue
		}
		if strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := stripQuotes(strings.TrimSpace(parts[1]))
		profiles[current] = append(profiles[current], Entry{Key: key, Value: val})
	}
	return profiles, scanner.Err()
}
