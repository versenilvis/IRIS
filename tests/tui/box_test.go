package tui

import (
	"strings"
	"testing"

	"github.com/Gaurav-Gosain/tuitest"
)

// box describes the suggestion box as it actually appears on screen.
type box struct {
	top    int
	bottom int
	items  int
	col    int
	rows   []string
}

func findBox(term *tuitest.Terminal) box {
	var b box
	b.top, b.bottom, b.col = -1, -1, -1
	for i, line := range strings.Split(term.Snapshot(), "\n") {
		trimmed := strings.TrimRight(line, " ")
		switch {
		case strings.Contains(trimmed, "╭"):
			b.top = i
			b.col = strings.Index(trimmed, "╭")
		case strings.Contains(trimmed, "╰"):
			b.bottom = i
		case strings.Contains(trimmed, "│"):
			b.items++
		default:
			continue
		}
		b.rows = append(b.rows, trimmed)
	}
	return b
}

// assertBoxIntact checks the box is a single whole rectangle. A box drawn
// against a stale cursor row loses its top rows to the shell's repaint, and one
// left behind by a line that stopped wrapping shows up as extra pieces.
func assertBoxIntact(t *testing.T, term *tuitest.Terminal, wantItems int, stage string) {
	t.Helper()
	b := findBox(term)

	switch {
	case b.top < 0:
		t.Errorf("%s: box has no top border\n%s", stage, screen(term))
	case b.bottom < 0:
		t.Errorf("%s: box has no bottom border\n%s", stage, screen(term))
	case b.items != wantItems:
		t.Errorf("%s: box shows %d item rows, want %d\n%s", stage, b.items, wantItems, screen(term))
	case b.bottom-b.top != wantItems+1:
		t.Errorf("%s: box spans rows %d..%d, not one contiguous block of %d\n%s",
			stage, b.top, b.bottom, wantItems+2, screen(term))
	}
}
