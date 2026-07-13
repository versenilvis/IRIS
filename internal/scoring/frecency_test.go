package scoring

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestFrecencyStore_RecordAndQueryLocal(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "history.db")
	store, err := NewFrecencyStore(dbPath)
	if err != nil {
		t.Fatalf("NewFrecencyStore failed: %v", err)
	}
	defer store.Close()

	cwd := "/home/user/project"
	_ = store.Record("git status", cwd)
	_ = store.Record("git status", cwd)
	_ = store.Record("git status", cwd)
	_ = store.Record("git commit -m 'test'", cwd)

	entries, err := store.QueryLocal(cwd, "git", 10)
	if err != nil {
		t.Fatalf("QueryLocal failed: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Cmd != "git status" || entries[0].Count != 3 {
		t.Errorf("expected top entry to be 'git status' with count 3, got %s (count %d)", entries[0].Cmd, entries[0].Count)
	}
}

func TestFrecencyStore_RawScoreDistribution(t *testing.T) {
	store := &FrecencyStore{}
	now := time.Now()

	oldHeavyScore := store.RawScore(5000, now.Add(-30*24*time.Hour))
	recentLightScore := store.RawScore(5, now.Add(-30*time.Minute))

	if oldHeavyScore <= 0 || recentLightScore <= 0 {
		t.Errorf("expected positive raw scores, got %f and %f", oldHeavyScore, recentLightScore)
	}
	if recentLightScore >= oldHeavyScore {
		t.Logf("recent light score (%f) vs old heavy score (%f)", recentLightScore, oldHeavyScore)
	}
}

func TestFrecencyStore_QueryGlobalDedupe(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "history.db")
	store, err := NewFrecencyStore(dbPath)
	if err != nil {
		t.Fatalf("NewFrecencyStore failed: %v", err)
	}
	defer store.Close()

	_ = store.Record("make build", "/repo/a")
	_ = store.Record("make build", "/repo/a")
	_ = store.Record("make build", "/repo/b")

	entries, err := store.QueryGlobal("make", 10)
	if err != nil {
		t.Fatalf("QueryGlobal failed: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 deduplicated entry, got %d", len(entries))
	}
	if entries[0].Count != 3 {
		t.Errorf("expected combined count 3 across workspaces, got %d", entries[0].Count)
	}
}

func TestFrecencyStore_Permissions(t *testing.T) {
	tmpRoot := t.TempDir()
	dbDir := filepath.Join(tmpRoot, "subdir", "iris")
	dbPath := filepath.Join(dbDir, "history.db")

	if err := os.MkdirAll(dbDir, 0755); err != nil {
		t.Fatalf("failed to make pre-existing dir: %v", err)
	}
	if err := os.WriteFile(dbPath, []byte{}, 0644); err != nil {
		t.Fatalf("failed to write dummy existing db file: %v", err)
	}

	store, err := NewFrecencyStore(dbPath)
	if err != nil {
		t.Fatalf("NewFrecencyStore failed: %v", err)
	}
	defer store.Close()

	dirInfo, err := os.Stat(dbDir)
	if err != nil {
		t.Fatalf("stat dbDir failed: %v", err)
	}
	if perm := dirInfo.Mode().Perm(); perm != 0700 {
		t.Errorf("expected directory permissions 0700, got %04o", perm)
	}

	fileInfo, err := os.Stat(dbPath)
	if err != nil {
		t.Fatalf("stat dbPath failed: %v", err)
	}
	if perm := fileInfo.Mode().Perm(); perm != 0600 {
		t.Errorf("expected database file permissions 0600, got %04o", perm)
	}
}
