package integration

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestRecordSessionCommand_MergeAndDeduplicate(t *testing.T) {
	sessionHistoryMu.Lock()
	origSessionHistory := sessionHistory
	sessionHistory = nil
	sessionHistoryMu.Unlock()

	mu.Lock()
	origHistoryCache := historyCache
	historyCache = nil
	mu.Unlock()

	t.Cleanup(func() {
		sessionHistoryMu.Lock()
		sessionHistory = origSessionHistory
		sessionHistoryMu.Unlock()

		mu.Lock()
		historyCache = origHistoryCache
		mu.Unlock()
	})

	RecordSessionCommand("git status")
	RecordSessionCommand("npm run dev")
	RecordSessionCommand("npm run dev") // duplicate subsequent command should be ignored
	RecordSessionCommand("git push origin fix/scoring")

	results, err := SearchHistory("", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) < 3 {
		t.Fatalf("expected at least 3 session commands in search results, got %d", len(results))
	}

	// newest session command must be results[0]
	if results[0].Cmd != "git push origin fix/scoring" {
		t.Errorf("expected results[0] to be 'git push origin fix/scoring', got %q", results[0].Cmd)
	}
	if results[1].Cmd != "npm run dev" {
		t.Errorf("expected results[1] to be 'npm run dev', got %q", results[1].Cmd)
	}
	if results[2].Cmd != "git status" {
		t.Errorf("expected results[2] to be 'git status', got %q", results[2].Cmd)
	}
}

func TestSearchHistory_Prefix(t *testing.T) {
	histFile := filepath.Join(t.TempDir(), "history")
	_ = os.WriteFile(histFile, []byte(""), 0600)
	t.Setenv("HISTFILE", histFile)

	sessionHistoryMu.Lock()
	origSessionHistory := sessionHistory
	sessionHistory = nil
	sessionHistoryMu.Unlock()

	mu.Lock()
	origHistoryCache := historyCache
	historyCache = nil
	mu.Unlock()

	t.Cleanup(func() {
		sessionHistoryMu.Lock()
		sessionHistory = origSessionHistory
		sessionHistoryMu.Unlock()

		mu.Lock()
		historyCache = origHistoryCache
		mu.Unlock()
	})

	RecordSessionCommand("npx tailwindcss -i input.css")
	for i := range 250 {
		RecordSessionCommand(fmt.Sprintf("git commit -m 'change %d'", i))
	}

	resN, err := SearchHistory("n", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(resN) == 0 {
		t.Fatal("expected results, got 0")
	}
	if resN[0].Cmd != "npx tailwindcss -i input.css" {
		t.Fatalf("expected prefix match 'npx tailwindcss -i input.css' at index 0, got %q", resN[0].Cmd)
	}
}
