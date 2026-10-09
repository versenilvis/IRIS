package tui

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLiveShellAliasesReceivedOverIPC(t *testing.T) {
	home := t.TempDir()
	for _, sub := range []string{".config/iris", ".local/share/iris", ".cache"} {
		if err := os.MkdirAll(filepath.Join(home, sub), 0o755); err != nil {
			t.Fatal(err)
		}
	}

	extra := `alias mydynamiccmd="echo dynamic_success"`
	if err := os.WriteFile(filepath.Join(home, ".zshrc.extra"), []byte(extra), 0o644); err != nil {
		t.Fatal(err)
	}

	term := startIn(t, home)
	defer func() { _ = term.Close() }()

	if err := term.Type("mydynamicc"); err != nil {
		t.Fatal(err)
	}
	if err := term.WaitForText("mydynamiccmd", 5*time.Second); err != nil {
		t.Fatalf("dynamic alias suggestion never appeared: %v\n%s", err, screen(term))
	}
}
