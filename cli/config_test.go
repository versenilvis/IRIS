package cli

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

	t.Setenv("HOME", tmpDir)
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	ConfigInitCmd.Run(ConfigInitCmd, []string{})

	configPath, err := config.ConfigPath()
	if err != nil {
		t.Fatalf("failed to get config path: %v", err)
	}
	if _, statErr := os.Stat(configPath); statErr != nil {
		t.Errorf("expected config file to be created at %s, but it was not", configPath)
	}
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("failed to read config file: %v", err)
	}
	if !strings.Contains(string(content), "shell-login = false") {
		t.Error("expected initialized config to include shell-login = false")
	}
	if !strings.Contains(string(content), `filter-exact-match = "auto"`) {
		t.Error("expected initialized config to include filter-exact-match = auto")
	}

	buf := new(bytes.Buffer)
	ConfigShowCmd.SetOut(buf)
	ConfigShowCmd.Run(ConfigShowCmd, []string{})
	if buf.Len() == 0 {
		t.Errorf("expected show command to output configuration")
	}
}
