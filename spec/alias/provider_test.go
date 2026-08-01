package alias

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCargoProvider(t *testing.T) {
	tempDir := t.TempDir()
	cargoDir := filepath.Join(tempDir, ".cargo")
	if err := os.Mkdir(cargoDir, 0755); err != nil {
		t.Fatal(err)
	}

	configToml := `
[alias]
b = "build"
c = "check"
t = ["test", "--", "--nocapture"]
`
	if err := os.WriteFile(filepath.Join(cargoDir, "config.toml"), []byte(configToml), 0644); err != nil {
		t.Fatal(err)
	}

	p := &CargoProvider{}
	aliases := p.parse(tempDir)

	expected := map[string]string{
		"b": "build",
		"c": "check",
		"t": "test -- --nocapture",
	}

	if len(aliases) != 3 {
		t.Fatalf("expected 3 aliases, got %d", len(aliases))
	}

	for _, a := range aliases {
		if expected[a.Name] != a.Expansion {
			t.Errorf("expected %s to map to %s, got %s", a.Name, expected[a.Name], a.Expansion)
		}
	}
}

func TestGitProviderParseOutput(t *testing.T) {
	outputWithScope := []byte("global\talias.recent !git for-each-ref --sort=committerdate --format=\"%(committerdate:relative) %(refname:short)\" refs/heads/ | tail -10\nlocal\talias.co checkout\n")
	p := &GitProvider{}
	entries := p.parseOutput(outputWithScope, true)

	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d: %v", len(entries), entries)
	}

	if entries[0].Name != "recent" {
		t.Errorf("expected name 'recent', got %q", entries[0].Name)
	}
	if entries[0].Scope != "global" {
		t.Errorf("expected scope 'global', got %q", entries[0].Scope)
	}
	if entries[0].Expansion != "!git for-each-ref --sort=committerdate --format=\"%(committerdate:relative) %(refname:short)\" refs/heads/ | tail -10" {
		t.Errorf("unexpected expansion: %q", entries[0].Expansion)
	}

	if entries[1].Name != "co" {
		t.Errorf("expected name 'co', got %q", entries[1].Name)
	}
	if entries[1].Scope != "local" {
		t.Errorf("expected scope 'local', got %q", entries[1].Scope)
	}
	if entries[1].Expansion != "checkout" {
		t.Errorf("expected expansion 'checkout', got %q", entries[1].Expansion)
	}
}

func TestGetProvider_UnregisteredCommand(t *testing.T) {
	p := GetProvider("unregistered_command_xyz")
	if p != nil {
		t.Errorf("expected nil for unregistered command, got %v", p)
	}
}

func TestResolveGitConfigPath_WalkUp(t *testing.T) {
	tempDir := t.TempDir()
	gitDir := filepath.Join(tempDir, ".git")
	_ = os.Mkdir(gitDir, 0755)
	configPath := filepath.Join(gitDir, "config")
	_ = os.WriteFile(configPath, []byte("[alias]\nst = status"), 0644)

	subDir := filepath.Join(tempDir, "a", "b", "c")
	_ = os.MkdirAll(subDir, 0755)

	resolved := resolveGitConfigPath(subDir)
	if resolved != configPath {
		t.Errorf("expected %q, got %q", configPath, resolved)
	}
}

func TestResolveGitConfigPath_WorktreeFile(t *testing.T) {
	tempDir := t.TempDir()
	mainGitDir := filepath.Join(tempDir, "mainrepo", ".git")
	_ = os.MkdirAll(mainGitDir, 0755)
	mainConfig := filepath.Join(mainGitDir, "config")
	_ = os.WriteFile(mainConfig, []byte("[alias]\nst = status"), 0644)

	wtDir := filepath.Join(tempDir, "wt")
	_ = os.MkdirAll(wtDir, 0755)
	wtGitFile := filepath.Join(wtDir, ".git")
	_ = os.WriteFile(wtGitFile, []byte("gitdir: "+mainGitDir), 0644)

	wtSubDir := filepath.Join(wtDir, "sub")
	_ = os.MkdirAll(wtSubDir, 0755)

	resolved := resolveGitConfigPath(wtSubDir)
	if resolved != mainConfig {
		t.Errorf("expected %q, got %q", mainConfig, resolved)
	}
}
