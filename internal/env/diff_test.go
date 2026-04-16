package env

import (
	"strings"
	"testing"
)

func TestDiff_Added(t *testing.T) {
	base := map[string]string{"FOO": "bar"}
	target := map[string]string{"FOO": "bar", "NEW_KEY": "value"}
	result := Diff(base, target)
	if _, ok := result.Added["NEW_KEY"]; !ok {
		t.Error("expected NEW_KEY to be in Added")
	}
	if len(result.Removed) != 0 || len(result.Changed) != 0 {
		t.Error("expected no removed or changed keys")
	}
}

func TestDiff_Removed(t *testing.T) {
	base := map[string]string{"FOO": "bar", "OLD_KEY": "val"}
	target := map[string]string{"FOO": "bar"}
	result := Diff(base, target)
	if _, ok := result.Removed["OLD_KEY"]; !ok {
		t.Error("expected OLD_KEY to be in Removed")
	}
}

func TestDiff_Changed(t *testing.T) {
	base := map[string]string{"FOO": "old"}
	target := map[string]string{"FOO": "new"}
	result := Diff(base, target)
	v, ok := result.Changed["FOO"]
	if !ok {
		t.Fatal("expected FOO to be in Changed")
	}
	if v[0] != "old" || v[1] != "new" {
		t.Errorf("unexpected changed values: %v", v)
	}
}

func TestDiff_NoChanges(t *testing.T) {
	base := map[string]string{"FOO": "bar"}
	target := map[string]string{"FOO": "bar"}
	result := Diff(base, target)
	if len(result.Added)+len(result.Removed)+len(result.Changed) != 0 {
		t.Error("expected no diff")
	}
}

func TestDiffResult_Format_MasksSecrets(t *testing.T) {
	base := map[string]string{}
	target := map[string]string{"API_SECRET": "supersecret"}
	result := Diff(base, target)
	formatted := result.Format(true)
	if strings.Contains(formatted, "supersecret") {
		t.Error("expected secret value to be masked")
	}
	if !strings.Contains(formatted, "***") {
		t.Error("expected masked placeholder")
	}
}

func TestDiffResult_Format_ShowsValues(t *testing.T) {
	base := map[string]string{}
	target := map[string]string{"APP_NAME": "envoy"}
	result := Diff(base, target)
	formatted := result.Format(false)
	if !strings.Contains(formatted, "envoy") {
		t.Error("expected value to be visible")
	}
}
