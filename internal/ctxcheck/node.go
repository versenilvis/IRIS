package ctxcheck

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

var nodePackageManagers = map[string]bool{
	"npm":  true,
	"pnpm": true,
	"yarn": true,
	"bun":  true,
}

var nodeBuiltins = map[string]bool{
	"install":   true,
	"i":         true,
	"add":       true,
	"ci":        true,
	"publish":   true,
	"pack":      true,
	"init":      true,
	"create":    true,
	"link":      true,
	"unlink":    true,
	"outdated":  true,
	"update":    true,
	"up":        true,
	"upgrade":   true,
	"audit":     true,
	"cache":     true,
	"config":    true,
	"help":      true,
	"version":   true,
	"login":     true,
	"logout":    true,
	"whoami":    true,
	"ping":      true,
	"doctor":    true,
	"exec":      true,
	"remove":    true,
	"rm":        true,
	"uninstall": true,
	"un":        true,
	"info":      true,
	"view":      true,
	"rebuild":   true,
	"prune":     true,
	"dedupe":    true,
	"why":       true,
	"explain":   true,
	"pm":        true,
	"set":       true,
	"get":       true,
}

type packageJSON struct {
	Scripts map[string]string `json:"scripts"`
}

func findNearestPackageJSON(cwd string) (string, error) {
	dir := filepath.Clean(cwd)
	for {
		p := filepath.Join(dir, "package.json")
		if fi, err := os.Stat(p); err == nil && !fi.IsDir() {
			return p, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", os.ErrNotExist
}

func readPackageScripts(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var pkg packageJSON
	if err := json.Unmarshal(data, &pkg); err != nil {
		return nil, err
	}
	return pkg.Scripts, nil
}

func ValidateNode(tokens []string, cwd string) Verdict {
	if len(tokens) == 0 {
		return Free
	}

	tool := tokens[0]
	if tool == "npx" || tool == "bunx" {
		return Free
	}
	if !nodePackageManagers[tool] {
		return Free
	}

	args := tokens[1:]

	// pnpm dlx is free
	if tool == "pnpm" && len(args) > 0 && args[0] == "dlx" {
		return Free
	}

	// ignore anything after --
	for idx, a := range args {
		if a == "--" {
			args = args[:idx]
			break
		}
	}

	var scriptName string
	isExplicitRun := false

	for _, a := range args {
		if a == "-w" || a == "--workspace" || a == "--workspaces" ||
			a == "--filter" || a == "--prefix" || a == "-C" || a == "--cwd" {
			return Unknown
		}
	}

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if strings.HasPrefix(arg, "-") {
			continue
		}

		if nodeBuiltins[arg] {
			return Free
		}

		if arg == "run" || arg == "run-script" {
			isExplicitRun = true
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				scriptName = args[i+1]
			}
			break
		}

		// bun test has a native test runner that does not require package.json
		if tool == "bun" && arg == "test" {
			return Free
		}

		// standard short script forms
		if arg == "test" || arg == "start" || arg == "stop" || arg == "restart" || (tool == "npm" && arg == "t") {
			scriptName = arg
			if arg == "t" {
				scriptName = "test"
			}
			isExplicitRun = true
			break
		}

		// for pnpm, yarn, bun: bare word might be a script or custom subcommand
		if tool != "npm" {
			scriptName = arg
			break
		}

		// for npm: unknown bare subcommand
		return Unknown
	}

	if scriptName == "" && !isExplicitRun {
		// bare npm / pnpm / yarn / bun without script argument
		return Free
	}

	path, err := getCachedPath("node", cwd, findNearestPackageJSON)
	if err != nil {
		return Invalid
	}

	data, err := getCachedManifest(path, func(p string) (any, error) {
		return readPackageScripts(p)
	})
	if err != nil {
		return Unknown
	}
	scripts, ok := data.(map[string]string)
	if !ok || scripts == nil {
		return Unknown
	}

	if _, ok := scripts[scriptName]; ok {
		return Valid
	}

	// npm with missing script is Invalid
	if tool == "npm" || isExplicitRun {
		return Invalid
	}

	// for pnpm/yarn/bun: bare word not in scripts and not built-in -> Unknown
	return Unknown
}
