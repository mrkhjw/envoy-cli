package env

import (
	"testing"
)

func TestFilter_ByPrefix(t *testing.T) {
	entries := map[string]string{
		"APP_NAME":   "myapp",
		"APP_PORT":   "8080",
		"DB_HOST":    "localhost",
	}
	res := Filter(entries, FilterOptions{Prefix: "APP_"})
	if len(res.Matched) != 2 {
		t.Errorf("expected 2 matched, got %d", len(res.Matched))
	}
	if res.Skipped != 1 {
		t.Errorf("expected 1 skipped, got %d", res.Skipped)
	}
}

func TestFilter_BySuffix(t *testing.T) {
	entries := map[string]string{
		"DB_HOST":  "localhost",
		"APP_HOST": "0.0.0.0",
		"APP_PORT": "8080",
	}
	res := Filter(entries, FilterOptions{Suffix: "_HOST"})
	if len(res.Matched) != 2 {
		t.Errorf("expected 2 matched, got %d", len(res.Matched))
	}
}

func TestFilter_ByKeys(t *testing.T) {
	entries := map[string]string{
		"APP_NAME": "myapp",
		"APP_PORT": "8080",
		"DB_HOST":  "localhost",
	}
	res := Filter(entries, FilterOptions{Keys: []string{"APP_NAME", "DB_HOST"}})
	if len(res.Matched) != 2 {
		t.Errorf("expected 2 matched, got %d", len(res.Matched))
	}
	if _, ok := res.Matched["APP_PORT"]; ok {
		t.Error("APP_PORT should not be in matched")
	}
}

func TestFilter_SecretsOnly(t *testing.T) {
	entries := map[string]string{
		"APP_NAME":    "myapp",
		"API_SECRET":  "abc123",
		"DB_PASSWORD": "pass",
	}
	res := Filter(entries, FilterOptions{SecretsOnly: true})
	if len(res.Matched) != 2 {
		t.Errorf("expected 2 secrets matched, got %d", len(res.Matched))
	}
	if _, ok := res.Matched["APP_NAME"]; ok {
		t.Error("APP_NAME should not be matched as secret")
	}
}

func TestFilter_EmptyOptions(t *testing.T) {
	entries := map[string]string{
		"A": "1",
		"B": "2",
	}
	res := Filter(entries, FilterOptions{})
	if len(res.Matched) != 2 {
		t.Errorf("expected all 2 entries matched, got %d", len(res.Matched))
	}
}
