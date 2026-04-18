package env

import (
	"strings"
	"testing"
)

func TestCompare_IdenticalMaps(t *testing.T) {
	env1 := map[string]string{"FOO": "bar", "BAZ": "qux"}
	env2 := map[string]string{"FOO": "bar", "BAZ": "qux"}

	r := Compare("a.env", "b.env", env1, env2)

	if len(r.InBoth) != 2 {
		t.Errorf("expected 2 identical keys, got %d", len(r.InBoth))
	}
	if len(r.Conflict) != 0 || len(r.OnlyIn1) != 0 || len(r.OnlyIn2) != 0 {
		t.Error("expected no conflicts or unique keys")
	}
}

func TestCompare_OnlyIn1(t *testing.T) {
	env1 := map[string]string{"FOO": "bar", "ONLY1": "val"}
	env2 := map[string]string{"FOO": "bar"}

	r := Compare("a.env", "b.env", env1, env2)

	if len(r.OnlyIn1) != 1 {
		t.Errorf("expected 1 key only in file1, got %d", len(r.OnlyIn1))
	}
	if _, ok := r.OnlyIn1["ONLY1"]; !ok {
		t.Error("expected ONLY1 in OnlyIn1")
	}
}

func TestCompare_OnlyIn2(t *testing.T) {
	env1 := map[string]string{"FOO": "bar"}
	env2 := map[string]string{"FOO": "bar", "ONLY2": "val"}

	r := Compare("a.env", "b.env", env1, env2)

	if len(r.OnlyIn2) != 1 {
		t.Errorf("expected 1 key only in file2, got %d", len(r.OnlyIn2))
	}
}

func TestCompare_Conflicts(t *testing.T) {
	env1 := map[string]string{"FOO": "bar"}
	env2 := map[string]string{"FOO": "baz"}

	r := Compare("a.env", "b.env", env1, env2)

	if len(r.Conflict) != 1 {
		t.Errorf("expected 1 conflict, got %d", len(r.Conflict))
	}
	if r.Conflict["FOO"] != [2]string{"bar", "baz"} {
		t.Errorf("unexpected conflict values: %v", r.Conflict["FOO"])
	}
}

func TestCompare_Summary_MasksSecrets(t *testing.T) {
	env1 := map[string]string{"SECRET_KEY": "abc123"}
	env2 := map[string]string{"SECRET_KEY": "xyz789"}

	r := Compare("a.env", "b.env", env1, env2)
	summary := r.Summary(true)

	if strings.Contains(summary, "abc123") || strings.Contains(summary, "xyz789") {
		t.Error("expected secret values to be masked in summary")
	}
	if !strings.Contains(summary, "***") {
		t.Error("expected masked placeholder in summary")
	}
}
