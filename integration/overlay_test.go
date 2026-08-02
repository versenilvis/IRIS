package integration

import (
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"
	"github.com/versenilvis/iris/internal/config"
	"github.com/versenilvis/iris/spec"
)

func TestRenderGhostText_CursorAtEnd(t *testing.T) {
	o := NewOverlay()
	items := []spec.Suggestion{
		{Cmd: "git checkout -b feature"},
	}
	o.UpdateItems(items)

	// case 1: cursor at end of buffer -> should render ghost text suffix
	out := o.RenderGhostText("git check", false, true)
	if !strings.Contains(out, "out -b feature") {
		t.Fatalf("Expected ghost text suffix 'out -b feature', got: %q", out)
	}
	if o.LastGhostLen == 0 {
		t.Fatalf("Expected LastGhostLen > 0, got %d", o.LastGhostLen)
	}

	// case 2: cursor moved left (cursorAtEnd == false) -> should clear ghost text
	outClear := o.RenderGhostText("git check", false, false)
	if strings.Contains(outClear, "out -b feature") {
		t.Fatalf("Expected ghost text to be hidden/cleared when cursor moved left, got: %q", outClear)
	}
	if o.LastGhostLen != 0 {
		t.Fatalf("Expected LastGhostLen == 0 after clearing, got %d", o.LastGhostLen)
	}
}

func TestGetGhostText(t *testing.T) {
	o := NewOverlay()
	items := []spec.Suggestion{
		{Cmd: "docker exec -it my-container bash"},
	}
	o.UpdateItems(items)

	// case 1: cursor at end
	ghost := o.GetGhostText("docker e", true)
	expected := "xec -it my-container bash"
	if ghost != expected {
		t.Fatalf("Expected %q, got %q", expected, ghost)
	}

	// case 2: cursor not at end (moved left)
	ghostLeft := o.GetGhostText("docker e", false)
	if ghostLeft != "" {
		t.Fatalf("Expected empty string when cursor not at end, got %q", ghostLeft)
	}

	// case 3: user navigated menu with Up/Down arrow -> should sync with highlighted item
	o.SetUserNavigated(true)
	ghostNav := o.GetGhostText("docker e", true)
	if ghostNav != expected {
		t.Fatalf("Expected %q when user navigated menu, got %q", expected, ghostNav)
	}
}

func TestGhostText_MenuSync(t *testing.T) {
	o := NewOverlay()
	items := []spec.Suggestion{
		{Cmd: "git checkout -b first"},
		{Cmd: "git checkout master"},
	}
	o.UpdateItems(items)

	// default item 0
	ghost0 := o.GetGhostText("git check", true)
	if ghost0 != "out -b first" {
		t.Fatalf("Expected 'out -b first', got %q", ghost0)
	}

	// move cursor down to item 1
	o.MoveCursor("down")
	ghost1 := o.GetGhostText("git check", true)
	if ghost1 != "out master" {
		t.Fatalf("Expected 'out master', got %q", ghost1)
	}

	out := o.RenderGhostText("git check", true, true)
	if !strings.Contains(out, "out master") {
		t.Fatalf("Expected RenderGhostText to render 'out master', got %q", out)
	}
}

func TestGhostText_Truncation(t *testing.T) {
	o := NewOverlay()
	longCmd := "git commit -m '" + strings.Repeat("a", 150) + "'"
	items := []spec.Suggestion{
		{Cmd: longCmd},
	}
	o.UpdateItems(items)
	o.SetPromptLen(10)

	// typed query length 105 -> total cursor col = 115, default width = 120 -> available cols = 5
	typedQuery := "git commit -m '" + strings.Repeat("a", 90)
	out := o.RenderGhostText(typedQuery, false, true)
	if !strings.Contains(out, "…") {
		t.Fatalf("Expected truncated ghost text with '…', got %q", out)
	}
}

func TestHideMenu_PreservesTypedQueryForAI(t *testing.T) {
	o := NewOverlay()
	o.HideMenu("git commit")

	if o.GetTypedQuery() != "git commit" {
		t.Fatalf("Expected TypedQuery to be preserved as 'git commit', got %q", o.GetTypedQuery())
	}

	aiSugg := spec.Suggestion{
		Cmd:        "git commit -m 'fix: test'",
		Desc:       "AI suggestion",
		Source:     "ai",
		Confidence: 85,
	}
	if !o.InjectAISuggestion(aiSugg) {
		t.Fatalf("Expected InjectAISuggestion to succeed after HideMenu")
	}
	if !o.IsVisible() || len(o.Items) == 0 || o.Items[0].Cmd != aiSugg.Cmd {
		t.Fatalf("Expected AI suggestion to be injected into Items[0] and Visible=true")
	}
}

