package tui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Gaurav-Gosain/tuitest"
)

const (
	ctrlLeft  = "\x1b[1;5D"
	ctrlRight = "\x1b[1;5C"
)

func wordKeyHome(t *testing.T) string {
	t.Helper()
	home := t.TempDir()
	for _, sub := range []string{".config/iris", ".local/share/iris", ".cache"} {
		if err := os.MkdirAll(filepath.Join(home, sub), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	extra := "bindkey '^[[1;5D' backward-word\nbindkey '^[[1;5C' forward-word\n"
	if err := os.WriteFile(filepath.Join(home, ".zshrc.extra"), []byte(extra), 0o644); err != nil {
		t.Fatal(err)
	}
	return home
}

func promptLine(t *testing.T, term *tuitest.Terminal) string {
	t.Helper()
	for _, line := range strings.Split(screen(term), "\n") {
		if strings.HasPrefix(line, "> ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "> "))
		}
	}
	return ""
}

// The reported failure: ctrl+arrow left iris' idea of the line out of step with
// the shell's, and the plain arrows afterwards did nothing at all.
func TestArrowsStillWorkAfterCtrlArrow(t *testing.T) {
	term := startIn(t, wordKeyHome(t))
	defer func() { _ = term.Close() }()

	if err := term.Type("echo hello world"); err != nil {
		t.Fatal(err)
	}
	if err := term.WaitStable(2 * time.Second); err != nil {
		t.Fatal(err)
	}

	for _, k := range []string{ctrlLeft, "\x1b[D", "\x1b[D"} {
		if err := term.SendKeys(k); err != nil {
			t.Fatal(err)
		}
		if err := term.WaitStable(time.Second); err != nil {
			t.Fatal(err)
		}
	}
	if err := term.Type("X"); err != nil {
		t.Fatal(err)
	}
	if err := term.WaitStable(2 * time.Second); err != nil {
		t.Fatal(err)
	}

	// ctrl+left parks on the 'w', two lefts step back over "o ", so X lands
	// between "hell" and "o world"
	if got := promptLine(t, term); got != "echo hellXo world" {
		t.Errorf("prompt = %q; want %q\n%s", got, "echo hellXo world", screen(term))
	}
}

func TestCtrlArrowComposesWithoutManglingTheLine(t *testing.T) {
	term := startIn(t, wordKeyHome(t))
	defer func() { _ = term.Close() }()

	if err := term.Type("echo alpha beta"); err != nil {
		t.Fatal(err)
	}
	if err := term.WaitStable(2 * time.Second); err != nil {
		t.Fatal(err)
	}

	for _, k := range []string{ctrlLeft, ctrlLeft, ctrlRight} {
		if err := term.SendKeys(k); err != nil {
			t.Fatal(err)
		}
		if err := term.WaitStable(time.Second); err != nil {
			t.Fatal(err)
		}
	}
	if err := term.Type("-"); err != nil {
		t.Fatal(err)
	}
	if err := term.WaitStable(2 * time.Second); err != nil {
		t.Fatal(err)
	}

	// zsh's forward-word stops at the start of the next word, so the third
	// motion lands on the "b" rather than after "alpha"
	if got := promptLine(t, term); got != "echo alpha -beta" {
		t.Errorf("prompt = %q; want %q\n%s", got, "echo alpha -beta", screen(term))
	}
}

// The menu belongs to the line, not to the cursor sitting at its end, so moving
// by a word has to leave it up and redrawn rather than tearing it down.
func TestCtrlArrowKeepsTheMenuUp(t *testing.T) {
	home := wordKeyHome(t)
	hist := ": 1700000000:0;echo hello world\n: 1700000001:0;echo hello there\n"
	if err := os.WriteFile(filepath.Join(home, ".zsh_history"), []byte(hist), 0o644); err != nil {
		t.Fatal(err)
	}
	term := startIn(t, home, "IRIS_CORE_MODE=history", "HISTFILE="+filepath.Join(home, ".zsh_history"))
	defer func() { _ = term.Close() }()

	if err := term.Type("echo hello"); err != nil {
		t.Fatal(err)
	}
	if err := term.WaitForText("Accept", 10*time.Second); err != nil {
		t.Fatalf("menu never appeared: %v\n%s", err, screen(term))
	}

	if err := term.SendKeys(ctrlLeft); err != nil {
		t.Fatal(err)
	}
	if err := term.WaitStable(2 * time.Second); err != nil {
		t.Fatal(err)
	}

	if got := screen(term); !strings.Contains(got, "Accept") {
		t.Errorf("menu gone after ctrl+left:\n%s", got)
	}
	if got := promptLine(t, term); got != "echo hello" {
		t.Errorf("prompt = %q; want %q\n%s", got, "echo hello", screen(term))
	}
}
