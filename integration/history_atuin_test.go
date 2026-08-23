package integration

import (
	"testing"
	"github.com/versenilvis/iris/internal/config"
)

func TestSearchHistory_AtuinSourceMapping(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Core.Atuin = 2
	cfg.Core.AtuinDBPath = "/tmp/does-not-exist.db"
	config.Init(cfg)

	sessionHistoryMu.Lock()
	origSessionHistory := sessionHistory
	sessionHistory = []string{"npm run build", "ls -l"}
	sessionHistoryMu.Unlock()

	mu.Lock()
	origHistoryCache := historyCache
	origAtuinCmds := atuinCmds
	historyCache = nil
	atuinCmds = []string{"git push", "ls -l"}
	mu.Unlock()

	t.Cleanup(func() {
		sessionHistoryMu.Lock()
		sessionHistory = origSessionHistory
		sessionHistoryMu.Unlock()

		mu.Lock()
		historyCache = origHistoryCache
		atuinCmds = origAtuinCmds
		mu.Unlock()
	})

	results, err := SearchHistory("", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) < 3 {
		t.Fatalf("expected at least 3 results, got %d", len(results))
	}

	sourceMap := make(map[string]string)
	for _, r := range results {
		sourceMap[r.Cmd] = r.Source
	}

	if source := sourceMap["ls -l"]; source != "atuin" {
		t.Errorf("expected 'ls -l' to have source 'atuin', got %q", source)
	}
	if source := sourceMap["npm run build"]; source != "session" {
		t.Errorf("expected 'npm run build' to have source 'session', got %q", source)
	}
	if source := sourceMap["git push"]; source != "atuin" {
		t.Errorf("expected 'git push' to have source 'atuin', got %q", source)
	}
}
