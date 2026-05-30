package tests

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/versenilvis/iris/commands/core"
	_ "github.com/versenilvis/iris/commands/dev"
)

// setupGitRepo creates a real git repo in a temp dir with:
// - local branches: main (HEAD), dev, feature/login, stable
// - remote branches: origin/main, origin/dev (written directly to .git/refs)
// - a tag: v1.0
// - a stash entry
func setupGitRepo(t *testing.T) (tmp string, cleanup func()) {
	t.Helper()

	tmp = t.TempDir()
	ctx := context.Background()

	run := func(args ...string) {
		t.Helper()
		out, err := exec.CommandContext(ctx, args[0], args[1:]...).CombinedOutput()
		if err != nil {
			t.Logf("git cmd %v: %s", args, out)
		}
	}

	run("git", "-C", tmp, "init", "--initial-branch=main")

	// use fallback for older git that doesn't support --initial-branch
	if _, err := os.Stat(filepath.Join(tmp, ".git", "refs", "heads", "main")); err != nil {
		run("git", "-C", tmp, "init")
	}

	run("git", "-C", tmp, "config", "user.email", "iris-test@example.com") // this is for ci/cd
	run("git", "-C", tmp, "config", "user.name", "Iris Test")              // this is for ci/cd

	// initial commit so branches can be created
	if err := os.WriteFile(filepath.Join(tmp, "file.go"), []byte("package main"), 0644); err != nil {
		t.Fatal(err)
	}
	run("git", "-C", tmp, "add", ".")
	run("git", "-C", tmp, "commit", "-m", "initial")

	// local branches (incl. slash branch to test tokenization)
	run("git", "-C", tmp, "branch", "dev")
	run("git", "-C", tmp, "branch", "stable")
	run("git", "-C", tmp, "branch", "feature/login")

	// tag
	run("git", "-C", tmp, "tag", "v1.0")

	// add a real remote in config
	run("git", "-C", tmp, "remote", "add", "origin", "https://github.com/versenilvis/iris.git")

	// write fake remote refs directly (no need for actual remote server)
	for _, ref := range []string{"main", "dev"} {
		dir := filepath.Join(tmp, ".git", "refs", "remotes", "origin")
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}
		// point them to the same commit as HEAD for simplicity
		headBytes, err := os.ReadFile(filepath.Join(tmp, ".git", "refs", "heads", "main"))
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, ref), headBytes, 0644); err != nil {
			t.Fatal(err)
		}
	}

	// stash entry
	if err := os.WriteFile(filepath.Join(tmp, "dirty.go"), []byte("dirty"), 0644); err != nil {
		t.Fatal(err)
	}
	run("git", "-C", tmp, "add", ".")
	run("git", "-C", tmp, "stash")

	// chdir into repo so generators can run git commands
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}

	cleanup = func() { _ = os.Chdir(oldWd) }
	return tmp, cleanup
}

