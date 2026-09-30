package scoring

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"testing"
	"time"

	"github.com/versenilvis/iris/internal/workspace"
)

func newTestStore(t *testing.T) *FrecencyStore {
	t.Helper()
	dbPath := filepath.Join(t.TempDir(), "history.db")
	store, err := NewFrecencyStore(dbPath)
	if err != nil {
		t.Fatalf("failed to create test store: %v", err)
	}
	t.Cleanup(func() {
		_ = store.Close()
	})
	return store
}

// 1. command in exact cwd beats higher-count command from different project
func TestCandidate_ExactCwdBeatsFarCwd(t *testing.T) {
	store := newTestStore(t)
	ctx := context.Background()

	cwdA := "/home/user/project_a"
	pidA := "/home/user/project_a"
	cwdB := "/home/user/project_b"

	// low count in cwdA
	if err := store.Record(ctx, "git status", cwdA, 0); err != nil {
		t.Fatalf("record failed: %v", err)
	}

	// high count in cwdB
	for i := 0; i < 500; i++ {
		_, err := store.db.ExecContext(ctx, `
INSERT INTO history_entries (cmd, cwd, project_id, count, last_used)
VALUES ('git fetch', ?, ?, 1, CURRENT_TIMESTAMP)
ON CONFLICT(cmd, cwd) DO UPDATE SET count = count + 1`, cwdB, cwdB)
		if err != nil {
			t.Fatalf("insert failed: %v", err)
		}
	}

	candidates := store.QueryHistoryCandidates(ctx, "git", cwdA, pidA)
	if len(candidates) < 2 {
		t.Fatalf("expected at least 2 candidates, got %d", len(candidates))
	}

	if candidates[0].Cmd != "git status" {
		t.Fatalf("expected git status to rank first, got %s (tier %d)", candidates[0].Cmd, candidates[0].Tier)
	}
	if candidates[0].Tier != 4 {
		t.Fatalf("expected tier 4 for exact cwd, got %d", candidates[0].Tier)
	}
	if candidates[1].Tier != 0 {
		t.Fatalf("expected tier 0 for different project, got %d", candidates[1].Tier)
	}
}

// 2. tier hierarchy: child is tier 3, parent is tier 2, sibling is tier 1
func TestCandidate_TierHierarchy(t *testing.T) {
	store := newTestStore(t)
	ctx := context.Background()

	tmpDir := t.TempDir()
	repoRoot := filepath.Join(tmpDir, "repo")
	_ = os.MkdirAll(filepath.Join(repoRoot, ".git"), 0755)
	backend := filepath.Join(repoRoot, "backend")
	_ = os.MkdirAll(backend, 0755)
	frontend := filepath.Join(repoRoot, "frontend")
	_ = os.MkdirAll(frontend, 0755)

	// insert commands in different folders within the same project
	if err := store.Record(ctx, "go test ./...", backend, 0); err != nil {
		t.Fatalf("record failed: %v", err)
	}
	if err := store.Record(ctx, "go build", repoRoot, 0); err != nil {
		t.Fatalf("record failed: %v", err)
	}
	if err := store.Record(ctx, "go run .", frontend, 0); err != nil {
		t.Fatalf("record failed: %v", err)
	}

	// standing at repo root
	fromRoot := store.QueryHistoryCandidates(ctx, "go", repoRoot, repoRoot)
	tierMapRoot := map[string]int{}
	for _, c := range fromRoot {
		tierMapRoot[c.Cmd] = c.Tier
	}
	if tierMapRoot["go build"] != 4 {
		t.Errorf("expected go build tier 4 at repo root, got %d", tierMapRoot["go build"])
	}
	if tierMapRoot["go test ./..."] != 3 {
		t.Errorf("expected go test tier 3 (descendant) at repo root, got %d", tierMapRoot["go test ./..."])
	}
	if tierMapRoot["go run ."] != 3 {
		t.Errorf("expected go run tier 3 (descendant) at repo root, got %d", tierMapRoot["go run ."])
	}

	// standing at repo/backend
	fromBackend := store.QueryHistoryCandidates(ctx, "go", backend, repoRoot)
	tierMapBackend := map[string]int{}
	for _, c := range fromBackend {
		tierMapBackend[c.Cmd] = c.Tier
	}
	if tierMapBackend["go test ./..."] != 4 {
		t.Errorf("expected go test tier 4 at repo/backend, got %d", tierMapBackend["go test ./..."])
	}
	if tierMapBackend["go build"] != 2 {
		t.Errorf("expected go build tier 2 (ancestor) at repo/backend, got %d", tierMapBackend["go build"])
	}
	if tierMapBackend["go run ."] != 1 {
		t.Errorf("expected go run tier 1 (sibling) at repo/backend, got %d", tierMapBackend["go run ."])
	}
}

