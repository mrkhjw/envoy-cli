package env

import (
	"strings"
	"testing"
)

func baseV1Entries() []Entry {
	return []Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "API_SECRET", Value: "old-secret"},
	}
}

func baseV2Entries() []Entry {
	return []Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "DB_HOST", Value: "db.prod"},
		{Key: "API_SECRET", Value: "new-secret"},
		{Key: "NEW_KEY", Value: "hello"},
	}
}

func TestVersionDiff_Added(t *testing.T) {
	r := VersionDiff("v1", baseV1Entries(), "v2", baseV2Entries())
	if len(r.Added) != 1 || r.Added[0].Key != "NEW_KEY" {
		t.Errorf("expected NEW_KEY to be added, got %+v", r.Added)
	}
}

func TestVersionDiff_Removed(t *testing.T) {
	v2 := []Entry{
		{Key: "APP_NAME", Value: "myapp"},
	}
	r := VersionDiff("v1", baseV1Entries(), "v2", v2)
	if len(r.Removed) != 2 {
		t.Errorf("expected 2 removed, got %d", len(r.Removed))
	}
}

func TestVersionDiff_Changed(t *testing.T) {
	r := VersionDiff("v1", baseV1Entries(), "v2", baseV2Entries())
	changedKeys := map[string]bool{}
	for _, e := range r.Changed {
		changedKeys[e.Key] = true
	}
	if !changedKeys["DB_HOST"] || !changedKeys["API_SECRET"] {
		t.Errorf("expected DB_HOST and API_SECRET to be changed, got %+v", r.Changed)
	}
}

func TestVersionDiff_NoChanges(t *testing.T) {
	r := VersionDiff("v1", baseV1Entries(), "v2", baseV1Entries())
	if len(r.Added)+len(r.Removed)+len(r.Changed) != 0 {
		t.Errorf("expected no changes")
	}
}

func TestVersionDiff_Format_MasksSecrets(t *testing.T) {
	r := VersionDiff("v1", baseV1Entries(), "v2", baseV2Entries())
	out := r.Format(true)
	if strings.Contains(out, "new-secret") {
		t.Errorf("expected secret to be masked, got: %s", out)
	}
	if !strings.Contains(out, "***") {
		t.Errorf("expected *** placeholder in output")
	}
}

func TestVersionDiff_Format_Header(t *testing.T) {
	r := VersionDiff("v1", baseV1Entries(), "v2", baseV2Entries())
	out := r.Format(false)
	if !strings.HasPrefix(out, "diff v1..v2") {
		t.Errorf("expected diff header, got: %s", out)
	}
}
