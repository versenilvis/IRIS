package js

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/versenilvis/iris/spec"
)

func NpmScriptGenerator(tokens []string, _ string, _ string) []spec.Suggestion {
	cwd := spec.GetCWD()
	data, err := os.ReadFile(filepath.Join(cwd, "package.json"))
	if err != nil {
		return []spec.Suggestion{
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
	var results []spec.Suggestion

	for _, name := range priority {
		if cmd, ok := pkg.Scripts[name]; ok {
			results = append(results, spec.Suggestion{Cmd: name, Desc: cmd})
			seen[name] = true
		}
	}

	for name, cmd := range pkg.Scripts {
		if !seen[name] {
			results = append(results, spec.Suggestion{Cmd: name, Desc: cmd})
		}
	}

	return results
}

func init() {
	spec.Register(&spec.Spec{
		Name:        "npm",
		Description: "node packages",
		Subcommands: []spec.Subcommand{
			{Name: "install", Description: "install packages", Options: []spec.Option{
				{Name: "--save-dev", Description: "save as devDependency"},
				{Name: "--save-exact", Description: "exact version"},
				{Name: "--legacy-peer-deps", Description: "skip peer deps"},
			}},
			{Name: "run", Description: "run script", Generator: NpmScriptGenerator},
			{Name: "test", Description: "run tests"},
			{Name: "init", Description: "create package.json", Options: []spec.Option{
				{Name: "-y", Description: "skip prompts"},
			}},
			{Name: "publish", Description: "publish package"},
			{Name: "update", Description: "update packages"},
			{Name: "uninstall", Description: "remove package"},
			{Name: "ls", Description: "list installed"},
			{Name: "audit", Description: "security audit", Options: []spec.Option{
				{Name: "fix", Description: "auto fix"},
			}},
			{Name: "ci", Description: "clean install from lockfile"},
			{Name: "pack", Description: "create tarball"},
			{Name: "link", Description: "symlink package"},
			{Name: "cache", Description: "manage cache", Subcommands: []spec.Subcommand{
				{Name: "clean", Description: "clear cache"},
				{Name: "verify", Description: "verify cache"},
			}},
		},
		Options: []spec.Option{
			{Name: "--prefix", Description: "set working directory"},
		},
	})
}
