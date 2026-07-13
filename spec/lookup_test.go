package spec

import (
	"strings"
	"sync"
	"testing"
)

func TestLookup(t *testing.T) {
	// Setup Registry
	Registry = make(map[string]*Spec)
	Register(&Spec{
		Name: "git",
		Subcommands: []Subcommand{
			{Name: "commit", Options: []Option{{Name: "--message"}}, MaxArgs: 1},
			{Name: "remote", Subcommands: []Subcommand{{Name: "add"}}},
		},
		Options: []Option{{Name: "--verbose"}},
	})

	// Setup Aliases
	ShellAliases = map[string]string{
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
			results := Lookup(tt.input)
			if len(results) < tt.minResults {
				t.Errorf("Lookup(%q) got %d results; want at least %d", tt.input, len(results), tt.minResults)
			}
			if tt.mustContain != "" {
				found := false
				for _, r := range results {
					if strings.Contains(r.Cmd, tt.mustContain) || strings.Contains(r.Desc, tt.mustContain) {
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

func TestLookup_NoFlagGateAndPriority(t *testing.T) {
	Registry["demo"] = &Spec{
		Name: "demo",
		Subcommands: []Subcommand{
			{Name: "sub1", Priority: 85},
		},
		Options: []Option{
			{Name: "--verbose", Priority: 70},
		},
	}

	results := Lookup("demo ")
	foundSub, foundOpt := false, false
	for _, r := range results {
		if strings.Contains(r.Cmd, "sub1") {
			foundSub = true
			if r.Priority != 85 {
				t.Errorf("expected subcommand Priority 85, got %d", r.Priority)
			}
		}
		if strings.Contains(r.Cmd, "--verbose") {
			foundOpt = true
			if r.Priority != 70 {
				t.Errorf("expected option Priority 70, got %d", r.Priority)
			}
		}
	}
	if !foundSub {
		t.Error("expected sub1 in lookup results")
	}
	if !foundOpt {
		t.Error("expected --verbose in lookup results even without typed dash")
	}
}

func TestLookupConcurrent(t *testing.T) {
	Registry = make(map[string]*Spec)
	Register(&Spec{
		Name: "git",
		Subcommands: []Subcommand{
			{Name: "commit", Options: []Option{{Name: "--message"}}, MaxArgs: 1},
		},
	})

	ShellAliases = map[string]string{
		"gca": "git commit -a",
	}

	var wg sync.WaitGroup
	const goroutines = 10
	const iterations = 50

	for range goroutines {
		wg.Go(func() {
			for range iterations {
				_ = Lookup("gca")
				_ = Lookup("git ")
			}
		})
	}
	wg.Wait()
}
