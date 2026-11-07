package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSpecInitCmd(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	buf := new(bytes.Buffer)
	SpecInitCmd.SetOut(buf)

	SpecInitCmd.Run(SpecInitCmd, []string{"mytool"})
	out := buf.String()
	if !strings.Contains(out, "initialized spec file at") {
		t.Fatalf("unexpected output: %s", out)
	}

	targetPath := filepath.Join(tmpDir, "iris", "specs", "mytool.jsonc")
	data, err := os.ReadFile(targetPath)
	if err != nil {
		t.Fatalf("failed to read created spec file: %v", err)
	}
	if !strings.Contains(string(data), `"name": "mytool"`) {
		t.Errorf("expected spec to contain mytool, got: %s", string(data))
	}

	buf.Reset()
	SpecInitCmd.Run(SpecInitCmd, []string{"mytool"})
	if !strings.Contains(buf.String(), "spec file already exists") {
		t.Errorf("expected already exists warning, got: %s", buf.String())
	}
}

func TestSpecPathCmd(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	buf := new(bytes.Buffer)
	SpecPathCmd.SetOut(buf)

	SpecPathCmd.Run(SpecPathCmd, []string{})
	out := strings.TrimSpace(buf.String())
	expected := filepath.Join(tmpDir, "iris", "specs")
	if out != expected {
		t.Fatalf("expected path %s, got: %s", expected, out)
	}
}

func TestSpecGenerateCmd_InvalidName(t *testing.T) {
	buf := new(bytes.Buffer)
	SpecGenerateCmd.SetOut(buf)
	SpecGenerateCmd.SetErr(buf)

	err := SpecGenerateCmd.RunE(SpecGenerateCmd, []string{"../escape"})
	if err == nil {
		t.Fatal("expected error for path traversal name")
	}
}

func TestSpecGenerateCmd_Success(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)
	t.Setenv("XDG_DATA_HOME", tmpDir)

	buf := new(bytes.Buffer)
	SpecGenerateCmd.SetOut(buf)
	SpecGenerateCmd.SetErr(buf)

	// generate spec for cal which exists on standard linux systems
	err := SpecGenerateCmd.RunE(SpecGenerateCmd, []string{"cal"})
	if err != nil {
		t.Fatalf("SpecGenerateCmd failed: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "generated spec for cal at") {
		t.Fatalf("unexpected output: %s", out)
	}

	cacheFile := filepath.Join(tmpDir, "iris", "specs", "cache", "cal.json")
	data, err := os.ReadFile(cacheFile)
	if err != nil {
		t.Fatalf("expected cache file %s: %v", cacheFile, err)
	}

	if !strings.Contains(string(data), `"name": "cal"`) {
		t.Errorf("expected spec to contain cal name, got: %s", string(data))
	}
}
