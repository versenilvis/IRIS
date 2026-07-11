package logger

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLogger(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "iris-log-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	logFilePath := filepath.Join(tempDir, "test.log")

	// test 1: default init sets level to info
	Init(logFilePath, false)
	Debugf("this debug msg should not be logged")
	Infof("this info msg should be logged")
	Close()

	data, err := os.ReadFile(logFilePath)
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}

	content := string(data)
	if strings.Contains(content, "this debug msg should not be logged") {
		t.Errorf("expected debug log to be skipped, got: %s", content)
	}
	if !strings.Contains(content, "this info msg should be logged") {
		t.Errorf("expected info log to be recorded, got: %s", content)
	}
	if !strings.Contains(content, "[INFO]") {
		t.Errorf("expected log to contain INFO tag, got: %s", content)
	}
	if !strings.Contains(content, "logger_test.go:") {
		t.Errorf("expected log to contain caller trace info, got: %s", content)
	}

	// test 2: override init with debug = true
	_ = os.Remove(logFilePath)
	Init(logFilePath, true)
	Debugf("this debug msg should now be logged")
	Close()

	data, err = os.ReadFile(logFilePath)
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}

	content = string(data)
	if !strings.Contains(content, "this debug msg should now be logged") {
		t.Errorf("expected debug log to be recorded, got: %s", content)
	}
	if !strings.Contains(content, "[DEBUG]") {
		t.Errorf("expected log to contain DEBUG tag, got: %s", content)
	}

	// test 3: test log rotation to .old
	// create a large file
	largeData := make([]byte, 6*1024*1024)
	err = os.WriteFile(logFilePath, largeData, 0644)
	if err != nil {
		t.Fatalf("failed to write large file: %v", err)
	}

	Init(logFilePath, false)
	Infof("new log after rotation")
	Close()

	// check if old file exists and is rotated
	oldPath := logFilePath + ".old"
	if _, err = os.Stat(oldPath); os.IsNotExist(err) {
		t.Errorf("expected rotated log file to exist at %s", oldPath)
	}

	data, err = os.ReadFile(logFilePath)
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}
	content = string(data)
	if !strings.Contains(content, "new log after rotation") {
		t.Errorf("expected new log to exist in the fresh log file, got: %s", content)
	}
}
