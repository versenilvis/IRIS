package root

import (
	"testing"

	"github.com/versenilvis/iris/spec"
)

func TestMergeResults(t *testing.T) {

	t.Run("Dedup exact match", func(t *testing.T) {
		// Mock history items that might conflict with specs
		res := MergeResults("git", "spec")
		seen := make(map[string]bool)
		for _, r := range res {
			if seen[r.Cmd] {
				t.Errorf("Duplicate suggestion found: %q", r.Cmd)
			}
			seen[r.Cmd] = true
		}
	})


	t.Run("Limit 100", func(t *testing.T) {
		res := MergeResults("a", "history")
		if len(res) > 100 {
			t.Errorf("Expected max 100 suggestions, got %d", len(res))
		}
	})

	t.Run("AI Suggestion Promotion", func(t *testing.T) {
		aiSugg := &spec.Suggestion{
			Cmd:        "git commit -m \"fix(auth): login bug\"",
			Desc:       "AI suggestion",
			Source:     "ai",
			Confidence: 85,
		}
		SetCurrentAISuggestion(aiSugg)
		defer SetCurrentAISuggestion(nil)

		res := MergeResults("git c", "history")
		if len(res) == 0 {
			t.Fatalf("Expected suggestions, got 0")
		}
		if res[0].Cmd != aiSugg.Cmd {
			t.Errorf("Expected AI suggestion at index 0, got %q (confidence %d)", res[0].Cmd, res[0].Confidence)
		}
	})
}
