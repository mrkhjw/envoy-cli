package env

import (
	"strings"
	"testing"
)

func baseStatsEntries() []Entry {
	return []Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "API_SECRET", Value: "s3cr3t"},
		{Key: "DB_PASSWORD", Value: "pass"},
		{Key: "DEBUG", Value: ""},
		{Key: "APP_NAME", Value: "duplicate"},
		{Comment: true, Raw: "# a comment"},
	}
}

func TestStats_Total(t *testing.T) {
	res := Stats(baseStatsEntries())
	if res.Total != 5 {
		t.Errorf("expected Total=5, got %d", res.Total)
	}
}

func TestStats_Secrets(t *testing.T) {
	res := Stats(baseStatsEntries())
	if res.Secrets != 2 {
		t.Errorf("expected Secrets=2, got %d", res.Secrets)
	}
}

func TestStats_Empty(t *testing.T) {
	res := Stats(baseStatsEntries())
	if res.Empty != 1 {
		t.Errorf("expected Empty=1, got %d", res.Empty)
	}
}

func TestStats_Comments(t *testing.T) {
	res := Stats(baseStatsEntries())
	if res.Comments != 1 {
		t.Errorf("expected Comments=1, got %d", res.Comments)
	}
}

func TestStats_UniqueKeys(t *testing.T) {
	res := Stats(baseStatsEntries())
	if res.Unique != 4 {
		t.Errorf("expected Unique=4, got %d", res.Unique)
	}
}

func TestStatsResult_Format(t *testing.T) {
	res := Stats(baseStatsEntries())
	out := res.Format()
	for _, s := range []string{"Total:", "Secrets:", "Empty:", "Comments:", "Unique Keys:"} {
		if !strings.Contains(out, s) {
			t.Errorf("Format() missing %q in %q", s, out)
		}
	}
}
