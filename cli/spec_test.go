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
