package integration

import (
	"strings"
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

func TestSearchHistory_PrefixRankedAboveSubstringAndAlias(t *testing.T) {
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

	RecordSessionCommand("git branch -M main")
	RecordSessionCommand("go mod init github.com/versenilvis/payments")
	RecordSessionCommand("go get -u github.com/gin-gonic/gin@latest")
	RecordSessionCommand("go get github.com/jackc/pgx/v5")
	RecordSessionCommand("z git")
	RecordSessionCommand("go get github.com/joho/godotenv")

	aliases := map[string]string{
		"g": "git",
	}

	results, err := SearchHistory("git", aliases)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) == 0 {
		t.Fatalf("expected results, got 0")
	}

	idxGitBranch, idxZGit, idxGoGet := -1, -1, -1
	for idx, r := range results {
		if r.Cmd == "git branch -M main" && idxGitBranch == -1 {
			idxGitBranch = idx
		}
		if r.Cmd == "z git" && idxZGit == -1 {
			idxZGit = idx
		}
		if strings.HasPrefix(r.Cmd, "go get github.com") && idxGoGet == -1 {
			idxGoGet = idx
		}
	}

	if idxGitBranch == -1 {
		t.Fatalf("expected 'git branch -M main' in results")
	}
	if idxZGit == -1 {
		t.Fatalf("expected 'z git' in results")
	}
	if idxGoGet == -1 {
		t.Fatalf("expected 'go get github.com...' in results")
	}

	if idxGitBranch >= idxZGit {
		t.Errorf("expected 'git branch -M main' (index %d) before 'z git' (index %d)", idxGitBranch, idxZGit)
	}
	if idxZGit >= idxGoGet {
		t.Errorf("expected 'z git' (index %d) before 'go get github.com' (index %d)", idxZGit, idxGoGet)
	}
}

func TestSearchHistory_FuzzyMatchPreserved(t *testing.T) {
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

	RecordSessionCommand("git commit -m \"initial commit\"")
	RecordSessionCommand("cargo build --release")

	results, err := SearchHistory("gcm", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) == 0 {
		t.Fatalf("expected fuzzy results for 'gcm', got 0")
	}

	if results[0].Cmd != "git commit -m \"initial commit\"" {
		t.Errorf("expected 'git commit -m \"initial commit\"' at index 0, got %q", results[0].Cmd)
	}
}
