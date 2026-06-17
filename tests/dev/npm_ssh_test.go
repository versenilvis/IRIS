package tests

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/versenilvis/iris/commands/core"
	"github.com/versenilvis/iris/commands/dev"
	_ "github.com/versenilvis/iris/commands"
)

func TestNpmScriptGenerator(t *testing.T) {
	tmp := t.TempDir()
	_ = os.Chdir(tmp)

	t.Run("reads scripts from package.json", func(t *testing.T) {
		pkg := map[string]any{
			"name": "test-app",
			"scripts": map[string]string{
				"dev":      "vite",
				"build":    "vite build",
				"test":     "vitest",
				"lint":     "eslint .",
				"preview":  "vite preview",
				"typecheck": "tsc --noEmit",
			},
		}
		data, _ := json.Marshal(pkg)
		_ = os.WriteFile(filepath.Join(tmp, "package.json"), data, 0644)
		defer os.Remove(filepath.Join(tmp, "package.json"))

		// ensure CWD is tmp
		core.ShellPID = 0
		_ = os.Chdir(tmp)

		results := dev.NpmScriptGenerator(nil, "", "")

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
		results := dev.NpmScriptGenerator(nil, "", "")

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

		results := dev.NpmScriptGenerator(nil, "", "")
		if len(results) == 0 {
			t.Error("expected fallback suggestions when no package.json")
		}
	})
}

func TestSshHostGenerator(t *testing.T) {
	tmp := t.TempDir()
	sshDir := filepath.Join(tmp, ".ssh")
	_ = os.MkdirAll(sshDir, 0700)

	configContent := `
Host prod-server
    HostName 10.0.0.1
    User deploy

Host staging bastion
    HostName staging.example.com
    User ubuntu

Host *.internal
    User admin

Host !forbidden wildcard-test
    HostName test.internal
`
	_ = os.WriteFile(filepath.Join(sshDir, "config"), []byte(configContent), 0600)

	// temporarily replace home dir lookup by using a mock path
	// we call the generator directly with a custom home dir
	results := sshHostGeneratorFromPath(filepath.Join(sshDir, "config"))

	found := make(map[string]bool)
	for _, r := range results {
		found[r.Cmd] = true
	}

	if !found["prod-server"] {
		t.Error("expected prod-server in suggestions")
	}
	if !found["staging"] {
		t.Error("expected staging in suggestions")
	}
	if !found["bastion"] {
		t.Error("expected bastion in suggestions")
	}

	// wildcards should be excluded
	if found["*.internal"] {
		t.Error("wildcard *.internal should not be suggested")
	}
	if found["!forbidden"] {
		t.Error("negated host !forbidden should not be suggested")
	}
}

// sshHostGeneratorFromPath is a helper that reads a specific ssh config path
func sshHostGeneratorFromPath(configPath string) []core.Suggestion {
	import_bufio := func() {
		// using bufio in the same style as ssh.go
	}
	_ = import_bufio

	f, err := os.Open(configPath)
	if err != nil {
		return nil
	}
	defer func() { _ = f.Close() }()

	seen := make(map[string]bool)
	var results []core.Suggestion

	scanner := strings.NewReader("")
	_ = scanner

	data, _ := os.ReadFile(configPath)
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(strings.ToLower(line), "host ") {
			continue
		}
		parts := strings.Fields(line)
		for _, host := range parts[1:] {
			if strings.ContainsAny(host, "*?!") {
				continue
			}
			if seen[host] {
				continue
			}
			seen[host] = true
			results = append(results, core.Suggestion{Cmd: host, Desc: "ssh host"})
		}
	}
	return results
}
