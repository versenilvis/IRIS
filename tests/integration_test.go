package tests

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/versenilvis/iris/commands/core"
	_ "github.com/versenilvis/iris/commands/fs" // Register z command
)

func TestIntegration_ZoxideMultiWord(t *testing.T) {
	// Create mock directory with spaces
	tmp := t.TempDir()
	targetDir := filepath.Join(tmp, "My Awesome Project")
	os.MkdirAll(targetDir, 0755)

	// Mock shell environment
	oldWd, _ := os.Getwd()
	os.Chdir(tmp)
	defer os.Chdir(oldWd)

	// Mock zoxide binary
	mockBinDir := t.TempDir()
	mockZoxide := filepath.Join(mockBinDir, "zoxide")
	script := "#!/bin/sh\necho \"" + targetDir + "\""
	os.WriteFile(mockZoxide, []byte(script), 0755)
	
	oldPath := os.Getenv("PATH")
	os.Setenv("PATH", mockBinDir+string(os.PathListSeparator)+oldPath)
	defer os.Setenv("PATH", oldPath)

	t.Run("z matches multi-word folder without quotes", func(t *testing.T) {
		// Simulating user typing "z My Awe"
		input := "z My Awe"
		results := core.Lookup(input)
		
		found := false
		for _, r := range results {
			// Expected result should be "z My Awesome Project/"
			if strings.Contains(r.Cmd, "My Awesome Project") {
				found = true
				// Check for word duplication bug: "z My My Awesome..."
				if strings.Contains(r.Cmd, "My My") {
					t.Errorf("Word duplication detected in suggestion: %s", r.Cmd)
				}
			}
		}
		if !found {
			t.Errorf("Could not find 'My Awesome Project' in suggestions for %q", input)
		}
	})
}
