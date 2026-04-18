package env

import (
	"fmt"
	"strings"
	"time"
)

// AuditEntry represents a single audit log record.
type AuditEntry struct {
	Timestamp time.Time
	Action    string
	Key       string
	File      string
	Masked    bool
}

// AuditLog holds a list of audit entries.
type AuditLog struct {
	Entries []AuditEntry
}

// Record adds a new entry to the audit log.
func (a *AuditLog) Record(action, key, file string) {
	a.Entries = append(a.Entries, AuditEntry{
		Timestamp: time.Now(),
		Action:    action,
		Key:       key,
		File:      file,
		Masked:    isSecret(key),
	})
}

// Format returns a human-readable audit log string.
func (a *AuditLog) Format() string {
	if len(a.Entries) == 0 {
		return "No audit entries."
	}
	var sb strings.Builder
	for _, e := range a.Entries {
		key := e.Key
		if e.Masked {
			key = key + " (secret)"
		}
		sb.WriteString(fmt.Sprintf("[%s] %s | key=%s | file=%s\n",
			e.Timestamp.Format(time.RFC3339),
			e.Action,
			key,
			e.File,
		))
	}
	return strings.TrimRight(sb.String(), "\n")
}

// AuditMap generates an audit log from a map of env vars for a given action.
func AuditMap(vars map[string]string, action, file string) *AuditLog {
	log := &AuditLog{}
	for k := range vars {
		log.Record(action, k, file)
	}
	return log
}
