package env

import (
	"strings"
	"testing"
)

func baseFreezeEntries() []Entry {
	return []Entry{
		{Key: "APP_NAME", Value: "myapp", RawLine: "APP_NAME=myapp"},
		{Key: "API_SECRET", Value: "s3cr3t", RawLine: "API_SECRET=s3cr3t"},
		{Key: "DEBUG", Value: "true", RawLine: "DEBUG=true"},
		{Key: "ALREADY", Value: "yes", RawLine: "ALREADY=yes #frozen"},
	}
}

func TestFreeze_AllKeys(t *testing.T) {
	entries := baseFreezeEntries()
	res := Freeze(entries, FreezeOption{})

	if len(res.Frozen) != 3 {
		t.Errorf("expected 3 frozen, got %d", len(res.Frozen))
	}
	if len(res.Skipped) != 1 {
		t.Errorf("expected 1 skipped, got %d", len(res.Skipped))
	}
	if res.Skipped[0] != "ALREADY" {
		t.Errorf("expected ALREADY skipped, got %s", res.Skipped[0])
	}
}

func TestFreeze_SpecificKeys(t *testing.T) {
	entries := baseFreezeEntries()
	res := Freeze(entries, FreezeOption{Keys: []string{"APP_NAME"}})

	if len(res.Frozen) != 1 || res.Frozen[0] != "APP_NAME" {
		t.Errorf("expected APP_NAME frozen, got %v", res.Frozen)
	}
}

func TestFreeze_DryRunDoesNotTag(t *testing.T) {
	entries := baseFreezeEntries()
	res := Freeze(entries, FreezeOption{DryRun: true})

	for _, e := range res.Entries {
		if e.Key == "APP_NAME" && strings.HasSuffix(e.RawLine, "#frozen") {
			t.Error("dry run should not modify RawLine")
		}
	}
	if len(res.Frozen) != 3 {
		t.Errorf("expected 3 reported frozen in dry run, got %d", len(res.Frozen))
	}
}

func TestFreeze_SkipsAlreadyFrozen(t *testing.T) {
	entries := baseFreezeEntries()
	res := Freeze(entries, FreezeOption{Keys: []string{"ALREADY"}})

	if len(res.Frozen) != 0 {
		t.Errorf("expected 0 frozen, got %d", len(res.Frozen))
	}
	if len(res.Skipped) != 1 {
		t.Errorf("expected 1 skipped, got %d", len(res.Skipped))
	}
}

func TestFreezeResult_Format_MasksSecrets(t *testing.T) {
	res := FreezeResult{
		Frozen:  []string{"API_SECRET"},
		Skipped: []string{},
		Entries: []Entry{
			{Key: "API_SECRET", Value: "s3cr3t", RawLine: "API_SECRET=s3cr3t #frozen"},
		},
	}
	out := res.Format(true)
	if strings.Contains(out, "s3cr3t") {
		t.Error("expected secret value to be masked")
	}
	if !strings.Contains(out, "***") {
		t.Error("expected *** placeholder in output")
	}
}
