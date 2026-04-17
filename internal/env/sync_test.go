package env

import (
	"os"
	"testing"
)

func TestSync_AddsNewKeys(t *testing.T) {
	dst := map[string]string{"APP_NAME": "myapp"}
	src := map[string]string{"APP_NAME": "myapp", "NEW_KEY": "value"}

	tmp, _ := os.CreateTemp("", "*.env")
	tmp.Close()
	defer os.Remove(tmp.Name())

	res, err := Sync(dst, src, tmp.Name(), SyncOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Applied) != 1 || res.Applied[0] != "NEW_KEY" {
		t.Errorf("expected NEW_KEY applied, got %v", res.Applied)
	}
	if len(res.Skipped) != 1 {
		t.Errorf("expected APP_NAME skipped, got %v", res.Skipped)
	}
}

func TestSync_OverwriteUpdatesExisting(t *testing.T) {
	dst := map[string]string{"DB_PASS": "old"}
	src := map[string]string{"DB_PASS": "new"}

	tmp, _ := os.CreateTemp("", "*.env")
	tmp.Close()
	defer os.Remove(tmp.Name())

	res, err := Sync(dst, src, tmp.Name(), SyncOptions{Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Applied) != 1 || res.Applied[0] != "DB_PASS" {
		t.Errorf("expected DB_PASS applied, got %v", res.Applied)
	}
}

func TestSync_DryRunDoesNotWrite(t *testing.T) {
	dst := map[string]string{"KEY": "val"}
	src := map[string]string{"NEW": "data"}

	tmp, _ := os.CreateTemp("", "*.env")
	tmp.Close()
	defer os.Remove(tmp.Name())

	_, err := Sync(dst, src, tmp.Name(), SyncOptions{DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	info, _ := os.Stat(tmp.Name())
	if info.Size() != 0 {
		t.Error("dry run should not write to file")
	}
}

func TestSync_NoOverwriteSkipsExisting(t *testing.T) {
	dst := map[string]string{"KEY": "original"}
	src := map[string]string{"KEY": "changed"}

	tmp, _ := os.CreateTemp("", "*.env")
	tmp.Close()
	defer os.Remove(tmp.Name())

	res, _ := Sync(dst, src, tmp.Name(), SyncOptions{Overwrite: false})
	if len(res.Skipped) != 1 || res.Skipped[0] != "KEY" {
		t.Errorf("expected KEY skipped, got %v", res.Skipped)
	}
}
