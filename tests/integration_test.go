package tests

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/versenilvis/iris/spec"
	_ "github.com/versenilvis/iris/commands"
)

func TestIntegration_ZoxideMultiWord(t *testing.T) {
	// Create mock directory with spaces
	tmp := t.TempDir()
	targetDir := filepath.Join(tmp, "My Awesome Project")
	_ = os.MkdirAll(targetDir, 0755)

	// Mock shell environment
	oldWd, _ := os.Getwd()
	_ = os.Chdir(tmp)
	defer func() { _ = os.Chdir(oldWd) }()

	// Mock zoxide binary
	mockBinDir := t.TempDir()
	mockZoxide := filepath.Join(mockBinDir, "zoxide")
	// Create mock zoxide binary
	script := "#!/bin/sh\necho \"" + targetDir + "\""
	_ = os.WriteFile(mockZoxide, []byte(script), 0755)
	t.Setenv("PATH", mockBinDir+string(os.PathListSeparator)+os.Getenv("PATH"))

	t.Run("z matches multi-word folder without quotes", func(t *testing.T) {
		// Simulating user typing "z My Awe"
		input := "z My Awe"
		results := spec.Lookup(input)
		
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
