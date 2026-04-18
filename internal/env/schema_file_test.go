package env

import (
	"os"
	"testing"
)

func writeTempSchema(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "schema-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(content)
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestLoadSchema_Basic(t *testing.T) {
	path := writeTempSchema(t, "APP_HOST required\nLOG_LEVEL optional default=info\n")
	schema, err := LoadSchema(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(schema) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(schema))
	}
	if !schema[0].Required {
		t.Error("expected APP_HOST to be required")
	}
	if schema[1].Default != "info" {
		t.Errorf("expected default info, got %s", schema[1].Default)
	}
}

func TestLoadSchema_IgnoresComments(t *testing.T) {
	path := writeTempSchema(t, "# comment\nDB_URL required\n")
	schema, err := LoadSchema(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(schema) != 1 || schema[0].Key != "DB_URL" {
		t.Errorf("unexpected schema: %v", schema)
	}
}

func TestLoadSchema_NotFound(t *testing.T) {
	_, err := LoadSchema("/nonexistent/schema.txt")
	if err == nil {
		t.Error("expected error for missing file")
	}
}
