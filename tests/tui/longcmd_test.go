package tui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Gaurav-Gosain/tuitest"
)

// hangReport is the command from the bug report, kept verbatim: it wraps onto
// three rows on a wide terminal and mixes quotes, box drawing runes and
// backslash escapes.
const hangReport = `sleep 1; { echo "═══ ps ═══"; ps -eo pid,ppid,pgid,stat,wchan:20,etimes,cmd | grep -Ei "iris|zsh|tmux" | grep -v grep; echo "═══ env of hung shell ═══"; for p in $(pgrep -x zsh); do echo "-- pid $p --"; tr "\0" "\n" < /proc/$p/environ 2>/dev/null | grep -E "^(IRIS_|TMUX|FF_SHOWN)"; done; echo "═══ dump ═══"; cat ~/iris-hang-dump.log; } > ~/iris-hang-report.txt 2>&1`

// TestLongCommandRendersIntact reproduces the reported screen: a wide terminal,
// the command wrapped onto three rows, and the menu open underneath.
func TestLongCommandRendersIntact(t *testing.T) {
	home := t.TempDir()
	for _, sub := range []string{".config/iris", ".local/share/iris", ".cache"} {
		if err := os.MkdirAll(filepath.Join(home, sub), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	var b strings.Builder
	for _, e := range []string{
		"su -", "ssh build-host", hangReport, "systemctl status",
		"sudo pacman -Syu", "sort -u notes.txt",
	} {
		b.WriteString(": 1700000000:0;" + e + "\n")
	}
	if err := os.WriteFile(filepath.Join(home, ".zsh_history"), []byte(b.String()), 0o644); err != nil {
		t.Fatal(err)
	}

	env := []string{
		"IRIS_CORE_MODE=history",
		"IRIS_UI_MAX_HEIGHT=6",
		"HISTFILE=" + filepath.Join(home, ".zsh_history"),
	}
	if os.Getenv("IRIS_TUI_NO_GHOST") != "" {
		env = append(env, "IRIS_UI_GHOST_TEXT=0")
	}
	term := startIn(t, home, env...)
	defer func() { _ = term.Close() }()

	if err := term.Type("s"); err != nil {
		t.Fatal(err)
	}
	if err := term.WaitForText("history", 10*time.Second); err != nil {
		t.Fatalf("menu never appeared: %v\n%s", err, screen(term))
	}
	if err := term.WaitStable(2 * time.Second); err != nil {
		t.Fatal(err)
	}

	// walk onto the long entry, which replaces the line and wraps it
	for range 4 {
		if err := term.SendKeys(tuitest.Down); err != nil {
			t.Fatal(err)
		}
		if err := term.WaitStable(2 * time.Second); err != nil {
			t.Fatal(err)
		}
		if strings.Contains(term.Snapshot(), "iris-hang-report.txt") {
			break
		}
	}
	assertBoxIntact(t, term, 6, "long command selected")
	assertCommandIntact(t, term, "long command selected")
}

// assertCommandIntact checks the shell's line still reads as it should. The
// overlay writes runs of spaces to erase its own ghost text, and one landing on
// the wrong row punches a hole straight through the command.
func assertCommandIntact(t *testing.T, term *tuitest.Terminal, stage string) {
	t.Helper()
	var typed strings.Builder
	for _, line := range strings.Split(term.Snapshot(), "\n") {
		if strings.ContainsAny(line, "╭╮╰╯│") {
			break
		}
		typed.WriteString(strings.TrimRight(line, " "))
	}
	flat := typed.String()
	for _, fragment := range []string{
		"ps -eo pid,ppid,pgid,stat,wchan:20,etimes,cmd",
		"~/iris-hang-report.txt 2>&1",
	} {
		if !strings.Contains(flat, fragment) {
			t.Errorf("%s: command line lost %q:\n%s", stage, fragment, screen(term))
			return
		}
	}
}
