package env

import (
	"os"
	"testing"
)

func writeTempProfileFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "profile-*.env")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(content)
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestLoadProfiles_Basic(t *testing.T) {
	path := writeTempProfileFile(t, "# @profile dev\nDB_HOST=localhost\n# @profile prod\nDB_HOST=prod.db\n")
	profiles, err := LoadProfiles(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(profiles["dev"]) != 1 || profiles["dev"][0].Value != "localhost" {
		t.Errorf("expected dev DB_HOST=localhost")
	}
	if len(profiles["prod"]) != 1 || profiles["prod"][0].Value != "prod.db" {
		t.Errorf("expected prod DB_HOST=prod.db")
	}
}

func TestLoadProfiles_DefaultProfile(t *testing.T) {
	path := writeTempProfileFile(t, "APP=base\n# @profile dev\nDEBUG=true\n")
	profiles, err := LoadProfiles(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(profiles["default"]) != 1 || profiles["default"][0].Key != "APP" {
		t.Errorf("expected default profile with APP")
	}
}

func TestProfile_NotFound(t *testing.T) {
	profiles := map[string][]Entry{"dev": {{Key: "X", Value: "1"}}}
	_, err := Profile(profiles, "staging")
	if err == nil {
		t.Error("expected error for missing profile")
	}
}

func TestProfile_Found(t *testing.T) {
	profiles := map[string][]Entry{
		"dev": {{Key: "DB_HOST", Value: "localhost"}, {Key: "SECRET_KEY", Value: "abc"}},
	}
	res, err := Profile(profiles, "dev")
	if err != nil {
		t.Fatal(err)
	}
	if res.Loaded != 2 {
		t.Errorf("expected 2 loaded, got %d", res.Loaded)
	}
}

func TestProfileResult_Format_MasksSecrets(t *testing.T) {
	res := ProfileResult{
		Profile: "dev",
		Loaded:  1,
		Entries: []Entry{{Key: "SECRET_KEY", Value: "supersecret"}},
	}
	out := res.Format(true)
	if contains(out, "supersecret") {
		t.Error("expected secret to be masked")
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsStr(s, sub))
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
