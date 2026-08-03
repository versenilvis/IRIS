package root

import (
	"os"
	"strings"
	"testing"

	"github.com/versenilvis/iris/internal/config"
)

func TestThemeInitCommand(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "iris-theme-cmd-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// UserConfigDir follows HOME on macOS and XDG_CONFIG_HOME on Unix.
	// Override both so the command cannot read or write the user's real theme.
	t.Setenv("HOME", tmpDir)
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	ThemeInitCmd.Run(ThemeInitCmd, []string{})

	themePath, err := config.ThemePath()
	if err != nil {
		t.Fatalf("failed to get theme path: %v", err)
	}
	if _, statErr := os.Stat(themePath); statErr != nil {
		t.Errorf("expected theme file to be created at %s, but it was not", themePath)
	}
	content, err := os.ReadFile(themePath)
	if err != nil {
		t.Fatalf("failed to read theme file: %v", err)
	}
	if !strings.Contains(string(content), "sel_text = \"#110f18\"") {
		t.Error("expected initialized theme to include sel_text = \"#110f18\"")
	}
}
