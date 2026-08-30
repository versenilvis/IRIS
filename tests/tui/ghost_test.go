package tui

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func ghostHome(t *testing.T) string {
	t.Helper()
	home := wordKeyHome(t)
	// add-zle-hook-widget only takes a widget, and iris' own init registers a
	// plain function, so on a bare zsh the line hook never fires. Register it
	// properly here: this is the setup these tests are about.
	extra, _ := os.ReadFile(filepath.Join(home, ".zshrc.extra"))
	extra = append(extra, []byte("zle -N _iris_send_lbuffer\nadd-zle-hook-widget line-pre-redraw _iris_send_lbuffer\n")...)
	if err := os.WriteFile(filepath.Join(home, ".zshrc.extra"), extra, 0o644); err != nil {
		t.Fatal(err)
	}
	hist := ": 1700000000:0;nvim ~/.config/iris/config.toml\n: 1700000001:0;nvim ~/.config/iris/theme.toml\n"
	if err := os.WriteFile(filepath.Join(home, ".zsh_history"), []byte(hist), 0o644); err != nil {
		t.Fatal(err)
	}
	return home
}

// The hint belongs after the end of the line. Reporting only the text left of
// the cursor made iris take that prefix for the whole line and draw the hint
// over everything the user still had to the right of it.
func TestGhostTextStaysOffWhenTheCursorIsMidLine(t *testing.T) {
	home := ghostHome(t)
	term := startIn(t, home, "IRIS_CORE_MODE=history", "HISTFILE="+filepath.Join(home, ".zsh_history"))
	defer func() { _ = term.Close() }()

	if err := term.Type("nvim ~/.con"); err != nil {
		t.Fatal(err)
	}
	if err := term.WaitForText("Accept", 10*time.Second); err != nil {
		t.Fatalf("no menu: %v\n%s", err, screen(term))
	}
	if err := term.Type("Z"); err != nil {
		t.Fatal(err)
	}
	if err := term.WaitStable(2 * time.Second); err != nil {
		t.Fatal(err)
	}

	// backspacing behind the Z leaves "nvim ~/.co" to the left of the cursor,
	// which is a live prefix of both history entries
	for _, k := range []string{"\x1b[D", "\x7f"} {
		if err := term.SendKeys(k); err != nil {
			t.Fatal(err)
		}
		if err := term.WaitStable(2 * time.Second); err != nil {
			t.Fatal(err)
		}
	}

	if got := promptLine(t, term); got != "nvim ~/.coZ" {
		t.Errorf("prompt = %q; want %q\n%s", got, "nvim ~/.coZ", screen(term))
	}
}

func TestMovingTheCursorLeavesNoGhostBehind(t *testing.T) {
	home := ghostHome(t)
	term := startIn(t, home, "IRIS_CORE_MODE=history", "HISTFILE="+filepath.Join(home, ".zsh_history"))
	defer func() { _ = term.Close() }()

	if err := term.Type("nvim ~/.c"); err != nil {
		t.Fatal(err)
	}
	if err := term.WaitForText("Accept", 10*time.Second); err != nil {
		t.Fatalf("no menu: %v\n%s", err, screen(term))
	}

	for _, k := range []string{"\x1b[D", ctrlLeft} {
		if err := term.SendKeys(k); err != nil {
			t.Fatal(err)
		}
		if err := term.WaitStable(2 * time.Second); err != nil {
			t.Fatal(err)
		}
		if got := promptLine(t, term); got != "nvim ~/.c" {
			t.Errorf("after %q prompt = %q; want %q\n%s", k, got, "nvim ~/.c", screen(term))
		}
	}
}
