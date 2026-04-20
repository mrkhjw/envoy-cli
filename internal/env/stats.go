package env

import "fmt"

// StatsResult holds summary statistics about an env map.
type StatsResult struct {
	Total    int
	Secrets  int
	Empty    int
	Comments int
	Unique   int
}

// Stats computes statistics over a slice of Entry.
func Stats(entries []Entry) StatsResult {
	seen := make(map[string]bool)
	result := StatsResult{}

	for _, e := range entries {
		if e.Comment {
			result.Comments++
			continue
		}
		result.Total++
		if e.Value == "" {
			result.Empty++
		}
		if isSecret(e.Key) {
			result.Secrets++
		}
		if !seen[e.Key] {
			seen[e.Key] = true
			result.Unique++
		}
	}
	return result
}

// Format returns a human-readable summary of StatsResult.
func (r StatsResult) Format() string {
	return fmt.Sprintf(
		"Total: %d | Secrets: %d | Empty: %d | Comments: %d | Unique Keys: %d",
		r.Total, r.Secrets, r.Empty, r.Comments, r.Unique,
	)
}
