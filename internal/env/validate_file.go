package env

import (
	"bufio"
	"fmt"
	"os"
)

// ValidateFile reads a .env file and runs validation on its contents.
func ValidateFile(path string) (*ValidationResult, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open file %q: %w", path, err)
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file %q: %w", path, err)
	}

	vars, err := ParseFile(path)
	if err != nil {
		return nil, err
	}

	return Validate(vars, lines), nil
}
