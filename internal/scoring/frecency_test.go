package scoring

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"os/exec"
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

func TestFrecencyStore_SubstringMatch(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "history.db")
	store, err := NewFrecencyStore(dbPath)
	if err != nil {
		t.Fatalf("NewFrecencyStore failed: %v", err)
	}
	defer store.Close()

	cwd := "/home/user/project"
	_ = store.Record(context.Background(), "sudo chown -R www-data:www-data /var/www/html", cwd, 0)
	_ = store.Record(context.Background(), "sudo systemctl restart nginx", cwd, 0)
	_ = store.Record(context.Background(), "chmod +x script.sh", cwd, 0)

	_ = store.Record(context.Background(), "echo hello world", cwd, 0)

	// Test QueryLocal with substring
	entries, err := store.QueryLocal(context.Background(), cwd, "cho", 10)
	if err != nil {
		t.Fatalf("QueryLocal failed: %v", err)
	}
	if len(entries) != 2 { // should match 'sudo chown...' and 'echo...' (has 'cho')
		t.Fatalf("expected 2 entries for 'cho', got %d", len(entries))
	}

	foundChown := false
	for _, e := range entries {
		if e.Cmd == "sudo chown -R www-data:www-data /var/www/html" {
			foundChown = true
			break
		}
	}
	if !foundChown {
		t.Errorf("expected 'sudo chown...' to be found with substring 'cho'")
	}

	// Test QueryGlobal with substring
	entriesGlobal, err := store.QueryGlobal(context.Background(), "chown", 10)
	if err != nil {
		t.Fatalf("QueryGlobal failed: %v", err)
	}
	if len(entriesGlobal) != 1 {
		t.Fatalf("expected 1 entry for 'chown', got %d", len(entriesGlobal))
	}
	if entriesGlobal[0].Cmd != "sudo chown -R www-data:www-data /var/www/html" {
		t.Errorf("expected 'sudo chown...' to be found with substring 'chown', got %s", entriesGlobal[0].Cmd)
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

func TestFrecencyStore_RecordExitCodeZeroVsNonZero(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "history.db")
	store, err := NewFrecencyStore(dbPath)
	if err != nil {
		t.Fatalf("NewFrecencyStore failed: %v", err)
	}
	defer store.Close()

	cwd := tmpDir

	// exitCode != 0 should not insert
	_ = store.Record(context.Background(), "failed cmd", cwd, 1)
	_ = store.RecordSequence(context.Background(), "prev", "failed next", cwd, 1)

	ctx := context.Background()

	var count int
	_ = store.db.QueryRowContext(ctx, "SELECT count(*) FROM history_entries WHERE cmd = 'failed cmd'").Scan(&count)
	if count != 0 {
		t.Fatalf("expected 0 entries for failed cmd, got %d", count)
	}

	var seqCount int
	_ = store.db.QueryRowContext(ctx, "SELECT count(*) FROM command_sequences WHERE next_cmd = 'failed next'").Scan(&seqCount)
	if seqCount != 0 {
		t.Fatalf("expected 0 sequence entries for failed next, got %d", seqCount)
	}

	// exitCode == 0 should insert and increment
	_ = store.Record(ctx, "success cmd", cwd, 0)
	_ = store.Record(ctx, "success cmd", cwd, 0)

	var successCount int
	_ = store.db.QueryRowContext(ctx, "SELECT count FROM history_entries WHERE cmd = 'success cmd'").Scan(&successCount)
	if successCount != 2 {
		t.Fatalf("expected count 2 for success cmd, got %d", successCount)
	}
}

