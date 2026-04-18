package env

import (
	"strings"
	"testing"
)

func TestConvert_EnvFormat(t *testing.T) {
	entries := map[string]string{"APP_NAME": "envoy"}
	res, err := Convert(entries, FormatEnv, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(res.Output, "APP_NAME=envoy") {
		t.Errorf("expected env format, got: %s", res.Output)
	}
	if res.Count != 1 {
		t.Errorf("expected count 1, got %d", res.Count)
	}
}

func TestConvert_ExportFormat(t *testing.T) {
	entries := map[string]string{"PORT": "8080"}
	res, err := Convert(entries, FormatExport, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(res.Output, "export PORT=8080") {
		t.Errorf("expected export format, got: %s", res.Output)
	}
}

func TestConvert_YAMLFormat(t *testing.T) {
	entries := map[string]string{"DB_HOST": "localhost"}
	res, err := Convert(entries, FormatYAML, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(res.Output, `db_host: "localhost"`) {
		t.Errorf("expected yaml format, got: %s", res.Output)
	}
}

func TestConvert_TOMLFormat(t *testing.T) {
	entries := map[string]string{"DB_PORT": "5432"}
	res, err := Convert(entries, FormatTOML, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(res.Output, `db_port = "5432"`) {
		t.Errorf("expected toml format, got: %s", res.Output)
	}
}

func TestConvert_UnknownFormat(t *testing.T) {
	entries := map[string]string{"KEY": "val"}
	_, err := Convert(entries, ConvertFormat("xml"), false)
	if err == nil {
		t.Error("expected error for unknown format")
	}
}

func TestConvert_MaskSecrets(t *testing.T) {
	entries := map[string]string{"API_SECRET": "supersecret", "APP_NAME": "envoy"}
	res, err := Convert(entries, FormatEnv, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(res.Output, "supersecret") {
		t.Error("expected secret to be masked")
	}
	if !strings.Contains(res.Output, "APP_NAME=envoy") {
		t.Error("expected non-secret to remain unmasked")
	}
}

func TestConvertResult_Format(t *testing.T) {
	r := ConvertResult{Format: FormatYAML, Count: 3}
	out := r.Format()
	if !strings.Contains(out, "3") || !strings.Contains(out, "yaml") {
		t.Errorf("unexpected format output: %s", out)
	}
}
