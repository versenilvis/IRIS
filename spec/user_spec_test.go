package spec

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestCleanJSONC(t *testing.T) {
	input := []byte(`{
		// line comment
		"name": "myapp",
		/* block comment */
		"url": "https://example.com/api",
		"items": [
			"one",
			"two", // inline comment
		],
		"nested": {
			"quote": "hello \"world\" // not a comment",
			"trailing": true,
		},
	}`)

	cleaned := CleanJSONC(input)

	var parsed map[string]any
	if err := json.Unmarshal(cleaned, &parsed); err != nil {
		t.Fatalf("CleanJSONC produced invalid JSON: %v\noutput:\n%s", err, string(cleaned))
	}

	if parsed["name"] != "myapp" {
		t.Errorf("expected name to be myapp, got %v", parsed["name"])
	}
	if parsed["url"] != "https://example.com/api" {
		t.Errorf("expected url to be preserved, got %v", parsed["url"])
	}
}

func TestLoadUserSpecBytes_Single(t *testing.T) {
	data := []byte(`{
		"name": "deployer",
		"description": "Deployment tool",
		"options": [
			{"name": "--dry-run", "description": "Simulate run"}
		]
	}`)

	specs, err := LoadUserSpecBytes(data, "fallback")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(specs) != 1 {
		t.Fatalf("expected 1 spec, got %d", len(specs))
	}
	if specs[0].Name != "deployer" {
		t.Errorf("expected name deployer, got %s", specs[0].Name)
	}
}

func TestLoadUserSpecBytes_DefaultName(t *testing.T) {
	data := []byte(`{
		"description": "Tool without explicit name",
		"generator": "file"
	}`)

	specs, err := LoadUserSpecBytes(data, "custom-cmd")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(specs) != 1 {
		t.Fatalf("expected 1 spec, got %d", len(specs))
	}
	if specs[0].Name != "custom-cmd" {
		t.Errorf("expected fallback name custom-cmd, got %s", specs[0].Name)
	}
	if specs[0].Generator == nil {
		t.Error("expected file generator to be resolved")
	}
}

func TestLoadUserSpecBytes_ArrayAndWrapper(t *testing.T) {
	arrayData := []byte(`[
		{"name": "tool1", "description": "First"},
		{"name": "tool2", "description": "Second"}
	]`)

	specs, err := LoadUserSpecBytes(arrayData, "")
	if err != nil {
		t.Fatalf("unexpected error parsing array: %v", err)
	}
	if len(specs) != 2 {
		t.Fatalf("expected 2 specs, got %d", len(specs))
	}

	wrapperData := []byte(`{
		"$schema": "./schema.json",
		"specs": [
			{"name": "wrapped1"},
			{"name": "wrapped2"}
		]
	}`)

	wrappedSpecs, err := LoadUserSpecBytes(wrapperData, "")
	if err != nil {
		t.Fatalf("unexpected error parsing wrapper: %v", err)
	}
	if len(wrappedSpecs) != 2 {
		t.Fatalf("expected 2 specs, got %d", len(wrappedSpecs))
	}
}

func TestUserSpec_RecursiveSubcommands(t *testing.T) {
	data := []byte(`{
		"name": "cluster",
		"subcommands": [
			{
				"name": "node",
				"subcommands": [
					{
						"name": "drain",
						"options": [
							{"name": "--force", "aliases": ["-f"], "description": "Force drain"}
						],
						"generator": "dir"
					}
				]
			}
		]
	}`)

	specs, err := LoadUserSpecBytes(data, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(specs) != 1 {
		t.Fatalf("expected 1 spec, got %d", len(specs))
	}

	cmd := specs[0]
	if len(cmd.Subcommands) != 1 {
		t.Fatalf("expected 1 child subcommand, got %d", len(cmd.Subcommands))
	}
	node := cmd.Subcommands[0]
	if node.Name != "node" || len(node.Subcommands) != 1 {
		t.Fatalf("unexpected node subcommand: %+v", node)
	}
	drain := node.Subcommands[0]
	if drain.Name != "drain" {
		t.Errorf("expected drain subcommand, got %s", drain.Name)
	}
	if len(drain.Options) != 2 {
		t.Fatalf("expected 2 options (including alias), got %d", len(drain.Options))
	}
	if drain.Generator == nil {
		t.Error("expected dir generator on leaf subcommand")
	}
}

func TestLoadUserSpecs_DirectoryAndReload(t *testing.T) {
	tempDir := t.TempDir()

	file1 := filepath.Join(tempDir, "first.json")
	if err := os.WriteFile(file1, []byte(`{"name": "first_cmd", "description": "First"}`), 0600); err != nil {
		t.Fatal(err)
	}

	file2 := filepath.Join(tempDir, "second.jsonc")
	if err := os.WriteFile(file2, []byte(`{
		// comment
		"name": "second_cmd",
		"aliases": ["sc2"],
		"description": "Second",
	}`), 0600); err != nil {
		t.Fatal(err)
	}

	ignoreFile := filepath.Join(tempDir, "ignore.txt")
	if err := os.WriteFile(ignoreFile, []byte(`{"name": "ignored"}`), 0600); err != nil {
		t.Fatal(err)
	}

	if err := LoadUserSpecs(tempDir); err != nil {
		t.Fatalf("LoadUserSpecs failed: %v", err)
	}

	if Registry["first_cmd"] == nil {
		t.Error("expected first_cmd to be registered")
	}
	if Registry["second_cmd"] == nil {
		t.Error("expected second_cmd to be registered")
	}
	if Registry["sc2"] == nil {
		t.Error("expected alias sc2 to be registered")
	}
	if Registry["ignored"] != nil {
		t.Error("ignored.txt should not be registered")
	}

	// test reload after file removal
	_ = os.Remove(file1)
	if err := LoadUserSpecs(tempDir); err != nil {
		t.Fatalf("LoadUserSpecs reload failed: %v", err)
	}
	if Registry["first_cmd"] != nil {
		t.Error("first_cmd should have been deregistered after file removal")
	}
	if Registry["second_cmd"] == nil {
		t.Error("second_cmd should remain registered")
	}
}
