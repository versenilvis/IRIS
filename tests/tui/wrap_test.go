package tui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Gaurav-Gosain/tuitest"
)

// seedHistory writes a zsh history file whose entries share a prefix long
// enough to wrap the prompt onto several rows while the menu stays populated.
func seedHistory(t *testing.T, home, prefix string) {
	t.Helper()
	var b strings.Builder
	for _, suffix := range []string{"alpha", "bravo", "charlie", "delta", "echo", "foxtrot"} {
		b.WriteString(": 1700000000:0;" + prefix + suffix + "\n")
	}
	if err := os.WriteFile(filepath.Join(home, ".zsh_history"), []byte(b.String()), 0o644); err != nil {
		t.Fatal(err)
	}
}

func startWithHistory(t *testing.T, prefix string) *tuitest.Terminal {
	t.Helper()
	home := t.TempDir()
	for _, sub := range []string{".config/iris", ".local/share/iris", ".cache"} {
		if err := os.MkdirAll(filepath.Join(home, sub), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	seedHistory(t, home, prefix)
	return startIn(t, home,
		"IRIS_CORE_MODE=history",
		"IRIS_UI_MAX_HEIGHT=5",
		"IRIS_UI_GHOST_TEXT=0",
		"HISTFILE="+filepath.Join(home, ".zsh_history"),
	)
}

// TestBoxStaysWholeWhenTheCommandWraps types past the width of the terminal so
// the input takes several rows, then deletes back across the boundary. The box
// hangs off the cursor, so both directions move it and anything the new
// position does not cover is left on screen.
func TestBoxStaysWholeWhenTheCommandWraps(t *testing.T) {
	prefix := "echo " + strings.Repeat("z", 300) + " "
	term := startWithHistory(t, prefix)
	defer func() { _ = term.Close() }()

	if err := term.Type(prefix[:40]); err != nil {
		t.Fatal(err)
	}
	if err := term.WaitForText("history", 10*time.Second); err != nil {
		t.Fatalf("menu never appeared: %v\n%s", err, screen(term))
	}
	if err := term.WaitStable(2 * time.Second); err != nil {
		t.Fatal(err)
	}
	assertBoxIntact(t, term, 5, "one row of input")

	// grow the input across two wrap boundaries
	if err := term.Type(prefix[40:250]); err != nil {
		t.Fatal(err)
	}
	if err := term.WaitStable(3 * time.Second); err != nil {
		t.Fatal(err)
	}
	assertBoxIntact(t, term, 5, "three rows of input")

	// and back to one row
	for range 210 {
		if err := term.SendKeys(tuitest.Backspace); err != nil {
			t.Fatal(err)
		}
	}
	if err := term.WaitStable(3 * time.Second); err != nil {
		t.Fatal(err)
	}
	assertBoxIntact(t, term, 5, "back to one row")
}
