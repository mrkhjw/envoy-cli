package env

import (
	"strings"
	"testing"
)

func TestFlatten_NoOptions(t *testing.T) {
	env := map[string]string{
		"FOO": "bar",
		"BAZ": "qux",
	}
	res := Flatten(env, FlattenOptions{})
	if len(res.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(res.Entries))
	}
	if res.Renamed != 0 {
		t.Errorf("expected 0 renamed, got %d", res.Renamed)
	}
}

func TestFlatten_AddsPrefix(t *testing.T) {
	env := map[string]string{
		"HOST": "localhost",
		"PORT": "5432",
	}
	res := Flatten(env, FlattenOptions{Prefix: "APP", Separator: "_"})
	for _, e := range res.Entries {
		if !strings.HasPrefix(e.Key, "APP_") {
			t.Errorf("expected key to have prefix APP_, got %s", e.Key)
		}
	}
	if res.Renamed != 2 {
		t.Errorf("expected 2 renamed, got %d", res.Renamed)
	}
}

func TestFlatten_SkipsPrefixIfAlreadyPresent(t *testing.T) {
	env := map[string]string{
		"APP_HOST": "localhost",
	}
	res := Flatten(env, FlattenOptions{Prefix: "APP", Separator: "_"})
	if res.Renamed != 0 {
		t.Errorf("expected 0 renamed (already prefixed), got %d", res.Renamed)
	}
	if res.Entries[0].Key != "APP_HOST" {
		t.Errorf("expected APP_HOST, got %s", res.Entries[0].Key)
	}
}

func TestFlatten_Uppercase(t *testing.T) {
	env := map[string]string{
		"db_host": "localhost",
		"db_port": "5432",
	}
	res := Flatten(env, FlattenOptions{Uppercase: true})
	for _, e := range res.Entries {
		if e.Key != strings.ToUpper(e.Key) {
			t.Errorf("expected uppercase key, got %s", e.Key)
		}
	}
}

func TestFlattenResult_Format_MasksSecrets(t *testing.T) {
	env := map[string]string{
		"API_SECRET": "supersecret",
		"HOST":       "localhost",
	}
	res := Flatten(env, FlattenOptions{})
	out := res.Format()
	if strings.Contains(out, "supersecret") {
		t.Error("expected secret to be masked in Format output")
	}
	if !strings.Contains(out, "localhost") {
		t.Error("expected non-secret value to appear in Format output")
	}
}
