package scoring

import (
	"testing"

	"github.com/versenilvis/iris/internal/workspace"
)

func TestApplyContextRules_MultiEcosystem(t *testing.T) {
	tests := []struct {
		name     string
		ws       workspace.WorkspaceInfo
		cmd      string
		expected int
	}{
		{
			name:     "git status inside git repo",
			ws:       workspace.WorkspaceInfo{HasGit: true},
			cmd:      "git status -s",
			expected: 40,
		},
		{
			name:     "git init inside git repo penalized",
			ws:       workspace.WorkspaceInfo{HasGit: true},
			cmd:      "git init",
			expected: -50,
		},
		{
			name:     "bun run dev inside node project",
			ws:       workspace.WorkspaceInfo{HasNodeProject: true},
			cmd:      "bun run dev",
			expected: 50,
		},
		{
			name:     "go test inside go project",
			ws:       workspace.WorkspaceInfo{HasGoProject: true},
			cmd:      "go test ./...",
			expected: 50,
		},
		{
			name:     "cargo check inside rust project",
			ws:       workspace.WorkspaceInfo{HasRustProject: true},
			cmd:      "cargo check",
			expected: 50,
		},
		{
			name:     "pytest inside python project",
			ws:       workspace.WorkspaceInfo{HasPythonProject: true},
			cmd:      "pytest -v",
			expected: 50,
		},
		{
			name:     "just build inside justfile project",
			ws:       workspace.WorkspaceInfo{HasJustfile: true},
			cmd:      "just build",
			expected: 50,
		},
		{
			name:     "kubectl get pods inside k8s workspace",
			ws:       workspace.WorkspaceInfo{HasK8s: true},
			cmd:      "kubectl get pods",
			expected: 40,
		},
		{
			name:     "unrelated command gets no bonus",
			ws:       workspace.WorkspaceInfo{HasGit: true},
			cmd:      "echo hello",
			expected: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := ApplyContextRules(tc.ws, tc.cmd)
			if got != tc.expected {
				t.Errorf("expected bonus %d, got %d for command %q", tc.expected, got, tc.cmd)
			}
		})
	}
}

func TestApplyContextRules_Clamping(t *testing.T) {
	ws := workspace.WorkspaceInfo{HasGit: true, HasNodeProject: true, HasMakefile: true}
	// Create rules that sum above 100 and below -100
	rules := []ContextRule{
		&SimpleContextRule{check: func(w workspace.WorkspaceInfo, c string) bool { return true }, bonus: 80},
		&SimpleContextRule{check: func(w workspace.WorkspaceInfo, c string) bool { return true }, bonus: 60},
	}
	got := ApplyCustomContextRules(ws, "any cmd", rules)
	if got != 100 {
		t.Errorf("expected clamp to 100, got %d", got)
	}

	negRules := []ContextRule{
		&SimpleContextRule{check: func(w workspace.WorkspaceInfo, c string) bool { return true }, bonus: -80},
		&SimpleContextRule{check: func(w workspace.WorkspaceInfo, c string) bool { return true }, bonus: -60},
	}
	gotNeg := ApplyCustomContextRules(ws, "any cmd", negRules)
	if gotNeg != -100 {
		t.Errorf("expected clamp to -100, got %d", gotNeg)
	}
}
