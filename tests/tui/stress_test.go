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

// TestWalkingDeepIntoTheListKeepsTheBoxWhole walks a long history whose entries
// swing between one row and three, the way a real history does. Each step
// rewrites the shell's line, so the box moves up and down the screen and every
// step is a chance to leave part of the old one behind.
func TestWalkingDeepIntoTheListKeepsTheBoxWhole(t *testing.T) {
	home := t.TempDir()
	for _, sub := range []string{".config/iris", ".local/share/iris", ".cache"} {
		if err := os.MkdirAll(filepath.Join(home, sub), 0o755); err != nil {
			t.Fatal(err)
		}
	}

	var b strings.Builder
	for i := range 40 {
		var entry string
		switch i % 4 {
		case 0:
			entry = fmt.Sprintf("echo short-%02d", i)
		case 1:
			entry = "echo " + strings.Repeat("w", 120) + fmt.Sprintf("-%02d", i)
		case 2:
			entry = fmt.Sprintf("echo mid-%02d ", i) + strings.Repeat("m", 40)
		case 3:
			entry = "echo " + strings.Repeat("v", 260) + fmt.Sprintf("-%02d", i)
		}
		b.WriteString(": 1700000000:0;" + entry + "\n")
	}
	if err := os.WriteFile(filepath.Join(home, ".zsh_history"), []byte(b.String()), 0o644); err != nil {
		t.Fatal(err)
	}

	term := startIn(t, home,
		"IRIS_CORE_MODE=history",
		"IRIS_UI_MAX_HEIGHT=6",
		"HISTFILE="+filepath.Join(home, ".zsh_history"),
	)
	defer func() { _ = term.Close() }()

	if err := term.Type("echo "); err != nil {
		t.Fatal(err)
	}
	if err := term.WaitForText("history", 10*time.Second); err != nil {
		t.Fatalf("menu never appeared: %v\n%s", err, screen(term))
	}
	if err := term.WaitStable(2 * time.Second); err != nil {
		t.Fatal(err)
	}

	// held down, not tapped: the keys arrive faster than the shell repaints
	for range 30 {
		if err := term.SendKeys(tuitest.Down); err != nil {
			t.Fatal(err)
		}
		time.Sleep(8 * time.Millisecond)
	}
	if err := term.WaitStable(3 * time.Second); err != nil {
		t.Fatal(err)
	}
	assertBoxIntact(t, term, 6, "after holding Down")
}

// TestShrinkingTheQueryLeavesNoFragment types until the line wraps, which puts
// the box two rows further down, then deletes back to one row. The list is
// rebuilt as the query changes, so the box that has to be erased is a different
// one from the box being drawn.
func TestShrinkingTheQueryLeavesNoFragment(t *testing.T) {
	home := t.TempDir()
	for _, sub := range []string{".config/iris", ".local/share/iris", ".cache"} {
		if err := os.MkdirAll(filepath.Join(home, sub), 0o755); err != nil {
			t.Fatal(err)
		}
	}

	prefix := "echo " + strings.Repeat("p", 200)
	var b strings.Builder
	for i := range 8 {
		fmt.Fprintf(&b, ": 1700000000:0;%s-tail-%02d\n", prefix, i)
	}
	if err := os.WriteFile(filepath.Join(home, ".zsh_history"), []byte(b.String()), 0o644); err != nil {
		t.Fatal(err)
	}

	term := startIn(t, home,
		"IRIS_CORE_MODE=history",
		"IRIS_UI_MAX_HEIGHT=6",
		"HISTFILE="+filepath.Join(home, ".zsh_history"),
	)
	defer func() { _ = term.Close() }()

	if err := term.Type(prefix[:150]); err != nil {
		t.Fatal(err)
	}
	if err := term.WaitForText("history", 10*time.Second); err != nil {
		t.Fatalf("menu never appeared: %v\n%s", err, screen(term))
	}
	if err := term.WaitStable(3 * time.Second); err != nil {
		t.Fatal(err)
	}
	assertBoxIntact(t, term, 6, "query wrapped onto two rows")

	// delete back onto a single row
	for range 130 {
		if err := term.SendKeys(tuitest.Backspace); err != nil {
			t.Fatal(err)
		}
		time.Sleep(6 * time.Millisecond)
	}
	if err := term.WaitStable(3 * time.Second); err != nil {
		t.Fatal(err)
	}
	assertBoxIntact(t, term, 6, "query back on one row")
}

// TestReloadDoesNotStrandTheBox reloads while the menu is on screen. The reload
// replaces the process, so the one that comes back has no record of the box the
// old one drew and cannot erase it.
func TestReloadDoesNotStrandTheBox(t *testing.T) {
	home := t.TempDir()
	for _, sub := range []string{".config/iris", ".local/share/iris", ".cache"} {
		if err := os.MkdirAll(filepath.Join(home, sub), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	var b strings.Builder
	for _, e := range []string{
		"nvim ~/.config/", "nvim ~/.config/iris/", "nvim ~/.config/opencode/",
		"nvim ~/.config/iris/config.toml", "nvim a.cxx", "nv a.go", "nv a.cpp",
	} {
		b.WriteString(": 1700000000:0;" + e + "\n")
	}
	if err := os.WriteFile(filepath.Join(home, ".zsh_history"), []byte(b.String()), 0o644); err != nil {
		t.Fatal(err)
	}

	term := startIn(t, home,
		"IRIS_CORE_MODE=history",
		"IRIS_UI_MAX_HEIGHT=6",
		"HISTFILE="+filepath.Join(home, ".zsh_history"),
	)
	defer func() { _ = term.Close() }()

	if err := term.Type("nvim ~/.config/"); err != nil {
		t.Fatal(err)
	}
	if err := term.WaitForText("history", 10*time.Second); err != nil {
		t.Fatalf("menu never appeared: %v\n%s", err, screen(term))
	}
	if err := term.WaitStable(2 * time.Second); err != nil {
		t.Fatal(err)
	}

	// clear the line and reload with the box still up
	if err := term.SendKeys(tuitest.Ctrl('u')); err != nil {
		t.Fatal(err)
	}
	if err := term.Type("iris reload"); err != nil {
		t.Fatal(err)
	}
	if err := term.SendKeys(tuitest.Enter); err != nil {
		t.Fatal(err)
	}
	if err := term.WaitStable(5 * time.Second); err != nil {
		t.Fatal(err)
	}

	if got := screen(term); strings.ContainsAny(got, "╭╮╰╯│") {
		t.Errorf("a box survived the reload:\n%s", got)
	}
}
