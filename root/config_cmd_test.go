package root

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestConfigCommands(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "iris-config-cmd-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	origConfigHome := os.Getenv("XDG_CONFIG_HOME")
	defer func() {
		_ = os.Setenv("XDG_CONFIG_HOME", origConfigHome)
	}()
	_ = os.Setenv("XDG_CONFIG_HOME", tmpDir)

	ConfigInitCmd.Run(ConfigInitCmd, []string{})

	configPath := filepath.Join(tmpDir, "iris", "config.toml")
	if _, err := os.Stat(configPath); err != nil {
		t.Errorf("expected config file to be created at %s, but it was not", configPath)
	}

	buf := new(bytes.Buffer)
	ConfigShowCmd.SetOut(buf)
	ConfigShowCmd.Run(ConfigShowCmd, []string{})
	if buf.Len() == 0 {
		t.Errorf("expected show command to output configuration")
	}
}
