package env

import (
	"fmt"
	"strings"
)

// TagEntry represents an env entry with associated tags.
type TagEntry struct {
	Key   string
	Value string
	Tags  []string
}

// TagResult holds the result of a tagging operation.
type TagResult struct {
	Tagged   []TagEntry
	Untagged []string
}

// TagOptions controls Tag behaviour.
type TagOptions struct {
	Tags     []string
	Keys     []string // if empty, tag all
	MaskSecrets bool
}

// Tag annotates env entries with the given tags as inline comments.
func Tag(entries []Entry, opts TagOptions) TagResult {
	tagSet := make(map[string]bool)
	for _, t := range opts.Tags {
		tagSet[strings.TrimSpace(t)] = true
	}

	keySet := make(map[string]bool)
	for _, k := range opts.Keys {
		keySet[strings.ToUpper(strings.TrimSpace(k))] = true
	}

	result := TagResult{}
	for _, e := range entries {
		norm := strings.ToUpper(e.Key)
		if len(keySet) > 0 && !keySet[norm] {
			result.Untagged = append(result.Untagged, e.Key)
			continue
		}
		val := e.Value
		if opts.MaskSecrets && isSecret(e.Key) {
			val = "****"
		}
		tagList := make([]string, 0, len(tagSet))
		for t := range tagSet {
			tagList = append(tagList, t)
		}
		result.Tagged = append(result.Tagged, TagEntry{Key: e.Key, Value: val, Tags: tagList})
	}
	return result
}

// Format returns a human-readable summary of the TagResult.
func (r TagResult) Format() string {
	var sb strings.Builder
	for _, te := range r.Tagged {
		sb.WriteString(fmt.Sprintf("%s=%s # tags:%s\n", te.Key, te.Value, strings.Join(te.Tags, ",")))
	}
	if len(r.Untagged) > 0 {
		sb.WriteString(fmt.Sprintf("untagged: %s\n", strings.Join(r.Untagged, ", ")))
	}
	return sb.String()
}
