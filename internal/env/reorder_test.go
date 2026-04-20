package env

import (
	"os"
	"testing"
)

func writeTempReorderEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "reorder-*.env")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(content)
	f.Close()
	return f.Name()
}

func TestReorder_PinsKeysFirst(t *testing.T) {
	entries := []Entry{
		{Key: "ZEBRA", Value: "z"},
		{Key: "ALPHA", Value: "a"},
		{Key: "BETA", Value: "b"},
	}
	res := Reorder(entries, ReorderOptions{Keys: []string{"BETA", "ALPHA"}})
	if res.Entries[0].Key != "BETA" || res.Entries[1].Key != "ALPHA" {
		t.Errorf("expected BETA then ALPHA first, got %v %v", res.Entries[0].Key, res.Entries[1].Key)
	}
	if res.Moved != 2 {
		t.Errorf("expected 2 moved, got %d", res.Moved)
	}
}

func TestReorder_RestSortedAlpha(t *testing.T) {
	entries := []Entry{
		{Key: "ZEBRA", Value: "z"},
		{Key: "MANGO", Value: "m"},
		{Key: "ALPHA", Value: "a"},
	}
	res := Reorder(entries, ReorderOptions{Keys: []string{"ZEBRA"}})
	if res.Entries[1].Key != "ALPHA" || res.Entries[2].Key != "MANGO" {
		t.Errorf("rest not sorted: %v %v", res.Entries[1].Key, res.Entries[2].Key)
	}
}

func TestReorder_DryRunDoesNotChangeLabel(t *testing.T) {
	entries := []Entry{{Key: "A", Value: "1"}, {Key: "B", Value: "2"}}
	res := Reorder(entries, ReorderOptions{Keys: []string{"B"}, DryRun: true})
	if !res.DryRun {
		t.Error("expected DryRun=true")
	}
	f := res.Format()
	if f[:9] != "[dry-run]" {
		t.Errorf("expected dry-run prefix, got %q", f)
	}
}

func TestReorder_MissingKeyIgnored(t *testing.T) {
	entries := []Entry{{Key: "A", Value: "1"}}
	res := Reorder(entries, ReorderOptions{Keys: []string{"MISSING", "A"}})
	if res.Entries[0].Key != "A" {
		t.Errorf("expected A first, got %v", res.Entries[0].Key)
	}
	if res.Moved != 1 {
		t.Errorf("expected 1 moved, got %d", res.Moved)
	}
}

func TestReorderFile_WritesOutput(t *testing.T) {
	src := writeTempReorderEnv(t, "ZEBRA=z\nALPHA=a\nBETA=b\n")
	dst := src + ".out"
	t.Cleanup(func() { os.Remove(src); os.Remove(dst) })

	_, err := ReorderFile(src, dst, ReorderOptions{Keys: []string{"BETA"}})
	if err != nil {
		t.Fatal(err)
	}
	out, _ := ParseFile(dst)
	if out[0].Key != "BETA" {
		t.Errorf("expected BETA first in output, got %v", out[0].Key)
	}
}
