package tests

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/versenilvis/iris/commands/js"
	"github.com/versenilvis/iris/spec"
)

func TestNpmScriptGenerator(t *testing.T) {
	tmp := t.TempDir()
	_ = os.Chdir(tmp)

	t.Run("reads scripts from package.json", func(t *testing.T) {
		pkg := map[string]any{
			"name": "test-app",
			"scripts": map[string]string{
				"dev":       "vite",
				"build":     "vite build",
				"test":      "vitest",
				"lint":      "eslint .",
				"preview":   "vite preview",
				"typecheck": "tsc --noEmit",
			},
		}
		data, _ := json.Marshal(pkg)
		_ = os.WriteFile(filepath.Join(tmp, "package.json"), data, 0644)
		defer os.Remove(filepath.Join(tmp, "package.json"))

		// ensure CWD is tmp
		spec.ShellPID = 0
		_ = os.Chdir(tmp)

		results := js.NpmScriptGenerator(nil, "", "")

		found := make(map[string]bool)
		for _, r := range results {
			found[r.Cmd] = true
		}

		for _, expected := range []string{"dev", "build", "test", "lint", "preview", "typecheck"} {
			if !found[expected] {
				t.Errorf("expected script '%s' in suggestions", expected)
			}
		}
	})

	t.Run("priority scripts come first", func(t *testing.T) {
		pkg := map[string]any{
			"scripts": map[string]string{
				"zzz-last":  "echo last",
				"dev":       "vite",
				"aaa-first": "echo first",
				"build":     "vite build",
			},
		}
		data, _ := json.Marshal(pkg)
		_ = os.WriteFile(filepath.Join(tmp, "package.json"), data, 0644)
		defer os.Remove(filepath.Join(tmp, "package.json"))

		_ = os.Chdir(tmp)
		results := js.NpmScriptGenerator(nil, "", "")

		if len(results) < 2 {
			t.Fatal("expected at least 2 results")
		}

		// dev should appear before zzz-last
		devIdx, zzzIdx := -1, -1
		for i, r := range results {
			if r.Cmd == "dev" {
				devIdx = i
			}
			if r.Cmd == "zzz-last" {
				zzzIdx = i
			}
		}
		if devIdx == -1 {
			t.Error("dev not found")
		}
		if zzzIdx == -1 {
			t.Error("zzz-last not found")
		}
		if devIdx > zzzIdx {
			t.Errorf("'dev' (idx %d) should come before 'zzz-last' (idx %d)", devIdx, zzzIdx)
		}
	})

	t.Run("fallback when no package.json", func(t *testing.T) {
		emptyDir := t.TempDir()
		_ = os.Chdir(emptyDir)
		defer func() { _ = os.Chdir(tmp) }()

		results := js.NpmScriptGenerator(nil, "", "")
		if len(results) == 0 {
			t.Error("expected fallback suggestions when no package.json")
		}
	})
}
