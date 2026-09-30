package tui

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/versenilvis/iris/integration"
	"github.com/versenilvis/iris/internal/ctxcheck"
	"github.com/versenilvis/iris/internal/scoring"
)

var testShells = []string{"zsh", "bash", "fish"}

// Case 1: Repo A has justfile with reload, cd to Repo B without justfile -> no ghost
func TestScopePrediction_RepoWithoutJustfile_NoGhost(t *testing.T) {
	for _, sh := range testShells {
		t.Run(sh, func(t *testing.T) {
			home := wordKeyHome(t)
			repoA := filepath.Join(home, "repoA")
			repoB := filepath.Join(home, "repoB")
			_ = os.MkdirAll(repoA, 0o755)
			_ = os.MkdirAll(repoB, 0o755)

			_ = os.WriteFile(filepath.Join(repoA, "justfile"), []byte("build:\n\techo build\nreload:\n\techo reload\n"), 0o644)

			dbPath := filepath.Join(home, ".local/share/iris/history.db")
			store, err := scoring.NewFrecencyStore(dbPath)
			if err != nil {
				t.Fatal(err)
			}
			ctx := context.Background()
			for range 5 {
				_ = store.Record(ctx, "just reload", repoA, 0)
			}
			_ = store.Close()

			term := startInDirShell(t, home, repoB, sh, "IRIS_CORE_MODE=history")
			defer func() { _ = term.Close() }()

			if err := term.Type("just "); err != nil {
				t.Fatal(err)
			}
			if err := term.WaitStable(2 * time.Second); err != nil {
				t.Fatal(err)
			}
			if got := screen(term); strings.Contains(got, "just reload") {
				t.Fatalf("expected no ghost text for 'just reload' in repoB, got:\n%s", got)
			}
		})
	}
}

// Case 2: Repo B has justfile with test, history has just test only in A -> ghost appears in B
func TestScopePrediction_RepoWithRecipe_GhostAppears(t *testing.T) {
	for _, sh := range testShells {
		t.Run(sh, func(t *testing.T) {
			home := wordKeyHome(t)
			repoA := filepath.Join(home, "repoA")
			repoB := filepath.Join(home, "repoB")
			_ = os.MkdirAll(repoA, 0o755)
			_ = os.MkdirAll(repoB, 0o755)

			_ = os.WriteFile(filepath.Join(repoA, "justfile"), []byte("build:\n\techo build\ntest:\n\techo test\n"), 0o644)
			_ = os.WriteFile(filepath.Join(repoB, "justfile"), []byte("build:\n\techo build-b\ntest:\n\techo test-b\n"), 0o644)

			dbPath := filepath.Join(home, ".local/share/iris/history.db")
			store, err := scoring.NewFrecencyStore(dbPath)
			if err != nil {
				t.Fatal(err)
			}
			ctx := context.Background()
			for range 5 {
				_ = store.Record(ctx, "just test", repoA, 0)
			}
			_ = store.Close()

			term := startInDirShell(t, home, repoB, sh, "IRIS_CORE_MODE=history")
			defer func() { _ = term.Close() }()

			if err := term.Type("just "); err != nil {
				t.Fatal(err)
			}
			if err := term.WaitStable(2 * time.Second); err != nil {
				t.Fatal(err)
			}
			if got := screen(term); !strings.Contains(got, integration.PredictionSymbol) || !strings.Contains(got, "just test") {
				t.Fatalf("expected ghost text for 'just test' in repoB, got:\n%s", got)
			}

			if err := term.SendKeys("\x1b[C"); err != nil {
				t.Fatal(err)
			}
			if err := term.WaitStable(2 * time.Second); err != nil {
				t.Fatal(err)
			}
			if got := promptLine(t, term); got != "just test" {
				t.Fatalf("expected expanded prompt 'just test', got %q", got)
			}
		})
	}
}

