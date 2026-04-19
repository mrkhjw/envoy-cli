package env

import (
	"strings"
	"testing"
)

func TestSearch_ByKey(t *testing.T) {
	env := map[string]string{"APP_HOST": "localhost", "APP_PORT": "8080", "DB_URL": "postgres://"}
	res := Search(env, SearchOptions{Key: "APP"})
	if len(res.Matches) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(res.Matches))
	}
}

func TestSearch_ByValue(t *testing.T) {
	env := map[string]string{"HOST": "localhost", "BACKUP": "localhost-backup"}
	res := Search(env, SearchOptions{Value: "localhost"})
	if len(res.Matches) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(res.Matches))
	}
}

func TestSearch_CaseSensitive(t *testing.T) {
	env := map[string]string{"APP_HOST": "localhost", "app_debug": "true"}
	res := Search(env, SearchOptions{Key: "APP", CaseSensitive: true})
	if len(res.Matches) != 1 {
		t.Fatalf("expected 1 match, got %d", len(res.Matches))
	}
}

func TestSearch_NoMatches(t *testing.T) {
	env := map[string]string{"HOST": "localhost"}
	res := Search(env, SearchOptions{Key: "MISSING"})
	if len(res.Matches) != 0 {
		t.Fatalf("expected 0 matches, got %d", len(res.Matches))
	}
	if res.Format(false) != "no matches found" {
		t.Errorf("unexpected format output")
	}
}

func TestSearch_MasksSecrets(t *testing.T) {
	env := map[string]string{"API_SECRET": "supersecret"}
	res := Search(env, SearchOptions{Key: "SECRET"})
	out := res.Format(true)
	if strings.Contains(out, "supersecret") {
		t.Errorf("expected secret to be masked")
	}
	if !strings.Contains(out, "***") {
		t.Errorf("expected *** in output")
	}
}

func TestSearch_KeyAndValue(t *testing.T) {
	env := map[string]string{"APP_HOST": "localhost", "APP_PORT": "8080"}
	res := Search(env, SearchOptions{Key: "APP", Value: "8080"})
	if len(res.Matches) != 1 {
		t.Fatalf("expected 1 match, got %d", len(res.Matches))
	}
	if res.Matches[0].Key != "APP_PORT" {
		t.Errorf("unexpected match key: %s", res.Matches[0].Key)
	}
}
