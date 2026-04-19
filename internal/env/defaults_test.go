package env

import (
	"strings"
	"testing"
)

func baseDefaultEntries() []Entry {
	return []Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "DB_PASS", Value: ""},
		{Key: "PORT", Value: "8080"},
	}
}

func TestDefaults_FillsEmptyValues(t *testing.T) {
	entries := baseDefaultEntries()
	defaults := map[string]string{"DB_PASS": "secret123"}
	res := Defaults(entries, defaults, false)
	for _, e := range res.Entries {
		if e.Key == "DB_PASS" && e.Value != "secret123" {
			t.Errorf("expected DB_PASS=secret123, got %s", e.Value)
		}
	}
	if len(res.Applied) != 1 || res.Applied[0] != "DB_PASS" {
		t.Errorf("expected Applied=[DB_PASS], got %v", res.Applied)
	}
}

func TestDefaults_SkipsNonEmptyWithoutOverwrite(t *testing.T) {
	entries := baseDefaultEntries()
	defaults := map[string]string{"APP_NAME": "other"}
	res := Defaults(entries, defaults, false)
	if len(res.Skipped) != 1 || res.Skipped[0] != "APP_NAME" {
		t.Errorf("expected Skipped=[APP_NAME], got %v", res.Skipped)
	}
	for _, e := range res.Entries {
		if e.Key == "APP_NAME" && e.Value != "myapp" {
			t.Errorf("expected APP_NAME unchanged, got %s", e.Value)
		}
	}
}

func TestDefaults_OverwriteReplacesExisting(t *testing.T) {
	entries := baseDefaultEntries()
	defaults := map[string]string{"PORT": "9090"}
	res := Defaults(entries, defaults, true)
	for _, e := range res.Entries {
		if e.Key == "PORT" && e.Value != "9090" {
			t.Errorf("expected PORT=9090, got %s", e.Value)
		}
	}
	if len(res.Applied) != 1 {
		t.Errorf("expected 1 applied, got %d", len(res.Applied))
	}
}

func TestDefaults_AddsNewKey(t *testing.T) {
	entries := baseDefaultEntries()
	defaults := map[string]string{"NEW_KEY": "newval"}
	res := Defaults(entries, defaults, false)
	found := false
	for _, e := range res.Entries {
		if e.Key == "NEW_KEY" && e.Value == "newval" {
			found = true
		}
	}
	if !found {
		t.Error("expected NEW_KEY to be added")
	}
}

func TestDefaultsResult_Format(t *testing.T) {
	res := DefaultsResult{
		Applied: []string{"DB_PASS"},
		Skipped: []string{"APP_NAME"},
	}
	out := res.Format()
	if !strings.Contains(out, "applied=1") || !strings.Contains(out, "skipped=1") {
		t.Errorf("unexpected format output: %s", out)
	}
}
