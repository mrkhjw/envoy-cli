package env

import (
	"os"
	"testing"
	"time"
)

func writeTempWatchEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "watch-*.env")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(content)
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestHashFile_Stable(t *testing.T) {
	path := writeTempWatchEnv(t, "KEY=value\n")
	h1, err := HashFile(path)
	if err != nil {
		t.Fatal(err)
	}
	h2, err := HashFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if h1 != h2 {
		t.Errorf("expected stable hash, got %s and %s", h1, h2)
	}
}

func TestHashFile_ChangesOnEdit(t *testing.T) {
	path := writeTempWatchEnv(t, "KEY=value\n")
	h1, _ := HashFile(path)
	os.WriteFile(path, []byte("KEY=changed\n"), 0644)
	h2, _ := HashFile(path)
	if h1 == h2 {
		t.Error("expected hash to change after file edit")
	}
}

func TestWatch_DetectsChange(t *testing.T) {
	path := writeTempWatchEnv(t, "KEY=original\n")

	done := make(chan struct{})
	results := make(chan WatchResult, 1)

	go Watch(path, 20*time.Millisecond, done, func(r WatchResult) {
		results <- r
	})

	time.Sleep(40 * time.Millisecond)
	os.WriteFile(path, []byte("KEY=updated\n"), 0644)

	select {
	case r := <-results:
		if !r.Changed {
			t.Error("expected change to be detected")
		}
	case <-time.After(300 * time.Millisecond):
		t.Error("timed out waiting for change")
	}
	close(done)
}

func TestWatchResult_Format_Changed(t *testing.T) {
	r := WatchResult{File: ".env", Changed: true, OldHash: "abcdef1234567890", NewHash: "1234567890abcdef"}
	out := r.Format()
	if out == "" {
		t.Error("expected non-empty format output")
	}
}

func TestWatchResult_Format_NoChange(t *testing.T) {
	r := WatchResult{File: ".env", Changed: false}
	out := r.Format()
	if out == "" {
		t.Error("expected non-empty format output")
	}
}
