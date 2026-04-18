package env

import (
	"strings"
	"testing"
)

func TestAuditLog_Record(t *testing.T) {
	log := &AuditLog{}
	log.Record("read", "APP_NAME", ".env")
	log.Record("read", "DB_PASSWORD", ".env")

	if len(log.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(log.Entries))
	}
	if log.Entries[0].Key != "APP_NAME" {
		t.Errorf("expected APP_NAME, got %s", log.Entries[0].Key)
	}
	if !log.Entries[1].Masked {
		t.Errorf("expected DB_PASSWORD to be masked")
	}
	if log.Entries[0].Masked {
		t.Errorf("expected APP_NAME to not be masked")
	}
}

func TestAuditLog_Format_Empty(t *testing.T) {
	log := &AuditLog{}
	out := log.Format()
	if out != "No audit entries." {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestAuditLog_Format_ContainsEntries(t *testing.T) {
	log := &AuditLog{}
	log.Record("write", "API_KEY", "prod.env")
	out := log.Format()

	if !strings.Contains(out, "write") {
		t.Errorf("expected 'write' in output")
	}
	if !strings.Contains(out, "API_KEY") {
		t.Errorf("expected 'API_KEY' in output")
	}
	if !strings.Contains(out, "(secret)") {
		t.Errorf("expected '(secret)' label for API_KEY")
	}
	if !strings.Contains(out, "prod.env") {
		t.Errorf("expected filename in output")
	}
}

func TestAuditMap(t *testing.T) {
	vars := map[string]string{
		"APP_ENV":   "production",
		"DB_SECRET": "hunter2",
	}
	log := AuditMap(vars, "export", ".env")
	if len(log.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(log.Entries))
	}
	for _, e := range log.Entries {
		if e.Action != "export" {
			t.Errorf("expected action 'export', got %s", e.Action)
		}
		if e.File != ".env" {
			t.Errorf("expected file '.env', got %s", e.File)
		}
	}
}
