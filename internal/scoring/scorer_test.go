package scoring

import (
	"testing"
	"time"

	"github.com/versenilvis/iris/internal/workspace"
	"github.com/versenilvis/iris/spec"
)

func TestScore_GitInitAndStatusInGitRepo(t *testing.T) {
	suggestions := []spec.Suggestion{
		{Cmd: "git init", Source: "spec"},
		{Cmd: "git status -s", Source: "spec"},
	}
	signals := SignalSet{
		Workspace: workspace.WorkspaceInfo{HasGit: true},
		Query:     "git",
	}

	scored := Score(suggestions, signals)
	if len(scored) != 2 {
		t.Fatalf("expected 2 scored suggestions, got %d", len(scored))
	}
	if scored[0].Cmd != "git status -s" {
		t.Errorf("expected 'git status -s' at top when inside git repo, got %s", scored[0].Cmd)
	}
	if scored[1].Cmd != "git init" {
		t.Errorf("expected 'git init' at bottom, got %s", scored[1].Cmd)
	}
	if scored[1].Breakdown.ContextBonus != -50 {
		t.Errorf("expected -50 penalty for git init, got %d", scored[1].Breakdown.ContextBonus)
	}
}

func TestScore_NormalizedFrecency(t *testing.T) {
	suggestions := []spec.Suggestion{
		{Cmd: "ls -la", Source: "history"},
		{Cmd: "git push", Source: "history"},
	}

	now := time.Now()
	signals := SignalSet{
		LocalFrecency: []FrecencyEntry{
			{Cmd: "ls -la", RawScore: 25000.0, LastUsed: now.Add(-30 * 24 * time.Hour)},
			{Cmd: "git push", RawScore: 500.0, LastUsed: now},
		},
	}

	scored := Score(suggestions, signals)
	var lsBreakdown, pushBreakdown ScoreBreakdown
	for _, s := range scored {
		switch s.Cmd {
		case "ls -la":
			lsBreakdown = s.Breakdown
		case "git push":
			pushBreakdown = s.Breakdown
		}
	}

	if lsBreakdown.Frecency != 100 {
		t.Errorf("expected max raw score to normalize to 100, got %d", lsBreakdown.Frecency)
	}
	if pushBreakdown.Frecency <= 0 || pushBreakdown.Frecency > 100 {
		t.Errorf("expected normalized frecency in (0, 100], got %d", pushBreakdown.Frecency)
	}
}

func TestScore_PrefixOverFuzzyMatch(t *testing.T) {
	suggestions := []spec.Suggestion{
		{Cmd: "make build", Source: "spec"}, // fuzzy/contains for 'bl'
		{Cmd: "block", Source: "spec"},      // prefix exact for 'bl'
	}
	signals := SignalSet{Query: "bl"}

	scored := Score(suggestions, signals)
	if len(scored) != 2 {
		t.Fatalf("expected 2 scored suggestions, got %d", len(scored))
	}
	if scored[0].Cmd != "block" {
		t.Errorf("expected prefix match 'block' to outscore fuzzy 'make build', got %s", scored[0].Cmd)
	}
}

func TestScore_AISuggestionConfidence(t *testing.T) {
	suggestions := []spec.Suggestion{
		{Cmd: "npm run custom-script", Source: "ai", Confidence: 85},
		{Cmd: "npm help", Source: "history"},
	}
	signals := SignalSet{
		Workspace: workspace.WorkspaceInfo{HasNodeProject: true},
		Query:     "npm",
	}

	scored := Score(suggestions, signals)
	if len(scored) < 1 {
		t.Fatalf("expected scored suggestions")
	}
	if scored[0].Cmd != "npm run custom-script" {
		t.Errorf("expected high-confidence AI suggestion with context bonus at top, got %s", scored[0].Cmd)
	}
}

func TestScore_UnsortedHistorySorting(t *testing.T) {
	suggestions := []spec.Suggestion{
		{Cmd: "cmdA", Source: "history"},
		{Cmd: "cmdB", Source: "history"},
		{Cmd: "cmdC", Source: "history"},
	}
	signals := SignalSet{
		LocalFrecency: []FrecencyEntry{
			{Cmd: "cmdC", RawScore: 100.0},
			{Cmd: "cmdA", RawScore: 50.0},
			{Cmd: "cmdB", RawScore: 10.0},
		},
	}

	scored := Score(suggestions, signals)
	if len(scored) != 3 {
		t.Fatalf("expected 3 items, got %d", len(scored))
	}
	if scored[0].Cmd != "cmdC" || scored[1].Cmd != "cmdA" || scored[2].Cmd != "cmdB" {
		t.Errorf("expected cmdC > cmdA > cmdB based on frecency, got %s, %s, %s", scored[0].Cmd, scored[1].Cmd, scored[2].Cmd)
	}
}

func TestBasePriorityFor_HistoryWithConfidence(t *testing.T) {
	s1 := spec.Suggestion{Source: "history"}
	if p := basePriorityFor(s1); p != 40 {
		t.Errorf("expected default history priority 40 when confidence unset, got %d", p)
	}

	s2 := spec.Suggestion{Source: "history", Confidence: 85}
	if p := basePriorityFor(s2); p != 85 {
		t.Errorf("expected history priority 85 when confidence is 85, got %d", p)
	}

	s3 := spec.Suggestion{Source: "history", Confidence: 150}
	if p := basePriorityFor(s3); p != 100 {
		t.Errorf("expected capped history priority 100 when confidence is > 100, got %d", p)
	}
}
