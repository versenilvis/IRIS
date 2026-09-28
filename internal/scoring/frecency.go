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

type FrecencyStore struct {
	db *sql.DB
	mu sync.Mutex
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

	store := &FrecencyStore{db: db}
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
    count     INTEGER DEFAULT 1,
    last_used TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(prev_cmd, next_cmd, cwd)
);

CREATE INDEX IF NOT EXISTS idx_sequences_prev_cwd ON command_sequences(prev_cmd, cwd);
`
	_, err := f.db.ExecContext(ctxTimeout, schema)
	return err
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

	f.mu.Lock()
	defer f.mu.Unlock()

	if ctx == nil {
		ctx = context.Background()
	}
	ctxTimeout, cancel := context.WithTimeout(ctx, 1000*time.Millisecond)
	defer cancel()

	var query string
	if exitCode == 0 {
		query = `
INSERT INTO history_entries (cmd, cwd, count, last_used)
VALUES (?, ?, 1, CURRENT_TIMESTAMP)
ON CONFLICT(cmd, cwd) DO UPDATE SET
    count = count + 1,
    last_used = CURRENT_TIMESTAMP;
`
	} else {
		query = `
INSERT INTO history_entries (cmd, cwd, count, last_used)
VALUES (?, ?, 0, CURRENT_TIMESTAMP)
ON CONFLICT(cmd, cwd) DO UPDATE SET
    last_used = CURRENT_TIMESTAMP;
`
	}
	_, err := f.db.ExecContext(ctxTimeout, query, cmd, cwd)
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

	f.mu.Lock()
	defer f.mu.Unlock()

	if ctx == nil {
		ctx = context.Background()
	}
	ctxTimeout, cancel := context.WithTimeout(ctx, 1000*time.Millisecond)
	defer cancel()

	var query string
	if nextExitCode == 0 {
		query = `
INSERT INTO command_transitions (prev_skeleton, next_skeleton, cwd, count, last_used)
VALUES (?, ?, ?, 1, CURRENT_TIMESTAMP)
ON CONFLICT(prev_skeleton, next_skeleton, cwd) DO UPDATE SET
    count = count + 1,
    last_used = CURRENT_TIMESTAMP;
`
	} else {
		query = `
INSERT INTO command_transitions (prev_skeleton, next_skeleton, cwd, count, last_used)
VALUES (?, ?, ?, 0, CURRENT_TIMESTAMP)
ON CONFLICT(prev_skeleton, next_skeleton, cwd) DO UPDATE SET
    last_used = CURRENT_TIMESTAMP;
`
	}
	_, err := f.db.ExecContext(ctxTimeout, query, prevSkeleton, nextSkeleton, cwd)
	return err
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

	f.mu.Lock()
	defer f.mu.Unlock()

	if ctx == nil {
		ctx = context.Background()
	}
	ctxTimeout, cancel := context.WithTimeout(ctx, 1000*time.Millisecond)
	defer cancel()

	var query string
	if nextExitCode == 0 {
		query = `
INSERT INTO command_sequences (prev_cmd, next_cmd, cwd, count, last_used)
VALUES (?, ?, ?, 1, CURRENT_TIMESTAMP)
ON CONFLICT(prev_cmd, next_cmd, cwd) DO UPDATE SET
    count = count + 1,
    last_used = CURRENT_TIMESTAMP;
`
	} else {
		query = `
INSERT INTO command_sequences (prev_cmd, next_cmd, cwd, count, last_used)
VALUES (?, ?, ?, 0, CURRENT_TIMESTAMP)
ON CONFLICT(prev_cmd, next_cmd, cwd) DO UPDATE SET
    last_used = CURRENT_TIMESTAMP;
`
	}
	_, err := f.db.ExecContext(ctxTimeout, query, prevCmd, nextCmd, cwd)
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

func (f *FrecencyStore) QueryTopHistoryByPrefix(ctx context.Context, prefix, cwd string) string {
	if f == nil {
		return ""
	}
	prefix = strings.TrimSpace(prefix)
	if prefix == "" {
		return ""
	}
	f.mu.Lock()
	defer f.mu.Unlock()

	var cmd string
	row := f.db.QueryRowContext(ctx, `
SELECT cmd
FROM history_entries
WHERE (cmd LIKE ? OR cmd LIKE ? OR cmd = ?) AND cmd != ?
ORDER BY CASE WHEN cwd = ? THEN 1 ELSE 0 END DESC, count DESC, last_used DESC
LIMIT 1
`, prefix+" %", prefix+"%", prefix, prefix, cwd)
	if err := row.Scan(&cmd); err == nil {
		return cmd
	}
	return ""
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
