package tests

import (
	"strings"
	"sync"
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
		{"Top-level", "gi", 1, "git"},
		{"Subcommand", "git ", 1, "git commit"},
		{"Alias expansion", "gca", 1, "git commit -a"},
		{"Alias with space", "ta", 1, "tmux a -t"},
		{"Deep subcommand", "git remote ", 1, "git remote add"},
		{"Option dedup", "git --verbose -", 0, ""},
		{"Flag with value ignore", "git --output=json ", 2, "git --output=json commit"},
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

func TestLookupConcurrent(t *testing.T) {
	core.Registry = make(map[string]*core.Spec)
	core.Register(&core.Spec{
		Name: "git",
		Subcommands: []core.Subcommand{
			{Name: "commit", Options: []core.Option{{Name: "--message"}}, MaxArgs: 1},
		},
	})

	core.ShellAliases = map[string]string{
		"gca": "git commit -a",
	}

	var wg sync.WaitGroup
	const goroutines = 10
	const iterations = 50

	for range goroutines {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range iterations {
				_ = core.Lookup("gca")
				_ = core.Lookup("git ")
			}
		}()
	}
	wg.Wait()
}
