package env

import "fmt"

// RotateFile reads an env file, rotates secret keys, and writes the result back.
func RotateFile(path string, opts RotateOptions) (RotateResult, error) {
	env, err := ParseFile(path)
	if err != nil {
		return RotateResult{}, fmt.Errorf("rotate: parse %s: %w", path, err)
	}

	rotated, result := Rotate(env, opts)

	if !opts.DryRun {
		if err := writeEnvFile(path, rotated); err != nil {
			return result, fmt.Errorf("rotate: write %s: %w", path, err)
		}
	}
	return result, nil
}