// 3. standing outside projects (pid empty): child project commands are tier 0 and blocked by gate
func TestCandidate_NonProjectDescendantBlocked(t *testing.T) {
	store := newTestStore(t)
	ctx := context.Background()

	irisCwd := "/home/user/dev/iris"
	devCwd := "/home/user/dev"

	if err := store.Record(ctx, "just reload", irisCwd, 0); err != nil {
		t.Fatalf("record failed: %v", err)
	}

	// standing in ~/dev with no project id
	candidates := store.QueryHistoryCandidates(ctx, "just", devCwd, "")
	if len(candidates) != 1 {
		t.Fatalf("expected 1 candidate, got %d", len(candidates))
	}
	c := candidates[0]
	if c.Tier != 0 {
		t.Fatalf("expected tier 0 when pid is empty, got %d", c.Tier)
	}

	// gate check
	allow := c.Tier > 0 || store.ScopeCount(ctx, c.Cmd) >= GlobalScopeThreshold
	if allow {
		t.Fatalf("expected just reload to be blocked by gate outside project")
	}
}

// 4. merge logic: count accumulates from local rows without duplicating global total
func TestCandidate_MergeCountAndGlobalScopes(t *testing.T) {
	store := newTestStore(t)
	ctx := context.Background()

	tmpDir := t.TempDir()
	repoRoot := filepath.Join(tmpDir, "repo")
	_ = os.MkdirAll(filepath.Join(repoRoot, ".git"), 0755)
	backend := filepath.Join(repoRoot, "backend")
	_ = os.MkdirAll(backend, 0755)
	otherRepo := filepath.Join(tmpDir, "other_repo")
	_ = os.MkdirAll(filepath.Join(otherRepo, ".git"), 0755)

	// child count 5
	for i := 0; i < 5; i++ {
		_ = store.Record(ctx, "make build", backend, 0)
	}
	// cwd count 1
	_ = store.Record(ctx, "make build", repoRoot, 0)
	// other repo count 10
	for i := 0; i < 10; i++ {
		_ = store.Record(ctx, "make build", otherRepo, 0)
	}

	pid := workspace.DetectRoot(repoRoot)
	candidates := store.QueryHistoryCandidates(ctx, "make", repoRoot, pid)
	if len(candidates) != 1 {
		t.Fatalf("expected 1 candidate, got %d", len(candidates))
	}
	c := candidates[0]

	// tier must be 4 (max of child 3 and exact 4)
	if c.Tier != 4 {
		t.Fatalf("expected tier 4, got %d", c.Tier)
	}
	// count must be 6 (local rows 5 + 1), not adding global 10
	if c.Count != 6 {
		t.Fatalf("expected count 6 (5+1), got %d", c.Count)
	}
	// scope count must be tracked from global distinct projects
	scopeCount := store.ScopeCount(ctx, c.Cmd)
	if scopeCount < 2 {
		t.Fatalf("expected scope count >= 2, got %d", scopeCount)
	}
}