// Case 3: Repo root and backend subdir share commands; common parent outside does not leak
func TestScopePrediction_ProjectSubdirSharing_NoParentLeak(t *testing.T) {
	for _, sh := range testShells {
		t.Run(sh, func(t *testing.T) {
			home := wordKeyHome(t)
			parent := filepath.Join(home, "projects")
			repo := filepath.Join(parent, "myrepo")
			backend := filepath.Join(repo, "backend")
			_ = os.MkdirAll(filepath.Join(repo, ".git"), 0o755)
			_ = os.MkdirAll(backend, 0o755)

			dbPath := filepath.Join(home, ".local/share/iris/history.db")
			store, err := scoring.NewFrecencyStore(dbPath)
			if err != nil {
				t.Fatal(err)
			}
			ctx := context.Background()
			for range 5 {
				_ = store.Record(ctx, "mycustomtool run", backend, 0)
				_ = store.Record(ctx, "projectapp start", repo, 0)
			}
			_ = store.Close()

			// isolate each terminal session in a subtest so tuitest cleanup runs before the next starts
			t.Run("repo_root", func(t *testing.T) {
				termRepo := startInDirShell(t, home, repo, sh, "IRIS_CORE_MODE=history")
				if err := termRepo.Type("mycustomtool "); err != nil {
					t.Fatal(err)
				}
				_ = termRepo.WaitStable(2 * time.Second)
				if got := screen(termRepo); !strings.Contains(got, "mycustomtool run") {
					t.Fatalf("expected 'mycustomtool run' in repo root, got:\n%s", got)
				}
			})

			t.Run("backend_subdir", func(t *testing.T) {
				termBackend := startInDirShell(t, home, backend, sh, "IRIS_CORE_MODE=history")
				if err := termBackend.Type("projectapp "); err != nil {
					t.Fatal(err)
				}
				_ = termBackend.WaitStable(2 * time.Second)
				if got := screen(termBackend); !strings.Contains(got, "projectapp start") {
					t.Fatalf("expected 'projectapp start' in backend subdir, got:\n%s", got)
				}
			})

			t.Run("parent_dir", func(t *testing.T) {
				termParent := startInDirShell(t, home, parent, sh, "IRIS_CORE_MODE=history")
				if err := termParent.Type("mycustomtool "); err != nil {
					t.Fatal(err)
				}
				_ = termParent.WaitStable(2 * time.Second)
				if got := screen(termParent); strings.Contains(got, "mycustomtool run") {
					t.Fatalf("expected no leak of 'mycustomtool run' in parent, got:\n%s", got)
				}
			})
		})
	}
}

// Case 4: Delete recipe from justfile in Repo A -> ghost disappears even at Tier 4
func TestScopePrediction_DeletedRecipe_GhostDisappears(t *testing.T) {
	for _, sh := range testShells {
		t.Run(sh, func(t *testing.T) {
			home := wordKeyHome(t)
			repoA := filepath.Join(home, "repoA")
			_ = os.MkdirAll(repoA, 0o755)
			justfilePath := filepath.Join(repoA, "justfile")
			_ = os.WriteFile(justfilePath, []byte("build:\n\techo build\nreload:\n\techo reload\n"), 0o644)

			dbPath := filepath.Join(home, ".local/share/iris/history.db")
			store, err := scoring.NewFrecencyStore(dbPath)
			if err != nil {
				t.Fatal(err)
			}
			ctx := context.Background()
			for range 5 {
				_ = store.Record(ctx, "just reload", repoA, 0)
			}
			_ = store.Close()

			term := startInDirShell(t, home, repoA, sh, "IRIS_CORE_MODE=history")
			defer func() { _ = term.Close() }()

			if err := term.Type("just "); err != nil {
				t.Fatal(err)
			}
			_ = term.WaitStable(2 * time.Second)
			if got := screen(term); !strings.Contains(got, "just reload") {
				t.Fatalf("expected ghost text before deletion, got:\n%s", got)
			}

			// delete recipe from justfile and invalidate cache
			_ = os.WriteFile(justfilePath, []byte("build:\n\techo build\n"), 0o644)
			ctxcheck.InvalidateCache()

			_ = term.SendKeys("\x15") // ctrl+u
			_ = term.WaitStable(1 * time.Second)
			if err := term.Type("just "); err != nil {
				t.Fatal(err)
			}
			_ = term.WaitStable(2 * time.Second)
			if got := screen(term); strings.Contains(got, "just reload") {
				t.Fatalf("expected ghost text to disappear after recipe deletion, got:\n%s", got)
			}
		})
	}
}

