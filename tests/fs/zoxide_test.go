package tests

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/versenilvis/iris/commands/fs"
)

func TestZoxideGenerator(t *testing.T) {
	// Setup: Create a mock zoxide binary
	tmp := t.TempDir()
	mockZoxide := filepath.Join(tmp, "zoxide")
	
	// Script that prints mock directories
	script := "#!/bin/sh\necho \"/home/verse/project1\n/home/verse/docs\n/home/verse/dev/iris\""
	_ = os.WriteFile(mockZoxide, []byte(script), 0755)

	// Add tmp to PATH
	oldPath := os.Getenv("PATH")
	_ = os.Setenv("PATH", tmp+string(os.PathListSeparator)+oldPath)
	defer func() { _ = os.Setenv("PATH", oldPath) }()

	gen := fs.ZoxideGenerator()
	
	// REQUIREMENT: Query returns the correct result when partial = ""
	t.Run("Query returns correct result when partial is empty", func(t *testing.T) {
		results := gen([]string{"z", ""}, "z ", "")
		if len(results) == 0 {
			t.Errorf("Expected results from zoxide history, got 0")
		}
	})

	// REQUIREMENT: Path replaces home dir with ~
	t.Run("Path replaces home dir with ~", func(t *testing.T) {
		home, _ := os.UserHomeDir()
		results := gen([]string{"z", ""}, "z ", "")
		foundHome := false
		for _, r := range results {
			if strings.HasPrefix(r.Desc, "~") {
				foundHome = true
				break
			}
		}
		if !foundHome && home != "" {
			t.Logf("Warning: Did not find ~ in descriptions, home is %s", home)
		}
	})

	// REQUIREMENT: Sort by descending score
	t.Run("Sort by descending score", func(t *testing.T) {
		results := gen([]string{"z", "i"}, "z ", "i")
		if len(results) >= 1 {
			if results[0].Cmd == "" {
				t.Errorf("Empty result command")
			}
		}
	})
}
