package env

import "fmt"

// WrapFile reads a .env file, applies wrap options, and writes the result back.
func WrapFile(path string, opts WrapOptions) (WrapResult, error) {
	entries, err := ParseFile(path)
	if err != nil {
		return WrapResult{}, fmt.Errorf("wrap: failed to parse %s: %w", path, err)
	}

	result := Wrap(entries, opts)

	if opts.DryRun {
		return result, nil
	}

	if err := writeEnvFile(path, result.Wrapped); err != nil {
		return WrapResult{}, fmt.Errorf("wrap: failed to write %s: %w", path, err)
	}

	return result, nil
}
