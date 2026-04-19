package env

import "fmt"

// ArchiveFile parses the given source .env file and archives it to dest.
func ArchiveFile(source, dest, label string) (ArchiveResult, error) {
	entries, err := ParseFile(source)
	if err != nil {
		return ArchiveResult{}, fmt.Errorf("parse error: %w", err)
	}

	return Archive(entries, dest, label)
}