func TestFrecencyStore_LegacyMigration(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "legacy_history.db")
	ctx := context.Background()

	// 1. Create a legacy database without project_id column
	rawDB, err := openRawLegacyDB(dbPath)
	if err != nil {
		t.Fatalf("failed to create raw legacy db: %v", err)
	}

	existingDir := filepath.Join(tmpDir, "repo")
	_ = os.MkdirAll(filepath.Join(existingDir, ".git"), 0755)

	nonExistingDir := filepath.Join(tmpDir, "deleted_folder")

	// Insert legacy rows: some count > 0, some count = 0
	if _, err = rawDB.ExecContext(ctx, "INSERT INTO history_entries (cmd, cwd, count) VALUES ('git status', ?, 5)", existingDir); err != nil {
		t.Fatalf("inserting git status failed: %v", err)
	}
	if _, err = rawDB.ExecContext(ctx, "INSERT INTO history_entries (cmd, cwd, count) VALUES ('failed cmd', ?, 0)", existingDir); err != nil {
		t.Fatalf("inserting failed cmd failed: %v", err)
	}
	if _, err = rawDB.ExecContext(ctx, "INSERT INTO history_entries (cmd, cwd, count) VALUES ('dead folder cmd', ?, 3)", nonExistingDir); err != nil {
		t.Fatalf("inserting dead folder cmd failed: %v", err)
	}
	if _, err = rawDB.ExecContext(ctx, "INSERT INTO command_sequences (prev_cmd, next_cmd, cwd, count) VALUES ('seq_failed', 'next', ?, 0)", existingDir); err != nil {
		t.Fatalf("inserting failed sequence failed: %v", err)
	}
	if _, err = rawDB.ExecContext(ctx, "INSERT INTO command_sequences (prev_cmd, next_cmd, cwd, count) VALUES ('seq_ok', 'next', ?, 2)", existingDir); err != nil {
		t.Fatalf("inserting ok sequence failed: %v", err)
	}
	_ = rawDB.Close()

	// 2. Open via NewFrecencyStore to trigger backup, migration, cleanup, and backfill
	store, err := NewFrecencyStore(dbPath)
	if err != nil {
		t.Fatalf("NewFrecencyStore on legacy db failed: %v", err)
	}

	// Verify backup file was created
	if _, errStat := os.Stat(dbPath + ".bak"); errStat != nil {
		t.Fatalf("expected backup file %s.bak to exist: %v", dbPath, errStat)
	}

	// Close store which waits for backfill to finish
	_ = store.Close()

	// 3. Inspect resulting database
	checkDB, err := openRawLegacyDB(dbPath)
	if err != nil {
		t.Fatalf("failed to reopen db: %v", err)
	}
	defer checkDB.Close()

	// count <= 0 rows must be deleted
	var failedCount int
	_ = checkDB.QueryRowContext(ctx, "SELECT count(*) FROM history_entries WHERE cmd = 'failed cmd'").Scan(&failedCount)
	if failedCount != 0 {
		t.Fatalf("expected failed cmd (count=0) to be deleted, found %d", failedCount)
	}

	// existingDir row must have project_id backfilled to existingDir
	var pid string
	err = checkDB.QueryRowContext(ctx, "SELECT project_id FROM history_entries WHERE cmd = 'git status'").Scan(&pid)
	if err != nil || pid != existingDir {
		t.Fatalf("expected project_id=%q for git status, got %q (err=%v)", existingDir, pid, err)
	}

	// nonExistingDir row must have project_id backfilled to empty string "" (not NULL)
	var deadPid *string
	err = checkDB.QueryRowContext(ctx, "SELECT project_id FROM history_entries WHERE cmd = 'dead folder cmd'").Scan(&deadPid)
	val := "<nil>"
	if deadPid != nil {
		val = *deadPid
	}
	if err != nil || deadPid == nil || *deadPid != "" {
		t.Fatalf("expected project_id='' for dead folder, got %q (err=%v)", val, err)
	}

	// sequence count <= 0 rows must be deleted
	var failedSeqCount int
	_ = checkDB.QueryRowContext(ctx, "SELECT count(*) FROM command_sequences WHERE prev_cmd = 'seq_failed'").Scan(&failedSeqCount)
	if failedSeqCount != 0 {
		t.Fatalf("expected failed sequence (count=0) to be deleted, found %d", failedSeqCount)
	}

	// sequence ok row must have project_id backfilled
	var seqPid string
	err = checkDB.QueryRowContext(ctx, "SELECT project_id FROM command_sequences WHERE prev_cmd = 'seq_ok'").Scan(&seqPid)
	if err != nil || seqPid != existingDir {
		t.Fatalf("expected project_id=%q for seq_ok, got %q (err=%v)", existingDir, seqPid, err)
	}
	_ = checkDB.Close()

	// 4. Reopen already-migrated database: verify .bak is not recreated
	_ = os.Remove(dbPath + ".bak")
	store2, err := NewFrecencyStore(dbPath)
	if err != nil {
		t.Fatalf("reopening migrated store failed: %v", err)
	}
	_ = store2.Close()
	if _, errStat := os.Stat(dbPath + ".bak"); !os.IsNotExist(errStat) {
		t.Fatalf("expected .bak not to be created on already-migrated database: %v", errStat)
	}
}

