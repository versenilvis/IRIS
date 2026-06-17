package dev

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/versenilvis/iris/commands/core"
)

func NpmScriptGenerator(tokens []string, _ string, _ string) []core.Suggestion {
	cwd := core.GetCWD()
	data, err := os.ReadFile(filepath.Join(cwd, "package.json"))
	if err != nil {
		return []core.Suggestion{
			{Cmd: "dev", Desc: "development server"},
			{Cmd: "build", Desc: "build for production"},
			{Cmd: "start", Desc: "start application"},
			{Cmd: "test", Desc: "run test suite"},
			{Cmd: "lint", Desc: "run linter"},
		}
	}

	var pkg struct {
		Scripts map[string]string `json:"scripts"`
	}
	if err := json.Unmarshal(data, &pkg); err != nil || len(pkg.Scripts) == 0 {
		return nil
	}

	priority := []string{"dev", "start", "build", "test", "lint", "preview", "typecheck", "format"}
	seen := make(map[string]bool)
	var results []core.Suggestion

	for _, name := range priority {
		if cmd, ok := pkg.Scripts[name]; ok {
			results = append(results, core.Suggestion{Cmd: name, Desc: cmd})
			seen[name] = true
		}
	}

	for name, cmd := range pkg.Scripts {
		if !seen[name] {
			results = append(results, core.Suggestion{Cmd: name, Desc: cmd})
		}
	}

	return results
}

func init() {
	core.Register(&core.Spec{
		Name:        "npm",
		Description: "node packages",
		Subcommands: []core.Subcommand{
			{Name: "install", Description: "install packages", Options: []core.Option{
				{Name: "--save-dev", Description: "save as devDependency"},
				{Name: "--save-exact", Description: "exact version"},
				{Name: "--legacy-peer-deps", Description: "skip peer deps"},
			}},
			{Name: "run", Description: "run script", Generator: NpmScriptGenerator},
			{Name: "test", Description: "run tests"},
			{Name: "init", Description: "create package.json", Options: []core.Option{
				{Name: "-y", Description: "skip prompts"},
			}},
			{Name: "publish", Description: "publish package"},
			{Name: "update", Description: "update packages"},
			{Name: "uninstall", Description: "remove package"},
			{Name: "ls", Description: "list installed"},
			{Name: "audit", Description: "security audit", Options: []core.Option{
				{Name: "fix", Description: "auto fix"},
			}},
			{Name: "ci", Description: "clean install from lockfile"},
			{Name: "pack", Description: "create tarball"},
			{Name: "link", Description: "symlink package"},
			{Name: "cache", Description: "manage cache", Subcommands: []core.Subcommand{
				{Name: "clean", Description: "clear cache"},
				{Name: "verify", Description: "verify cache"},
			}},
		},
		Options: []core.Option{
			{Name: "--prefix", Description: "set working directory"},
		},
	})
}
