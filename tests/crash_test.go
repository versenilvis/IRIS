package tests

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/versenilvis/iris/root"
)

func TestWriteCrashLog(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "iris-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	origHome := os.Getenv("HOME")
	origCache := os.Getenv("XDG_CACHE_HOME")
	defer func() {
		_ = os.Setenv("HOME", origHome)
		_ = os.Setenv("XDG_CACHE_HOME", origCache)
	}()
	_ = os.Setenv("HOME", tmpDir)
	_ = os.Setenv("XDG_CACHE_HOME", filepath.Join(tmpDir, ".cache"))

	testErr := "test panic message"
	root.WriteCrashLog(testErr)

	dir := filepath.Join(tmpDir, ".cache", "iris", "crashes")
	files, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("failed to read crashes dir: %v", err)
	}
	if len(files) != 1 {
		t.Fatalf("expected exactly 1 crash log file, got %d", len(files))
	}

	logPath := filepath.Join(dir, files[0].Name())
	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("failed to read crash log: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "=== IRIS CRASH ") {
		t.Errorf("expected header in log, got: %s", content)
	}
	if !strings.Contains(content, "panic: test panic message") {
		t.Errorf("expected panic message in log, got: %s", content)
	}
	if !strings.Contains(content, "version:") {
		t.Errorf("expected version in log, got: %s", content)
	}
}

func TestCrashLogCommand(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "iris-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	origHome := os.Getenv("HOME")
	origCache := os.Getenv("XDG_CACHE_HOME")
	defer func() {
		_ = os.Setenv("HOME", origHome)
		_ = os.Setenv("XDG_CACHE_HOME", origCache)
	}()
	_ = os.Setenv("HOME", tmpDir)
	_ = os.Setenv("XDG_CACHE_HOME", filepath.Join(tmpDir, ".cache"))

	var buf bytes.Buffer
	root.CrashCmd.SetOut(&buf)
	root.CrashCmd.SetArgs([]string{})
	root.ClearLog = false

	root.CrashCmd.Run(root.CrashCmd, []string{})
	if !strings.Contains(buf.String(), "no crash log found") {
		t.Errorf("expected 'no crash log found', got: %q", buf.String())
	}

	root.WriteCrashLog("mock error")
	buf.Reset()
	root.CrashCmd.Run(root.CrashCmd, []string{})
	if !strings.Contains(buf.String(), "crash_") || !strings.Contains(buf.String(), ".log") {
		t.Errorf("expected crash log path, got: %q", buf.String())
	}

	buf.Reset()
	root.ClearLog = true
	root.CrashCmd.Run(root.CrashCmd, []string{})
	if !strings.Contains(buf.String(), "crash log cleared") {
		t.Errorf("expected 'crash log cleared', got: %q", buf.String())
	}

	dir := filepath.Join(tmpDir, ".cache", "iris", "crashes")
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		t.Errorf("expected crashes directory to be deleted, but it exists")
	}
}
