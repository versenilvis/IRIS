package scoring

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/versenilvis/iris/internal/config"
	"github.com/versenilvis/iris/internal/workspace"
	_ "modernc.org/sqlite"
)

type FrecencyEntry struct {
	Cmd      string
	Cwd      string
	Count    int
	LastUsed time.Time
	RawScore float64
}

type TransitionEntry struct {
	PrevSkeleton string
	NextSkeleton string
	Cwd          string
	Count        int
	LastUsed     time.Time
}

type SequenceEntry struct {
	PrevCmd  string
	NextCmd  string
	Cwd      string
	Count    int
	LastUsed time.Time
}

type Candidate struct {
	Cmd      string
	Tier     int
	Count    int
	LastUsed time.Time
}

const (
	LocalLimit           = 200
	GlobalLimit          = 100
	MaxCandidates        = 30
	GlobalScopeThreshold = 3
)

type FrecencyStore struct {
	db         *sql.DB
	mu         sync.Mutex
	bgWg       sync.WaitGroup
	dbPath     string
	backupOnce sync.Once
}

func (f *FrecencyStore) backupDatabase(ctx context.Context) error {
	if f.dbPath == "" || f.dbPath == ":memory:" || f.db == nil {
		return nil
	}
	bakPath := f.dbPath + ".bak"

	// create empty destination file with 0600 permissions
	fBak, err := os.OpenFile(bakPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o600)
	if err != nil {
		return fmt.Errorf("failed to create backup file %s: %w", bakPath, err)
	}
	_ = fBak.Close()

	// VACUUM INTO writes an atomic, WAL-consistent copy into the empty file
	_, err = f.db.ExecContext(ctx, "VACUUM INTO ?", bakPath)
	if err == nil {
		return nil
	}

	// fallback: truncate WAL checkpoint, then copy file bytes
	_, _ = f.db.ExecContext(ctx, "PRAGMA wal_checkpoint(TRUNCATE)")
	data, errRead := os.ReadFile(f.dbPath)
	if errRead != nil {
		_ = os.Remove(bakPath)
		return fmt.Errorf("backup failed via VACUUM INTO (%w) and fallback read (%w)", err, errRead)
	}
	if errWrite := os.WriteFile(bakPath, data, 0o600); errWrite != nil {
		_ = os.Remove(bakPath)
		return fmt.Errorf("backup failed via VACUUM INTO (%w) and fallback write (%w)", err, errWrite)
	}
	return nil
}

func NewFrecencyStore(dbPath string) (*FrecencyStore, error) {
	if dbPath == "" {
		var err error
		dbPath, err = config.HistoryDBPath()
		if err != nil {
			return nil, err
		}
	}

	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create directory for history.db: %w", err)
	}
	_ = os.Chmod(dir, 0700)

	if f, err := os.OpenFile(dbPath, os.O_CREATE, 0600); err == nil {
		_ = f.Close()
	}
	_ = os.Chmod(dbPath, 0600)

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite database: %w", err)
	}
	db.SetMaxOpenConns(1)

	store := &FrecencyStore{db: db, dbPath: dbPath}
	if err := store.initSchema(context.Background()); err != nil {
		_ = db.Close()
		return nil, err
	}
	_ = os.Chmod(dbPath, 0600)
	go store.BootstrapSequences(context.Background(), "", "")

	return store, nil
}

func (f *FrecencyStore) configureSQLite(ctx context.Context) error {
	_, err := f.db.ExecContext(ctx, "PRAGMA journal_mode = WAL; PRAGMA busy_timeout = 5000;")
	return err
}

func (f *FrecencyStore) initSchema(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}
	ctxTimeout, cancel := context.WithTimeout(ctx, 2000*time.Millisecond)
	defer cancel()

	if err := f.configureSQLite(ctxTimeout); err != nil {
		return err
	}

	schema := `
CREATE TABLE IF NOT EXISTS history_entries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    cmd TEXT NOT NULL,
    cwd TEXT NOT NULL,
    project_id TEXT DEFAULT NULL,
    count INTEGER DEFAULT 1,
    last_used TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(cmd, cwd)
);

CREATE INDEX IF NOT EXISTS idx_history_cwd_cmd ON history_entries(cwd, cmd);

CREATE TABLE IF NOT EXISTS command_transitions (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    prev_skeleton TEXT NOT NULL,
    next_skeleton TEXT NOT NULL,
    cwd           TEXT NOT NULL,
    count         INTEGER DEFAULT 1,
    last_used     TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(prev_skeleton, next_skeleton, cwd)
);

CREATE INDEX IF NOT EXISTS idx_transitions_prev_cwd ON command_transitions(prev_skeleton, cwd);

CREATE TABLE IF NOT EXISTS command_sequences (
    id        INTEGER PRIMARY KEY AUTOINCREMENT,
    prev_cmd  TEXT NOT NULL,
    next_cmd  TEXT NOT NULL,
    cwd       TEXT NOT NULL,
    project_id TEXT DEFAULT NULL,
    count     INTEGER DEFAULT 1,
    last_used TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(prev_cmd, next_cmd, cwd)
);

CREATE INDEX IF NOT EXISTS idx_sequences_prev_cwd ON command_sequences(prev_cmd, cwd);
`
	if _, err := f.db.ExecContext(ctxTimeout, schema); err != nil {
		return err
	}

	addedHist, err := f.addColumnIfNotExists(ctxTimeout, "history_entries", "project_id", "TEXT DEFAULT NULL")
	if err != nil {
		return err
	}
	addedSeq, err := f.addColumnIfNotExists(ctxTimeout, "command_sequences", "project_id", "TEXT DEFAULT NULL")
	if err != nil {
		return err
	}

	indexSQL := `
CREATE INDEX IF NOT EXISTS idx_history_project_cmd ON history_entries(project_id, cmd);
CREATE INDEX IF NOT EXISTS idx_sequences_project_prev ON command_sequences(project_id, prev_cmd);
`
	if _, err := f.db.ExecContext(ctxTimeout, indexSQL); err != nil {
		return err
	}

	if addedHist || addedSeq {
		cleanupSQL := `
DELETE FROM history_entries WHERE count <= 0;
DELETE FROM command_sequences WHERE count <= 0;
`
		if _, err := f.db.ExecContext(ctxTimeout, cleanupSQL); err != nil {
			return err
		}
	}

	f.bgWg.Add(1)
	go func() {
		defer f.bgWg.Done()
		f.backfillProjectIDs()
	}()
	return nil
}