func TestRenderMatchedTitle_ASCII(t *testing.T) {
	out := renderMatchedTitle("git status", "git", false, 40)
	stripped := ansi.Strip(out)
	if stripped != "git status"+strings.Repeat(" ", 30) {
		t.Fatalf("unexpected content: %q", stripped)
	}
	if !strings.Contains(out, "git") {
		t.Fatalf("missing match segment: %q", out)
	}
}

func TestRenderMatchedTitle_WideRunes(t *testing.T) {
	out := renderMatchedTitle("日記を表示する", "日記", false, 20)
	stripped := ansi.Strip(out)
	if !strings.HasPrefix(stripped, "日記を表示する") {
		t.Fatalf("content mangled: %q", stripped)
	}
	// split must occur after 日記 (width 4), not mid-grapheme
	if !strings.Contains(stripped, "を") {
		t.Fatalf("lost remainder: %q", stripped)
	}
}

func TestRenderMatchedTitle_CaseInsensitive(t *testing.T) {
	out := renderMatchedTitle("GitHub CLI", "git", false, 40)
	stripped := ansi.Strip(out)
	if !strings.HasPrefix(stripped, "GitHub CLI") {
		t.Fatalf("case-insensitive match mangled content: %q", stripped)
	}
}

func TestMenuItemRowsHonorsMaxHeight(t *testing.T) {
	// ui.max-height was parsed and validated but never read: the box was fixed
	// at 6 rows however it was configured, unlike ui.max-width.
	tests := []struct {
		name      string
		maxHeight int
		want      int
	}{
		{"configured", 15, 15 - borderLines},
		{"minimum", 3, 1},
		{"zero falls back", 0, defaultMaxItems},
		{"out of range falls back", 999, defaultMaxItems},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.DefaultConfig()
			cfg.UI.MaxHeight = tt.maxHeight
			config.Init(cfg)

			if got := menuItemRows(); got != tt.want {
				t.Errorf("menuItemRows() with max-height=%d = %d; want %d", tt.maxHeight, got, tt.want)
			}
		})
	}
}

func TestClearRowsCoversShrunkBox(t *testing.T) {
	// A box that shrinks -- max-height is hot-reloaded -- must still be erased
	// at its old height, or the surplus rows stay on screen.
	cfg := config.DefaultConfig()
	cfg.UI.MaxHeight = 20
	config.Init(cfg)
	lastDrawnLines.Store(int32(menuItemRows() + borderLines))
	tall := clearRows()

	cfg = config.DefaultConfig()
	cfg.UI.MaxHeight = 5
	config.Init(cfg)

	if got := clearRows(); got != tall {
		t.Errorf("clearRows() after shrink = %d; want %d (the height still on screen)", got, tall)
	}
}

func TestScrolloffIsSymmetric(t *testing.T) {
	// The window kept a row of context above the highlight but none below, so
	// paging down pinned it to the last visible row while paging up did not.
	cfg := config.DefaultConfig()
	cfg.UI.MaxHeight = 8 // 6 item rows
	config.Init(cfg)

	items := make([]spec.Suggestion, 20)
	for i := range items {
		items[i] = spec.Suggestion{Cmd: string(rune('a' + i))}
	}

	o := NewOverlay()
	o.SetQueryAndItems("q", items)

	window := menuItemRows()

	// Walk down to the bottom, then back up, asserting the highlight never sits
	// flush against an edge while there are still items past it.
	for range len(items) - 1 {
		o.MoveCursor("down")
		o.Render()
		if o.Cursor < len(items)-1 && o.Cursor == o.StartIdx+window-1 {
			t.Fatalf("moving down: cursor %d pinned to last visible row (start %d, window %d)", o.Cursor, o.StartIdx, window)
		}
	}
	for range len(items) - 1 {
		o.MoveCursor("up")
		o.Render()
		if o.Cursor > 0 && o.Cursor == o.StartIdx {
			t.Fatalf("moving up: cursor %d pinned to first visible row (start %d, window %d)", o.Cursor, o.StartIdx, window)
		}
	}
}
