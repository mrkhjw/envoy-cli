package env

import (
	"strings"
	"testing"
)

var baseVerifyEntries = []Entry{
	{Key: "APP_NAME", Value: "envoy"},
	{Key: "DB_PASSWORD", Value: "secret123"},
	{Key: "PORT", Value: "8080"},
}

func TestVerify_AllMatch(t *testing.T) {
	opts := VerifyOptions{
		Expected: map[string]string{
			"APP_NAME": "envoy",
			"PORT":     "8080",
		},
	}
	res := Verify(baseVerifyEntries, opts)
	if len(res.Matched) != 2 {
		t.Errorf("expected 2 matched, got %d", len(res.Matched))
	}
	if !res.OK() {
		t.Error("expected OK()=true")
	}
}

func TestVerify_Mismatch(t *testing.T) {
	opts := VerifyOptions{
		Expected: map[string]string{
			"PORT": "9090",
		},
	}
	res := Verify(baseVerifyEntries, opts)
	if len(res.Mismatch) != 1 || res.Mismatch[0] != "PORT" {
		t.Errorf("expected PORT in Mismatch, got %v", res.Mismatch)
	}
	if res.OK() {
		t.Error("expected OK()=false")
	}
}

func TestVerify_Missing(t *testing.T) {
	opts := VerifyOptions{
		Expected: map[string]string{
			"UNKNOWN_KEY": "value",
		},
	}
	res := Verify(baseVerifyEntries, opts)
	if len(res.Missing) != 1 || res.Missing[0] != "UNKNOWN_KEY" {
		t.Errorf("expected UNKNOWN_KEY in Missing, got %v", res.Missing)
	}
	if res.OK() {
		t.Error("expected OK()=false")
	}
}

func TestVerify_EmptyExpected(t *testing.T) {
	opts := VerifyOptions{Expected: map[string]string{}}
	res := Verify(baseVerifyEntries, opts)
	if !res.OK() {
		t.Error("expected OK()=true for empty expected map")
	}
}

func TestVerifyResult_Format_ContainsSummary(t *testing.T) {
	res := VerifyResult{
		Matched:  []string{"APP_NAME"},
		Mismatch: []string{"PORT"},
		Missing:  []string{"SECRET_KEY"},
	}
	out := res.Format()
	if !strings.Contains(out, "1 matched") {
		t.Errorf("expected '1 matched' in output: %s", out)
	}
	if !strings.Contains(out, "MISMATCH: PORT") {
		t.Errorf("expected 'MISMATCH: PORT' in output: %s", out)
	}
	if !strings.Contains(out, "MISSING:  SECRET_KEY") {
		t.Errorf("expected 'MISSING:  SECRET_KEY' in output: %s", out)
	}
}