func (f *FrecencyStore) addColumnIfNotExists(ctx context.Context, table, column, colDef string) (bool, error) {
	rows, err := f.db.QueryContext(ctx, fmt.Sprintf("PRAGMA table_info(%s)", table))
	if err != nil {
		return false, err
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		var cid int
		var name, ctype string
		var notnull, pk int
		var dfltValue interface{}
		if scanErr := rows.Scan(&cid, &name, &ctype, &notnull, &dfltValue, &pk); scanErr == nil {
			if strings.EqualFold(name, column) {
				return false, nil
			}
		}
	}
	if rowsErr := rows.Err(); rowsErr != nil {
		return false, rowsErr
	}
	var backupErr error
	f.backupOnce.Do(func() {
		backupErr = f.backupDatabase(ctx)
	})
	if backupErr != nil {
		return false, fmt.Errorf("schema migration aborted: failed to create database backup: %w", backupErr)
	}
	_, err = f.db.ExecContext(ctx, fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", table, column, colDef))
	if err != nil {
		return false, err
	}
	return true, nil
}

func (f *FrecencyStore) backfillProjectIDs() {
	if f == nil || f.db == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	rows, err := f.db.QueryContext(ctx, "SELECT DISTINCT cwd FROM history_entries WHERE project_id IS NULL")
	if err == nil {
		defer func() { _ = rows.Close() }()
		var cwds []string
		for rows.Next() {
			var d string
			if errScan := rows.Scan(&d); errScan == nil && d != "" {
				cwds = append(cwds, d)
			}
		}
		if rowsErr := rows.Err(); rowsErr == nil {
			for _, d := range cwds {
				norm := workspace.Normalize(d)
				pid := ""
				if _, statErr := os.Stat(norm); statErr == nil {
					pid = workspace.ProjectID(workspace.DetectRoot(norm))
				}
				f.mu.Lock()
				_, _ = f.db.ExecContext(ctx, "UPDATE history_entries SET project_id = ? WHERE cwd = ? AND project_id IS NULL", pid, d)
				f.mu.Unlock()
			}
		}
	}

	seqRows, seqErr := f.db.QueryContext(ctx, "SELECT DISTINCT cwd FROM command_sequences WHERE project_id IS NULL")
	if seqErr == nil {
		defer func() { _ = seqRows.Close() }()
		var cwds []string
		for seqRows.Next() {
			var d string
			if errScan := seqRows.Scan(&d); errScan == nil && d != "" {
				cwds = append(cwds, d)
			}
		}
		if seqRowsErr := seqRows.Err(); seqRowsErr == nil {
			for _, d := range cwds {
				norm := workspace.Normalize(d)
				pid := ""
				if _, statErr := os.Stat(norm); statErr == nil {
					pid = workspace.ProjectID(workspace.DetectRoot(norm))
				}
				f.mu.Lock()
				_, _ = f.db.ExecContext(ctx, "UPDATE command_sequences SET project_id = ? WHERE cwd = ? AND project_id IS NULL", pid, d)
				f.mu.Unlock()
			}
		}
	}
}

func (f *FrecencyStore) Record(ctx context.Context, cmd, cwd string, exitCode int) error {
	if f == nil {
		return nil
	}
	cmd = strings.TrimSpace(cmd)
	cwd = strings.TrimSpace(cwd)
	if cmd == "" || cwd == "" {
		return nil
	}
	if exitCode != 0 {
		return nil
	}

	normCwd := workspace.Normalize(cwd)
	projectID := workspace.ProjectID(workspace.DetectRoot(normCwd))

	f.mu.Lock()
	defer f.mu.Unlock()

	if ctx == nil {
		ctx = context.Background()
	}
	ctxTimeout, cancel := context.WithTimeout(ctx, 1000*time.Millisecond)
	defer cancel()

	query := `
INSERT INTO history_entries (cmd, cwd, project_id, count, last_used)
VALUES (?, ?, ?, 1, CURRENT_TIMESTAMP)
ON CONFLICT(cmd, cwd) DO UPDATE SET
    project_id = COALESCE(NULLIF(excluded.project_id, ''), project_id),
    count = count + 1,
    last_used = CURRENT_TIMESTAMP;
`
	_, err := f.db.ExecContext(ctxTimeout, query, cmd, cwd, projectID)
	return err
}

func (f *FrecencyStore) RecordTransition(ctx context.Context, prevSkeleton, nextSkeleton, cwd string, nextExitCode int) error {
	if f == nil {
		return nil
	}
	prevSkeleton = strings.TrimSpace(prevSkeleton)
	nextSkeleton = strings.TrimSpace(nextSkeleton)
	cwd = strings.TrimSpace(cwd)
	if prevSkeleton == "" || nextSkeleton == "" || cwd == "" {
		return nil
	}
	if nextExitCode != 0 {
		return nil
	}

	normCwd := workspace.Normalize(cwd)
	f.mu.Lock()
	defer f.mu.Unlock()

	if ctx == nil {
		ctx = context.Background()
	}
	ctxTimeout, cancel := context.WithTimeout(ctx, 1000*time.Millisecond)
	defer cancel()

	query := `
INSERT INTO command_transitions (prev_skeleton, next_skeleton, cwd, count, last_used)
VALUES (?, ?, ?, 1, CURRENT_TIMESTAMP)
ON CONFLICT(prev_skeleton, next_skeleton, cwd) DO UPDATE SET
    count = count + 1,
    last_used = CURRENT_TIMESTAMP;
`
	_, err := f.db.ExecContext(ctxTimeout, query, prevSkeleton, nextSkeleton, normCwd)
	return err
}

var navCommands = map[string]bool{
	"cd":    true,
	"z":     true,
	"zi":    true,
	"j":     true,
	"pushd": true,
	"popd":  true,
}

// avoid repeating directory jumps after arriving at destination
func IsNavCommand(cmd string) bool {
	fields := strings.Fields(cmd)
	if len(fields) == 0 {
		return false
	}
	return navCommands[fields[0]]
}

func (f *FrecencyStore) RecordSequence(ctx context.Context, prevCmd, nextCmd, cwd string, nextExitCode int) error {
	if f == nil {
		return nil
	}
	prevCmd = strings.TrimSpace(prevCmd)
	nextCmd = strings.TrimSpace(nextCmd)
	cwd = strings.TrimSpace(cwd)
	if prevCmd == "" || nextCmd == "" || cwd == "" {
		return nil
	}
	if nextExitCode != 0 {
		return nil
	}
	if IsNavCommand(nextCmd) && strings.EqualFold(prevCmd, nextCmd) {
		return nil
	}

	normCwd := workspace.Normalize(cwd)
	projectID := workspace.ProjectID(workspace.DetectRoot(normCwd))

	f.mu.Lock()
	defer f.mu.Unlock()

	if ctx == nil {
		ctx = context.Background()
	}
	ctxTimeout, cancel := context.WithTimeout(ctx, 1000*time.Millisecond)
	defer cancel()

	query := `
INSERT INTO command_sequences (prev_cmd, next_cmd, cwd, project_id, count, last_used)
VALUES (?, ?, ?, ?, 1, CURRENT_TIMESTAMP)
ON CONFLICT(prev_cmd, next_cmd, cwd) DO UPDATE SET
    project_id = COALESCE(NULLIF(excluded.project_id, ''), project_id),
    count = count + 1,
    last_used = CURRENT_TIMESTAMP;
`
	_, err := f.db.ExecContext(ctxTimeout, query, prevCmd, nextCmd, cwd, projectID)
	return err
}

func (f *FrecencyStore) QuerySequencesWithFallback(ctx context.Context, prevCmd, cwd string) ([]SequenceEntry, bool) {
	if f == nil {
		return nil, false
	}
	prevCmd = strings.TrimSpace(prevCmd)
	cwd = strings.TrimSpace(cwd)
	if prevCmd == "" {
		return nil, false
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	if ctx == nil {
		ctx = context.Background()
	}
	ctxTimeout, cancel := context.WithTimeout(ctx, 1000*time.Millisecond)
	defer cancel()

	var localEntries []SequenceEntry
	rows, err := f.db.QueryContext(ctxTimeout, `
SELECT prev_cmd, next_cmd, cwd, count, last_used
FROM command_sequences
WHERE prev_cmd = ? AND cwd = ? AND count > 0
ORDER BY count DESC
`, prevCmd, cwd)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var prev, next, rCwd string
			var count int
			var lastUsedRaw string
			if err := rows.Scan(&prev, &next, &rCwd, &count, &lastUsedRaw); err == nil {
				t, _ := parseTimestamp(lastUsedRaw)
				localEntries = append(localEntries, SequenceEntry{
					PrevCmd:  prev,
					NextCmd:  next,
					Cwd:      rCwd,
					Count:    count,
					LastUsed: t,
				})
			}
		}
		if rowErr := rows.Err(); rowErr != nil {
			localEntries = nil
		}
	}
	if len(localEntries) > 0 {
		return localEntries, true
	}

	var globalEntries []SequenceEntry
	gRows, gErr := f.db.QueryContext(ctxTimeout, `
SELECT prev_cmd, next_cmd, SUM(count) as total_count, MAX(last_used) as max_last_used
FROM command_sequences
WHERE prev_cmd = ? AND count > 0
GROUP BY next_cmd
ORDER BY total_count DESC
`, prevCmd)
	if gErr == nil {
		defer gRows.Close()
		for gRows.Next() {
			var prev, next string
			var count int
			var lastUsedRaw string
			if err := gRows.Scan(&prev, &next, &count, &lastUsedRaw); err == nil {
				t, _ := parseTimestamp(lastUsedRaw)
				globalEntries = append(globalEntries, SequenceEntry{
					PrevCmd:  prev,
					NextCmd:  next,
					Cwd:      "",
					Count:    count,
					LastUsed: t,
				})
			}
		}
		if gRowErr := gRows.Err(); gRowErr != nil {
			globalEntries = nil
		}
	}
	if len(globalEntries) > 0 {
		return globalEntries, false
	}

	return nil, false
}

