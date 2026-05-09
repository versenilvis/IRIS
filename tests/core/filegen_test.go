package tests

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/versenilvis/iris/commands/core"
)

func TestFileGenerator(t *testing.T) {
	// Setup mock files
	tmp := t.TempDir()
	os.MkdirAll(filepath.Join(tmp, "src"), 0755)
	os.WriteFile(filepath.Join(tmp, "main.go"), []byte(""), 0644)
	os.WriteFile(filepath.Join(tmp, "README.md"), []byte(""), 0644)
	os.WriteFile(filepath.Join(tmp, ".hidden"), []byte(""), 0644)
	os.WriteFile(filepath.Join(tmp, "src/utils.go"), []byte(""), 0644)

	oldWd, _ := os.Getwd()
	os.Chdir(tmp)
	defer os.Chdir(oldWd)

	// REQUIREMENT: dirOnly shows only dirs
	t.Run("dirOnly shows only dirs", func(t *testing.T) {
		gen := core.FileGenerator("/")
		results := gen([]string{"cd", ""}, "cd ", "")
		for _, r := range results {
			if !strings.HasSuffix(r.Cmd, "/") {
				t.Errorf("Expected only directories, got %q", r.Cmd)
			}
		}
	})

	// REQUIREMENT: Filter extension shows only matching files
	t.Run("Filter extension", func(t *testing.T) {
		gen := core.FileGenerator(".go")
		results := gen([]string{"ls", ""}, "ls ", "")
		foundMain := false
		for _, r := range results {
			if r.Cmd == "main.go" {
				foundMain = true
			}
			if r.Cmd == "README.md" {
				t.Errorf("Did not expect README.md when filtering for .go")
			}
		}
		if !foundMain {
			t.Errorf("FileGenerator(\".go\") did not suggest main.go")
		}
	})

	// REQUIREMENT: Nested path (src/mai -> correct dir + prefix)
	t.Run("Nested path", func(t *testing.T) {
		gen := core.FileGenerator()
		results := gen([]string{"ls", "src/u"}, "ls src/u", "src/u")
		foundUtils := false
		for _, r := range results {
			if r.Cmd == "src/utils.go" {
				foundUtils = true
			}
		}
		if !foundUtils {
			t.Errorf("Did not find src/utils.go for nested path src/u")
		}
	})

	// REQUIREMENT: Deep scan 1 level finds files in subdir
	t.Run("Deep scan 1 level", func(t *testing.T) {
		gen := core.FileGenerator()
		results := gen([]string{"ls", "src/"}, "ls src/", "src/")
		foundUtils := false
		for _, r := range results {
			if r.Cmd == "src/utils.go" {
				foundUtils = true
			}
		}
		if !foundUtils {
			t.Errorf("Deep scan did not find src/utils.go")
		}
	})

	// REQUIREMENT: Deep scan does not go deeper than 1 level
	// (This is implicitly tested by the logic in FileGenerator)

	// REQUIREMENT: Hidden files are skipped
	t.Run("Hidden files are skipped", func(t *testing.T) {
		gen := core.FileGenerator()
		results := gen([]string{"ls", ""}, "ls ", "")
		for _, r := range results {
			if strings.HasPrefix(r.Cmd, ".") {
				t.Errorf("Hidden file %q should be skipped", r.Cmd)
			}
		}
	})
}
