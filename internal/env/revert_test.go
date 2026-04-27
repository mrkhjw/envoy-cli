package env

import (
	"strings"
	"testing"
)

var baseRevertCurrent = []Entry{
	{Key: "APP_NAME", Value: "new-app"},
	{Key: "DB_PASSWORD", Value: "newpass"},
	{Key: "DEBUG", Value: "true"},
	{Key: "VERSION", Value: "2.0"},
}

var baseRevertBaseline = []Entry{
	{Key: "APP_NAME", Value: "original-app"},
	{Key: "DB_PASSWORD", Value: "oldpass"},
	{Key: "DEBUG", Value: "false"},
}

func TestRevert_AllKeys(t *testing.T) {
	out, res := Revert(baseRevertCurrent, baseRevertBaseline, RevertOptions{Overwrite: true})
	if len(res.Reverted) != 3 {
		t.Fatalf("expected 3 reverted, got %d", len(res.Reverted))
	}
	if len(res.Skipped) != 1 {
		t.Fatalf("expected 1 skipped (VERSION not in baseline), got %d", len(res.Skipped))
	}
	m := make(map[string]string)
	for _, e := range out {
		m[e.Key] = e.Value
	}
	if m["APP_NAME"] != "original-app" {
		t.Errorf("expected original-app, got %s", m["APP_NAME"])
	}
	if m["VERSION"] != "2.0" {
		t.Errorf("VERSION should be unchanged, got %s", m["VERSION"])
	}
}

func TestRevert_SpecificKeys(t *testing.T) {
	_, res := Revert(baseRevertCurrent, baseRevertBaseline, RevertOptions{
		Keys:      []string{"APP_NAME"},
		Overwrite: true,
	})
	if len(res.Reverted) != 1 || res.Reverted[0] != "APP_NAME" {
		t.Errorf("expected only APP_NAME reverted, got %v", res.Reverted)
	}
}

func TestRevert_DryRunDoesNotChange(t *testing.T) {
	out, res := Revert(baseRevertCurrent, baseRevertBaseline, RevertOptions{
		DryRun:    true,
		Overwrite: true,
	})
	if !res.DryRun {
		t.Error("expected DryRun to be true")
	}
	for _, e := range out {
		if e.Key == "APP_NAME" && e.Value != "new-app" {
			t.Errorf("dry-run should not change value, got %s", e.Value)
		}
	}
}

func TestRevert_SkipsWhenValueUnchanged(t *testing.T) {
	current := []Entry{
		{Key: "DEBUG", Value: "false"},
	}
	_, res := Revert(current, baseRevertBaseline, RevertOptions{})
	if len(res.Skipped) != 1 {
		t.Errorf("expected 1 skipped (same value), got %d", len(res.Skipped))
	}
	if len(res.Reverted) != 0 {
		t.Errorf("expected 0 reverted, got %d", len(res.Reverted))
	}
}

func TestRevert_Format(t *testing.T) {
	_, res := Revert(baseRevertCurrent, baseRevertBaseline, RevertOptions{Overwrite: true})
	f := res.Format()
	if !strings.Contains(f, "Reverted:") {
		t.Error("expected Format to contain 'Reverted:'")
	}
	if !strings.Contains(f, "APP_NAME") {
		t.Error("expected Format to mention APP_NAME")
	}
}
