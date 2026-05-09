package core_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/versenilvis/iris/commands/core"
)

func TestFileGenerator(t *testing.T) {
	// Create a temp directory structure
	tmp := t.TempDir()
	
	os.Mkdir(filepath.Join(tmp, ".git"), 0755)
	os.Mkdir(filepath.Join(tmp, "src"), 0755)
	os.Mkdir(filepath.Join(tmp, "docs"), 0755)
	os.WriteFile(filepath.Join(tmp, "src", "main.go"), []byte("package main"), 0644)
	os.WriteFile(filepath.Join(tmp, "src", "utils.go"), []byte("package core"), 0644)
	os.WriteFile(filepath.Join(tmp, "README.md"), []byte("# Readme"), 0644)

	// Helper to change CWD for test
	oldWd, _ := os.Getwd()
	os.Chdir(tmp)
	defer os.Chdir(oldWd)

	t.Run("dirOnly shows only dirs", func(t *testing.T) {
		gen := core.FileGenerator("/")
		results := gen([]string{"cd", ""}, "cd ", "")
		for _, r := range results {
			if r.Desc != "directory" {
				t.Errorf("FileGenerator(\"/\") suggested a non-directory: %s", r.Cmd)
			}
		}
	})

	t.Run("Filter extension", func(t *testing.T) {
		gen := core.FileGenerator(".go")
		results := gen([]string{"go", "src/"}, "go src/", "src/")
		foundMain := false
		for _, r := range results {
			if r.Cmd == "src/main.go" {
				foundMain = true
			}
		}
		if !foundMain {
			t.Errorf("FileGenerator(\".go\") did not suggest src/main.go")
		}
	})

	t.Run("Hidden files are skipped", func(t *testing.T) {
		gen := core.FileGenerator()
		results := gen([]string{"ls", ""}, "ls ", "")
		for _, r := range results {
			if strings.HasPrefix(filepath.Base(r.Cmd), ".") {
				t.Errorf("FileGenerator suggested hidden file: %s", r.Cmd)
			}
		}
	})
}
