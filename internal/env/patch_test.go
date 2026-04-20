package env

import (
	"strings"
	"testing"
)

var basePatchEntries = []Entry{
	{Key: "APP_NAME", Value: "myapp"},
	{Key: "DB_PASSWORD", Value: "secret"},
	{Key: "PORT", Value: "8080"},
}

func TestPatch_SetExisting(t *testing.T) {
	ops := []PatchOp{{Op: "set", Key: "PORT", Value: "9090"}}
	out, res := Patch(basePatchEntries, ops, false)
	if len(res.Applied) != 1 {
		t.Fatalf("expected 1 applied, got %d", len(res.Applied))
	}
	for _, e := range out {
		if e.Key == "PORT" && e.Value != "9090" {
			t.Errorf("expected PORT=9090, got %s", e.Value)
		}
	}
}

func TestPatch_SetNew(t *testing.T) {
	ops := []PatchOp{{Op: "set", Key: "NEW_KEY", Value: "val"}}
	out, res := Patch(basePatchEntries, ops, false)
	if len(res.Applied) != 1 {
		t.Fatalf("expected 1 applied")
	}
	found := false
	for _, e := range out {
		if e.Key == "NEW_KEY" {
			found = true
		}
	}
	if !found {
		t.Error("NEW_KEY not added")
	}
}

func TestPatch_Delete(t *testing.T) {
	ops := []PatchOp{{Op: "delete", Key: "PORT"}}
	out, res := Patch(basePatchEntries, ops, false)
	if len(res.Applied) != 1 {
		t.Fatalf("expected 1 applied")
	}
	for _, e := range out {
		if e.Key == "PORT" {
			t.Error("PORT should be deleted")
		}
	}
}

func TestPatch_DeleteMissing(t *testing.T) {
	ops := []PatchOp{{Op: "delete", Key: "MISSING"}}
	_, res := Patch(basePatchEntries, ops, false)
	if len(res.Skipped) != 1 {
		t.Errorf("expected 1 skipped, got %d", len(res.Skipped))
	}
}

func TestPatch_Rename(t *testing.T) {
	ops := []PatchOp{{Op: "rename", Key: "PORT", NewKey: "HTTP_PORT"}}
	out, res := Patch(basePatchEntries, ops, false)
	if len(res.Applied) != 1 {
		t.Fatalf("expected 1 applied")
	}
	for _, e := range out {
		if e.Key == "PORT" {
			t.Error("old key should not exist")
		}
	}
}

func TestPatch_DryRun(t *testing.T) {
	ops := []PatchOp{{Op: "set", Key: "PORT", Value: "9999"}}
	out, res := Patch(basePatchEntries, ops, true)
	if !res.DryRun {
		t.Error("expected dry run")
	}
	for _, e := range out {
		if e.Key == "PORT" && e.Value == "9999" {
			t.Error("dry run should not modify")
		}
	}
}

func TestPatchResult_Format(t *testing.T) {
	res := PatchResult{Applied: []string{"set PORT"}, Skipped: []string{"MISSING"}, DryRun: true}
	out := res.Format()
	if !strings.Contains(out, "dry-run") {
		t.Error("expected dry-run in output")
	}
	if !strings.Contains(out, "set PORT") {
		t.Error("expected applied entry")
	}
}
