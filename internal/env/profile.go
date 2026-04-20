package env

import (
	"fmt"
	"sort"
	"strings"
)

// ProfileEntry represents a named env profile (e.g. dev, staging, prod)
type ProfileEntry struct {
	Name    string
	Entries []Entry
}

// ProfileResult holds the result of a profile switch/load operation
type ProfileResult struct {
	Profile string
	Loaded  int
	Entries []Entry
}

func (r ProfileResult) Format(mask bool) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Profile: %s (%d keys)\n", r.Profile, r.Loaded)
	entries := r.Entries
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})
	for _, e := range entries {
		val := e.Value
		if mask && isSecret(e.Key) {
			val = "***"
		}
		fmt.Fprintf(&sb, "  %s=%s\n", e.Key, val)
	}
	return strings.TrimRight(sb.String(), "\n")
}

// Profile filters entries from a map of named profiles by profile name.
// It returns a ProfileResult with the matched entries.
func Profile(profiles map[string][]Entry, name string) (ProfileResult, error) {
	entries, ok := profiles[name]
	if !ok {
		available := make([]string, 0, len(profiles))
		for k := range profiles {
			available = append(available, k)
		}
		sort.Strings(available)
		return ProfileResult{}, fmt.Errorf("profile %q not found; available: %s", name, strings.Join(available, ", "))
	}
	return ProfileResult{
		Profile: name,
		Loaded:  len(entries),
		Entries: entries,
	}, nil
}
