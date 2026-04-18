package env

import (
	"os"
	"testing"
)

func TestInject_SetsEnvVars(t *testing.T) {
	os.Unsetenv("INJECT_FOO")
	os.Unsetenv("INJECT_BAR")
	vars := map[string]string{"INJECT_FOO": "hello", "INJECT_BAR": "world"}
	result := Inject(vars, false)
	if len(result.Injected) != 2 {
		t.Errorf("expected 2 injected, got %d", len(result.Injected))
	}
	if os.Getenv("INJECT_FOO") != "hello" {
		t.Error("expected INJECT_FOO=hello")
	}
}

func TestInject_SkipsExistingWithoutOverwrite(t *testing.T) {
	os.Setenv("INJECT_EXISTING", "original")
	defer os.Unsetenv("INJECT_EXISTING")
	vars := map[string]string{"INJECT_EXISTING": "new"}
	result := Inject(vars, false)
	if len(result.Skipped) != 1 {
		t.Errorf("expected 1 skipped, got %d", len(result.Skipped))
	}
	if os.Getenv("INJECT_EXISTING") != "original" {
		t.Error("expected value to remain original")
	}
}

func TestInject_OverwritesExisting(t *testing.T) {
	os.Setenv("INJECT_OVR", "old")
	defer os.Unsetenv("INJECT_OVR")
	vars := map[string]string{"INJECT_OVR": "new"}
	result := Inject(vars, true)
	if len(result.Injected) != 1 {
		t.Errorf("expected 1 injected, got %d", len(result.Injected))
	}
	if os.Getenv("INJECT_OVR") != "new" {
		t.Error("expected value to be new")
	}
}

func TestInjectFile_Valid(t *testing.T) {
	os.Unsetenv("HELLO")
	f := writeTempEnv(t, "HELLO=world\n")
	result, err := InjectFile(f, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Injected) != 1 {
		t.Errorf("expected 1 injected, got %d", len(result.Injected))
	}
	if os.Getenv("HELLO") != "world" {
		t.Error("expected HELLO=world")
	}
}

func TestInjectFile_NotFound(t *testing.T) {
	_, err := InjectFile("/nonexistent/.env", false)
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestInjectResult_Format_Empty(t *testing.T) {
	r := InjectResult{}
	out := r.Format()
	if out != "Nothing to inject.\n" {
		t.Errorf("unexpected output: %q", out)
	}
}