// 5. scope count distinction: 3 non-project folders pass gate, single project blocked outside
func TestCandidate_ScopeCountDisambiguation(t *testing.T) {
	store := newTestStore(t)
	ctx := context.Background()

	// ls run in 3 independent non-project directories
	_ = store.Record(ctx, "ls -la", "/tmp/d1", 0)
	_ = store.Record(ctx, "ls -la", "/tmp/d2", 0)
	_ = store.Record(ctx, "ls -la", "/tmp/d3", 0)

	// project-only command
	_ = store.Record(ctx, "just reload", "/home/user/project", 0)

	// check from an unrelated directory /tmp/d4 with pid ""
	lsCandidates := store.QueryHistoryCandidates(ctx, "ls", "/tmp/d4", "")
	if len(lsCandidates) != 1 {
		t.Fatalf("expected 1 ls candidate, got %d", len(lsCandidates))
	}
	lsCand := lsCandidates[0]
	lsScope := store.ScopeCount(ctx, lsCand.Cmd)
	if lsScope < 3 {
		t.Fatalf("expected ls scope count >= 3, got %d", lsScope)
	}
	if lsCand.Tier == 0 && lsScope < GlobalScopeThreshold {
		t.Fatalf("expected ls to pass gate with 3 scopes")
	}

	justCandidates := store.QueryHistoryCandidates(ctx, "just", "/tmp/d4", "")
	if len(justCandidates) != 1 {
		t.Fatalf("expected 1 just candidate, got %d", len(justCandidates))
	}
	justCand := justCandidates[0]
	justScope := store.ScopeCount(ctx, justCand.Cmd)
	if justScope != 1 {
		t.Fatalf("expected just scope count 1, got %d", justScope)
	}
	if justCand.Tier > 0 || justScope >= GlobalScopeThreshold {
		t.Fatalf("expected single-project command to be blocked outside")
	}
}

// 6. count = 0 filtered, literal matching for _ and %, case-sensitive prefix
func TestCandidate_LiteralPrefixAndCaseSensitivity(t *testing.T) {
	store := newTestStore(t)
	ctx := context.Background()
	cwd := "/home/user/test"

	// insert count 0 entry
	_, _ = store.db.ExecContext(ctx, "INSERT INTO history_entries (cmd, cwd, count) VALUES ('test_fail', ?, 0)", cwd)

	_ = store.Record(ctx, "docker_ps", cwd, 0)
	_ = store.Record(ctx, "docker%ps", cwd, 0)
	_ = store.Record(ctx, "docker-compose", cwd, 0)
	_ = store.Record(ctx, "git status", cwd, 0)

	// count = 0 not returned
	if cands := store.QueryHistoryCandidates(ctx, "test", cwd, cwd); len(cands) != 0 {
		t.Fatalf("expected 0 candidates for count=0, got %d", len(cands))
	}

	// docker_ should not match docker-compose via wildcard
	candsUnder := store.QueryHistoryCandidates(ctx, "docker_", cwd, cwd)
	if len(candsUnder) != 1 || candsUnder[0].Cmd != "docker_ps" {
		t.Fatalf("expected only docker_ps for docker_, got %v", candsUnder)
	}

	// docker% should not match docker_ps via wildcard
	candsPercent := store.QueryHistoryCandidates(ctx, "docker%", cwd, cwd)
	if len(candsPercent) != 1 || candsPercent[0].Cmd != "docker%ps" {
		t.Fatalf("expected only docker%%ps for docker%%, got %v", candsPercent)
	}

	// case sensitivity: Git must not match git status
	candsCase := store.QueryHistoryCandidates(ctx, "Git", cwd, cwd)
	if len(candsCase) != 0 {
		t.Fatalf("expected 0 candidates for Git, got %v", candsCase)
	}
}

