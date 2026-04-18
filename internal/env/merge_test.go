package env

import (
	"strings"
	"testing"
)

func TestMerge_AddsNewKeys(t *testing.T) {
	dst := map[string]string{"A": "1"}
	src := map[string]string{"B": "2"}
	out, res := Merge(dst, src, false)
	if out["B"] != "2" {
		t.Errorf("expected B=2")
	}
	if len(res.Added) != 1 || res.Added[0] != "B" {
		t.Errorf("expected B in Added")
	}
}

func TestMerge_SkipsExistingWithoutOverwrite(t *testing.T) {
	dst := map[string]string{"A": "1"}
	src := map[string]string{"A": "99"}
	out, res := Merge(dst, src, false)
	if out["A"] != "1" {
		t.Errorf("expected A unchanged")
	}
	if len(res.Skipped) != 1 {
		t.Errorf("expected A skipped")
	}
}

func TestMerge_OverwriteUpdatesExisting(t *testing.T) {
	dst := map[string]string{"A": "1"}
	src := map[string]string{"A": "99"}
	out, res := Merge(dst, src, true)
	if out["A"] != "99" {
		t.Errorf("expected A=99")
	}
	if len(res.Updated) != 1 {
		t.Errorf("expected A updated")
	}
}

func TestMerge_SameValueNotUpdated(t *testing.T) {
	dst := map[string]string{"A": "1"}
	src := map[string]string{"A": "1"}
	_, res := Merge(dst, src, true)
	if len(res.Updated) != 0 {
		t.Errorf("expected no update when value same")
	}
	if len(res.Skipped) != 1 {
		t.Errorf("expected A skipped")
	}
}

func TestMergeResult_Format(t *testing.T) {
	res := MergeResult{
		Added:   []string{"NEW"},
		Updated: []string{"CHANGED"},
		Skipped: []string{"OLD"},
	}
	out := res.Format()
	if !strings.Contains(out, "+ NEW") {
		t.Errorf("expected + NEW in output")
	}
	if !strings.Contains(out, "~ CHANGED") {
		t.Errorf("expected ~ CHANGED in output")
	}
	if !strings.Contains(out, "= OLD (skipped)") {
		t.Errorf("expected = OLD (skipped) in output")
	}
}
