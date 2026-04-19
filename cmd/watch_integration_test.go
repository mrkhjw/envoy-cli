package cmd

import (
	"bytes"
	"os"
	"strings"
	"syscall"
	"testing"
	"time"
)

func TestWatchCmd_OutputsWatchingMessage(t *testing.T) {
	path := writeTempWatchCmdEnv(t, "API_KEY=secret\n")

	buf := &bytes.Buffer{}
	watchCmd.SetOut(buf)
	watchCmd.SetErr(buf)
	watchCmd.SetArgs([]string{"--interval", "30", path})

	errCh := make(chan error, 1)
	go func() {
		errCh <- watchCmd.Execute()
	}()

	time.Sleep(80 * time.Millisecond)
	syscall.Kill(syscall.Getpid(), syscall.SIGINT)

	select {
	case <-errCh:
	case <-time.After(1 * time.Second):
		t.Error("watch command did not exit after SIGINT")
		return
	}

	output := buf.String()
	if !strings.Contains(output, "Watching") {
		t.Errorf("expected 'Watching' in output, got: %s", output)
	}
}

func TestWatchCmd_DetectsFileChange(t *testing.T) {
	path := writeTempWatchCmdEnv(t, "DB_URL=original\n")

	buf := &bytes.Buffer{}
	watchCmd.SetOut(buf)
	watchCmd.SetErr(buf)
	watchCmd.SetArgs([]string{"--interval", "30", path})

	errCh := make(chan error, 1)
	go func() {
		errCh <- watchCmd.Execute()
	}()

	time.Sleep(60 * time.Millisecond)
	os.WriteFile(path, []byte("DB_URL=updated\n"), 0644)
	time.Sleep(80 * time.Millisecond)
	syscall.Kill(syscall.Getpid(), syscall.SIGINT)

	select {
	case <-errCh:
	case <-time.After(1 * time.Second):
		t.Error("watch command did not exit")
		return
	}

	output := buf.String()
	if !strings.Contains(output, "changed") {
		t.Errorf("expected change notification in output, got: %s", output)
	}
}