func TestFrecencyStore_LegacyMigration_WALConsistency(t *testing.T) {
	if os.Getenv("TEST_SUBPROCESS_WAL") == "1" {
		dbPath := os.Getenv("TEST_WAL_DBPATH")
		rawDB, err := openRawLegacyDB(dbPath)
		if err != nil {
			os.Exit(1)
		}
		ctx := context.Background()
		if _, err = rawDB.ExecContext(ctx, "PRAGMA journal_mode = WAL;"); err != nil {
			os.Exit(2)
		}
		for i := range 20 {
			cmd := fmt.Sprintf("cmd_%02d", i)
			if _, err = rawDB.ExecContext(ctx, "INSERT INTO history_entries (cmd, cwd, count) VALUES (?, ?, 1)", cmd, filepath.Dir(dbPath)); err != nil {
				os.Exit(3)
			}
		}
		// Exit immediately without calling rawDB.Close() to leave uncheckpointed WAL frames on disk
		os.Exit(0)
	}

	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "history.db")
	ctx := context.Background()

	cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=TestFrecencyStore_LegacyMigration_WALConsistency")
	cmd.Env = append(os.Environ(), "TEST_SUBPROCESS_WAL=1", "TEST_WAL_DBPATH="+dbPath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("subprocess failed: %v, out: %s", err, out)
	}

	// Ensure WAL file exists with uncheckpointed frames
	if fi, errStat := os.Stat(dbPath + "-wal"); errStat != nil || fi.Size() == 0 {
		t.Fatalf("expected non-empty WAL file before migration: %v", errStat)
	}

	// 3. Open via NewFrecencyStore to trigger migration and backup
	store, err := NewFrecencyStore(dbPath)
	if err != nil {
		t.Fatalf("NewFrecencyStore failed: %v", err)
	}
	_ = store.Close()

	// 4. Open the .bak file independently and verify integrity and exact row count
	bakPath := dbPath + ".bak"
	bakDB, err := openRawLegacyDB(bakPath)
	if err != nil {
		t.Fatalf("failed to open backup file: %v", err)
	}
	defer bakDB.Close()

	var integrity string
	if err = bakDB.QueryRowContext(ctx, "PRAGMA integrity_check;").Scan(&integrity); err != nil || integrity != "ok" {
		t.Fatalf("expected integrity_check=ok on backup, got %q (err=%v)", integrity, err)
	}

	var rowCount int
	if err = bakDB.QueryRowContext(ctx, "SELECT count(*) FROM history_entries;").Scan(&rowCount); err != nil {
		t.Fatalf("failed to count rows in backup: %v", err)
	}
	if rowCount != 20 {
		t.Fatalf("expected 20 rows in backup from uncheckpointed WAL, got %d", rowCount)
	}

	// Verify backup file permissions are 0600
	if fi, errStat := os.Stat(bakPath); errStat == nil {
		if fi.Mode().Perm() != 0o600 {
			t.Fatalf("expected 0600 permissions on backup, got %v", fi.Mode().Perm())
		}
	}
}

