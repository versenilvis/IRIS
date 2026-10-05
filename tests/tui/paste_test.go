package tui

import (
	"strings"
	"testing"
	"time"
)

func TestBracketedPasteRendersCleanly(t *testing.T) {
	home := predictionHome(t)
	term := startIn(t, home, "IRIS_CORE_MODE=history")
	defer func() { _ = term.Close() }()

	// bracketed paste sequence enclosing a command
	pasteSeq := "\x1b[200~echo 'hello world'\x1b[201~"
	if err := term.SendKeys(pasteSeq); err != nil {
		t.Fatal(err)
	}
	if err := term.WaitStable(2 * time.Second); err != nil {
		t.Fatal(err)
	}

	got := promptLine(t, term)
	if !strings.Contains(got, "echo 'hello world'") {
		t.Fatalf("prompt = %q; want to contain %q\nscreen:\n%s", got, "echo 'hello world'", screen(term))
	}
}

func TestBracketedPasteWithTabDoesNotTriggerCompletion(t *testing.T) {
	home := predictionHome(t)
	term := startIn(t, home, "IRIS_CORE_MODE=history")
	defer func() { _ = term.Close() }()

	// tab inside paste should be forwarded as raw character, not iris completion
	pasteSeq := "\x1b[200~echo\t'pasted'\x1b[201~"
	if err := term.SendKeys(pasteSeq); err != nil {
		t.Fatal(err)
	}
	if err := term.WaitStable(2 * time.Second); err != nil {
		t.Fatal(err)
	}

	got := promptLine(t, term)
	if !strings.Contains(got, "pasted") {
		t.Fatalf("prompt = %q; want to contain 'pasted'\nscreen:\n%s", got, screen(term))
	}
}
