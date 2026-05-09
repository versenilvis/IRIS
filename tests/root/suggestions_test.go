package tests

import (
	"testing"

	"github.com/versenilvis/iris/root"
)

func TestMergeResults(t *testing.T) {
	// REQUIREMENT: Dedup exact match
	t.Run("Dedup exact match", func(t *testing.T) {
		// Mock history items that might conflict with specs
		res := root.MergeResults("git", "spec")
		seen := make(map[string]bool)
		for _, r := range res {
			if seen[r.Cmd] {
				t.Errorf("Duplicate suggestion found: %q", r.Cmd)
			}
			seen[r.Cmd] = true
		}
	})

	// REQUIREMENT: Limit 100
	t.Run("Limit 100", func(t *testing.T) {
		res := root.MergeResults("a", "history")
		if len(res) > 100 {
			t.Errorf("Expected max 100 suggestions, got %d", len(res))
		}
	})
}
