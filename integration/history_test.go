package integration

import (
	"testing"
)

func TestRecordSessionCommand_MergeAndDeduplicate(t *testing.T) {
	// reset sessionHistory
	sessionHistoryMu.Lock()
	sessionHistory = nil
	sessionHistoryMu.Unlock()

	mu.Lock()
	historyCache = nil
	mu.Unlock()

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
