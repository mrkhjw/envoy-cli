package env

import (
	"strings"
	"testing"
)

func baseTagEntries() []Entry {
	return []Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "DB_PASSWORD", Value: "secret123"},
		{Key: "PORT", Value: "8080"},
	}
}

func TestTag_AllKeys(t *testing.T) {
	entries := baseTagEntries()
	result := Tag(entries, TagOptions{Tags: []string{"prod", "v2"}})
	if len(result.Tagged) != 3 {
		t.Fatalf("expected 3 tagged, got %d", len(result.Tagged))
	}
	if len(result.Untagged) != 0 {
		t.Fatalf("expected 0 untagged, got %d", len(result.Untagged))
	}
}

func TestTag_SpecificKeys(t *testing.T) {
	entries := baseTagEntries()
	result := Tag(entries, TagOptions{Tags: []string{"infra"}, Keys: []string{"PORT"}})
	if len(result.Tagged) != 1 {
		t.Fatalf("expected 1 tagged, got %d", len(result.Tagged))
	}
	if result.Tagged[0].Key != "PORT" {
		t.Errorf("expected PORT, got %s", result.Tagged[0].Key)
	}
	if len(result.Untagged) != 2 {
		t.Fatalf("expected 2 untagged, got %d", len(result.Untagged))
	}
}

func TestTag_MasksSecrets(t *testing.T) {
	entries := baseTagEntries()
	result := Tag(entries, TagOptions{Tags: []string{"secure"}, MaskSecrets: true})
	for _, te := range result.Tagged {
		if isSecret(te.Key) && te.Value != "****" {
			t.Errorf("expected masked value for %s, got %s", te.Key, te.Value)
		}
	}
}

func TestTag_EmptyTags(t *testing.T) {
	entries := baseTagEntries()
	result := Tag(entries, TagOptions{Tags: []string{}})
	if len(result.Tagged) != 3 {
		t.Fatalf("expected 3 tagged, got %d", len(result.Tagged))
	}
	for _, te := range result.Tagged {
		if len(te.Tags) != 0 {
			t.Errorf("expected no tags for %s", te.Key)
		}
	}
}

func TestTagResult_Format(t *testing.T) {
	entries := baseTagEntries()
	result := Tag(entries, TagOptions{Tags: []string{"prod"}, Keys: []string{"APP_NAME"}})
	out := result.Format()
	if !strings.Contains(out, "APP_NAME") {
		t.Error("expected APP_NAME in output")
	}
	if !strings.Contains(out, "tags:prod") {
		t.Error("expected tags:prod in output")
	}
	if !strings.Contains(out, "untagged") {
		t.Error("expected untagged section in output")
	}
}
