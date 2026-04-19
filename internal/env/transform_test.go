package env

import (
	"strings"
	"testing"
)

var baseTransformEntries = []Entry{
	{Key: "app_name", Value: "myapp"},
	{Key: "DB_HOST", Value: "  localhost  "},
	{Key: "api_secret", Value: "abc123"},
}

func TestTransform_UppercaseKeys(t *testing.T) {
	res := Transform(baseTransformEntries, TransformOpts{UppercaseKeys: true})
	for _, e := range res.Entries {
		if e.Key != strings.ToUpper(e.Key) {
			t.Errorf("expected uppercase key, got %s", e.Key)
		}
	}
	if res.Modified != 2 { // DB_HOST already uppercase
		t.Errorf("expected 2 modified, got %d", res.Modified)
	}
}

func TestTransform_LowercaseKeys(t *testing.T) {
	res := Transform(baseTransformEntries, TransformOpts{LowercaseKeys: true})
	for _, e := range res.Entries {
		if e.Key != strings.ToLower(e.Key) {
			t.Errorf("expected lowercase key, got %s", e.Key)
		}
	}
}

func TestTransform_TrimValues(t *testing.T) {
	res := Transform(baseTransformEntries, TransformOpts{TrimValues: true})
	for _, e := range res.Entries {
		if e.Value != strings.TrimSpace(e.Value) {
			t.Errorf("expected trimmed value for %s", e.Key)
		}
	}
	if res.Modified != 1 {
		t.Errorf("expected 1 modified, got %d", res.Modified)
	}
}

func TestTransform_SpecificKeys(t *testing.T) {
	res := Transform(baseTransformEntries, TransformOpts{
		UppercaseKeys: true,
		Keys: []string{"app_name"},
	})
	if res.Entries[0].Key != "APP_NAME" {
		t.Errorf("expected APP_NAME, got %s", res.Entries[0].Key)
	}
	if res.Entries[1].Key != "DB_HOST" {
		t.Errorf("expected DB_HOST unchanged, got %s", res.Entries[1].Key)
	}
	if res.Modified != 1 {
		t.Errorf("expected 1 modified, got %d", res.Modified)
	}
}

func TestTransform_Format_MasksSecrets(t *testing.T) {
	res := Transform(baseTransformEntries, TransformOpts{})
	out := res.Format(true)
	if strings.Contains(out, "abc123") {
		t.Error("expected secret to be masked")
	}
}

func TestTransform_NoOpts(t *testing.T) {
	res := Transform(baseTransformEntries, TransformOpts{})
	if res.Modified != 0 {
		t.Errorf("expected 0 modified, got %d", res.Modified)
	}
}
