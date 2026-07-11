package tests

import (
	"os"
	"path/filepath"
	"testing"

	_ "github.com/versenilvis/iris/commands/runner"
	"github.com/versenilvis/iris/spec"
)

// Verify that the just command generator parses recipes and returns nil on read errors
func TestJustGenerator(t *testing.T) {
	tmp := t.TempDir()
	content := []byte("# build project\nbuild:\n\techo build\n")
	_ = os.WriteFile(filepath.Join(tmp, "justfile"), content, 0644)

	oldWd, _ := os.Getwd()
	_ = os.Chdir(tmp)
	defer func() { _ = os.Chdir(oldWd) }()

	s := spec.Registry["just"]
	if s == nil || s.Generator == nil {
		t.Fatalf("expected just spec with generator to be registered in Registry")
	}

	res := s.Generator([]string{"just", ""}, "just ", "")
	if len(res) != 1 || res[0].Cmd != "build" || res[0].Desc != "build project" {
		t.Fatalf("expected recipe build with comment, got %v", res)
	}

	// Verify missing file returns nil
	_ = os.Remove(filepath.Join(tmp, "justfile"))
	resMissing := s.Generator([]string{"just", ""}, "just ", "")
	if resMissing != nil {
		t.Fatalf("expected nil when justfile cannot be read, got %v", resMissing)
	}
}
