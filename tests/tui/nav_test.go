package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Gaurav-Gosain/tuitest"
)

// TestNavigatingOntoALongEntryKeepsTheBoxWhole walks the history menu with the
// arrow keys across entries of very different lengths. Selecting one rewrites
// the shell's line, so the input jumps between one row and several and the box
// moves with the cursor.
func TestNavigatingOntoALongEntryKeepsTheBoxWhole(t *testing.T) {
	home := t.TempDir()
	for _, sub := range []string{".config/iris", ".local/share/iris", ".cache"} {
		if err := os.MkdirAll(filepath.Join(home, sub), 0o755); err != nil {
			t.Fatal(err)
		}
	}

	long := "echo " + strings.Repeat("L", 260)
	entries := []string{
		"echo short-one",
		long,
		"echo short-two",
		"echo " + strings.Repeat("M", 130),
		"echo short-three",
	}
	var b strings.Builder
	for _, e := range entries {
		b.WriteString(": 1700000000:0;" + e + "\n")
	}
	if err := os.WriteFile(filepath.Join(home, ".zsh_history"), []byte(b.String()), 0o644); err != nil {
		t.Fatal(err)
	}

	term := startIn(t, home,
		"IRIS_CORE_MODE=history",
		"IRIS_UI_MAX_HEIGHT=5",
		"HISTFILE="+filepath.Join(home, ".zsh_history"),
	)
	defer func() { _ = term.Close() }()

	if err := term.Type("echo "); err != nil {
		t.Fatal(err)
	}
	if err := term.WaitForText("history", 10*time.Second); err != nil {
		t.Fatalf("menu never appeared: %v\n%s", err, screen(term))
	}

	for step := range 6 {
		if err := term.SendKeys(tuitest.Down); err != nil {
			t.Fatal(err)
		}
		if err := term.WaitStable(2 * time.Second); err != nil {
			t.Fatal(err)
		}
		assertBoxIntact(t, term, 5, fmt.Sprintf("after Down x%d", step+1))
		if t.Failed() {
			return
		}
	}
}

// TestBoxHoldsItsColumnWhileNavigating walks the list and checks the box does
// not slide sideways. Each step rewrites the shell's line to the selected
// entry, and entries differ in length, so a box that follows the cursor jumps
// to a new column under the entry the user is trying to read.
func TestBoxHoldsItsColumnWhileNavigating(t *testing.T) {
	home := t.TempDir()
	for _, sub := range []string{".config/iris", ".local/share/iris", ".cache"} {
		if err := os.MkdirAll(filepath.Join(home, sub), 0o755); err != nil {
			t.Fatal(err)
		}
	}

	var b strings.Builder
	for _, e := range []string{
		"nvim ~/.config/iris/config.toml",
		"nv ~/.zshrc",
		"nvim ~/.config/opencode/opencode.json",
		"nvim x",
		"nvim ~/.local/share/iris/history.db",
	} {
		b.WriteString(": 1700000000:0;" + e + "\n")
	}
	if err := os.WriteFile(filepath.Join(home, ".zsh_history"), []byte(b.String()), 0o644); err != nil {
		t.Fatal(err)
	}

	term := startIn(t, home,
		"IRIS_CORE_MODE=history",
		"IRIS_UI_MAX_HEIGHT=5",
		"HISTFILE="+filepath.Join(home, ".zsh_history"),
	)
	defer func() { _ = term.Close() }()

	if err := term.Type("nv"); err != nil {
		t.Fatal(err)
	}
	if err := term.WaitForText("history", 10*time.Second); err != nil {
		t.Fatalf("menu never appeared: %v\n%s", err, screen(term))
	}
	if err := term.WaitStable(2 * time.Second); err != nil {
		t.Fatal(err)
	}
	want := findBox(term).col

	for step := range 4 {
		if err := term.SendKeys(tuitest.Down); err != nil {
			t.Fatal(err)
		}
		if err := term.WaitStable(2 * time.Second); err != nil {
			t.Fatal(err)
		}
		if got := findBox(term).col; got != want {
			t.Fatalf("after Down x%d the box moved from column %d to %d:\n%s",
				step+1, want, got, screen(term))
		}
	}
}