// Case 5: cd to existing directory appears (Valid), cd to deleted directory does not (Invalid)
func TestScopePrediction_CdDestination_ValidAndInvalid(t *testing.T) {
	for _, sh := range testShells {
		t.Run(sh, func(t *testing.T) {
			home := wordKeyHome(t)
			targetDir := filepath.Join(home, "target-folder")
			foreignDir := filepath.Join(home, "foreign")
			workDir := filepath.Join(home, "work")
			_ = os.MkdirAll(targetDir, 0o755)
			_ = os.MkdirAll(foreignDir, 0o755)
			_ = os.MkdirAll(workDir, 0o755)

			cdCmd := "cd " + targetDir
			dbPath := filepath.Join(home, ".local/share/iris/history.db")
			store, err := scoring.NewFrecencyStore(dbPath)
			if err != nil {
				t.Fatal(err)
			}
			ctx := context.Background()
			for range 5 {
				_ = store.Record(ctx, cdCmd, foreignDir, 0)
			}
			_ = store.Close()

			term := startInDirShell(t, home, workDir, sh, "IRIS_CORE_MODE=history")
			defer func() { _ = term.Close() }()

			prefix := "cd " + filepath.Join(home, "target-")
			if err := term.Type(prefix); err != nil {
				t.Fatal(err)
			}
			_ = term.WaitStable(2 * time.Second)
			if got := screen(term); !strings.Contains(got, cdCmd) {
				t.Fatalf("expected ghost text for valid cd destination, got:\n%s", got)
			}

			// remove target directory
			_ = os.RemoveAll(targetDir)
			ctxcheck.InvalidateCache()

			_ = term.SendKeys("\x15") // ctrl+u
			_ = term.WaitStable(1 * time.Second)
			if err := term.Type(prefix); err != nil {
				t.Fatal(err)
			}
			_ = term.WaitStable(2 * time.Second)
			if got := screen(term); strings.Contains(got, cdCmd) {
				t.Fatalf("expected no ghost text for deleted cd destination, got:\n%s", got)
			}
		})
	}
}

// Case 6: Compound cd && cmd only appears at Tier 4
func TestScopePrediction_CompoundCd_Tier4Only(t *testing.T) {
	for _, sh := range testShells {
		t.Run(sh, func(t *testing.T) {
			home := wordKeyHome(t)
			dirA := filepath.Join(home, "dirA")
			dirB := filepath.Join(home, "dirB")
			_ = os.MkdirAll(filepath.Join(dirA, "sub"), 0o755)
			_ = os.MkdirAll(filepath.Join(dirB, "sub"), 0o755)

			compound := "cd sub && just test"
			dbPath := filepath.Join(home, ".local/share/iris/history.db")
			store, err := scoring.NewFrecencyStore(dbPath)
			if err != nil {
				t.Fatal(err)
			}
			ctx := context.Background()
			for range 5 {
				_ = store.Record(ctx, compound, dirA, 0)
			}
			_ = store.Close()

			// isolate each directory check in a subtest so pty cleanup finishes cleanly
			t.Run("dirA", func(t *testing.T) {
				termA := startInDirShell(t, home, dirA, sh, "IRIS_CORE_MODE=history")
				if err := termA.Type("cd sub"); err != nil {
					t.Fatal(err)
				}
				_ = termA.WaitStable(2 * time.Second)
				if got := screen(termA); !strings.Contains(got, compound) {
					t.Fatalf("expected compound cd command at Tier 4 in dirA, got:\n%s", got)
				}
			})

			t.Run("dirB", func(t *testing.T) {
				termB := startInDirShell(t, home, dirB, sh, "IRIS_CORE_MODE=history")
				if err := termB.Type("cd sub"); err != nil {
					t.Fatal(err)
				}
				_ = termB.WaitStable(2 * time.Second)
				if got := screen(termB); strings.Contains(got, compound) {
					t.Fatalf("expected compound cd command to be blocked at Tier 0 in dirB, got:\n%s", got)
				}
			})
		})
	}
}

