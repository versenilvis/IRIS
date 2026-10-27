package tui

import (
	"fmt"
	"os"
	"os/exec"
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

func TestKittyKeyboardNavigation(t *testing.T) {
	home := t.TempDir()
	for _, sub := range []string{".config/iris", ".local/share/iris", ".cache"} {
		if err := os.MkdirAll(filepath.Join(home, sub), 0o755); err != nil {
			t.Fatal(err)
		}
	}

	var b strings.Builder
	for _, e := range []string{
		"nvim ~/.config/iris/config.toml",
		"nvim ~/.zshrc",
		"nvim test.txt",
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

	if err := term.Type("nvim "); err != nil {
		t.Fatal(err)
	}
	if err := term.WaitForText("history", 10*time.Second); err != nil {
		t.Fatalf("menu never appeared: %v\n%s", err, screen(term))
	}

	// navigate down using kitty keyboard protocol sequence
	if err := term.SendKeys("\x1b[57353u"); err != nil {
		t.Fatal(err)
	}
	if err := term.WaitStable(2 * time.Second); err != nil {
		t.Fatal(err)
	}
	if got := screen(term); !strings.Contains(got, "2/3") {
		t.Fatalf("after kitty down want 2/3 selection:\n%s", got)
	}

	// navigate up using kitty keyboard protocol sequence
	if err := term.SendKeys("\x1b[57352u"); err != nil {
		t.Fatal(err)
	}
	if err := term.WaitStable(2 * time.Second); err != nil {
		t.Fatal(err)
	}
	if got := screen(term); !strings.Contains(got, "1/3") {
		t.Fatalf("after kitty up want 1/3 selection:\n%s", got)
	}

	// toggle menu using kitty shift+tab sequence
	if err := term.SendKeys("\x1b[9;2u"); err != nil {
		t.Fatal(err)
	}
	if err := term.WaitStable(2 * time.Second); err != nil {
		t.Fatal(err)
	}
	if got := screen(term); strings.Contains(got, "history") {
		t.Fatalf("menu should be hidden after shift+tab:\n%s", got)
	}
}

func startInFish(t *testing.T, home string, extraEnv ...string) *tuitest.Terminal {
	t.Helper()
	if _, err := exec.LookPath("fish"); err != nil {
		t.Skip("fish not installed")
	}
	bin := binary(t)
	confDir := filepath.Join(home, ".config", "fish")
	_ = os.MkdirAll(confDir, 0o755)
	fishrc := "set -g fish_greeting ''\nset -g fish_man_page_completion_suppress 1\nfunction fish_prompt\n    echo '> '\nend\neval (" + bin + " init fish | string collect)\n"
	if err := os.WriteFile(filepath.Join(confDir, "config.fish"), []byte(fishrc), 0o644); err != nil {
		t.Fatal(err)
	}

	cacheDir, _ := os.MkdirTemp("", "fish-cache-*")
	t.Cleanup(func() {
		_ = os.RemoveAll(cacheDir)
	})

	env := []string{
		"HOME=" + home,
		"XDG_CONFIG_HOME=" + filepath.Join(home, ".config"),
		"XDG_DATA_HOME=" + filepath.Join(home, ".local/share"),
		"XDG_CACHE_HOME=" + cacheDir,
		"SHELL=/usr/bin/fish",
		"IRIS_ACTIVE_SHELL=fish",
		"PATH=" + os.Getenv("PATH"),
		"TERM=xterm-ghostty",
	}
	env = append(env, extraEnv...)

	term := tuitest.StartT(t, []string{bin, "--shell", "fish"},
		tuitest.WithSize(cols, rows),
		tuitest.WithEnv(env...),
		tuitest.WithDir(home),
	)

	if err := term.WaitForText(">", 20*time.Second); err != nil {
		t.Fatalf("iris never reached a prompt in fish: %v\n%s", err, term.Snapshot())
	}
	return term
}

func TestFishNavigationWithKittyKeys(t *testing.T) {
	home := t.TempDir()
	for _, sub := range []string{".config/iris", ".local/share/iris"} {
		if err := os.MkdirAll(filepath.Join(home, sub), 0o755); err != nil {
			t.Fatal(err)
		}
	}

	term := startInFish(t, home,
		"IRIS_CORE_MODE=spec",
		"IRIS_UI_MAX_HEIGHT=5",
	)
	defer func() {
		_ = term.SendKeys("\x03exit\n")
		time.Sleep(100 * time.Millisecond)
		_ = term.Close()
	}()

	if err := term.Type("git "); err != nil {
		t.Fatal(err)
	}
	if err := term.WaitForText("Accept", 10*time.Second); err != nil {
		t.Fatalf("menu never appeared: %v\n%s", err, screen(term))
	}

	// navigate down using kitty keyboard sequence in fish
	if err := term.SendKeys("\x1b[57353u"); err != nil {
		t.Fatal(err)
	}
	if err := term.WaitStable(2 * time.Second); err != nil {
		t.Fatal(err)
	}
	if got := screen(term); !strings.Contains(got, "2/") {
		t.Fatalf("after kitty down in fish want 2/.. selection:\n%s", got)
	}

	// navigate up using kitty keyboard sequence in fish
	if err := term.SendKeys("\x1b[57352u"); err != nil {
		t.Fatal(err)
	}
	if err := term.WaitStable(2 * time.Second); err != nil {
		t.Fatal(err)
	}
	if got := screen(term); !strings.Contains(got, "1/") {
		t.Fatalf("after kitty up in fish want 1/.. selection:\n%s", got)
	}
}


