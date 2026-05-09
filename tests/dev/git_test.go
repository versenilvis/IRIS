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

func TestGitSuggestions(t *testing.T) {
	// Setup a real git repo in temp dir for testing branch generators
	tmp := t.TempDir()
	oldWd, _ := os.Getwd()
	_ = os.Chdir(tmp)
	defer func() { _ = os.Chdir(oldWd) }()

	ctx := context.Background()
	// Initialize git repo
	_ = exec.CommandContext(ctx, "git", "init").Run()

	_ = exec.CommandContext(ctx, "git", "config", "user.email", "iris-test@example.com").Run() // this is for ci/cd
	_ = exec.CommandContext(ctx, "git", "config", "user.name", "Iris Test").Run()              // this is for ci/cd

	_ = os.WriteFile(filepath.Join(tmp, "file.go"), []byte("package main"), 0644)
	_ = exec.CommandContext(ctx, "git", "add", ".").Run()
	_ = exec.CommandContext(ctx, "git", "commit", "-m", "initial").Run()

	// Create branches
	_ = exec.CommandContext(ctx, "git", "branch", "feature/login").Run()
	_ = exec.CommandContext(ctx, "git", "branch", "dev").Run()
	_ = exec.CommandContext(ctx, "git", "tag", "v1.0").Run()

	// Setup stash
	_ = os.WriteFile(filepath.Join(tmp, "dirty.go"), []byte("dirty"), 0644)
	_ = exec.CommandContext(ctx, "git", "add", ".").Run()
	_ = exec.CommandContext(ctx, "git", "stash").Run()

	t.Run("git top-level", func(t *testing.T) {
		res := core.Lookup("git ")
		if len(res) < 10 {
			t.Errorf("Expected many git subcommands, got %d", len(res))
		}
	})

	t.Run("git tag -d show tags", func(t *testing.T) {
		res := core.Lookup("git tag -d ")
		found := false
		for _, r := range res {
			if strings.Contains(r.Cmd, "v1.0") {
				found = true
			}
		}
		if !found {
			t.Error("git tag -d should suggest existing tags")
		}
	})

	t.Run("git push HEAD options", func(t *testing.T) {
		// git push origin HEAD --force -> should show --force
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

	t.Run("git push upstream options", func(t *testing.T) {
		// git push -u origin -> show branches
		res := core.Lookup("git push -u origin ")
		found := false
		for _, r := range res {
			if strings.Contains(r.Cmd, "dev") {
				found = true
			}
		}
		if !found {
			t.Error("git push -u origin should suggest branches")
		}
	})

	t.Run("git reset options", func(t *testing.T) {
		// git reset --soft origin/main -> should be accepted (just testing lookup doesn't crash)
		_ = core.Lookup("git reset --soft origin/main ")

		// git reset HEAD -> show files
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

	t.Run("git checkout -b no suggest", func(t *testing.T) {
		// git checkout -b -> should NOT suggest branches
		res := core.Lookup("git checkout -b ")
		for _, r := range res {
			if strings.Contains(r.Cmd, "dev") {
				t.Error("git checkout -b should not suggest existing branches")
			}
		}
	})

	t.Run("git switch -c no suggest", func(t *testing.T) {
		res := core.Lookup("git switch -c ")
		for _, r := range res {
			if strings.Contains(r.Cmd, "dev") {
				t.Error("git switch -c should not suggest existing branches")
			}
		}
	})

	t.Run("stash entries", func(t *testing.T) {
		res := core.Lookup("git stash pop ")
		found := false
		for _, r := range res {
			if strings.Contains(r.Cmd, "stash@{0}") {
				found = true
			}
		}
		if !found {
			t.Error("git stash pop should suggest stash@{0}")
		}
	})

	t.Run("not a git repo", func(t *testing.T) {
		emptyDir := t.TempDir()
		_ = os.Chdir(emptyDir)
		// Should not crash
		_ = core.Lookup("git status ")
		_ = os.Chdir(tmp)
	})

	t.Run("active branch filter", func(t *testing.T) {
		// find current branch
		out, _ := exec.CommandContext(ctx, "git", "rev-parse", "--abbrev-ref", "HEAD").Output()
		current := strings.TrimSpace(string(out))

		res := core.Lookup("git checkout ")
		for _, r := range res {
			if strings.Contains(r.Cmd, current) && !strings.Contains(r.Cmd, "remotes/") {
				t.Errorf("Should not suggest active branch '%s'", current)
			}
		}
	})
}