// Case 7: Ghost on empty query (sequence) filtered by same validator logic
func TestScopePrediction_EmptyQuerySequence_Filtered(t *testing.T) {
	for _, sh := range testShells {
		t.Run(sh, func(t *testing.T) {
			ctx := context.Background()

			t.Run("repoB_blocked", func(t *testing.T) {
				homeB := wordKeyHome(t)
				repoB := filepath.Join(homeB, "repoB")
				_ = os.MkdirAll(repoB, 0o755)

				dbPathB := filepath.Join(homeB, ".local/share/iris/history.db")
				storeB, err := scoring.NewFrecencyStore(dbPathB)
				if err != nil {
					t.Fatal(err)
				}
				_ = storeB.RecordSequence(ctx, "echo ready", "just reload", repoB, 0)
				_ = storeB.Close()

				termB := startInDirShell(t, homeB, repoB, sh, "IRIS_CORE_MODE=history")
				_ = termB.Type("echo ready\n")
				_ = termB.WaitStable(2 * time.Second)
				if got := screen(termB); strings.Contains(got, "just reload") {
					t.Fatalf("expected empty-query sequence 'just reload' blocked in repoB, got:\n%s", got)
				}
			})

			t.Run("repoA_allowed", func(t *testing.T) {
				homeA := wordKeyHome(t)
				repoA := filepath.Join(homeA, "repoA")
				_ = os.MkdirAll(repoA, 0o755)
				_ = os.WriteFile(filepath.Join(repoA, "justfile"), []byte("reload:\n\techo reloading\n"), 0o644)

				dbPathA := filepath.Join(homeA, ".local/share/iris/history.db")
				storeA, err := scoring.NewFrecencyStore(dbPathA)
				if err != nil {
					t.Fatal(err)
				}
				_ = storeA.RecordSequence(ctx, "echo ready", "just reload", repoA, 0)
				_ = storeA.Close()

				termA := startInDirShell(t, homeA, repoA, sh, "IRIS_CORE_MODE=history")
				_ = termA.Type("echo ready\n")
				_ = termA.WaitStable(2 * time.Second)
				if got := screen(termA); !strings.Contains(got, "just reload") {
					t.Fatalf("expected empty-query sequence 'just reload' shown in repoA, got:\n%s", got)
				}
			})
		})
	}
}

// Case 8: Timeout on slow directory stat returns within deadline without hanging
func TestScopePrediction_Timeout_ReturnsWithinDeadline(t *testing.T) {
	for _, sh := range testShells {
		t.Run(sh, func(t *testing.T) {
			home := wordKeyHome(t)
			slowDir := filepath.Join(home, "slow-dir")
			foreignDir := filepath.Join(home, "foreign-dir")
			_ = os.MkdirAll(slowDir, 0o755)
			_ = os.MkdirAll(foreignDir, 0o755)
			_ = os.WriteFile(filepath.Join(slowDir, "justfile"), []byte("build:\n\techo build\n"), 0o644)

			dbPath := filepath.Join(home, ".local/share/iris/history.db")
			store, err := scoring.NewFrecencyStore(dbPath)
			if err != nil {
				t.Fatal(err)
			}
			ctx := context.Background()
			_ = store.Record(ctx, "just build", slowDir, 0)
			_ = store.Record(ctx, "just foreign", foreignDir, 0)
			_ = store.Close()

			ctxcheck.SetTestStatDelayHook(func(dir string) {
				if strings.Contains(dir, "slow-dir") {
					time.Sleep(50 * time.Millisecond)
				}
			})
			defer ctxcheck.SetTestStatDelayHook(nil)

			// verify direct validator timeout behavior: cuts off at ~15ms, returns Unknown
			t0 := time.Now()
			v1 := ctxcheck.ValidateWithTimeout("just build", slowDir, ctxcheck.Posix, 15*time.Millisecond)
			d1 := time.Since(t0)
			if v1 != ctxcheck.Unknown {
				t.Fatalf("expected Unknown on slow dir timeout, got %v", v1)
			}
			if d1 > 100*time.Millisecond {
				t.Fatalf("validation did not cut off within reasonable budget for 15ms deadline, took %v", d1)
			}

			// verify slowDir is remembered: second call returns Unknown immediately (< 20ms)
			t1 := time.Now()
			v2 := ctxcheck.ValidateWithTimeout("just build", slowDir, ctxcheck.Posix, 15*time.Millisecond)
			d2 := time.Since(t1)
			if v2 != ctxcheck.Unknown || d2 > 20*time.Millisecond {
				t.Fatalf("expected immediate Unknown from remembered slow dir, took %v with %v", d2, v2)
			}

			// in terminal: Tier 4 gets through (Unknown allowed at Tier 4), Tier 0 does not
			start := time.Now()
			term := startInDirShell(t, home, slowDir, sh, "IRIS_CORE_MODE=history")
			defer func() { _ = term.Close() }()

			if err := term.Type("just b"); err != nil {
				t.Fatal(err)
			}
			if err := term.WaitStable(2 * time.Second); err != nil {
				t.Fatal(err)
			}
			if got := screen(term); !strings.Contains(got, "just build") {
				t.Fatalf("expected Tier 4 ghost in slow dir, got:\n%s", got)
			}

			_ = term.SendKeys("\x15") // ctrl+u
			_ = term.WaitStable(1 * time.Second)
			if err := term.Type("just f"); err != nil {
				t.Fatal(err)
			}
			_ = term.WaitStable(2 * time.Second)
			if got := screen(term); strings.Contains(got, "just foreign") {
				t.Fatalf("expected Tier 0 ghost blocked in slow dir, got:\n%s", got)
			}

			elapsed := time.Since(start)
			if elapsed > 10*time.Second {
				t.Fatalf("render hung or took too long: %v", elapsed)
			}
		})
	}
}

