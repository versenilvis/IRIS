package core_test

import (
	"testing"

	"github.com/versenilvis/iris/commands/core"
	"github.com/versenilvis/iris/integration/shell"
)

type mockAdapter struct {
	shell.BashAdapter
}

func (m *mockAdapter) ScanAliases() map[string]string {
	return map[string]string{
		"gca": "git commit -a",
		"ta":  "tmux a -t",
	}
}

func TestLookup(t *testing.T) {
	// Use mock adapter
	shell.Current = &mockAdapter{}
	
	core.Register(&core.Spec{
		Name:        "git",
		Description: "git command",
		Subcommands: []core.Subcommand{
			{Name: "commit", Description: "commit changes"},
			{Name: "remote", Description: "manage remotes", Subcommands: []core.Subcommand{
				{Name: "add", Description: "add remote"},
			}},
		},
		Options: []core.Option{
			{Name: "--verbose", Description: "verbose output"},
		},
	})

	tests := []struct {
		name     string
		input    string
		minCount int
		checkCmd string
	}{
		{"Top-level suggestions", "gi", 1, "git"},
		{"Subcommand suggestions", "git ", 2, "git commit"},
		{"Alias expansion", "gca", 1, "git commit -a"},
		{"Deep subcommand", "git remote ", 1, "git remote add"},
		{"Option dedup", "git --verbose -", 0, ""}, 
		{"Flag with value ignore", "git --output=json ", 2, "git --output=json commit"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := core.Lookup(tt.input)
			if len(results) < tt.minCount {
				t.Errorf("Lookup(%q) returned %d results; want at least %d", tt.input, len(results), tt.minCount)
			}
			if tt.checkCmd != "" {
				found := false
				for _, r := range results {
					if r.Cmd == tt.checkCmd {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Lookup(%q) did not suggest %q", tt.input, tt.checkCmd)
				}
			}
		})
	}
}
