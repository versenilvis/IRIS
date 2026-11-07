package synthesize

import (
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/versenilvis/iris/spec"
)

func TestSynthesize_Combined(t *testing.T) {
	manDir, err := filepath.Abs("testdata/man1")
	if err != nil {
		t.Fatal(err)
	}
	bashDir, err := filepath.Abs("testdata/completions")
	if err != nil {
		t.Fatal(err)
	}

	synth := New(
		WithManDir(manDir),
		WithBashCompletionDir(bashDir),
	)

	result, err := synth.Synthesize("toycmd")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Name != "toycmd" {
		t.Errorf("got name %q, want 'toycmd'", result.Name)
	}
	if result.Description != "a demonstration command for iris" {
		t.Errorf("got description %q", result.Description)
	}

	rawJSON, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("failed to marshal JSON: %v", err)
	}

	specs, err := spec.LoadUserSpecBytes(rawJSON, "")
	if err != nil {
		t.Fatalf("LoadUserSpecBytes failed on synthesized spec: %v", err)
	}
	if len(specs) != 1 {
		t.Fatalf("expected 1 spec, got %d", len(specs))
	}
	loaded := specs[0]
	if loaded.Name != "toycmd" {
		t.Errorf("loaded spec name mismatch: got %s", loaded.Name)
	}

	options := make(map[string]spec.Option)
	for _, o := range loaded.Options {
		options[o.Name] = o
	}

	if _, exists := options["--file"]; !exists {
		t.Error("expected --file in loaded spec")
	}
	if _, exists := options["-f"]; !exists {
		t.Error("expected alias -f in loaded spec options")
	}
	if _, exists := options["--dry-run"]; !exists {
		t.Error("expected --dry-run from bash in loaded spec")
	}
	if _, exists := options["-q"]; !exists {
		t.Error("expected -q from man in loaded spec")
	}
}

func TestSynthesize_MissingCommand(t *testing.T) {
	manDir, err := filepath.Abs("testdata/man1")
	if err != nil {
		t.Fatal(err)
	}
	bashDir, err := filepath.Abs("testdata/completions")
	if err != nil {
		t.Fatal(err)
	}

	synth := New(
		WithManDir(manDir),
		WithBashCompletionDir(bashDir),
	)

	_, err = synth.Synthesize("nonexistent_command_12345")
	if err == nil {
		t.Fatal("expected error for nonexistent command")
	}
}

func TestSynthesize_InvalidCommandName(t *testing.T) {
	synth := New()
	for _, invalid := range []string{"", "../bad", "foo/bar", "with space"} {
		_, err := synth.Synthesize(invalid)
		if err == nil {
			t.Errorf("expected error for invalid command %q", invalid)
		}
	}
}
