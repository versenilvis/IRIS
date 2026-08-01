package spec

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/versenilvis/iris/spec/alias"
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

	partialResults := Lookup("demo ver")
	foundPartialOpt := false
	for _, r := range partialResults {
		if strings.Contains(r.Cmd, "--verbose") {
			foundPartialOpt = true
			break
		}
	}
	if !foundPartialOpt {
		t.Error("expected --verbose when partial query is 'ver' (trimmed dash match)")
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

func TestLookup_AliasFileGenerator(t *testing.T) {
	ResetRegistry()
	Register(&Spec{
		Name:      "nvim",
		Generator: FileGenerator(),
	})
	// alias "nv" -> "nvim", single token with trailing space
	ShellAliases = map[string]string{"nv": "nvim"}

	results := Lookup("nv ")
	if len(results) == 0 {
		t.Errorf("expected file suggestions for alias 'nv ' -> 'nvim', got none")
	}
	for _, r := range results {
		if !strings.HasPrefix(r.Cmd, "nvim ") {
			t.Errorf("expected suggestion to start with 'nvim ', got %q", r.Cmd)
		}
	}
}

func TestLookup_NvimFileGenerator(t *testing.T) {
	ResetRegistry()
	Register(&Spec{
		Name:      "nvim",
		Generator: FileGenerator(),
	})
	results := Lookup("nvim ")
	if len(results) == 0 {
		t.Errorf("expected file suggestions for 'nvim ', got none")
	}
}

func TestLookup_OptionAndFilePriority(t *testing.T) {
	ResetRegistry()
	Register(&Spec{
		Name:      "nvim",
		Generator: FileGenerator(),
		Options: []Option{
			{Name: "-c", Description: "Execute cmd"},
			{Name: "--cmd", Description: "Execute cmd before config"},
		},
	})

	// When query is 'nvim ', files should be prioritized over flags
	resultsEmpty := Lookup("nvim ")
	if len(resultsEmpty) == 0 {
		t.Fatalf("expected results for 'nvim ', got 0")
	}
	// first result should be a file or dir (Priority 50), not option (Priority 10)
	if strings.HasPrefix(resultsEmpty[0].Cmd, "nvim -") {
		t.Errorf("expected file/dir as top result for 'nvim ', got %q", resultsEmpty[0].Cmd)
	}

	// When query is 'nvim -', flags should be prioritized (Priority 80)
	resultsDash := Lookup("nvim -")
	if len(resultsDash) == 0 {
		t.Fatalf("expected results for 'nvim -', got 0")
	}
	if !strings.HasPrefix(resultsDash[0].Cmd, "nvim -") {
		t.Errorf("expected option as top result for 'nvim -', got %q", resultsDash[0].Cmd)
	}
}

func TestLookup_NestedDirectoryTrailingSpace(t *testing.T) {
	ResetRegistry()
	Register(&Spec{
		Name:      "cat",
		Generator: FileGenerator(),
	})

	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "subdir")
	_ = os.Mkdir(subDir, 0755)
	testFile := filepath.Join(subDir, "hello.txt")
	_ = os.WriteFile(testFile, []byte("hi"), 0644)

	query := "cat " + subDir + "/ "
	results := Lookup(query)
	if len(results) == 0 {
		t.Fatalf("expected results for %q, got 0", query)
	}
	found := false
	for _, r := range results {
		if strings.Contains(r.Cmd, "hello.txt") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected hello.txt in results, got %v", results)
	}
}

type mockGitProvider struct{}

func (m *mockGitProvider) ToolName() string { return "git" }
func (m *mockGitProvider) GetAliases(cwd string) []alias.AliasEntry {
	return []alias.AliasEntry{
		{
			Name:      "co",
			Expansion: "checkout",
			Scope:     "global",
		},
		{
			Name:      "recent",
			Expansion: "!git for-each-ref --sort=committerdate --format='%(committerdate:relative) %(refname:short)' refs/heads/ | tail -10",
			Scope:     "global",
		},
		{
			Name:      "local-alias",
			Expansion: "status",
			Scope:     "local",
		},
	}
}

func TestLookup_GitRecentAlias(t *testing.T) {
	ResetRegistry()
	Register(&Spec{
		Name: "git",
		Subcommands: []Subcommand{
			{Name: "commit"},
			{Name: "checkout", Options: []Option{{Name: "-b"}}},
		},
	})
	alias.Register(&mockGitProvider{})

	// Test Type 1: Standard subcommand alias 'git co' -> 'git checkout'
	resultsCo := Lookup("git co ")
	foundCo := false
	for _, r := range resultsCo {
		if strings.Contains(r.Cmd, "-b") {
			foundCo = true
			break
		}
	}
	if !foundCo {
		t.Errorf("expected '-b' option after expanding 'git co ', got %v", resultsCo)
	}

	// Test Type 2: Shell pipeline alias 'git rec' -> 'git recent'
	resultsRec := Lookup("git rec")
	foundRec := false
	for _, r := range resultsRec {
		if strings.Contains(r.Cmd, "git recent") {
			foundRec = true
			if r.Priority != 70 {
				t.Errorf("expected global alias priority 70, got %d", r.Priority)
			}
			break
		}
	}
	if !foundRec {
		t.Errorf("expected 'git recent' suggestion for 'git rec', got %v", resultsRec)
	}

	// Test Local Scope Priority (Priority = 85)
	resultsLocal := Lookup("git local")
	foundLocal := false
	for _, r := range resultsLocal {
		if strings.Contains(r.Cmd, "git local-alias") {
			foundLocal = true
			if r.Priority != 85 {
				t.Errorf("expected local alias priority 85, got %d", r.Priority)
			}
			break
		}
	}
	if !foundLocal {
		t.Errorf("expected 'git local-alias' suggestion, got %v", resultsLocal)
	}

	// Test Type 2 expansion with trailing space: 'git recent ' should expand without returning 0/nil
	resultsRecentSpace := Lookup("git recent ")
	if len(resultsRecentSpace) == 0 {
		t.Errorf("expected non-empty suggestions for 'git recent ', got 0 (possibly returned nil due to invalid subcommand/exclamation mark)")
	}
}

func TestLookup_RealGitProvider(t *testing.T) {
	ResetRegistry()
	Register(&Spec{
		Name: "git",
	})
	cwd, _ := os.Getwd()
	provider := &alias.GitProvider{}
	aliases := provider.GetAliases(cwd)
	t.Logf("Real Git aliases: %v", aliases)
	alias.Register(provider)

	results := Lookup("git rec")
	t.Logf("Lookup('git rec') results: %v", results)
	found := false
	for _, r := range results {
		if strings.Contains(r.Cmd, "recent") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected 'git recent' in results for real GitProvider, got %v", results)
	}
}

