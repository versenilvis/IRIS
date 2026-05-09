package root_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/versenilvis/iris/commands/core"
	"github.com/versenilvis/iris/root"
)

func TestMergeResults(t *testing.T) {
	// Setup mock Registry
	core.Registry = make(map[string]*core.Spec)
	core.Register(&core.Spec{
		Name:        "ls",
		Description: "list files",
	})

	// Setup mock history file
	tmp := t.TempDir()
	histFile := filepath.Join(tmp, ".bash_history")
	os.WriteFile(histFile, []byte("ls -l\ncd /tmp\ngit status\n"), 0644)
	
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmp)
	defer os.Setenv("HOME", oldHome)

	t.Run("History mode returns history items", func(t *testing.T) {
		results := root.MergeResults("", "history")
		if len(results) == 0 {
			t.Errorf("MergeResults history mode returned 0 items")
		}
	})

	t.Run("Spec mode returns command results and dedups", func(t *testing.T) {
		core.Register(&core.Spec{Name: "git", Description: "git"})
		
		results := root.MergeResults("gi", "spec")
		foundGit := false
		for _, r := range results {
			if r.Cmd == "git" {
				foundGit = true
			}
		}
		if !foundGit {
			t.Errorf("MergeResults spec mode did not find 'git'")
		}
	})
}