func (f *FrecencyStore) GetLatestHistoryEntry(ctx context.Context) (string, string) {
	if f == nil {
		return "", ""
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	var cmd, cwd string
	row := f.db.QueryRowContext(ctx, "SELECT cmd, cwd FROM history_entries ORDER BY last_used DESC LIMIT 1")
	if err := row.Scan(&cmd, &cwd); err == nil {
		return cmd, cwd
	}
	return "", ""
}

type localRow struct {
	cmd   string
	cwd   string
	pid   string
	count int
	last  time.Time
}

type globalRow struct {
	cmd   string
	total int
	last  time.Time
}

func tierOf(rowCwd, rowPID, cwd, pid string) int {
	if rowCwd == cwd {
		return 4
	}
	// descendant or ancestor checks only apply within the same project
	if pid == "" || rowPID != pid {
		return 0
	}
	switch {
	case isUnder(rowCwd, cwd):
		return 3
	case isUnder(cwd, rowCwd):
		return 2
	}
	return 1
}

func isUnder(child, parent string) bool {
	if parent == "" || parent == child {
		return false
	}
	return strings.HasPrefix(child, strings.TrimSuffix(parent, "/")+"/")
}

func rank(local []localRow, global []globalRow, cwd, pid string) []Candidate {
	m := map[string]*Candidate{}
	for _, r := range local {
		t := tierOf(r.cwd, r.pid, cwd, pid)
		c, ok := m[r.cmd]
		if !ok {
			c = &Candidate{Cmd: r.cmd}
			m[r.cmd] = c
		}
		c.Tier = max(c.Tier, t)
		c.Count += r.count
		if r.last.After(c.LastUsed) {
			c.LastUsed = r.last
		}
	}
	for _, g := range global {
		if _, ok := m[g.cmd]; ok {
			continue
		}
		m[g.cmd] = &Candidate{Cmd: g.cmd, Tier: 0, Count: g.total, LastUsed: g.last}
	}
	out := make([]Candidate, 0, len(m))
	for _, c := range m {
		out = append(out, *c)
	}
	sort.Slice(out, func(i, j int) bool {
		a, b := out[i], out[j]
		if a.Tier != b.Tier {
			return a.Tier > b.Tier
		}
		if a.Count != b.Count {
			return a.Count > b.Count
		}
		if !a.LastUsed.Equal(b.LastUsed) {
			return a.LastUsed.After(b.LastUsed)
		}
		return a.Cmd < b.Cmd
	})
	if len(out) > MaxCandidates {
		out = out[:MaxCandidates]
	}
	return out
}

func prefixUpperBound(p string) string {
	if p == "" {
		return ""
	}
	b := []byte(p)
	for i := len(b) - 1; i >= 0; i-- {
		if b[i] < 255 {
			b[i]++
			return string(b[:i+1])
		}
	}
	return ""
}

func (f *FrecencyStore) QueryHistoryCandidates(ctx context.Context, prefix, cwd, pid string) []Candidate {
	if f == nil || prefix == "" {
		return nil
	}
	cwd = strings.TrimSpace(cwd)
	pid = strings.TrimSpace(pid)

	if ctx == nil {
		ctx = context.Background()
	}
	ctxTimeout, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	var local []localRow
	var global []globalRow

	upper := prefixUpperBound(prefix)

	var localSQL string
	var localArgs []interface{}
	if pid != "" {
		if upper != "" {
			localSQL = `
SELECT cmd, cwd, COALESCE(project_id,''), count, last_used
FROM history_entries
WHERE count > 0 AND cwd = ? AND cmd >= ? AND cmd < ? AND instr(cmd, ?) = 1 AND cmd != ?
UNION
SELECT cmd, cwd, COALESCE(project_id,''), count, last_used
FROM history_entries
WHERE count > 0 AND project_id = ? AND cmd >= ? AND cmd < ? AND instr(cmd, ?) = 1 AND cmd != ?
ORDER BY count DESC LIMIT ?
`
			localArgs = []interface{}{cwd, prefix, upper, prefix, prefix, pid, prefix, upper, prefix, prefix, LocalLimit}
		} else {
			localSQL = `
SELECT cmd, cwd, COALESCE(project_id,''), count, last_used
FROM history_entries
WHERE count > 0 AND cwd = ? AND instr(cmd, ?) = 1 AND cmd != ?
UNION
SELECT cmd, cwd, COALESCE(project_id,''), count, last_used
FROM history_entries
WHERE count > 0 AND project_id = ? AND instr(cmd, ?) = 1 AND cmd != ?
ORDER BY count DESC LIMIT ?
`
			localArgs = []interface{}{cwd, prefix, prefix, pid, prefix, prefix, LocalLimit}
		}
	} else {
		if upper != "" {
			localSQL = `
SELECT cmd, cwd, COALESCE(project_id,''), count, last_used
FROM history_entries
WHERE count > 0 AND cwd = ? AND cmd >= ? AND cmd < ? AND instr(cmd, ?) = 1 AND cmd != ?
ORDER BY count DESC LIMIT ?
`
			localArgs = []interface{}{cwd, prefix, upper, prefix, prefix, LocalLimit}
		} else {
			localSQL = `
SELECT cmd, cwd, COALESCE(project_id,''), count, last_used
FROM history_entries
WHERE count > 0 AND cwd = ? AND instr(cmd, ?) = 1 AND cmd != ?
ORDER BY count DESC LIMIT ?
`
			localArgs = []interface{}{cwd, prefix, prefix, LocalLimit}
		}
	}

	var globalSQL string
	var globalArgs []interface{}
	if upper != "" {
		globalSQL = `
SELECT cmd, SUM(count), MAX(last_used)
FROM history_entries
WHERE count > 0 AND cmd >= ? AND cmd < ? AND instr(cmd, ?) = 1 AND cmd != ?
GROUP BY cmd ORDER BY SUM(count) DESC LIMIT ?
`
		globalArgs = []interface{}{prefix, upper, prefix, prefix, GlobalLimit}
	} else {
		globalSQL = `
SELECT cmd, SUM(count), MAX(last_used)
FROM history_entries
WHERE count > 0 AND instr(cmd, ?) = 1 AND cmd != ?
GROUP BY cmd ORDER BY SUM(count) DESC LIMIT ?
`
		globalArgs = []interface{}{prefix, prefix, GlobalLimit}
	}

	func() {
		f.mu.Lock()
		defer f.mu.Unlock()

		if rows, err := f.db.QueryContext(ctxTimeout, localSQL, localArgs...); err == nil {
			defer func() { _ = rows.Close() }()
			for rows.Next() {
				var cmd, rCwd, rPid, lastRaw string
				var count int
				if scanErr := rows.Scan(&cmd, &rCwd, &rPid, &count, &lastRaw); scanErr == nil {
					t, _ := parseTimestamp(lastRaw)
					local = append(local, localRow{
						cmd:   cmd,
						cwd:   rCwd,
						pid:   rPid,
						count: count,
						last:  t,
					})
				}
			}
			if rowErr := rows.Err(); rowErr != nil {
				local = nil
			}
		}

		if gRows, err := f.db.QueryContext(ctxTimeout, globalSQL, globalArgs...); err == nil {
			defer func() { _ = gRows.Close() }()
			for gRows.Next() {
				var cmd, lastRaw string
				var total int
				if scanErr := gRows.Scan(&cmd, &total, &lastRaw); scanErr == nil {
					t, _ := parseTimestamp(lastRaw)
					global = append(global, globalRow{
						cmd:   cmd,
						total: total,
						last:  t,
					})
				}
			}
			if gRowErr := gRows.Err(); gRowErr != nil {
				global = nil
			}
		}
	}()

	return rank(local, global, cwd, pid)
}

func (f *FrecencyStore) QuerySequenceCandidates(ctx context.Context, prevCmd, prefix, cwd, pid string) []Candidate {
	if f == nil || prevCmd == "" {
		return nil
	}
	prevCmd = strings.TrimSpace(prevCmd)
	cwd = strings.TrimSpace(cwd)
	pid = strings.TrimSpace(pid)

	if ctx == nil {
		ctx = context.Background()
	}
	ctxTimeout, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	var local []localRow
	var global []globalRow

	var localSQL string
	var localArgs []interface{}
	var globalSQL string
	var globalArgs []interface{}

	if prefix == "" {
		if pid != "" {
			localSQL = `
SELECT next_cmd, cwd, COALESCE(project_id,''), count, last_used
FROM command_sequences
WHERE count > 0 AND prev_cmd = ? AND cwd = ?
UNION
SELECT next_cmd, cwd, COALESCE(project_id,''), count, last_used
FROM command_sequences
WHERE count > 0 AND prev_cmd = ? AND project_id = ?
ORDER BY count DESC LIMIT ?
`
			localArgs = []interface{}{prevCmd, cwd, prevCmd, pid, LocalLimit}
		} else {
			localSQL = `
SELECT next_cmd, cwd, COALESCE(project_id,''), count, last_used
FROM command_sequences
WHERE count > 0 AND prev_cmd = ? AND cwd = ?
ORDER BY count DESC LIMIT ?
`
			localArgs = []interface{}{prevCmd, cwd, LocalLimit}
		}

		globalSQL = `
SELECT next_cmd, SUM(count), MAX(last_used)
FROM command_sequences
WHERE count > 0 AND prev_cmd = ?
GROUP BY next_cmd ORDER BY SUM(count) DESC LIMIT ?
`
		globalArgs = []interface{}{prevCmd, GlobalLimit}
	} else {
		if pid != "" {
			localSQL = `
SELECT next_cmd, cwd, COALESCE(project_id,''), count, last_used
FROM command_sequences
WHERE count > 0 AND prev_cmd = ? AND cwd = ? AND instr(next_cmd, ?) = 1 AND next_cmd != ?
UNION
SELECT next_cmd, cwd, COALESCE(project_id,''), count, last_used
FROM command_sequences
WHERE count > 0 AND prev_cmd = ? AND project_id = ? AND instr(next_cmd, ?) = 1 AND next_cmd != ?
ORDER BY count DESC LIMIT ?
`
			localArgs = []interface{}{prevCmd, cwd, prefix, prefix, prevCmd, pid, prefix, prefix, LocalLimit}
		} else {
			localSQL = `
SELECT next_cmd, cwd, COALESCE(project_id,''), count, last_used
FROM command_sequences
WHERE count > 0 AND prev_cmd = ? AND cwd = ? AND instr(next_cmd, ?) = 1 AND next_cmd != ?
ORDER BY count DESC LIMIT ?
`
			localArgs = []interface{}{prevCmd, cwd, prefix, prefix, LocalLimit}
		}

		globalSQL = `
SELECT next_cmd, SUM(count), MAX(last_used)
FROM command_sequences
WHERE count > 0 AND prev_cmd = ? AND instr(next_cmd, ?) = 1 AND next_cmd != ?
GROUP BY next_cmd ORDER BY SUM(count) DESC LIMIT ?
`
		globalArgs = []interface{}{prevCmd, prefix, prefix, GlobalLimit}
	}

	func() {
		f.mu.Lock()
		defer f.mu.Unlock()

		if rows, err := f.db.QueryContext(ctxTimeout, localSQL, localArgs...); err == nil {
			defer func() { _ = rows.Close() }()
			for rows.Next() {
				var nextCmd, rCwd, rPid, lastRaw string
				var count int
				if scanErr := rows.Scan(&nextCmd, &rCwd, &rPid, &count, &lastRaw); scanErr == nil {
					if IsNavCommand(nextCmd) && strings.EqualFold(nextCmd, prevCmd) {
						continue
					}
					t, _ := parseTimestamp(lastRaw)
					local = append(local, localRow{
						cmd:   nextCmd,
						cwd:   rCwd,
						pid:   rPid,
						count: count,
						last:  t,
					})
				}
			}
			if rowErr := rows.Err(); rowErr != nil {
				local = nil
			}
		}

		if gRows, err := f.db.QueryContext(ctxTimeout, globalSQL, globalArgs...); err == nil {
			defer func() { _ = gRows.Close() }()
			for gRows.Next() {
				var nextCmd, lastRaw string
				var total int
				if scanErr := gRows.Scan(&nextCmd, &total, &lastRaw); scanErr == nil {
					if IsNavCommand(nextCmd) && strings.EqualFold(nextCmd, prevCmd) {
						continue
					}
					t, _ := parseTimestamp(lastRaw)
					global = append(global, globalRow{
						cmd:   nextCmd,
						total: total,
						last:  t,
					})
				}
			}
			if gRowErr := gRows.Err(); gRowErr != nil {
				global = nil
			}
		}
	}()

	return rank(local, global, cwd, pid)
}

func (f *FrecencyStore) QueryTopHistoryByPrefix(ctx context.Context, prefix, cwd string) string {
	candidates := f.QueryHistoryCandidates(ctx, prefix, cwd, "")
	if len(candidates) > 0 {
		return candidates[0].Cmd
	}
	return ""
}

func (f *FrecencyStore) ScopeCount(ctx context.Context, cmd string) int {
	if f == nil || cmd == "" {
		return 0
	}
	f.mu.Lock()
	defer f.mu.Unlock()

	if ctx == nil {
		ctx = context.Background()
	}
	ctxTimeout, cancel := context.WithTimeout(ctx, 50*time.Millisecond)
	defer cancel()

	var scopes int
	row := f.db.QueryRowContext(ctxTimeout, `
SELECT COUNT(DISTINCT COALESCE(NULLIF(project_id,''), cwd))
FROM history_entries
WHERE count > 0 AND cmd = ?`, cmd)
	_ = row.Scan(&scopes)
	return scopes
}

func (f *FrecencyStore) BootstrapSequences(ctx context.Context, historyPath, defaultCwd string) {
	if f == nil {
		return
	}
	f.mu.Lock()
	var count int
	_ = f.db.QueryRowContext(ctx, "SELECT count(*) FROM command_sequences").Scan(&count)
	f.mu.Unlock()
	if count >= 20 {
		return
	}

	if historyPath == "" {
		home, _ := os.UserHomeDir()
		candidates := []string{
			filepath.Join(home, ".zsh_history"),
			filepath.Join(home, ".bash_history"),
			filepath.Join(home, ".local/share/fish/fish_history"),
		}
		for _, p := range candidates {
			if _, err := os.Stat(p); err == nil {
				historyPath = p
				break
			}
		}
	}
	if historyPath == "" {
		return
	}

	data, err := os.ReadFile(historyPath)
	if err != nil {
		return
	}

	rawLines := strings.Split(string(data), "\n")
	var cmds []string
	for _, l := range rawLines {
		l = strings.TrimSpace(l)
		if l == "" {
			continue
		}
		if strings.HasPrefix(l, ": ") {
			if idx := strings.Index(l, ";"); idx != -1 {
				l = strings.TrimSpace(l[idx+1:])
			}
		}
		if l != "" && len(l) < 300 {
			cmds = append(cmds, l)
		}
	}
	if len(cmds) < 2 {
		return
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	tx, err := f.db.BeginTx(ctx, nil)
	if err != nil {
		return
	}
	stmt, err := tx.PrepareContext(ctx, `
INSERT INTO command_sequences (prev_cmd, next_cmd, cwd, count, last_used)
VALUES (?, ?, ?, 1, CURRENT_TIMESTAMP)
ON CONFLICT(prev_cmd, next_cmd, cwd) DO UPDATE SET
    count = command_sequences.count + 1,
    last_used = CURRENT_TIMESTAMP;
`)
	if err != nil {
		_ = tx.Rollback()
		return
	}
	defer stmt.Close()

	if defaultCwd == "" {
		defaultCwd, _ = os.UserHomeDir()
	}

	start := max(0, len(cmds)-2000)
	for i := start; i < len(cmds)-1; i++ {
		prev := cmds[i]
		next := cmds[i+1]
		if prev != next && !strings.Contains(prev, "\n") && !strings.Contains(next, "\n") {
			_, _ = stmt.ExecContext(ctx, prev, next, defaultCwd)
		}
	}
	_ = tx.Commit()
}

func (f *FrecencyStore) QueryTransitionsWithFallback(ctx context.Context, prevSkeleton, cwd string) ([]TransitionEntry, bool) {
	if f == nil {
		return nil, false
	}
	prevSkeleton = strings.TrimSpace(prevSkeleton)
	cwd = strings.TrimSpace(cwd)
	if prevSkeleton == "" {
		return nil, false
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	if ctx == nil {
		ctx = context.Background()
	}
	ctxTimeout, cancel := context.WithTimeout(ctx, 1000*time.Millisecond)
	defer cancel()

	// Phase 1: Local query with depth fallback
	parts := strings.Fields(prevSkeleton)
	for len(parts) > 0 {
		key := strings.Join(parts, " ")
		var loopEntries []TransitionEntry
		func() {
			rows, err := f.db.QueryContext(ctxTimeout, `
SELECT prev_skeleton, next_skeleton, cwd, count, last_used
FROM command_transitions
WHERE prev_skeleton = ? AND cwd = ? AND count > 0
ORDER BY count DESC
`, key, cwd)
			if err == nil {
				defer rows.Close()
				for rows.Next() {
					var prev, next, rCwd string
					var count int
					var lastUsedRaw string
					if err := rows.Scan(&prev, &next, &rCwd, &count, &lastUsedRaw); err == nil {
						t, _ := parseTimestamp(lastUsedRaw)
						loopEntries = append(loopEntries, TransitionEntry{
							PrevSkeleton: prev,
							NextSkeleton: next,
							Cwd:          rCwd,
							Count:        count,
							LastUsed:     t,
						})
					}
				}
				if rowErr := rows.Err(); rowErr != nil {
					loopEntries = nil
				}
			}
		}()
		if len(loopEntries) > 0 {
			return loopEntries, true
		}
		parts = parts[:len(parts)-1]
	}

	// Phase 2: Global query with depth fallback
	parts = strings.Fields(prevSkeleton)
	for len(parts) > 0 {
		key := strings.Join(parts, " ")
		var loopEntries []TransitionEntry
		func() {
			rows, err := f.db.QueryContext(ctxTimeout, `
SELECT prev_skeleton, next_skeleton, SUM(count) as total_count, MAX(last_used) as max_last_used
FROM command_transitions
WHERE prev_skeleton = ? AND count > 0
GROUP BY next_skeleton
ORDER BY total_count DESC
`, key)
			if err == nil {
				defer rows.Close()
				for rows.Next() {
					var prev, next string
					var count int
					var lastUsedRaw string
					if err := rows.Scan(&prev, &next, &count, &lastUsedRaw); err == nil {
						t, _ := parseTimestamp(lastUsedRaw)
						loopEntries = append(loopEntries, TransitionEntry{
							PrevSkeleton: prev,
							NextSkeleton: next,
							Cwd:          "",
							Count:        count,
							LastUsed:     t,
						})
					}
				}
				if gRowErr := rows.Err(); gRowErr != nil {
					loopEntries = nil
				}
			}
		}()
		if len(loopEntries) > 0 {
			return loopEntries, false
		}
		parts = parts[:len(parts)-1]
	}

	return nil, false
}

func (f *FrecencyStore) RawScore(count int, lastUsed time.Time) float64 {
	if count <= 0 {
		return 0
	}
	age := max(time.Since(lastUsed), 0)

	var weight float64
	switch {
	case age <= time.Hour:
		weight = 100.0
	case age <= 24*time.Hour:
		weight = 50.0
	case age <= 7*24*time.Hour:
		weight = 20.0
	case age <= 30*24*time.Hour:
		weight = 5.0
	default:
		weight = 1.0
	}

	return float64(count) * weight
}

func (f *FrecencyStore) QueryLocal(ctx context.Context, cwd, prefix string, limit int) ([]FrecencyEntry, error) {
	if f == nil {
		return nil, nil
	}
	if limit <= 0 {
		limit = 50
	}
	f.mu.Lock()
	defer f.mu.Unlock()

	if ctx == nil {
		ctx = context.Background()
	}
	ctxTimeout, cancel := context.WithTimeout(ctx, 1000*time.Millisecond)
	defer cancel()

	var rows *sql.Rows
	var err error
	if prefix != "" {
		rows, err = f.db.QueryContext(ctxTimeout, `SELECT cmd, cwd, count, last_used FROM history_entries WHERE cwd = ? AND cmd LIKE ?`, cwd, "%"+prefix+"%")
	} else {
		rows, err = f.db.QueryContext(ctxTimeout, `SELECT cmd, cwd, count, last_used FROM history_entries WHERE cwd = ?`, cwd)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []FrecencyEntry
	for rows.Next() {
		var cmd, rCwd string
		var count int
		var lastUsedRaw string
		if err := rows.Scan(&cmd, &rCwd, &count, &lastUsedRaw); err != nil {
			continue
		}
		t, err := parseTimestamp(lastUsedRaw)
		if err != nil {
			t = time.Now()
		}
		entries = append(entries, FrecencyEntry{
			Cmd:      cmd,
			Cwd:      rCwd,
			Count:    count,
			LastUsed: t,
			RawScore: f.RawScore(count, t),
		})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	sort.SliceStable(entries, func(i, j int) bool {
		return entries[i].RawScore > entries[j].RawScore
	})

	if len(entries) > limit {
		entries = entries[:limit]
	}
	return entries, nil
}

func (f *FrecencyStore) QueryGlobal(ctx context.Context, prefix string, limit int) ([]FrecencyEntry, error) {
	if f == nil {
		return nil, nil
	}
	if limit <= 0 {
		limit = 50
	}
	f.mu.Lock()
	defer f.mu.Unlock()

	if ctx == nil {
		ctx = context.Background()
	}
	ctxTimeout, cancel := context.WithTimeout(ctx, 1000*time.Millisecond)
	defer cancel()

	var rows *sql.Rows
	var err error
	if prefix != "" {
		rows, err = f.db.QueryContext(ctxTimeout, `SELECT cmd, cwd, count, last_used FROM history_entries WHERE cmd LIKE ?`, "%"+prefix+"%")
	} else {
		rows, err = f.db.QueryContext(ctxTimeout, `SELECT cmd, cwd, count, last_used FROM history_entries`)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dedupe := make(map[string]*FrecencyEntry)
	for rows.Next() {
		var cmd, rCwd string
		var count int
		var lastUsedRaw string
		if err := rows.Scan(&cmd, &rCwd, &count, &lastUsedRaw); err != nil {
			continue
		}
		t, err := parseTimestamp(lastUsedRaw)
		if err != nil {
			t = time.Now()
		}
		score := f.RawScore(count, t)
		if existing, found := dedupe[cmd]; found {
			existing.Count += count
			existing.RawScore += score
			if t.After(existing.LastUsed) {
				existing.LastUsed = t
				existing.Cwd = rCwd
			}
		} else {
			dedupe[cmd] = &FrecencyEntry{
				Cmd:      cmd,
				Cwd:      rCwd,
				Count:    count,
				LastUsed: t,
				RawScore: score,
			}
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	var entries []FrecencyEntry
	for _, entry := range dedupe {
		entries = append(entries, *entry)
	}

	sort.SliceStable(entries, func(i, j int) bool {
		return entries[i].RawScore > entries[j].RawScore
	})

	if len(entries) > limit {
		entries = entries[:limit]
	}
	return entries, nil
}

func (f *FrecencyStore) Close() error {
	if f == nil {
		return nil
	}
	f.bgWg.Wait()
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.db != nil {
		return f.db.Close()
	}
	return nil
}

func parseTimestamp(s string) (time.Time, error) {
	if t, err := time.Parse("2006-01-02 15:04:05", s); err == nil {
		return t, nil
	}
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}
	if t, err := time.Parse("2006-01-02 15:04:05.999999999-07:00", s); err == nil {
		return t, nil
	}
	return time.Parse("2006-01-02", s)
}

var (
	globalFrecencyStore *FrecencyStore
	globalFrecencyMu    sync.Mutex
)

func GetFrecencyStore() (*FrecencyStore, error) {
	globalFrecencyMu.Lock()
	defer globalFrecencyMu.Unlock()

	if globalFrecencyStore != nil {
		return globalFrecencyStore, nil
	}

	store, err := NewFrecencyStore("")
	if err != nil {
		return nil, err
	}
	globalFrecencyStore = store
	return globalFrecencyStore, nil
}

// CloseGlobalFrecencyStore safely closes the singleton database connection.
// This is primarily used in testing to prevent goroutine leaks from the DB connectionOpener.
func CloseGlobalFrecencyStore() {
	globalFrecencyMu.Lock()
	defer globalFrecencyMu.Unlock()
	
	if globalFrecencyStore != nil {
		_ = globalFrecencyStore.Close()
		globalFrecencyStore = nil
	}
}
