package env

import "fmt"

// VerifyResult holds the outcome of a verification check.
type VerifyResult struct {
	Matched  []string
	Mismatch []string
	Missing  []string
}

// VerifyOptions controls how verification behaves.
type VerifyOptions struct {
	Expected map[string]string // key -> expected value
	MaskSecrets bool
}

// Verify checks that entries in the provided slice match the expected values.
// Keys absent from entries are reported as Missing; differing values as Mismatch.
func Verify(entries []Entry, opts VerifyOptions) VerifyResult {
	actual := make(map[string]string)
	for _, e := range entries {
		if e.Key != "" {
			actual[e.Key] = e.Value
		}
	}

	var result VerifyResult
	for k, expectedVal := range opts.Expected {
		gotVal, ok := actual[k]
		if !ok {
			result.Missing = append(result.Missing, k)
			continue
		}
		if gotVal == expectedVal {
			result.Matched = append(result.Matched, k)
		} else {
			result.Mismatch = append(result.Mismatch, k)
		}
	}
	return result
}

// Format returns a human-readable summary of the VerifyResult.
func (r VerifyResult) Format() string {
	out := fmt.Sprintf("verified: %d matched, %d mismatched, %d missing",
		len(r.Matched), len(r.Mismatch), len(r.Missing))
	for _, k := range r.Mismatch {
		out += fmt.Sprintf("\n  MISMATCH: %s", k)
	}
	for _, k := range r.Missing {
		out += fmt.Sprintf("\n  MISSING:  %s", k)
	}
	return out
}

// OK returns true when there are no mismatches and no missing keys.
func (r VerifyResult) OK() bool {
	return len(r.Mismatch) == 0 && len(r.Missing) == 0
}
