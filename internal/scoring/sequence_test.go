package scoring

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/versenilvis/iris/spec"
)

func TestFrecencyStore_RecordSequenceAndQuery(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "history.db")
	store, err := NewFrecencyStore(dbPath)
	if err != nil {
		t.Fatalf("failed to create frecency store: %v", err)
	}
	defer store.Close()

	ctx := context.Background()
	cwdA := "/home/user/projectA"
	cwdB := "/home/user/projectB"

	// Record sequence: git add . -> git commit -m "feat"
	_ = store.RecordSequence(ctx, "git add .", "git commit -m \"feat\"", cwdA, 0)
	_ = store.RecordSequence(ctx, "git add .", "git commit -m \"feat\"", cwdA, 0)
	// Failed execution should not increment count
	_ = store.RecordSequence(ctx, "git add .", "git commit -m \"broken\"", cwdA, 1)

	// Local query in cwdA
	entries, isLocal := store.QuerySequencesWithFallback(ctx, "git add .", cwdA)
	if !isLocal {
		t.Errorf("expected isLocal to be true for cwdA")
	}
	if len(entries) == 0 {
		t.Fatalf("expected sequence entries, got 0")
	}
	if entries[0].NextCmd != "git commit -m \"feat\"" || entries[0].Count != 2 {
		t.Errorf("expected count=2 for feat commit, got %+v", entries[0])
	}

	// Global fallback in cwdB
	entriesGlobal, isLocalGlobal := store.QuerySequencesWithFallback(ctx, "git add .", cwdB)
	if isLocalGlobal {
		t.Errorf("expected isLocal to be false for cwdB fallback")
	}
	if len(entriesGlobal) == 0 || entriesGlobal[0].NextCmd != "git commit -m \"feat\"" {
		t.Errorf("expected global fallback to find feat commit, got %+v", entriesGlobal)
	}
}

func TestScore_SequencePredictionPriority(t *testing.T) {
	suggestions := []spec.Suggestion{
		{Cmd: "git status", Source: "spec"},
		{Cmd: "git commit -m \"feat\"", Source: "history"},
	}

	signals := SignalSet{
		Query: "git",
		SequenceEntries: []SequenceEntry{
			{NextCmd: "git commit -m \"feat\"", Count: 10},
		},
		SequenceIsLocal: true,
	}

	scored := Score(suggestions, signals)
	if len(scored) != 2 {
		t.Fatalf("expected 2 scored suggestions, got %d", len(scored))
	}
	if scored[0].Cmd != "git commit -m \"feat\"" {
		t.Errorf("expected predicted sequence command at top, got %s", scored[0].Cmd)
	}
	if scored[0].Breakdown.Transition != 100 {
		t.Errorf("expected transition score 100 for exact sequence match, got %d", scored[0].Breakdown.Transition)
	}
}

func TestIsNavCommand(t *testing.T) {
	for _, cmd := range []string{"z po", "cd /tmp", "pushd dir", "popd", "j proj", "zi"} {
		if !IsNavCommand(cmd) {
			t.Errorf("expected IsNavCommand true for %q", cmd)
		}
	}
	for _, cmd := range []string{"air", "just test", "git status", "ls", "echo cd"} {
		if IsNavCommand(cmd) {
			t.Errorf("expected IsNavCommand false for %q", cmd)
		}
	}
}

func TestFrecencyStore_NavSelfLoopIgnored(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "history.db")
	store, err := NewFrecencyStore(dbPath)
	if err != nil {
		t.Fatalf("failed to create frecency store: %v", err)
	}
	defer store.Close()

	ctx := context.Background()
	cwd := "/home/user/ai-post"

	// navigation self-loop must not be recorded
	_ = store.RecordSequence(ctx, "z po", "z po", cwd, 0)
	_ = store.RecordSequence(ctx, "cd foo", "cd foo", cwd, 0)
	_ = store.RecordSequence(ctx, "z po", "air", cwd, 0)

	cands := store.QuerySequenceCandidates(ctx, "z po", "", cwd, cwd)
	if len(cands) == 0 {
		t.Fatalf("expected candidate, got 0")
	}
	if cands[0].Cmd != "air" {
		t.Fatalf("expected 'air' to be top candidate, got %q", cands[0].Cmd)
	}
	for _, c := range cands {
		if c.Cmd == "z po" {
			t.Fatalf("unexpected nav self-loop 'z po' in candidates")
		}
	}
}