// 7. sequence candidates: exact cwd beats far cwd, and empty query returns ordered sequence
func TestCandidate_SequenceCandidates(t *testing.T) {
	store := newTestStore(t)
	ctx := context.Background()

	cwdA := "/home/user/repo_a"
	pidA := "/home/user/repo_a"
	cwdB := "/home/user/repo_b"

	_ = store.RecordSequence(ctx, "git add .", "git commit -m \"local\"", cwdA, 0)

	for i := 0; i < 500; i++ {
		_, err := store.db.ExecContext(ctx, `
INSERT INTO command_sequences (prev_cmd, next_cmd, cwd, project_id, count, last_used)
VALUES ('git add .', 'git push origin main', ?, ?, 1, CURRENT_TIMESTAMP)
ON CONFLICT(prev_cmd, next_cmd, cwd) DO UPDATE SET count = count + 1`, cwdB, cwdB)
		if err != nil {
			t.Fatalf("insert failed: %v", err)
		}
	}

	// empty prefix should predict next command with exact cwd ranking highest
	candsEmpty := store.QuerySequenceCandidates(ctx, "git add .", "", cwdA, pidA)
	if len(candsEmpty) < 2 {
		t.Fatalf("expected >= 2 sequence candidates, got %d", len(candsEmpty))
	}
	if candsEmpty[0].Cmd != "git commit -m \"local\"" {
		t.Fatalf("expected local commit to rank first, got %s (tier %d)", candsEmpty[0].Cmd, candsEmpty[0].Tier)
	}

	// prefix query: local commit must rank first
	candsPrefix := store.QuerySequenceCandidates(ctx, "git add .", "git c", cwdA, pidA)
	if len(candsPrefix) == 0 || candsPrefix[0].Cmd != "git commit -m \"local\"" {
		t.Fatalf("expected git commit local to rank first for git c prefix, got %v", candsPrefix)
	}

	// sequence outside project with pid empty is tier 0
	candsOutside := store.QuerySequenceCandidates(ctx, "git add .", "", "/home/user/dev", "")
	if len(candsOutside) > 0 && candsOutside[0].Tier != 0 {
		t.Fatalf("expected tier 0 for sequence outside project, got %d", candsOutside[0].Tier)
	}
}

func TestPrefixUpperBound_EdgeCases(t *testing.T) {
	// case a: empty prefix
	if upper := prefixUpperBound(""); upper != "" {
		t.Fatalf("expected empty upper bound for empty string, got %q", upper)
	}

	// case b: trailing 0xFF byte and all 0xFF bytes
	if upper := prefixUpperBound("abc\xff"); upper != "abd" {
		t.Fatalf("expected 'abd' for 'abc\\xff', got %q", upper)
	}
	if upper := prefixUpperBound("\xff\xff"); upper != "" {
		t.Fatalf("expected '' for all 0xFF bytes, got %q", upper)
	}

	// case c: multi-byte UTF-8 characters
	cafeUpper := prefixUpperBound("café")
	if cafeUpper <= "café" {
		t.Fatalf("expected upper bound > café, got %q", cafeUpper)
	}
	tiengUpper := prefixUpperBound("tiếng")
	if tiengUpper <= "tiếng" {
		t.Fatalf("expected upper bound > tiếng, got %q", tiengUpper)
	}
	jpUpper := prefixUpperBound("こんにちは")
	if jpUpper <= "こんにちは" {
		t.Fatalf("expected upper bound > こんにちは, got %q", jpUpper)
	}

	// end-to-end DB query test with UTF-8 and 0xFF byte
	store := newTestStore(t)
	ctx := context.Background()
	cwd := "/home/user/utf8"

	_ = store.Record(ctx, "tiếng việt nam", cwd, 0)
	_ = store.Record(ctx, "café au lait", cwd, 0)
	_ = store.Record(ctx, "こんにちは世界", cwd, 0)
	_ = store.Record(ctx, "binary\xffspecial", cwd, 0)

	_ = store.Record(ctx, "\xff\xffspecial", cwd, 0)

	candsTieng := store.QueryHistoryCandidates(ctx, "tiếng", cwd, cwd)
	if len(candsTieng) != 1 || candsTieng[0].Cmd != "tiếng việt nam" {
		t.Fatalf("expected 'tiếng việt nam', got %v", candsTieng)
	}

	candsCafe := store.QueryHistoryCandidates(ctx, "café", cwd, cwd)
	if len(candsCafe) != 1 || candsCafe[0].Cmd != "café au lait" {
		t.Fatalf("expected 'café au lait', got %v", candsCafe)
	}

	candsJp := store.QueryHistoryCandidates(ctx, "こんにちは", cwd, cwd)
	if len(candsJp) != 1 || candsJp[0].Cmd != "こんにちは世界" {
		t.Fatalf("expected 'こんにちは世界', got %v", candsJp)
	}

	candsBin := store.QueryHistoryCandidates(ctx, "binary\xff", cwd, cwd)
	if len(candsBin) != 1 || candsBin[0].Cmd != "binary\xffspecial" {
		t.Fatalf("expected 'binary\\xffspecial', got %v", candsBin)
	}

	candsAllFF := store.QueryHistoryCandidates(ctx, "\xff\xff", cwd, cwd)
	if len(candsAllFF) != 1 || candsAllFF[0].Cmd != "\xff\xffspecial" {
		t.Fatalf("expected '\\xff\\xffspecial', got %v", candsAllFF)
	}
}

