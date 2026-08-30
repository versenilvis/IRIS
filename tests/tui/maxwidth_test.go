package tui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Gaurav-Gosain/tuitest"
)

// boxWidth measures the drawn box from its top border.
func boxWidth(t *testing.T, term *tuitest.Terminal) int {
	t.Helper()
	b := findBox(term)
	if b.top < 0 || len(b.rows) == 0 {
		t.Fatalf("no box on screen:\n%s", screen(term))
	}
	top := b.rows[0]
	start := strings.Index(top, "╭")
	end := strings.LastIndex(top, "╮")
	if start < 0 || end < 0 {
		t.Fatalf("no top border on screen:\n%s", screen(term))
	}
	return len([]rune(top[start:])) - len([]rune(top[end:])) + 1
}

func startWithMaxWidth(t *testing.T, value string) *tuitest.Terminal {
	t.Helper()
	home := t.TempDir()
	for _, sub := range []string{".config/iris", ".local/share/iris", ".cache"} {
		if err := os.MkdirAll(filepath.Join(home, sub), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	cfg := "version = 1\n\n[ui]\nmax-width = " + value + "\n"
	if err := os.WriteFile(filepath.Join(home, ".config/iris/config.toml"), []byte(cfg), 0o644); err != nil {
		t.Fatal(err)
	}
	term := startIn(t, home)
	if err := term.Type("nvi"); err != nil {
		t.Fatal(err)
	}
	if err := term.WaitForText("Accept", 10*time.Second); err != nil {
		t.Fatalf("menu never appeared: %v\n%s", err, screen(term))
	}
	return term
}

// A percentage has to scale with the terminal; a column count must not.
func TestMaxWidthAcceptsColumnsAndPercentages(t *testing.T) {
	cases := []struct {
		value string
		want  int
	}{
		{"90", 90},
		{`"80%"`, cols * 80 / 100},
		{`"50%"`, cols * 50 / 100},
		{"0", 76}, // unset keeps the overlay's own default
	}
	for _, c := range cases {
		t.Run(c.value, func(t *testing.T) {
			term := startWithMaxWidth(t, c.value)
			defer func() { _ = term.Close() }()
			if got := boxWidth(t, term); got != c.want {
				t.Errorf("max-width = %s drew a box %d columns wide; want %d\n%s", c.value, got, c.want, screen(term))
			}
		})
	}
}
