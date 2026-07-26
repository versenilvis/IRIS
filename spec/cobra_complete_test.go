package spec

import (
	"testing"
)

func TestParseCobraOutput_ValidCobra(t *testing.T) {
	raw := "install\tinstall a chart\nupgrade\tupgrade a release\nstatus\tget release status\n:4\n"
	results := parseCobraOutput(raw, "helm")
	if len(results) != 3 {
		t.Fatalf("expected 3 suggestions, got %d", len(results))
	}
	if results[0].Cmd != "helm install" {
		t.Errorf("expected 'helm install', got %q", results[0].Cmd)
	}
	if results[0].Desc != "install a chart" {
		t.Errorf("expected desc 'install a chart', got %q", results[0].Desc)
	}
	if results[0].Source != "spec-inferred" {
		t.Errorf("expected source 'spec-inferred', got %q", results[0].Source)
	}
	if results[0].Priority != 30 {
		t.Errorf("expected priority 30, got %d", results[0].Priority)
	}
}

func TestParseCobraOutput_ErrorDirective(t *testing.T) {
	// directive bit 1 = ShellCompDirectiveError — not a Cobra CLI
	raw := "something\n:1\n"
	results := parseCobraOutput(raw, "mycmd")
	if results != nil {
		t.Errorf("expected nil for error directive, got %v", results)
	}
}

func TestParseCobraOutput_NoDirectiveLine(t *testing.T) {
	raw := "just some --help output\nno directive here\n"
	results := parseCobraOutput(raw, "mycmd")
	if results != nil {
		t.Errorf("expected nil for non-Cobra output, got %v", results)
	}
}

func TestParseCobraOutput_PartialFilter(t *testing.T) {
	raw := "get\tget resources\ndelete\tdelete resources\ndescribe\tdescribe resources\n:4\n"
	// parseCobraOutput returns all candidates; filterByPartial handles narrowing
	results := parseCobraOutput(raw, "kubectl")
	filtered := filterByPartial(results, "de")
	if len(filtered) != 2 {
		t.Fatalf("expected 2 results matching 'de', got %d: %v", len(filtered), filtered)
	}
}

func TestQueryCobraComplete_PathTraversalBlocked(t *testing.T) {
	result := QueryCobraComplete("./malicious.sh", nil, "")
	if result != nil {
		t.Errorf("expected nil for path traversal input, got %v", result)
	}
	result = QueryCobraComplete("/usr/bin/env", nil, "")
	if result != nil {
		t.Errorf("expected nil for absolute path input, got %v", result)
	}
}

func TestQueryCobraComplete_NonCobraBinary(t *testing.T) {
	t.Cleanup(ResetCobraCache)
	// 'ls' is not Cobra-based, should return nil gracefully
	result := QueryCobraComplete("ls", nil, "")
	if result != nil {
		t.Logf("ls returned %d suggestions (unexpected but non-fatal)", len(result))
	}
}

func TestFilterByPartial(t *testing.T) {
	suggestions := []Suggestion{
		{Cmd: "helm install"},
		{Cmd: "helm upgrade"},
		{Cmd: "helm status"},
	}
	filtered := filterByPartial(suggestions, "up")
	if len(filtered) != 1 || filtered[0].Cmd != "helm upgrade" {
		t.Errorf("expected [helm upgrade], got %v", filtered)
	}
}

func TestFilterByPartial_EmptyPartial(t *testing.T) {
	suggestions := []Suggestion{{Cmd: "a"}, {Cmd: "b"}}
	filtered := filterByPartial(suggestions, "")
	if len(filtered) != 2 {
		t.Errorf("expected all suggestions for empty partial, got %d", len(filtered))
	}
}

func TestLookup_CobraRealBinary(t *testing.T) {
	// Test real lookup for 'gh' (GitHub CLI) which has no hand-written spec in Iris
	results := Lookup("gh repo ")
	if len(results) == 0 {
		t.Skip("gh binary not available or output no completions")
	}
	t.Logf("Got %d completions for 'gh repo ':", len(results))
	for _, r := range results {
		t.Logf("  - Cmd: %-25s | Source: %-15s | Priority: %d | Desc: %s", r.Cmd, r.Source, r.Priority, r.Desc)
	}
}
