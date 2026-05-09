package tests

import (
	"strings"
	"testing"

	"github.com/versenilvis/iris/commands/core"
)

func TestLookup(t *testing.T) {
	// Setup Registry
	core.Registry = make(map[string]*core.Spec)
	core.Register(&core.Spec{
		Name: "git",
		Subcommands: []core.Subcommand{
			{Name: "commit", Options: []core.Option{{Name: "--message"}}, MaxArgs: 1},
			{Name: "remote", Subcommands: []core.Subcommand{{Name: "add"}}},
		},
		Options: []core.Option{{Name: "--verbose"}},
	})

	// Setup Aliases
	core.ShellAliases = map[string]string{
		"gca": "git commit -a",
		"ta":  "tmux a -t",
	}

	tests := []struct {
		name        string
		input       string
		minResults  int
		mustContain string
	}{
		// REQUIREMENT: Token 1, no trailing space -> top-level suggestions
		{"Top-level", "gi", 1, "git"},
		// REQUIREMENT: Token 1, with trailing space -> subcommand suggestions
		{"Subcommand", "git ", 1, "git commit"},
		// REQUIREMENT: Alias expansion (gca -> git commit -a)
		{"Alias expansion", "gca", 1, "git commit -a"},
		// REQUIREMENT: Alias value with space (ta -> tmux a -t)
		{"Alias with space", "ta", 1, "tmux a -t"},
		// REQUIREMENT: Subcommand depth 2+ (git remote add)
		{"Deep subcommand", "git remote ", 1, "git remote add"},
		// REQUIREMENT: Option dedup (do not suggest --verbose if already typed)
		{"Option dedup", "git --verbose -", 0, ""}, 
		// REQUIREMENT: --flag=value does not count into argCount
		{"Flag with value ignore", "git --output=json ", 2, "git --output=json commit"},
		// REQUIREMENT: Unknown root command -> nil
		{"Unknown root command", "unknowncmd ", 0, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := core.Lookup(tt.input)
			if len(results) < tt.minResults {
				t.Errorf("Lookup(%q) got %d results; want at least %d", tt.input, len(results), tt.minResults)
			}
			if tt.mustContain != "" {
				found := false
				for _, r := range results {
					if strings.Contains(r.Cmd, tt.mustContain) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Lookup(%q) results did not contain %q", tt.input, tt.mustContain)
				}
			}
		})
	}
}
