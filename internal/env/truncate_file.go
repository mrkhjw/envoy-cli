package env

import "fmt"

// TruncateFile reads a .env file, truncates long values, and writes the result
// back to the same path unless DryRun is set.
func TruncateFile(path string, opts TruncateOptions) (TruncateResult, error) {
	entries, err := ParseFile(path)
	if err != nil {
		return TruncateResult{}, fmt.Errorf("truncate: %w", err)
	}

	res := Truncate(entries, opts)

	if !opts.DryRun {
		if err := writeEnvFile(path, res.Entries); err != nil {
			return TruncateResult{}, fmt.Errorf("truncate write: %w", err)
		}
	}

	return res, nil
}
