package env

import (
	"fmt"
	"strings"
)

// LookupResult holds the result of a key lookup.
type LookupResult struct {
	Key     string
	Value   string
	Found   bool
	Masked  bool
}

// LookupOptions controls lookup behaviour.
type LookupOptions struct {
	CaseSensitive bool
	MaskSecrets   bool
}

// Lookup finds a single key in a list of entries and returns its value.
func Lookup(entries []Entry, key string, opts LookupOptions) LookupResult {
	for _, e := range entries {
		if e.IsComment || e.Key == "" {
			continue
		}
		match := false
		if opts.CaseSensitive {
			match = e.Key == key
		} else {
			match = strings.EqualFold(e.Key, key)
		}
		if match {
			val := e.Value
			masked := false
			if opts.MaskSecrets && isSecret(e.Key) {
				val = "***"
				masked = true
			}
			return LookupResult{Key: e.Key, Value: val, Found: true, Masked: masked}
		}
	}
	return LookupResult{Key: key, Found: false}
}

// Format returns a human-readable string for a LookupResult.
func (r LookupResult) Format() string {
	if !r.Found {
		return fmt.Sprintf("key %q not found", r.Key)
	}
	if r.Masked {
		return fmt.Sprintf("%s=*** (masked)", r.Key)
	}
	return fmt.Sprintf("%s=%s", r.Key, r.Value)
}
