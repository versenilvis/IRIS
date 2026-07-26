package scoring

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestCollectSignals(t *testing.T) {
	tmpDir := t.TempDir()
	_ = os.WriteFile(filepath.Join(tmpDir, "package.json"), []byte("{}"), 0644)

	dbPath := filepath.Join(tmpDir, "history.db")
	store, err := NewFrecencyStore(dbPath)
	if err != nil {
		t.Fatalf("NewFrecencyStore failed: %v", err)
	}
	defer store.Close()

	_ = store.Record(context.Background(), "npm run dev", tmpDir, 0)
	_ = store.Record(context.Background(), "npm test", "/other/dir", 0)
	_ = store.RecordTransition(context.Background(), "git checkout", "npm run dev", tmpDir, 0)

	signals := CollectSignals(context.Background(), tmpDir, "npm", "npm", store, "git checkout")

	if !signals.Workspace.HasNodeProject {
		t.Error("expected HasNodeProject to be true in collected signals")
	}
	if len(signals.LocalFrecency) != 1 || signals.LocalFrecency[0].Cmd != "npm run dev" {
		t.Errorf("expected local frecency to contain 'npm run dev', got %v", signals.LocalFrecency)
	}
	if len(signals.GlobalFrecency) != 2 {
		t.Errorf("expected global frecency to contain 2 entries, got %d", len(signals.GlobalFrecency))
	}
	if !signals.TransitionIsLocal || len(signals.TransitionEntries) != 1 || signals.TransitionEntries[0].NextSkeleton != "npm run dev" {
		t.Errorf("expected transition entry 'npm run dev', got %v (isLocal=%v)", signals.TransitionEntries, signals.TransitionIsLocal)
	}
}
