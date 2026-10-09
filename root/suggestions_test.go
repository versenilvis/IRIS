package root

import (
	"testing"

	"github.com/versenilvis/iris/internal/config"
	"github.com/versenilvis/iris/spec"
	"github.com/versenilvis/iris/spec/alias"
)

func TestMergeResultsExactMatch(t *testing.T) {
	originalConfig := config.Get()
	const name = "irisexacttest"
	originalSpec, existed := spec.Registry[name]
	t.Cleanup(func() {
		config.Init(originalConfig)
		if existed {
			spec.Registry[name] = originalSpec
		} else {
			delete(spec.Registry, name)
		}
	})
	spec.Registry[name] = &spec.Spec{
		Name: name,
		Subcommands: []spec.Subcommand{
			{Name: "runner", Priority: 100},
			{Name: "run", Priority: 1},
		},
	}

	for _, mode := range []string{"spec", "history"} {
		for _, settings := range []struct {
			label            string
			autoExecute      bool
			filterExactMatch config.FilterExactMatchMode
			wantExact        bool
		}{
			{"default-off", false, config.FilterExactMatchAuto, false},
			{"default-on", true, config.FilterExactMatchAuto, true},
			{"force-retain-with-auto-off", false, config.FilterExactMatchOff, true},
			{"force-filter-with-auto-on", true, config.FilterExactMatchOn, false},
			{"force-retain-with-auto-on", true, config.FilterExactMatchOff, true},
			{"force-filter-with-auto-off", false, config.FilterExactMatchOn, false},
		} {
			t.Run(mode+"/"+settings.label, func(t *testing.T) {
				cfg := config.DefaultConfig()
				cfg.Core.AutoExecute = settings.autoExecute
				cfg.Core.FilterExactMatch = settings.filterExactMatch
				cfg.UI.MaxSuggestions = 1
				config.Init(cfg)

				results := MergeResults(name+" run", mode)
				want := name + " runner"
				if settings.wantExact {
					want = name + " run"
				}
				if len(results) != 1 || results[0].Cmd != want {
					t.Fatalf("MergeResults() = %v; want only %q", results, want)
				}
			})
		}
	}
}

func TestMergeResultsExactMatchToolAlias(t *testing.T) {
	original := config.Get()
	t.Cleanup(func() { config.Init(original) })
	alias.Register(&mockGitProvider{})
	t.Cleanup(alias.Reset)
	for _, tt := range []struct {
		name   string
		filter config.FilterExactMatchMode
		want   bool
	}{
		{"auto-preserves-alias", config.FilterExactMatchAuto, true},
		{"explicit-filter", config.FilterExactMatchOn, false},
		{"explicit-retain", config.FilterExactMatchOff, true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.DefaultConfig()
			cfg.Core.FilterExactMatch = tt.filter
			config.Init(cfg)
			found := false
			for _, s := range MergeResults("git co", "spec") {
				if s.Cmd == "git co" {
					found = true
				}
			}
			if found != tt.want {
				t.Fatalf("exact alias present = %v; want %v", found, tt.want)
			}
		})
	}
}

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
		if res[0].Source != "ai" {
			t.Errorf("Expected promoted source 'ai', got %q", res[0].Source)
		}
	})

	t.Run("Tool Alias Preserved When Matching Query", func(t *testing.T) {
		alias.Register(&mockGitProvider{})
		defer alias.Reset()
		res := MergeResults("git co", "default")
		foundCo := false
		for _, r := range res {
			if r.Cmd == "git co" {
				foundCo = true
				break
			}
		}
		if !foundCo {
			t.Errorf("expected 'git co' alias suggestion in results for query 'git co', got %v", res)
		}
	})
}

func TestPrevRecordedCommandState(t *testing.T) {
	origRegistry := spec.Registry
	t.Cleanup(func() {
		spec.Registry = origRegistry
	})

	spec.ResetRegistry()
	spec.Register(&spec.Spec{
		Name: "git",
		Subcommands: []spec.Subcommand{
			{Name: "checkout"},
		},
	})

	setPrevRecordedInfo("git checkout feature-abc", "/repo/dir")
	skel, cwd := getPrevRecordedInfo()
	if skel != "git checkout" || cwd != "/repo/dir" {
		t.Errorf("expected skel='git checkout' and cwd='/repo/dir', got skel=%q, cwd=%q", skel, cwd)
	}
	if gotSkel := getPrevSkeleton(); gotSkel != "git checkout" {
		t.Errorf("expected getPrevSkeleton()='git checkout', got %q", gotSkel)
	}

	setPrevRecordedInfo("", "")
	if gotSkel := getPrevSkeleton(); gotSkel != "" {
		t.Errorf("expected empty getPrevSkeleton() after reset, got %q", gotSkel)
	}
}
