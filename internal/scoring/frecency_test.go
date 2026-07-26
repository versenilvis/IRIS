package scoring

import (
	"context"
	"errors"
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
	_ = store.Record(context.Background(), "git status", cwd, 0)
	_ = store.Record(context.Background(), "git status", cwd, 0)
	_ = store.Record(context.Background(), "git status", cwd, 0)
	_ = store.Record(context.Background(), "git commit -m 'test'", cwd, 0)

	entries, err := store.QueryLocal(context.Background(), cwd, "git", 10)
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

	_ = store.Record(context.Background(), "make build", "/repo/a", 0)
	_ = store.Record(context.Background(), "make build", "/repo/a", 0)
	_ = store.Record(context.Background(), "make build", "/repo/b", 0)

	entries, err := store.QueryGlobal(context.Background(), "make", 10)
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

func TestFrecencyStore_SQLiteConfigurationAndContext(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "history.db")
	store, err := NewFrecencyStore(dbPath)
	if err != nil {
		t.Fatalf("NewFrecencyStore failed: %v", err)
	}
	defer store.Close()

	var journalMode string
	if qErr := store.db.QueryRowContext(context.Background(), "PRAGMA journal_mode;").Scan(&journalMode); qErr != nil {
		t.Fatalf("failed to query journal_mode: %v", qErr)
	}
	if journalMode != "wal" {
		t.Errorf("expected journal_mode 'wal', got '%s'", journalMode)
	}

	var busyTimeout int
	if qErr := store.db.QueryRowContext(context.Background(), "PRAGMA busy_timeout;").Scan(&busyTimeout); qErr != nil {
		t.Fatalf("failed to query busy_timeout: %v", qErr)
	}
	if busyTimeout != 5000 {
		t.Errorf("expected busy_timeout 5000, got %d", busyTimeout)
	}

	ctxCanceled, cancel := context.WithCancel(context.Background())
	cancel()

	err = store.Record(ctxCanceled, "git status", tmpDir, 0)
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled from Record with canceled context, got %v", err)
	}
}

func TestFrecencyStore_NilReceiver(t *testing.T) {
	var nilStore *FrecencyStore
	if err := nilStore.Record(context.Background(), "cmd", "cwd", 0); err != nil {
		t.Errorf("expected nil error on nil store Record, got %v", err)
	}
	if entries, err := nilStore.QueryLocal(context.Background(), "cwd", "", 10); err != nil || entries != nil {
		t.Errorf("expected nil entries and nil error on nil store QueryLocal, got %v, %v", entries, err)
	}
	if entries, err := nilStore.QueryGlobal(context.Background(), "", 10); err != nil || entries != nil {
		t.Errorf("expected nil entries and nil error on nil store QueryGlobal, got %v, %v", entries, err)
	}
	if err := nilStore.Close(); err != nil {
		t.Errorf("expected nil error on nil store Close, got %v", err)
	}
}

func TestFrecencyStore_ExitCodeBehavior(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "history.db")
	store, err := NewFrecencyStore(dbPath)
	if err != nil {
		t.Fatalf("NewFrecencyStore failed: %v", err)
	}
	defer store.Close()

	cwd := "/home/user/test"
	_ = store.Record(context.Background(), "grep foo", cwd, 0) // count=1
	_ = store.Record(context.Background(), "grep foo", cwd, 1) // count unchanged (1)

	entries, _ := store.QueryLocal(context.Background(), cwd, "grep", 10)
	if len(entries) != 1 || entries[0].Count != 1 {
		t.Errorf("expected grep count to be 1 after non-zero exit code, got %v", entries)
	}

	_ = store.RecordTransition(context.Background(), "git checkout", "git status", cwd, 0)
	_ = store.RecordTransition(context.Background(), "git checkout", "git status", cwd, 1)

	transitions, isLocal := store.QueryTransitionsWithFallback(context.Background(), "git checkout", cwd)
	if !isLocal || len(transitions) != 1 || transitions[0].Count != 1 {
		t.Errorf("expected transition count 1 after non-zero exit code, got %v, isLocal=%v", transitions, isLocal)
	}
}

func TestFrecencyStore_TransitionCwdIsolationAndDepthFallback(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "history.db")
	store, err := NewFrecencyStore(dbPath)
	if err != nil {
		t.Fatalf("NewFrecencyStore failed: %v", err)
	}
	defer store.Close()

	projectA := "/repo/a"
	projectB := "/repo/b"

	_ = store.RecordTransition(context.Background(), "git checkout", "npm run dev", projectA, 0)
	_ = store.RecordTransition(context.Background(), "git checkout", "go test", projectB, 0)
	_ = store.RecordTransition(context.Background(), "git checkout", "go test", projectB, 0)

	// query in project B should return go test (Local) and not npm run dev
	transB, isLocalB := store.QueryTransitionsWithFallback(context.Background(), "git checkout", projectB)
	if !isLocalB || len(transB) != 1 || transB[0].NextSkeleton != "go test" {
		t.Errorf("expected local transition 'go test' for project B, got %v (isLocal=%v)", transB, isLocalB)
	}

	// query in project C (no local data) should fallback to Global (returning both aggregated)
	projectC := "/repo/c"
	transC, isLocalC := store.QueryTransitionsWithFallback(context.Background(), "git checkout", projectC)
	if isLocalC || len(transC) != 2 {
		t.Errorf("expected global transitions for project C, got %v (isLocal=%v)", transC, isLocalC)
	}
	if transC[0].NextSkeleton != "go test" {
		t.Errorf("expected global top transition to be 'go test' (count 2), got %s", transC[0].NextSkeleton)
	}

	// depth fallback test: query deep skeleton with no exact match should fallback to shallower prefix
	_ = store.RecordTransition(context.Background(), "git remote", "git fetch", projectA, 0)
	transDeep, isLocalDeep := store.QueryTransitionsWithFallback(context.Background(), "git remote add", projectA)
	if !isLocalDeep || len(transDeep) != 1 || transDeep[0].NextSkeleton != "git fetch" {
		t.Errorf("expected depth fallback to 'git fetch' from 'git remote', got %v", transDeep)
	}
}
