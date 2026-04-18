package env

import (
	"bufio"
	"os"
	"strings"
)

// StripResult holds the outcome of a strip operation.
type StripResult struct {
	Removed []string
	Kept    int
}

func (r StripResult) Format() string {
	var sb strings.Builder
	if len(r.Removed) == 0 {
		sb.WriteString("No keys removed.\n")
		return sb.String()
	}
	sb.WriteString("Removed keys:\n")
	for _, k := range r.Removed {
		sb.WriteString("  - " + k + "\n")
	}
	return sb.String()
}

// Strip removes keys from entries that match the given keys list.
func Strip(entries []Entry, keys []string) ([]Entry, StripResult) {
	keySet := make(map[string]bool, len(keys))
	for _, k := range keys {
		keySet[strings.ToUpper(k)] = true
	}

	var result []Entry
	var removed []string

	for _, e := range entries {
		if keySet[strings.ToUpper(e.Key)] {
			removed = append(removed, e.Key)
		} else {
			result = append(result, e)
		}
	}

	return result, StripResult{Removed: removed, Kept: len(result)}
}

// StripFile removes specified keys from a .env file in place.
func StripFile(path string, keys []string) (StripResult, error) {
	entries, err := ParseFile(path)
	if err != nil {
		return StripResult{}, err
	}

	kept, result := Strip(entries, keys)

	f, err := os.Create(path)
	if err != nil {
		return StripResult{}, err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for _, e := range kept {
		_, err := w.WriteString(e.Key + "=" + e.Value + "\n")
		if err != nil {
			return StripResult{}, err
		}
	}

	return result, w.Flush()
}