func TestFrecencyStore_DoNotOverwriteProjectIDWithEmpty(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "history.db")
	store, err := NewFrecencyStore(dbPath)
	if err != nil {
		t.Fatalf("NewFrecencyStore failed: %v", err)
	}
	defer store.Close()

	ctx := context.Background()
	repoDir := filepath.Join(tmpDir, "repo")
	_ = os.MkdirAll(filepath.Join(repoDir, ".git"), 0755)

	// Record with valid project_id
	_ = store.Record(ctx, "npm test", repoDir, 0)

	var initialPID string
	_ = store.db.QueryRowContext(ctx, "SELECT project_id FROM history_entries WHERE cmd = 'npm test'").Scan(&initialPID)
	if initialPID != repoDir {
		t.Fatalf("expected initial project_id=%q, got %q", repoDir, initialPID)
	}

	// Now delete .git to simulate a transient detection failure
	_ = os.RemoveAll(filepath.Join(repoDir, ".git"))

	// Record again - project_id must NOT be overwritten with ""
	_ = store.Record(ctx, "npm test", repoDir, 0)

	var finalPID string
	_ = store.db.QueryRowContext(ctx, "SELECT project_id FROM history_entries WHERE cmd = 'npm test'").Scan(&finalPID)
	if finalPID != repoDir {
		t.Fatalf("project_id was clobbered! got %q, want %q", finalPID, repoDir)
	}
}

func TestFrecencyStore_CwdNormalization(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "history.db")
	store, err := NewFrecencyStore(dbPath)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	defer store.Close()

	ctx := context.Background()
	testDir := filepath.Join(tmpDir, "myrepo")
	_ = os.MkdirAll(testDir, 0755)

	// test Record with trailing slash and query without
	dirWithSlash := testDir + "/"
	if recErr := store.Record(ctx, "git status", dirWithSlash, 0); recErr != nil {
		t.Fatalf("record failed: %v", recErr)
	}
	entries, qErr := store.QueryLocal(ctx, testDir, "git", 10)
	if qErr != nil || len(entries) != 1 || entries[0].Cmd != "git status" {
		t.Fatalf("expected 1 entry from QueryLocal, got %v (err: %v)", entries, qErr)
	}

	// query candidate exact cwd match across trailing slash difference
	candidates := store.QueryHistoryCandidates(ctx, "git", testDir, "")
	if len(candidates) != 1 || candidates[0].Tier != 4 {
		t.Fatalf("expected tier 4 candidate, got %v", candidates)
	}

	// test RecordSequence and query with trailing slash mismatch
	if err := store.RecordSequence(ctx, "git status", "git diff", dirWithSlash, 0); err != nil {
		t.Fatalf("record sequence failed: %v", err)
	}
	seqEntries, ok := store.QuerySequencesWithFallback(ctx, "git status", testDir)
	if !ok || len(seqEntries) != 1 || seqEntries[0].NextCmd != "git diff" {
		t.Fatalf("expected sequence entry, got %v", seqEntries)
	}

	// test QueryTransitionsWithFallback across trailing slash difference
	if err := store.RecordTransition(ctx, "git status", "git commit", dirWithSlash, 0); err != nil {
		t.Fatalf("record transition failed: %v", err)
	}
	transitions, ok := store.QueryTransitionsWithFallback(ctx, "git status", testDir)
	if !ok || len(transitions) != 1 || transitions[0].NextSkeleton != "git commit" {
		t.Fatalf("expected transition entry, got %v", transitions)
	}

	// test tierOf with trailing slashes
	if tier := tierOf(dirWithSlash, "proj", testDir, "proj"); tier != 4 {
		t.Fatalf("expected tier 4 for slash mismatch, got %d", tier)
	}
}

func openRawLegacyDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	legacySchema := `
CREATE TABLE IF NOT EXISTS history_entries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    cmd TEXT NOT NULL,
    cwd TEXT NOT NULL,
    count INTEGER DEFAULT 1,
    last_used TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(cmd, cwd)
);

CREATE TABLE IF NOT EXISTS command_sequences (
    id        INTEGER PRIMARY KEY AUTOINCREMENT,
    prev_cmd  TEXT NOT NULL,
    next_cmd  TEXT NOT NULL,
    cwd       TEXT NOT NULL,
    count     INTEGER DEFAULT 1,
    last_used TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(prev_cmd, next_cmd, cwd)
);
`
	_, err = db.ExecContext(context.Background(), legacySchema)
	return db, err
}
