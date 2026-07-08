package tests

import (
	"strings"
	"testing"

	"github.com/versenilvis/iris/integration"
	"github.com/versenilvis/iris/spec"
)

func TestRenderGhostText_CursorAtEnd(t *testing.T) {
	o := integration.NewOverlay()
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
	o := integration.NewOverlay()
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
	o := integration.NewOverlay()
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
	o := integration.NewOverlay()
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
	o := integration.NewOverlay()
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