// Case 9: Scope gate threshold: Free command requires ScopeCount >= 3 at Tier 0
func TestScopePrediction_ScopeGate_Threshold3(t *testing.T) {
	for _, sh := range testShells {
		t.Run(sh, func(t *testing.T) {
			home := wordKeyHome(t)
			projA := filepath.Join(home, "projA")
			projB := filepath.Join(home, "projB")
			projC := filepath.Join(home, "projC")
			projD := filepath.Join(home, "projD")
			_ = os.MkdirAll(filepath.Join(projA, ".git"), 0o755)
			_ = os.MkdirAll(filepath.Join(projB, ".git"), 0o755)
			_ = os.MkdirAll(filepath.Join(projC, ".git"), 0o755)
			_ = os.MkdirAll(filepath.Join(projD, ".git"), 0o755)

			dbPath := filepath.Join(home, ".local/share/iris/history.db")
			store, err := scoring.NewFrecencyStore(dbPath)
			if err != nil {
				t.Fatal(err)
			}
			ctx := context.Background()

			// Part 1: command recorded only in projA (scope count = 1)
			freeCmd := "customrunner build"
			for range 5 {
				_ = store.Record(ctx, freeCmd, projA, 0)
			}
			_ = store.Close()

			t.Run("part1_single_scope", func(t *testing.T) {
				// standing in projB (Tier 0, scope = 1 < 3) -> no ghost
				termB := startInDirShell(t, home, projB, sh, "IRIS_CORE_MODE=history")
				if typeErr := termB.Type("customrunner "); typeErr != nil {
					t.Fatal(typeErr)
				}
				_ = termB.WaitStable(2 * time.Second)
				if got := screen(termB); strings.Contains(got, freeCmd) {
					t.Fatalf("expected no ghost for 1-scope command in projB, got:\n%s", got)
				}
			})

			// part 2: record same command in projB and projC -> now scope count = 3
			store, err = scoring.NewFrecencyStore(dbPath)
			if err != nil {
				t.Fatal(err)
			}
			_ = store.Record(ctx, freeCmd, projB, 0)
			_ = store.Record(ctx, freeCmd, projC, 0)
			_ = store.Close()

			t.Run("part2_multi_scope", func(t *testing.T) {
				// standing in projD (Tier 0, scope = 3 >= 3) -> ghost appears
				termD := startInDirShell(t, home, projD, sh, "IRIS_CORE_MODE=history")
				if typeErr := termD.Type("customrunner "); typeErr != nil {
					t.Fatal(typeErr)
				}
				_ = termD.WaitStable(2 * time.Second)
				if got := screen(termD); !strings.Contains(got, freeCmd) {
					t.Fatalf("expected ghost for 3-scope command in projD, got:\n%s", got)
				}
			})

			// part 3: command in 3 non-project directories (tests COALESCE(project_id, cwd))
			dir1 := filepath.Join(home, "standalone1")
			dir2 := filepath.Join(home, "standalone2")
			dir3 := filepath.Join(home, "standalone3")
			dir4 := filepath.Join(home, "standalone4")
			_ = os.MkdirAll(dir1, 0o755)
			_ = os.MkdirAll(dir2, 0o755)
			_ = os.MkdirAll(dir3, 0o755)
			_ = os.MkdirAll(dir4, 0o755)

			standaloneCmd := "dirtool deploy"
			store, err = scoring.NewFrecencyStore(dbPath)
			if err != nil {
				t.Fatal(err)
			}
			_ = store.Record(ctx, standaloneCmd, dir1, 0)
			_ = store.Record(ctx, standaloneCmd, dir2, 0)
			_ = store.Record(ctx, standaloneCmd, dir3, 0)
			_ = store.Close()

			t.Run("part3_standalone", func(t *testing.T) {
				// standing in dir4 (Tier 0, scope = 3 non-project directories) -> ghost appears
				termDir4 := startInDirShell(t, home, dir4, sh, "IRIS_CORE_MODE=history")
				if typeErr := termDir4.Type("dirtool "); typeErr != nil {
					t.Fatal(typeErr)
				}
				_ = termDir4.WaitStable(2 * time.Second)
				if got := screen(termDir4); !strings.Contains(got, standaloneCmd) {
					t.Fatalf("expected ghost for standalone command with 3 cwd scopes, got:\n%s", got)
				}
			})
		})
	}
}