// 8. benchmark 100k rows with p95 latency under 3ms
func TestCandidate_Benchmark100k(t *testing.T) {
	store := newTestStore(t)
	ctx := context.Background()

	tools := []string{
		"git", "npm", "cargo", "docker", "python", "just", "kubectl", "curl", "make", "node",
		"go", "yarn", "pnpm", "ls", "cd", "vim", "grep", "tar", "ssh", "echo",
	}

	// bulk insert 100k unique rows distributed across tools and projects
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("begin tx failed: %v", err)
	}
	stmt, err := tx.PrepareContext(ctx, `
INSERT INTO history_entries (cmd, cwd, project_id, count, last_used)
VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP)`)
	if err != nil {
		t.Fatalf("prepare failed: %v", err)
	}
	defer func() { _ = stmt.Close() }()

	for i := 0; i < 100000; i++ {
		tool := tools[i%len(tools)]
		cmd := fmt.Sprintf("%s action_%06d arg", tool, i)
		cwd := fmt.Sprintf("/home/user/project_%d/sub", i%50)
		pid := fmt.Sprintf("/home/user/project_%d", i%50)
		if _, execErr := stmt.ExecContext(ctx, cmd, cwd, pid, (i%10)+1); execErr != nil {
			t.Fatalf("exec insert failed: %v", execErr)
		}
	}

	seqStmt, err := tx.PrepareContext(ctx, `
INSERT INTO command_sequences (prev_cmd, next_cmd, cwd, project_id, count, last_used)
VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP)`)
	if err != nil {
		t.Fatalf("prepare seq failed: %v", err)
	}
	defer func() { _ = seqStmt.Close() }()

	for i := 0; i < 20000; i++ {
		tool := tools[i%len(tools)]
		prev := fmt.Sprintf("%s prev_%03d", tool, i%100)
		next := fmt.Sprintf("%s next_%06d", tool, i)
		cwd := fmt.Sprintf("/home/user/project_%d/sub", i%50)
		pid := fmt.Sprintf("/home/user/project_%d", i%50)
		if _, execErr := seqStmt.ExecContext(ctx, prev, next, cwd, pid, (i%5)+1); execErr != nil {
			t.Fatalf("exec seq insert failed: %v", execErr)
		}
	}

	if err := tx.Commit(); err != nil {
		t.Fatalf("commit failed: %v", err)
	}

	// warm up
	_ = store.QueryHistoryCandidates(ctx, "git", "/home/user/project_1/sub", "/home/user/project_1")

	// measure queries across short prefixes: 'g', 'n', 'git', 'npm'
	shortPrefixes := []string{"g", "n", "git", "npm"}
	iterations := len(shortPrefixes) * 25
	latencies := make([]time.Duration, iterations)
	for i := 0; i < iterations; i++ {
		prefix := shortPrefixes[i%len(shortPrefixes)]
		start := time.Now()
		_ = store.QueryHistoryCandidates(ctx, prefix, "/home/user/project_1/sub", "/home/user/project_1")
		latencies[i] = time.Since(start)
	}

	sort.Slice(latencies, func(i, j int) bool {
		return latencies[i] < latencies[j]
	})
	p95 := latencies[int(float64(iterations)*0.95)]
	t.Logf("100k rows short prefix ('g','n','git','npm') QueryHistoryCandidates p95: %v (p50: %v)", p95, latencies[iterations/2])

	// test empty prefix for sequences
	seqLatencies := make([]time.Duration, 50)
	for i := 0; i < 50; i++ {
		start := time.Now()
		_ = store.QuerySequenceCandidates(ctx, "git prev_000", "", "/home/user/project_1/sub", "/home/user/project_1")
		seqLatencies[i] = time.Since(start)
	}
	sort.Slice(seqLatencies, func(i, j int) bool {
		return seqLatencies[i] < seqLatencies[j]
	})
	p95Seq := seqLatencies[int(float64(len(seqLatencies))*0.95)]
	t.Logf("empty prefix QuerySequenceCandidates p95 latency: %v (p50: %v)", p95Seq, seqLatencies[25])
}

