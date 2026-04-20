package env

import (
	"testing"
)

func baseEnvMapEntries() []Entry {
	return []Entry{
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "API_KEY", Value: "secret123"},
		{Key: "PORT", Value: "8080"},
		{IsComment: true, Raw: "# a comment"},
		{Key: "", Value: "", Raw: ""},
	}
}

func TestBuildEnvMap_Total(t *testing.T) {
	entries := baseEnvMapEntries()
	result := BuildEnvMap(entries)
	if result.Total != 3 {
		t.Errorf("expected 3 non-comment keys, got %d", result.Total)
	}
}

func TestBuildEnvMap_KeysSorted(t *testing.T) {
	entries := baseEnvMapEntries()
	result := BuildEnvMap(entries)
	expected := []string{"API_KEY", "DB_HOST", "PORT"}
	for i, k := range result.Keys {
		if k != expected[i] {
			t.Errorf("key[%d]: expected %q, got %q", i, expected[i], k)
		}
	}
}

func TestToMap_BasicConversion(t *testing.T) {
	entries := baseEnvMapEntries()
	m := ToMap(entries)
	if m["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", m["DB_HOST"])
	}
	if m["PORT"] != "8080" {
		t.Errorf("expected PORT=8080, got %q", m["PORT"])
	}
	if _, ok := m[""]; ok {
		t.Error("empty key should not be included in map")
	}
}

func TestToMap_ExcludesComments(t *testing.T) {
	entries := baseEnvMapEntries()
	m := ToMap(entries)
	if len(m) != 3 {
		t.Errorf("expected 3 entries, got %d", len(m))
	}
}

func TestFromMap_ReturnsSortedEntries(t *testing.T) {
	m := map[string]string{
		"ZEBRA": "last",
		"ALPHA": "first",
		"MIDDLE": "mid",
	}
	entries := FromMap(m)
	if entries[0].Key != "ALPHA" {
		t.Errorf("expected ALPHA first, got %q", entries[0].Key)
	}
	if entries[2].Key != "ZEBRA" {
		t.Errorf("expected ZEBRA last, got %q", entries[2].Key)
	}
}

func TestFromMap_EmptyMap(t *testing.T) {
	entries := FromMap(map[string]string{})
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}