func TestGitSuggestions(t *testing.T) {
	tmp, cleanup := setupGitRepo(t)
	defer cleanup()

	t.Run("git top-level", func(t *testing.T) {
		res := core.Lookup("git ")
		if len(res) < 10 {
			t.Errorf("expected many git subcommands, got %d", len(res))
		}
	})

	t.Run("tag -d shows tags", func(t *testing.T) {
		res := core.Lookup("git tag -d ")
		found := false
		for _, r := range res {
			if strings.Contains(r.Cmd, "v1.0") {
				found = true
			}
		}
		if !found {
			t.Error("git tag -d should suggest v1.0")
		}
	})

	t.Run("push HEAD options", func(t *testing.T) {
		res := core.Lookup("git push origin HEAD --")
		found := false
		for _, r := range res {
			if strings.Contains(r.Cmd, "--force") {
				found = true
			}
		}
		if !found {
			t.Error("git push origin HEAD -- should suggest --force")
		}
	})

	t.Run("push -u origin suggests branches", func(t *testing.T) {
		res := core.Lookup("git push -u origin ")
		found := false
		for _, r := range res {
			if strings.Contains(r.Cmd, "dev") {
				found = true
			}
		}
		if !found {
			t.Errorf("git push -u origin should suggest branches, got: %v", res)
		}
	})

	t.Run("push origin suggests active branch", func(t *testing.T) {
		ctx := context.Background()
		out, err := exec.CommandContext(ctx, "git", "rev-parse", "--abbrev-ref", "HEAD").Output()
		if err != nil {
			t.Skip("can't determine HEAD branch")
		}
		activeBranch := strings.TrimSpace(string(out))
		res := core.Lookup("git push origin ")
		found := false
		for _, r := range res {
			parts := strings.Fields(r.Cmd)
			if len(parts) > 0 && parts[len(parts)-1] == activeBranch {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("git push origin should suggest active branch '%s'", activeBranch)
		}
		if len(res) > 0 {
			parts := strings.Fields(res[0].Cmd)
			if len(parts) == 0 || parts[len(parts)-1] != activeBranch {
				t.Errorf("expected active branch '%s' to be first suggestion, got: %s", activeBranch, res[0].Cmd)
			}
		}
	})

	t.Run("push origin no duplicate branches", func(t *testing.T) {
		res := core.Lookup("git push origin ")
		seen := make(map[string]int)
		for _, r := range res {
			parts := strings.Fields(r.Cmd)
			if len(parts) == 0 {
				continue
			}
			branch := parts[len(parts)-1]
			seen[branch]++
			if seen[branch] > 1 {
				t.Errorf("duplicate branch suggestion: %s", branch)
			}
		}
	})


	t.Run("branch with slash is suggested correctly", func(t *testing.T) {
		res := core.Lookup("git checkout ")
		found := false
		for _, r := range res {
			if strings.Contains(r.Cmd, "feature/login") {
				found = true
			}
		}
		if !found {
			t.Error("git checkout should suggest 'feature/login'")
		}
	})

	t.Run("remote branches suggested for push", func(t *testing.T) {
		res := core.Lookup("git push origin ")
		cmdStr := ""
		for _, r := range res {
			cmdStr += r.Cmd + " "
		}
		// should have at least dev or main from branch list
		if !strings.Contains(cmdStr, "dev") && !strings.Contains(cmdStr, "main") {
			t.Errorf("git push origin should suggest local branches, got: %s", cmdStr)
		}
	})

	t.Run("active branch not suggested for checkout", func(t *testing.T) {
		// find actual active branch
		ctx := context.Background()
		out, err := exec.CommandContext(ctx, "git", "rev-parse", "--abbrev-ref", "HEAD").Output()
		if err != nil {
			t.Skip("can't determine HEAD branch")
		}
		activeBranch := strings.TrimSpace(string(out))

		res := core.Lookup("git checkout ")
		for _, r := range res {
			// the suggestion should not contain the active branch as a standalone word
			parts := strings.Fields(r.Cmd)
			for _, p := range parts {
				if p == activeBranch {
					t.Errorf("git checkout should not suggest active branch '%s', got: %s", activeBranch, r.Cmd)
				}
			}
		}
	})

	t.Run("checkout -b no suggest", func(t *testing.T) {
		res := core.Lookup("git checkout -b ")
		for _, r := range res {
			if strings.Contains(r.Cmd, "dev") {
				t.Error("git checkout -b should not suggest existing branches")
			}
		}
	})

	t.Run("switch -c no suggest", func(t *testing.T) {
		res := core.Lookup("git switch -c ")
		for _, r := range res {
			if strings.Contains(r.Cmd, "dev") {
				t.Error("git switch -c should not suggest existing branches")
			}
		}
	})

	t.Run("stash variants suggest entries", func(t *testing.T) {
		for _, cmd := range []string{"apply", "drop", "pop"} {
			res := core.Lookup("git stash " + cmd + " ")
			found := false
			for _, r := range res {
				if strings.Contains(r.Cmd, "stash@{0}") {
					found = true
				}
			}
			if !found {
				t.Errorf("git stash %s should suggest stash@{0}", cmd)
			}
		}
	})

	t.Run("remote subcommands suggest remotes", func(t *testing.T) {
		for _, cmd := range []string{"remove", "rename", "set-url"} {
			res := core.Lookup("git remote " + cmd + " ")
			found := false
			for _, r := range res {
				// origin is our fake remote
				if strings.Contains(r.Cmd, "origin") {
					found = true
				}
			}
			if !found {
				t.Errorf("git remote %s should suggest origin", cmd)
			}
		}
	})

	t.Run("not a git repo no crash", func(t *testing.T) {
		emptyDir := t.TempDir()
		_ = os.Chdir(emptyDir)
		defer func() { _ = os.Chdir(tmp) }()
		_ = core.Lookup("git status ")
	})

	t.Run("reset options", func(t *testing.T) {
		_ = core.Lookup("git reset --soft origin/main ")
		res := core.Lookup("git reset HEAD ")
		found := false
		for _, r := range res {
			if strings.Contains(r.Cmd, "file.go") {
				found = true
			}
		}
		if !found {
			t.Error("git reset HEAD should suggest file.go")
		}
	})
}
