package env

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
)

const (
	charsetAlpha   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	charsetNumeric = "0123456789"
	charsetSpecial = "!@#$%^&*"
	charsetAll     = charsetAlpha + charsetNumeric + charsetSpecial
)

// GenerateOptions controls how values are generated.
type GenerateOptions struct {
	Length  int
	Format  string // "hex", "alphanumeric", "full"
	Keys    []string
	DryRun  bool
}

// GenerateResult holds the outcome of a Generate call.
type GenerateResult struct {
	Generated map[string]string
	Skipped   []string
}

func (r GenerateResult) Format() string {
	var sb strings.Builder
	for k, v := range r.Generated {
		sb.WriteString(fmt.Sprintf("  generated: %s=%s\n", k, v))
	}
	for _, k := range r.Skipped {
		sb.WriteString(fmt.Sprintf("  skipped:   %s\n", k))
	}
	return sb.String()
}

// Generate creates random values for the specified keys in entries.
func Generate(entries []Entry, opts GenerateOptions) ([]Entry, GenerateResult, error) {
	if opts.Length <= 0 {
		opts.Length = 32
	}

	target := map[string]bool{}
	for _, k := range opts.Keys {
		target[strings.ToUpper(k)] = true
	}

	result := GenerateResult{Generated: map[string]string{}}
	out := make([]Entry, 0, len(entries))

	for _, e := range entries {
		if len(target) > 0 && !target[strings.ToUpper(e.Key)] {
			out = append(out, e)
			continue
		}
		if !isSecret(e.Key) {
			result.Skipped = append(result.Skipped, e.Key)
			out = append(out, e)
			continue
		}
		val, err := generateValue(opts.Length, opts.Format)
		if err != nil {
			return nil, result, fmt.Errorf("generate value for %s: %w", e.Key, err)
		}
		result.Generated[e.Key] = val
		if !opts.DryRun {
			e.Value = val
		}
		out = append(out, e)
	}
	return out, result, nil
}

func generateValue(length int, format string) (string, error) {
	switch format {
	case "hex", "":
		b := make([]byte, (length+1)/2)
		if _, err := rand.Read(b); err != nil {
			return "", err
		}
		return hex.EncodeToString(b)[:length], nil
	case "alphanumeric":
		return randomString(length, charsetAlpha+charsetNumeric)
	case "full":
		return randomString(length, charsetAll)
	default:
		return "", fmt.Errorf("unknown format: %s", format)
	}
}

func randomString(length int, charset string) (string, error) {
	b := make([]byte, length)
	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		b[i] = charset[n.Int64()]
	}
	return string(b), nil
}
