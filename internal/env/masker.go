package env

import "strings"

// MaskedValue is the string used to replace secret values.
const MaskedValue = "********"

// MaskSecrets returns a copy of the given env map with secret values masked.
func MaskSecrets(vars map[string]string) map[string]string {
	masked := make(map[string]string, len(vars))
	for k, v := range vars {
		if isSecret(k) {
			masked[k] = MaskedValue
		} else {
			masked[k] = v
		}
	}
	return masked
}

// MaskLine masks the value portion of a KEY=VALUE line if the key is a secret.
func MaskLine(line string) string {
	line = strings.TrimSpace(line)
	if line == "" || strings.HasPrefix(line, "#") {
		return line
	}
	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return line
	}
	key := strings.TrimSpace(parts[0])
	if isSecret(key) {
		return key + "=" + MaskedValue
	}
	return line
}
