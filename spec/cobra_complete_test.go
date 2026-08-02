package spec

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
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
		t.Fatalf("expected nil for non-Cobra binary 'ls', got %v", result)
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

func TestBuildCobraCacheKey(t *testing.T) {
	key1 := buildCobraCacheKey("gh", []string{"repo", "foo bar"}, "baz")
	key2 := buildCobraCacheKey("gh", []string{"repo", "foo", "bar"}, "baz")
	key3 := buildCobraCacheKey("gh", []string{"repo", "foo bar"}, "other")

	if key1 == key2 {
		t.Errorf("expected quoted argument key1 and key2 to be distinct, but were equal")
	}
	if key1 == key3 {
		t.Errorf("expected key1 and key3 with different partials to be distinct, but were equal")
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

func TestLookup_CobraKubectl(t *testing.T) {
	results := Lookup("kubectl get ")
	if len(results) == 0 {
		t.Skip("kubectl binary not available or output no completions")
	}
	t.Logf("Got %d completions for 'kubectl get ':", len(results))
	for i, r := range results {
		if i >= 10 {
			t.Logf("  ... and %d more", len(results)-10)
			break
		}
		t.Logf("  - Cmd: %-30s | Source: %-15s | Priority: %d | Desc: %s", r.Cmd, r.Source, r.Priority, r.Desc)
	}
}

func TestNewProbeCmd_Setsid(t *testing.T) {
	cmd := newProbeCmd(context.Background(), "true", nil)
	if cmd.SysProcAttr == nil || !cmd.SysProcAttr.Setsid {
		t.Fatalf("expected probe command to set Setsid: true, got %+v", cmd.SysProcAttr)
	}
}

// TestNewProbeCmd_NoControllingTerminal writes a small script
// and runs it to check
// if newProbeCmd actually denies it access to tty
func TestNewProbeCmd_NoControllingTerminal(t *testing.T) {
	dir := t.TempDir()
	script := "#!/bin/sh\nexec 3<>/dev/tty"

	binPath := filepath.Join(dir, "ttyprobe")
	if err := os.WriteFile(binPath, []byte(script), 0o755); err != nil {
		t.Fatalf("failed to write probe script: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	runErr := newProbeCmd(ctx, binPath, nil).Run()

	var exitErr *exec.ExitError
	if !errors.As(runErr, &exitErr) {
		t.Fatalf("expected probe script to exit nonzero after failing to open /dev/tty, got %v", runErr)
	}
}

func TestLookup_CobraGolangciLint(t *testing.T) {
	results := Lookup("golangci-lint ")
	if len(results) == 0 {
		t.Skip("golangci-lint binary not available or output no completions")
	}
	t.Logf("Got %d completions for 'golangci-lint ':", len(results))
	for _, r := range results {
		t.Logf("  - Cmd: %-30s | Source: %-15s | Priority: %d | Desc: %s", r.Cmd, r.Source, r.Priority, r.Desc)
	}
}
