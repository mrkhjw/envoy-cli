package env

import (
	"fmt"
	"os"
	"strings"
)

// RenameResult holds the result of a rename operation.
type RenameResult struct {
	OldKey  string
	NewKey  string
	Renamed bool
}

func (r RenameResult) Format() string {
	if !r.Renamed {
		return fmt.Sprintf("key %q not found, nothing renamed", r.OldKey)
	}
	return fmt.Sprintf("renamed %q -> %q", r.OldKey, r.NewKey)
}

// Rename renames oldKey to newKey in the provided env map.
// Returns an error if newKey already exists and overwrite is false.
func Rename(entries []EnvEntry, oldKey, newKey string, overwrite bool) ([]EnvEntry, RenameResult, error) {
	newKeyUpper := strings.ToUpper(newKey)
	oldKeyUpper := strings.ToUpper(oldKey)

	oldIdx := -1
	newIdx := -1
	for i, e := range entries {
		if strings.ToUpper(e.Key) == oldKeyUpper {
			oldIdx = i
		}
		if strings.ToUpper(e.Key) == newKeyUpper {
			newIdx = i
		}
	}

	if oldIdx == -1 {
		return entries, RenameResult{OldKey: oldKey, NewKey: newKey, Renamed: false}, nil
	}

	if newIdx != -1 && !overwrite {
		return entries, RenameResult{}, fmt.Errorf("key %q already exists; use --overwrite to replace", newKey)
	}

	updated := make([]EnvEntry, 0, len(entries))
	for i, e := range entries {
		if i == newIdx && newIdx != -1 {
			continue // remove old newKey entry
		}
		if i == oldIdx {
			e.Key = newKey
		}
		updated = append(updated, e)
	}

	return updated, RenameResult{OldKey: oldKey, NewKey: newKey, Renamed: true}, nil
}

// RenameFile renames a key in a .env file and writes the result back.
func RenameFile(path, oldKey, newKey string, overwrite bool) (RenameResult, error) {
	entries, err := ParseFile(path)
	if err != nil {
		return RenameResult{}, err
	}

	updated, result, err := Rename(entries, oldKey, newKey, overwrite)
	if err != nil {
		return RenameResult{}, err
	}

	if !result.Renamed {
		return result, nil
	}

	var sb strings.Builder
	for _, e := range updated {
		sb.WriteString(fmt.Sprintf("%s=%s\n", e.Key, e.Value))
	}

	if err := os.WriteFile(path, []byte(sb.String()), 0644); err != nil {
		return RenameResult{}, err
	}

	return result, nil
}
