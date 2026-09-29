package tui

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/versenilvis/iris/integration"
	"github.com/versenilvis/iris/internal/scoring"
)

func predictionHome(t *testing.T) string {
	t.Helper()
	home := wordKeyHome(t)
	extra, _ := os.ReadFile(filepath.Join(home, ".zshrc.extra"))
	extra = append(extra, []byte("zle -N _iris_send_lbuffer\nadd-zle-hook-widget line-pre-redraw _iris_send_lbuffer\n")...)
	if err := os.WriteFile(filepath.Join(home, ".zshrc.extra"), extra, 0o644); err != nil {
		t.Fatal(err)
	}

	dbPath := filepath.Join(home, ".local/share/iris/history.db")
	store, err := scoring.NewFrecencyStore(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = store.Close() }()

	ctx := context.Background()
	_ = store.Record(ctx, "just reload", home, 0)
	_ = store.Record(ctx, "git add .", home, 0)
	_ = store.RecordSequence(ctx, "git add .", "git commit", home, 0)
	return home
}

func TestPredictionShowsHintAndExpandsOnRightArrow(t *testing.T) {
	home := predictionHome(t)
	term := startIn(t, home, "IRIS_CORE_MODE=history")
	defer func() { _ = term.Close() }()

	if err := term.Type("just "); err != nil {
		t.Fatal(err)
	}
	if err := term.WaitStable(2 * time.Second); err != nil {
		t.Fatal(err)
	}

	if got := screen(term); !strings.Contains(got, integration.PredictionSymbol) || !strings.Contains(got, "just reload") {
		t.Fatalf("expected prediction hint with %q and 'just reload', got screen:\n%s", integration.PredictionSymbol, got)
	}

	if err := term.SendKeys("\x1b[C"); err != nil {
		t.Fatal(err)
	}
	if err := term.WaitStable(2 * time.Second); err != nil {
		t.Fatal(err)
	}

	if got := promptLine(t, term); got != "just reload" {
		t.Fatalf("prompt = %q; want %q\nscreen:\n%s", got, "just reload", screen(term))
	}
}

func TestPredictionUnrelatedInputDoesNotShowOrExpand(t *testing.T) {
	home := predictionHome(t)
	term := startIn(t, home, "IRIS_CORE_MODE=history")
	defer func() { _ = term.Close() }()

	if err := term.Type("jar "); err != nil {
		t.Fatal(err)
	}
	if err := term.WaitStable(2 * time.Second); err != nil {
		t.Fatal(err)
	}

	if got := screen(term); strings.Contains(got, integration.PredictionSymbol) || strings.Contains(got, "just reload") {
		t.Fatalf("expected no prediction hint for unrelated input 'jar ', got screen:\n%s", got)
	}

	if err := term.SendKeys("\x1b[C"); err != nil {
		t.Fatal(err)
	}
	if err := term.WaitStable(2 * time.Second); err != nil {
		t.Fatal(err)
	}

	if got := promptLine(t, term); got != "jar" {
		t.Fatalf("prompt = %q; want 'jar'\nscreen:\n%s", got, screen(term))
	}
}

func TestPredictionRightArrowExpandsPredictionWhileMenuIsOpen(t *testing.T) {
	home := predictionHome(t)
	term := startIn(t, home, "IRIS_CORE_MODE=history")
	defer func() { _ = term.Close() }()

	if err := term.Type("just"); err != nil {
		t.Fatal(err)
	}
	if err := term.WaitForText("Accept", 10*time.Second); err != nil {
		t.Fatalf("menu did not appear: %v\n%s", err, screen(term))
	}

	// right arrow must expand prediction ("just reload"), not act like Tab
	if err := term.SendKeys("\x1b[C"); err != nil {
		t.Fatal(err)
	}
	if err := term.WaitStable(2 * time.Second); err != nil {
		t.Fatal(err)
	}

	if got := promptLine(t, term); got != "just reload" {
		t.Fatalf("prompt = %q; want 'just reload'\nscreen:\n%s", got, screen(term))
	}
}
