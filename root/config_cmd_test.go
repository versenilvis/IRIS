package root

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/versenilvis/iris/internal/config"
)

func TestConfigCommands(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "iris-config-cmd-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// UserConfigDir follows HOME on macOS and XDG_CONFIG_HOME on Unix.
	// Override both so the command cannot read or write the user's real config.
	t.Setenv("HOME", tmpDir)
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	ConfigInitCmd.Run(ConfigInitCmd, []string{})

	configPath, err := config.ConfigPath()
	if err != nil {
		t.Fatalf("failed to get config path: %v", err)
	}
	if _, err := os.Stat(configPath); err != nil {
		t.Errorf("expected config file to be created at %s, but it was not", configPath)
	}
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("failed to read config file: %v", err)
	}
	if !strings.Contains(string(content), "login = false") {
		t.Error("expected initialized config to include login = false")
	}

	buf := new(bytes.Buffer)
	ConfigShowCmd.SetOut(buf)
	ConfigShowCmd.Run(ConfigShowCmd, []string{})
	if buf.Len() == 0 {
		t.Errorf("expected show command to output configuration")
	}
}
