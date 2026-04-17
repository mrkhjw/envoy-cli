package env

import (
	"strings"
	"testing"
)

func TestExport_ShellFormat(t *testing.T) {
	vars := map[string]string{"APP_NAME": "myapp"}
	out, err := Export(vars, FormatShell, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "export APP_NAME=") {
		t.Errorf("expected shell export, got: %s", out)
	}
}

func TestExport_DockerFormat(t *testing.T) {
	vars := map[string]string{"PORT": "8080"}
	out, err := Export(vars, FormatDocker, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "-e PORT=8080") {
		t.Errorf("expected docker flag, got: %s", out)
	}
}

func TestExport_JSONFormat(t *testing.T) {
	vars := map[string]string{"KEY": "val"}
	out, err := Export(vars, FormatJSON, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "\"KEY\"") {
		t.Errorf("expected JSON key, got: %s", out)
	}
}

func TestExport_UnknownFormat(t *testing.T) {
	vars := map[string]string{"K": "v"}
	_, err := Export(vars, ExportFormat("xml"), false)
	if err == nil {
		t.Error("expected error for unknown format")
	}
}

func TestExport_MaskSecrets(t *testing.T) {
	vars := map[string]string{"SECRET_KEY": "supersecret", "APP": "myapp"}
	out, err := Export(vars, FormatShell, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(out, "supersecret") {
		t.Errorf("expected secret to be masked, got: %s", out)
	}
	if !strings.Contains(out, "myapp") {
		t.Errorf("expected non-secret value to remain, got: %s", out)
	}
}

func TestExportFile_Valid(t *testing.T) {
	f := writeTempEnv(t, "APP=myapp\nPORT=9000\n")
	out, err := ExportFile(f, FormatShell, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "APP") {
		t.Errorf("expected APP in output, got: %s", out)
	}
}

func TestExportFile_NotFound(t *testing.T) {
	_, err := ExportFile("/nonexistent/.env", FormatShell, false)
	if err == nil {
		t.Error("expected error for missing file")
	}
}
