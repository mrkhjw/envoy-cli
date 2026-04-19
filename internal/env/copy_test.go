package env

import (
	"strings"
	"testing"
)

func baseCopySrc() []Entry {
	return []Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "SECRET_KEY", Value: "s3cr3t"},
		{Key: "DEBUG", Value: "true"},
	}
}

func baseCopyDst() []Entry {
	return []Entry{
		{Key: "APP_NAME", Value: "oldapp"},
		{Key: "PORT", Value: "8080"},
	}
}

func TestCopy_AllKeys(t *testing.T) {
	src := baseCopySrc()
	dst := baseCopyDst()
	out, res := Copy(src, dst, CopyOptions{Overwrite: true})
	if len(res.Copied) != 3 {
		t.Errorf("expected 3 copied, got %d", len(res.Copied))
	}
	m := make(map[string]string)
	for _, e := range out {
		m[e.Key] = e.Value
	}
	if m["APP_NAME"] != "myapp" {
		t.Errorf("expected APP_NAME=myapp, got %s", m["APP_NAME"])
	}
}

func TestCopy_SkipsExistingWithoutOverwrite(t *testing.T) {
	src := baseCopySrc()
	dst := baseCopyDst()
	_, res := Copy(src, dst, CopyOptions{Overwrite: false})
	for _, k := range res.Skipped {
		if k == "APP_NAME" {
			return
		}
	}
	t.Error("expected APP_NAME to be skipped")
}

func TestCopy_SpecificKeys(t *testing.T) {
	src := baseCopySrc()
	dst := baseCopyDst()
	_, res := Copy(src, dst, CopyOptions{Keys: []string{"DEBUG"}, Overwrite: true})
	if len(res.Copied) != 1 || res.Copied[0] != "DEBUG" {
		t.Errorf("expected only DEBUG copied, got %v", res.Copied)
	}
}

func TestCopy_DryRun(t *testing.T) {
	src := baseCopySrc()
	dst := baseCopyDst()
	out, res := Copy(src, dst, CopyOptions{Overwrite: true, DryRun: true})
	if len(res.Copied) == 0 {
		t.Error("expected copied entries in dry run result")
	}
	for _, e := range out {
		if e.Key == "SECRET_KEY" {
			t.Error("dry run should not write SECRET_KEY to dst")
		}
	}
}

func TestCopyResult_Format(t *testing.T) {
	r := CopyResult{
		Copied:  []string{"A", "B"},
		Skipped: []string{"C"},
	}
	out := r.Format(false)
	if !strings.Contains(out, "2 copied") {
		t.Errorf("expected '2 copied' in output, got: %s", out)
	}
	if !strings.Contains(out, "1 skipped") {
		t.Errorf("expected '1 skipped' in output, got: %s", out)
	}
}