func TestBenchmark_RealDB_Comparison(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("no home dir")
	}
	realPath := filepath.Join(home, ".local/share/iris/history.db")
	if _, statErr := os.Stat(realPath); statErr != nil {
		t.Skip("real history.db not found")
	}

	realStore, storeErr := NewFrecencyStore(realPath)
	if storeErr != nil {
		t.Fatalf("open real store: %v", storeErr)
	}
	defer func() { _ = realStore.Close() }()

	ctx := context.Background()
	prefixes := []string{"g", "n", "git", "npm"}

	t.Log("=== REAL DB MEASUREMENTS ===")
	for _, p := range prefixes {
		u := prefixUpperBound(p)

		// 1. old global query (with count distinct)
		start := time.Now()
		func() {
			rOld, qErr := realStore.db.QueryContext(ctx, `
SELECT cmd, SUM(count), MAX(last_used),
       COUNT(DISTINCT COALESCE(NULLIF(project_id,''), cwd))
FROM history_entries
WHERE count > 0 AND cmd >= ? AND cmd < ? AND instr(cmd, ?) = 1 AND cmd != ?
GROUP BY cmd ORDER BY SUM(count) DESC LIMIT 100`, p, u, p, p)
			if qErr == nil {
				defer func() { _ = rOld.Close() }()
				for rOld.Next() {
				}
				_ = rOld.Err()
			}
		}()
		dOld := time.Since(start)

		// 2. new global query (without count distinct)
		var topCmd string
		start = time.Now()
		func() {
			rNew, qErr := realStore.db.QueryContext(ctx, `
SELECT cmd, SUM(count), MAX(last_used)
FROM history_entries
WHERE count > 0 AND cmd >= ? AND cmd < ? AND instr(cmd, ?) = 1 AND cmd != ?
GROUP BY cmd ORDER BY SUM(count) DESC LIMIT 100`, p, u, p, p)
			if qErr == nil {
				defer func() { _ = rNew.Close() }()
				if rNew.Next() {
					var total int
					var lastRaw string
					_ = rNew.Scan(&topCmd, &total, &lastRaw)
				}
				for rNew.Next() {
				}
				_ = rNew.Err()
			}
		}()
		dNew := time.Since(start)

		// 3. lazy scope count for topCmd
		dScope := time.Duration(0)
		if topCmd != "" {
			start = time.Now()
			_ = realStore.ScopeCount(ctx, topCmd)
			dScope = time.Since(start)
		}

		t.Logf("real DB prefix %-4q: old=%v, new=%v, lazy_scope=%v (cmd: %s)", p, dOld, dNew, dScope, topCmd)
	}
}
