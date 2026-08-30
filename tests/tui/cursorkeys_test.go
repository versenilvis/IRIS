package tui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// Walking the cursor to the far left empties the text before it, which is not
// the same as an empty line: the menu belongs to the line, so it has to stay.
func TestMovingToTheStartOfTheLineKeepsTheMenu(t *testing.T) {
	for _, tc := range []struct{ name, key string }{
		{"left", "\x1b[D"},
		{"ctrl-left", ctrlLeft},
	} {
		t.Run(tc.name, func(t *testing.T) {
			home := wordKeyHome(t)
			hist := ": 1700000000:0;git commit --amend\n"
			if err := os.WriteFile(filepath.Join(home, ".zsh_history"), []byte(hist), 0o644); err != nil {
				t.Fatal(err)
			}
			term := startIn(t, home, "HISTFILE="+filepath.Join(home, ".zsh_history"))
			defer func() { _ = term.Close() }()

			if err := term.Type("git com"); err != nil {
				t.Fatal(err)
			}
			if err := term.WaitForText("Accept", 10*time.Second); err != nil {
				t.Fatalf("no menu: %v\n%s", err, screen(term))
			}

			for range 8 {
				if err := term.SendKeys(tc.key); err != nil {
					t.Fatal(err)
				}
				if err := term.WaitStable(time.Second); err != nil {
					t.Fatal(err)
				}
			}

			if got := screen(term); !strings.Contains(got, "Accept") {
				t.Errorf("menu gone after walking left to the start:\n%s", got)
			}
			if got := promptLine(t, term); got != "git com" {
				t.Errorf("prompt = %q; want %q\n%s", got, "git com", screen(term))
			}
		})
	}
}

// Holding an arrow down repeats it faster than the throttled redraw can put the
// menu back, so tearing the menu down on every press leaves it gone for as long
// as the key is held. Spamming the same key looks fine, because the gaps between
// presses are long enough for the redraw to land.
func TestHoldingLeftKeepsTheMenuOnScreen(t *testing.T) {
	home := wordKeyHome(t)
	hist := ": 1700000000:0;git commit --amend\n"
	if err := os.WriteFile(filepath.Join(home, ".zsh_history"), []byte(hist), 0o644); err != nil {
		t.Fatal(err)
	}
	term := startIn(t, home, "HISTFILE="+filepath.Join(home, ".zsh_history"))
	defer func() { _ = term.Close() }()

	if err := term.Type("git com"); err != nil {
		t.Fatal(err)
	}
	if err := term.WaitForText("Accept", 10*time.Second); err != nil {
		t.Fatalf("no menu: %v\n%s", err, screen(term))
	}

	// sample inside the redraw gap, where a torn-down menu shows up
	const presses = 24
	seen := 0
	for range presses {
		if err := term.SendKeys("\x1b[D"); err != nil {
			t.Fatal(err)
		}
		time.Sleep(6 * time.Millisecond)
		if strings.Contains(term.Snapshot(), "Accept") {
			seen++
		}
		time.Sleep(19 * time.Millisecond)
	}

	if seen < presses/2 {
		t.Errorf("menu visible in only %d/%d samples while left was held", seen, presses)
	}
	if err := term.WaitStable(3 * time.Second); err != nil {
		t.Fatal(err)
	}
	if got := promptLine(t, term); got != "git com" {
		t.Errorf("prompt = %q; want %q\n%s", got, "git com", screen(term))
	}
}
