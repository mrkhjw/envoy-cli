package env

import (
	"os"
	"testing"
)

func writeTempEnvForValidate(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "envoy-validate-*.env")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestValidateFile_Valid(t *testing.T) {
	path := writeTempEnvForValidate(t, "APP=hello\nPORT=3000\n")
	result, err := ValidateFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Valid() {
		t.Errorf("expected valid result, got: %s", result.Summary())
	}
}

func TestValidateFile_Duplicate(t *testing.T) {
	path := writeTempEnvForValidate(t, "KEY=a\nKEY=b\n")
	result, err := ValidateFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Valid() {
		t.Error("expected validation failure for duplicate key")
	}
}

func TestValidateFile_NotFound(t *testing.T) {
	_, err := ValidateFile("/nonexistent/path/.env")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestValidateFile_EmptyFile(t *testing.T) {
	path := writeTempEnvForValidate(t, "")
	result, err := ValidateFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Valid() {
		t.Errorf("expected empty file to be valid, got: %s", result.Summary())
	}
}
